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
	s.AccessToken = "8235020a22677db9b33c43fd29c9d10800b203eb"
	return s, nil
}
