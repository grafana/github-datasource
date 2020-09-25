package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// Settings represents the Datasource options in Grafana
type Settings struct {
	AccessToken    string `json:"accessToken"`
	GithubURL      string `json:"githubUrl"`
	CachingEnabled bool   `json:"cachingEnabled"`
}

// LoadSettings converts the DataSourceInLoadSettings to usable Github settings
func LoadSettings(settings backend.DataSourceInstanceSettings) (Settings, error) {
	s := Settings{}
	if err := json.Unmarshal(settings.JSONData, &s); err != nil {
		return Settings{}, err
	}

	if val, ok := settings.DecryptedSecureJSONData["accessToken"]; ok {
		s.AccessToken = val
	}

	return s, nil
}
