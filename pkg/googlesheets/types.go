package googlesheets

import (
	"google.golang.org/api/sheets/v4"
)

// QueryModel represents a spreadsheet query.
type QueryModel struct {
	Spreadsheet          string `json:"spreadsheet"`
	Range                string `json:"range"`
	CacheDurationSeconds int    `json:"cacheDurationSeconds"`
	UseTimeFilter        bool   `json:"useTimeFilter"`
}

// GoogleSheetConfig contains Google Sheets API authentication properties.
type GoogleSheetConfig struct {
	AuthType string `json:"authType"` // jwt | key
	APIKey   string `json:"apiKey"`
	JWT      string `json:"jwt"`
}

type client interface {
	GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error)
}
