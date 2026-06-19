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
	GitHubURL string `json:"githubUrl,omitempty"`
	// GitHubPlan: Not used in the backend. Adding here for frontend parity
	GitHubPlan     string `json:"githubPlan,omitempty"`
	CachingEnabled bool   `json:"cachingEnabled,omitempty"`
	// Auth type related settings
	SelectedAuthType AuthType `json:"selectedAuthType,omitempty"`
	// personal-access-token auth related settings
	AccessToken string
	// github-app auth related settings
	AppId               string `json:"appId,omitempty"` // legacy config and provisioning, appId is stored as either number or string. But string should be desired format. So custom UnmarshalJSON function added to support this
	AppIdInt64          int64
	InstallationId      string `json:"installationId,omitempty"` // legacy config and provisioning, appId is stored as either number or string. But string should be desired format. So custom UnmarshalJSON function added to support this
	InstallationIdInt64 int64
	PrivateKey          string
}

// UnmarshalJSON decodes the settings while tolerating the appId and installationId
// fields being stored either as JSON strings (e.g. "1111") or, for legacy configs,
// as JSON numbers (e.g. 1111). Both are normalized to their string representation.
func (s *Settings) UnmarshalJSON(data []byte) error {
	type Alias Settings
	aux := struct {
		AppId          json.RawMessage `json:"appId,omitempty"`
		InstallationId json.RawMessage `json:"installationId,omitempty"`
		*Alias
	}{Alias: (*Alias)(s)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.AppId = rawMessageToString(aux.AppId)
	s.InstallationId = rawMessageToString(aux.InstallationId)
	return nil
}

func LoadSettings(settings backend.DataSourceInstanceSettings) (s Settings, err error) {
	if err := json.Unmarshal(settings.JSONData, &s); err != nil {
		return s, err
	}
	if s.SelectedAuthType == AuthTypeGithubApp {
		if s.AppIdInt64, err = stringToInt64(s.AppId, "app id"); err != nil {
			return s, err
		}
		if s.InstallationIdInt64, err = stringToInt64(s.InstallationId, "installation id"); err != nil {
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

func rawMessageToString(r json.RawMessage) string {
	return strings.Trim(string(r), `"`)
}

func stringToInt64(v string, m string) (out int64, err error) {
	out, err = strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing %s", m)
	}
	return out, nil
}
