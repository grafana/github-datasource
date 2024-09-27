package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Settings struct {
	SelectedAuthType string `json:"selectedAuthType"`
	AccessToken      string `json:"accessToken"`
	PrivateKey       string `json:"privateKey"`
	AppId            string `json:"appId"`
	InstallationId   string `json:"installationId"`
	GitHubURL        string `json:"githubUrl"`
	CachingEnabled   bool   `json:"cachingEnabled"`
}

func LoadSettings(settings backend.DataSourceInstanceSettings) (Settings, error) {
	s := Settings{}
	if err := json.Unmarshal(settings.JSONData, &s); err != nil {
		return Settings{}, err
	}

	if val, ok := settings.DecryptedSecureJSONData["accessToken"]; ok {
		s.AccessToken = val
	}

	if val, ok := settings.DecryptedSecureJSONData["privateKey"]; ok {
		s.PrivateKey = val
	}

	// Data sources created before the auth type was introduced will have an accessToken but no auth type.
	// In this case, we default to personal access token.
	if s.AccessToken != "" && s.SelectedAuthType == "" {
		s.SelectedAuthType = "personal-access-token"
	}

	return s, nil
}
