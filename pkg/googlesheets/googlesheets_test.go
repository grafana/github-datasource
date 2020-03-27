package googlesheets

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/sheets/v4"
)

type fakeClient struct {
}

func (f *fakeClient) GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error) {
	return loadTestSheet("./testdata/mixed-data.json")
}

func loadTestSheet(path string) (*sheets.Spreadsheet, error) {
	jsonBody, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sheet sheets.Spreadsheet
	if err := json.Unmarshal(jsonBody, &sheet); err != nil {
		return nil, err
	}

	return &sheet, nil
}

func TestGooglesheets(t *testing.T) {
	t.Run("getUniqueColumnName", func(t *testing.T) {
		t.Run("name is appended with number if not unique", func(t *testing.T) {
			columns := map[string]bool{"header": true, "name": true}
			name := getUniqueColumnName("header", 1, columns)
			assert.Equal(t, "header1", name)
		})

		t.Run("name becomes Field + column index if header row is empty", func(t *testing.T) {
			columns := map[string]bool{}
			name := getUniqueColumnName("", 3, columns)
			assert.Equal(t, "Field 4", name)
		})
	})

	t.Run("getSheetData", func(t *testing.T) {
		client := &fakeClient{}

		t.Run("spreadsheets get cached", func(t *testing.T) {
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			qm := QueryModel{Range: "A1:O", Spreadsheet: "someid", CacheDurationSeconds: 10}
			require.Equal(t, 0, gsd.Cache.ItemCount())

			_, meta, err := gsd.getSheetData(client, &qm)
			require.NoError(t, err)

			assert.False(t, meta["hit"].(bool))
			assert.Equal(t, 1, gsd.Cache.ItemCount())

			_, meta, err = gsd.getSheetData(client, &qm)
			require.NoError(t, err)
			assert.True(t, meta["hit"].(bool))
			assert.Equal(t, 1, gsd.Cache.ItemCount())
		})

		t.Run("spreadsheets don't get cached if CacheDurationSeconds is 0", func(t *testing.T) {
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			qm := QueryModel{Range: "A1:O", Spreadsheet: "someid", CacheDurationSeconds: 0}
			require.Equal(t, 0, gsd.Cache.ItemCount())

			_, meta, err := gsd.getSheetData(client, &qm)
			require.NoError(t, err)

			assert.False(t, meta["hit"].(bool))
			assert.Equal(t, 0, gsd.Cache.ItemCount())
		})
	})

	t.Run("transformSheetToDataFrame", func(t *testing.T) {
		sheet, err := loadTestSheet("./testdata/mixed-data.json")
		require.NoError(t, err)

		gsd := &GoogleSheets{
			Cache:  cache.New(300*time.Second, 50*time.Second),
			Logger: log.New(),
		}
		qm := QueryModel{Range: "A1:O", Spreadsheet: "someid", CacheDurationSeconds: 10}

		meta := make(map[string]interface{})
		frame, err := gsd.transformSheetToDataFrame(sheet.Sheets[0].Data[0], meta, "ref1", &qm)
		require.NoError(t, err)
		require.Equal(t, "ref1", frame.Name)

		t.Run("no of columns match", func(t *testing.T) {
			assert.Equal(t, 16, len(frame.Fields))
		})

		t.Run("no of rows matches field length", func(t *testing.T) {
			for _, field := range frame.Fields {
				assert.Equal(t, len(sheet.Sheets[0].Data[0].RowData)-1, field.Len())
			}
		})

		t.Run("meta is populated correctly", func(t *testing.T) {
			assert.Equal(t, qm.Spreadsheet, meta["spreadsheetId"])
			assert.Equal(t, qm.Range, meta["range"])
		})

		t.Run("meta warnings field is populated correctly", func(t *testing.T) {
			warnings := meta["warnings"].([]string)
			assert.Equal(t, 3, len(warnings))
			assert.Equal(t, "Multiple data types found in column \"MixedDataTypes\". Using string data type", warnings[0])
			assert.Equal(t, "Multiple units found in column \"MixedUnits\". Formatted value will be used", warnings[1])
			assert.Equal(t, "Multiple units found in column \"Mixed currencies\". Formatted value will be used", warnings[2])
			//assert.Equal(t, "Multiple data types found in column \"MixedUnits\". Using string data type", warnings[2])
		})
	})

	t.Run("query single cell", func(t *testing.T) {
		sheet, err := loadTestSheet("./testdata/single-cell.json")
		require.NoError(t, err)

		gsd := &GoogleSheets{
			Cache:  cache.New(300*time.Second, 50*time.Second),
			Logger: log.New(),
		}
		qm := QueryModel{Range: "A2", Spreadsheet: "someid", CacheDurationSeconds: 10}

		meta := make(map[string]interface{})
		frame, err := gsd.transformSheetToDataFrame(sheet.Sheets[0].Data[0], meta, "ref1", &qm)
		require.NoError(t, err)
		require.Equal(t, "ref1", frame.Name)

		t.Run("single field", func(t *testing.T) {
			assert.Equal(t, 1, len(frame.Fields))
		})

		t.Run("single row", func(t *testing.T) {
			for _, field := range frame.Fields {
				assert.Equal(t, 1, field.Len())
			}
		})

		t.Run("single value", func(t *testing.T) {
			strVal, ok := frame.Fields[0].At(0).(*string)
			require.True(t, ok)
			require.NotNil(t, strVal)
			assert.Equal(t, "🌭", *strVal)
		})
	})

	t.Run("column id formatting", func(t *testing.T) {
		require.Equal(t, "A", getExcelColumnName(1))
		require.Equal(t, "B", getExcelColumnName(2))
		require.Equal(t, "AH", getExcelColumnName(34))
		require.Equal(t, "BN", getExcelColumnName(66))
		require.Equal(t, "ZW", getExcelColumnName(699))
		require.Equal(t, "AJIL", getExcelColumnName(24582))
	})
}
