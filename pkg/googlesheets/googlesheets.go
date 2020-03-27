package googlesheets

import (
	"fmt"
	"strings"
	"time"

	"context"

	"github.com/araddon/dateparse"
	"github.com/davecgh/go-spew/spew"
	cd "github.com/grafana/google-sheets-datasource/pkg/googlesheets/columndefinition"
	gc "github.com/grafana/google-sheets-datasource/pkg/googlesheets/googleclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/patrickmn/go-cache"
	"google.golang.org/api/sheets/v4"
)

// GoogleSheets provides an interface to the Google Sheets API.
type GoogleSheets struct {
	Cache  *cache.Cache
	Logger log.Logger
}

// Query queries a spreadsheet and returns a corresponding data frame.
func (gs *GoogleSheets) Query(ctx context.Context, refID string, qm *QueryModel, config *GoogleSheetConfig, timeRange backend.TimeRange) (*data.Frame, error) {
	client, err := gc.New(ctx, gc.NewAuth(config.APIKey, config.AuthType, config.JWT))
	if err != nil {
		return nil, fmt.Errorf("unable to create Google API client: %w", err)
	}

	// This result may be cached
	data, meta, err := gs.getSheetData(client, qm)
	if err != nil {
		return nil, err
	}

	frame, err := gs.transformSheetToDataFrame(data, meta, refID, qm)
	if frame != nil && qm.UseTimeFilter {
		timeIndex := findTimeField(frame)
		if timeIndex >= 0 {
			frame, err = frame.FilterRowsByField(timeIndex, func(i interface{}) (bool, error) {
				val, ok := i.(*time.Time)
				if !ok {
					return false, fmt.Errorf("invalid time column: %s", spew.Sdump(i))
				}
				if val == nil || val.Before(timeRange.From) || val.After(timeRange.To) {
					return false, nil
				}
				return true, nil
			})
		}
	}
	return frame, err
}

// GetSpreadsheets gets spreadsheets from the Google API.
func (gs *GoogleSheets) GetSpreadsheets(ctx context.Context, config *GoogleSheetConfig) (map[string]string, error) {
	client, err := gc.New(ctx, gc.NewAuth(config.APIKey, config.AuthType, config.JWT))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	files, err := client.GetSpreadsheetFiles()
	if err != nil {
		return nil, err
	}

	fileNames := map[string]string{}
	for _, i := range files {
		fileNames[i.Id] = i.Name
	}

	return fileNames, nil
}

// getSheetData gets grid data corresponding to a spreadsheet.
func (gs *GoogleSheets) getSheetData(client client, qm *QueryModel) (*sheets.GridData, map[string]interface{}, error) {
	cacheKey := qm.Spreadsheet + qm.Range
	if item, expires, found := gs.Cache.GetWithExpiration(cacheKey); found && qm.CacheDurationSeconds > 0 {
		return item.(*sheets.GridData), map[string]interface{}{
			"hit":     true,
			"expires": expires.Unix(),
		}, nil
	}

	result, err := client.GetSpreadsheet(qm.Spreadsheet, qm.Range, true)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get spreadsheet: %w", err)
	}

	if result.Properties.TimeZone != "" {
		loc, err := time.LoadLocation(result.Properties.TimeZone)
		if err != nil {
			return nil, nil, fmt.Errorf("error while loading timezone: %w", err)
		}
		time.Local = loc
	}

	data := result.Sheets[0].Data[0]
	if qm.CacheDurationSeconds > 0 {
		gs.Cache.Set(cacheKey, data, time.Duration(qm.CacheDurationSeconds)*time.Second)
	}

	return data, map[string]interface{}{"hit": false}, nil
}

func getExcelColumnName(columnNumber int) string {
	dividend := columnNumber
	columnName := ""
	var modulo int

	for dividend > 0 {
		modulo = ((dividend - 1) % 26)
		fmt.Printf("MOD %d\n", modulo)
		columnName = string(65+modulo) + columnName
		dividend = ((dividend - modulo) / 26)
	}

	return columnName
}

func (gs *GoogleSheets) transformSheetToDataFrame(sheet *sheets.GridData, meta map[string]interface{}, refID string, qm *QueryModel) (*data.Frame, error) {
	fields := []*data.Field{}
	columns, start := getColumnDefinitions(sheet.RowData)
	warnings := []string{}

	for _, column := range columns {
		var field *data.Field
		// fname := getExcelColumnName(column.ColumnIndex + int(sheet.StartColumn))
		fname := column.Header

		switch column.GetType() {
		case "TIME":
			field = data.NewField(fname, nil, make([]*time.Time, len(sheet.RowData)-start))
		case "NUMBER":
			field = data.NewField(fname, nil, make([]*float64, len(sheet.RowData)-start))
		case "STRING":
			field = data.NewField(fname, nil, make([]*string, len(sheet.RowData)-start))
		default:
			return nil, fmt.Errorf("unknown column type: %s", column.GetType())
		}

		field.Config = &data.FieldConfig{
			Unit:  column.GetUnit(),
			Title: column.Header,
		}

		if column.HasMixedTypes() {
			warning := fmt.Sprintf("Multiple data types found in column %q. Using string data type", column.Header)
			warnings = append(warnings, warning)
			gs.Logger.Warn(warning)
		}

		if column.HasMixedUnits() {
			warning := fmt.Sprintf("Multiple units found in column %q. Formatted value will be used", column.Header)
			warnings = append(warnings, warning)
			gs.Logger.Warn(warning)
		}

		fields = append(fields, field)
	}

	frame := data.NewFrame(refID, // TODO: shoud set the name from metadata
		fields...,
	)
	frame.RefID = refID

	for rowIndex := start; rowIndex < len(sheet.RowData); rowIndex++ {
		for columnIndex, cellData := range sheet.RowData[rowIndex].Values {
			if columnIndex >= len(columns) {
				continue
			}

			// Skip any empty values
			if "" == cellData.FormattedValue {
				continue
			}

			switch columns[columnIndex].GetType() {
			case "TIME":
				time, err := dateparse.ParseLocal(cellData.FormattedValue)
				if err != nil {
					warnings = append(warnings, fmt.Sprintf("Error while parsing date at row %d in column %q",
						rowIndex, columns[columnIndex].Header))
				} else {
					frame.Fields[columnIndex].Set(rowIndex-start, &time)
				}
			case "NUMBER":
				if cellData.EffectiveValue != nil {
					frame.Fields[columnIndex].Set(rowIndex-start, &cellData.EffectiveValue.NumberValue)
				}
			case "STRING":
				frame.Fields[columnIndex].Set(rowIndex-start, &cellData.FormattedValue)
			}
		}
	}

	meta["warnings"] = warnings
	meta["spreadsheetId"] = qm.Spreadsheet
	meta["range"] = qm.Range
	frame.Meta = &data.QueryResultMeta{Custom: meta}
	gs.Logger.Debug("frame.Meta: %s", spew.Sdump(frame.Meta))

	return frame, nil
}

func getUniqueColumnName(formattedName string, columnIndex int, columns map[string]bool) string {
	name := formattedName
	if name == "" {
		name = fmt.Sprintf("Field %d", columnIndex+1)
	}

	nameExist := true
	counter := 1
	for nameExist {
		if _, exist := columns[name]; exist {
			name = fmt.Sprintf("%s%d", formattedName, counter)
			counter++
		} else {
			nameExist = false
		}
	}

	return name
}

func getColumnDefinitions(rows []*sheets.RowData) ([]*cd.ColumnDefinition, int) {
	columns := []*cd.ColumnDefinition{}
	columnMap := map[string]bool{}
	headerRow := rows[0].Values

	start := 0
	if len(rows) > 1 {
		start = 1
		for columnIndex, headerCell := range headerRow {
			name := getUniqueColumnName(strings.TrimSpace(headerCell.FormattedValue), columnIndex, columnMap)
			columnMap[name] = true
			columns = append(columns, cd.New(name, columnIndex))
		}
	} else {
		for columnIndex := range headerRow {
			name := getUniqueColumnName("", columnIndex, columnMap)
			columnMap[name] = true
			columns = append(columns, cd.New(name, columnIndex))
		}
	}

	// Check the types for each column
	for rowIndex := start; rowIndex < len(rows); rowIndex++ {
		for _, column := range columns {
			if column.ColumnIndex < len(rows[rowIndex].Values) {
				column.CheckCell(rows[rowIndex].Values[column.ColumnIndex])
			}
		}
	}

	return columns, start
}
