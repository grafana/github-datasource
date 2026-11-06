package models_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type settingsLogger struct {
	log.Logger
	warnings []map[string]any
}

func (l *settingsLogger) Warn(_ string, args ...any) {
	fields := map[string]any{}
	for i := 0; i < len(args); i += 2 {
		fields[args[i].(string)] = args[i+1]
	}
	l.warnings = append(l.warnings, fields)
}

func TestSettingsIDConversionLogging(t *testing.T) {
	previousLogger := log.DefaultLogger
	t.Cleanup(func() { log.DefaultLogger = previousLogger })
	for _, tc := range []struct {
		name       string
		jsonData   string
		appID      string
		installID  string
		warnFields []string
	}{
		{name: "strings", jsonData: `{"appId":"1111","installationId":"2222"}`, appID: "1111", installID: "2222"},
		{name: "numbers", jsonData: `{"appId":1111,"installationId":2222}`, appID: "1111", installID: "2222", warnFields: []string{"appId", "installationId"}},
		{name: "mixed with large ID", jsonData: `{"appId":9007199254740993,"installationId":"2222"}`, appID: "9007199254740993", installID: "2222", warnFields: []string{"appId"}},
		{name: "missing", jsonData: `{}`},
		{name: "null", jsonData: `{"appId":null,"installationId":null}`, appID: "null", installID: "null"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := &settingsLogger{Logger: log.NewNullLogger()}
			log.DefaultLogger = logger
			var settings models.Settings
			require.NoError(t, json.Unmarshal([]byte(tc.jsonData), &settings))
			assert.Equal(t, tc.appID, settings.AppId)
			assert.Equal(t, tc.installID, settings.InstallationId)
			require.Len(t, logger.warnings, len(tc.warnFields))
			for i, field := range tc.warnFields {
				assert.Equal(t, map[string]any{
					"field": field, "from": "number", "to": "string", "outcome": "coerced",
				}, logger.warnings[i])
			}
		})
	}
}

func TestLoadSettings(t *testing.T) {
	tests := []struct {
		name              string
		settings          backend.DataSourceInstanceSettings
		jsonData          json.RawMessage
		decryptedJsonData map[string]string
		want              models.Settings
		wantErr           error
	}{
		{
			name: "valid config should not throw error for pat authentication",
			jsonData: []byte(`{
				"githubUrl" 		: 	"https://foo.com"
			}`),
			decryptedJsonData: map[string]string{"accessToken": "foo"},
			want: models.Settings{
				GitHubURL:        "https://foo.com",
				SelectedAuthType: models.AuthTypePAT,
				AccessToken:      "foo",
			},
		},
		{
			name: "valid config should not throw error for github app type authentication",
			jsonData: []byte(`{
				"githubUrl" 		: 	"https://foo.com",
				"selectedAuthType" 	: 	"github-app",
				"appId" 			: 	"1111",
				"installationId" 	:	"2222"
			}`),
			decryptedJsonData: map[string]string{"privateKey": "foo"},
			want: models.Settings{
				GitHubURL:           "https://foo.com",
				SelectedAuthType:    models.AuthTypeGithubApp,
				AppId:               "1111",
				AppIdInt64:          1111,
				InstallationId:      "2222",
				InstallationIdInt64: 2222,
				PrivateKey:          "foo",
			},
		},
		{
			name: "valid config should not throw error for github app type authentication - passed as numbers",
			jsonData: []byte(`{
				"githubUrl" 		: 	"https://foo.com",
				"selectedAuthType" 	: 	"github-app",
				"appId" 			: 	1111,
				"installationId" 	:	2222
			}`),
			decryptedJsonData: map[string]string{"privateKey": "foo"},
			want: models.Settings{
				GitHubURL:           "https://foo.com",
				SelectedAuthType:    models.AuthTypeGithubApp,
				AppId:               "1111",
				AppIdInt64:          1111,
				InstallationId:      "2222",
				InstallationIdInt64: 2222,
				PrivateKey:          "foo",
			},
		},
		{
			name: "invalid config should throw error for github app type authentication - app id passed as string literals",
			jsonData: []byte(`{
				"githubUrl" 		: 	"https://foo.com",
				"selectedAuthType" 	: 	"github-app",
				"appId" 			: 	"1111xyz",
				"installationId" 	:	"2222"
			}`),
			decryptedJsonData: map[string]string{"privateKey": "foo"},
			wantErr:           errors.New("error parsing app id"),
		},
		{
			name: "invalid config should throw error for github app type authentication - installation id passed as string literals",
			jsonData: []byte(`{
				"githubUrl" 		: 	"https://foo.com",
				"selectedAuthType" 	: 	"github-app",
				"appId" 			: 	"1111",
				"installationId" 	:	"2222xyz"
			}`),
			decryptedJsonData: map[string]string{"privateKey": "foo"},
			wantErr:           errors.New("error parsing installation id"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.jsonData == nil {
				tt.jsonData = []byte(`{}`)
			}
			got, err := models.LoadSettings(backend.DataSourceInstanceSettings{JSONData: tt.jsonData, DecryptedSecureJSONData: tt.decryptedJsonData})
			if tt.wantErr != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			require.Nil(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}
