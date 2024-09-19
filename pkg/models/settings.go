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
	GithubURL        string `json:"githubUrl"`
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

	return s, nil
}
