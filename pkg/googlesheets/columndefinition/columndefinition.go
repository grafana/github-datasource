package columndefinition

import (
	"strings"

	"google.golang.org/api/sheets/v4"
)

// ColumnType is the set of possible column types
type ColumnType string

const (
	// ColumTypeTime is the TIME type
	ColumTypeTime ColumnType = "TIME"
	// ColumTypeNumber is the NUMBER type
	ColumTypeNumber = "NUMBER"
	// ColumTypeString is the STRING type
	ColumTypeString = "STRING"
)

// ColumnDefinition represents a spreadsheet column definition.
type ColumnDefinition struct {
	Header      string
	ColumnIndex int
	types       map[ColumnType]bool
	units       map[string]bool
}

// New creates a new ColumnDefinition.
func New(header string, index int) *ColumnDefinition {
	return &ColumnDefinition{
		Header:      header,
		ColumnIndex: index,
		types:       map[ColumnType]bool{},
		units:       map[string]bool{},
	}
}

// CheckCell checks a ColumnDefinition's cell.
func (cd *ColumnDefinition) CheckCell(cell *sheets.CellData) {
	cd.checkType(cell)
	cd.checkUnit(cell)
}

// GetType gets the type of a ColumnDefinition.
func (cd *ColumnDefinition) GetType() ColumnType {
	if len(cd.types) == 1 {
		for columnType := range cd.types {
			return columnType
		}
	}

	//The column has mixed or no data types - fallback to string
	return ColumTypeString
}

// GetUnit gets the unit of a ColumnDefinition.
func (cd *ColumnDefinition) GetUnit() string {
	if len(cd.units) == 1 {
		for unit := range cd.units {
			return unit
		}
	}

	return ""
}

// HasMixedTypes returns whether a ColumnDefinition has mixed types.
func (cd *ColumnDefinition) HasMixedTypes() bool {
	return len(cd.types) > 1
}

// HasMixedUnits returns whether a ColumnDefinition has mixed units.
func (cd *ColumnDefinition) HasMixedUnits() bool {
	return len(cd.units) > 1
}

func (cd *ColumnDefinition) checkType(cell *sheets.CellData) {
	if cell == nil || cell.FormattedValue == "" {
		return
	}

	// Has a number value (will not detect 0)
	hasNumberValue := cell.UserEnteredValue != nil && cell.UserEnteredValue.NumberValue != 0
	hasNumberFormat := cell.UserEnteredFormat != nil && cell.UserEnteredFormat.NumberFormat != nil

	if hasNumberFormat {
		if cell.UserEnteredFormat.NumberFormat.Type == "DATE" || cell.UserEnteredFormat.NumberFormat.Type == "DATE_TIME" {
			cd.types["TIME"] = true
			return
		}
	}

	if hasNumberFormat || hasNumberValue || "0" == cell.FormattedValue {
		cd.types["NUMBER"] = true
	} else {
		cd.types["STRING"] = true
	}
}

var unitMappings = map[string]string{
	"$":   "currencyUSD",
	"£":   "currencyGBP",
	"€":   "currencyEUR",
	"¥":   "currencyJPY",
	"₽":   "currencyRUB",
	"₴":   "currencyUAH",
	"R$":  "currencyBRL",
	"kr.": "currencyDKK",
	"kr":  "currencySEK",
	"czk": "currencyCZK",
	"CHF": "currencyCHF",
	"PLN": "currencyPLN",
	"฿":   "currencyBTC",
	"R":   "currencyZAR",
	"₹":   "currencyINR",
	"₩":   "currencyKRW",
}

// A lot more that can be done/improved here. For example it should be possible to extract
// the number of decimals from the pattern. Read more here: https://developers.google.com/sheets/api/guides/formats
func (cd *ColumnDefinition) checkUnit(cellData *sheets.CellData) {
	if cellData == nil || cellData.UserEnteredFormat == nil || cellData.UserEnteredFormat.NumberFormat == nil {
		return
	}

	switch cellData.UserEnteredFormat.NumberFormat.Type {
	case "NUMBER":
		for unit, unitID := range unitMappings {
			if strings.Contains(cellData.UserEnteredFormat.NumberFormat.Pattern, unit) {
				cd.units[unitID] = true
			}
		}
	case "PERCENT":
		cd.units["percent"] = true
	case "CURRENCY":
		for unit, unitID := range unitMappings {
			if strings.Contains(cellData.FormattedValue, unit) {
				cd.units[unitID] = true
			}
		}
	}
}
