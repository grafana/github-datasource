package models_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				AppId:               []byte(`"1111"`),
				AppIdInt64:          1111,
				InstallationId:      []byte(`"2222"`),
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
				AppId:               []byte(`1111`),
				AppIdInt64:          1111,
				InstallationId:      []byte(`2222`),
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
