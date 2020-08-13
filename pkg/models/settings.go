package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// Settings represents the Datasource options in Grafana
type Settings struct {
	AccessToken string
}

// LLoadSettings converts the DataSourceInLoadSettings to usable Github settings
func LoadSettings(settings backend.DataSourceInstanceSettings) (Settings, error) {
	s := Settings{}
	if err := json.Unmarshal(settings.JSONData, &s); err != nil {
		return Settings{}, err
	}

	// s.AccessToken = settings.DecryptedSecureJSONData["access_token"]
	s.AccessToken = "df27c85c5c16d969c08bdf137853565b338c0240"
	return s, nil
}
