package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type AuthType string

const (
	AuthTypePAT       AuthType = "personal-access-token"
	AuthTypeGithubApp AuthType = "github-app"
)

type Settings struct {
	// General settings
	GitHubURL      string `json:"githubUrl,omitempty"`
	CachingEnabled bool   `json:"cachingEnabled,omitempty"`
	// Auth type related settings
	SelectedAuthType AuthType `json:"selectedAuthType,omitempty"`
	// personal-access-token auth related settings
	AccessToken string
	// github-app auth related settings
	AppId               json.RawMessage `json:"appId,omitempty"`
	AppIdInt64          int64
	InstallationId      json.RawMessage `json:"installationId,omitempty"`
	InstallationIdInt64 int64
	PrivateKey          string
}

func LoadSettings(settings backend.DataSourceInstanceSettings) (s Settings, err error) {
	if err := json.Unmarshal(settings.JSONData, &s); err != nil {
		return s, err
	}
	if s.SelectedAuthType == AuthTypeGithubApp {
		if s.AppIdInt64, err = rawMessageToInt64(s.AppId, "app id"); err != nil {
			return s, err
		}
		if s.InstallationIdInt64, err = rawMessageToInt64(s.InstallationId, "installation id"); err != nil {
			return s, err
		}
		if val, ok := settings.DecryptedSecureJSONData["privateKey"]; ok {
			s.PrivateKey = val
		}
	}
	if val, ok := settings.DecryptedSecureJSONData["accessToken"]; ok {
		s.AccessToken = val
	}
	// Data sources created before the auth type was introduced will have an accessToken but no auth type.
	// In this case, we default to personal access token.
	if s.AccessToken != "" && s.SelectedAuthType == "" {
		s.SelectedAuthType = AuthTypePAT
	}
	return s, nil
}

func rawMessageToInt64(r json.RawMessage, m string) (out int64, err error) {
	out, err = strconv.ParseInt(strings.ReplaceAll(string(r), `"`, ""), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing %s", m)
	}
	return out, nil
}
