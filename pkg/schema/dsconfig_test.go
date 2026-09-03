package schema_test

import (
	_ "embed"
	"testing"

	"github.com/grafana/dsconfig/schema"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/pluginschema"
	"k8s.io/kube-openapi/pkg/spec3"

	"github.com/grafana/github-datasource/pkg/models"
)

//go:embed dsconfig.json
var configSchemaJSON []byte

//go:generate go test -run TestPlugin -generateArtifacts
func TestPlugin(t *testing.T) {
	secureKeys := []string{}
	for _, k := range models.SecureJsonDataKeys {
		secureKeys = append(secureKeys, string(k))
	}
	schema.RunPluginTests(t, schema.PluginUnderTest{
		ID:                models.PluginID,
		ConfigSchemaJSON:  configSchemaJSON,
		SettingsJSONModel: models.Settings{},
		SecureKeys:        secureKeys,
		SettingsExamples: &pluginschema.SettingsExamples{
			Examples: map[string]*spec3.Example{
				"": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "Default configuration",
						Description: "The defaults a new datasource starts with: GitHub.com (Free, Pro & Team) and personal access token authentication. The token still has to be supplied before the datasource can load.",
						Value: map[string]any{
							"jsonData": map[string]any{
								"githubPlan":       "github-basic",
								"selectedAuthType": string(models.AuthTypePAT),
							},
							"secureJsonData": map[string]any{
								"accessToken": "",
							},
						},
					},
				},
				"personalAccessToken": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "Personal access token on GitHub.com",
						Description: "Authenticate against GitHub.com with a fine grained personal access token. Grant read-only repository permissions plus 'Metadata'; the plugin never writes. Do not set jsonData.githubUrl — any non-empty value makes the backend treat the instance as on-prem Enterprise Server.",
						Value: map[string]any{
							"jsonData": map[string]any{
								"githubPlan":       "github-basic",
								"selectedAuthType": string(models.AuthTypePAT),
							},
							"secureJsonData": map[string]any{
								"accessToken": "REPLACE_WITH_ACCESS_TOKEN",
							},
						},
					},
				},
				"githubApp": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "GitHub App authentication",
						Description: "Authenticate as a GitHub App installation. appId and installationId are numeric values (accepted as JSON strings or numbers) and privateKey is the complete PEM text including the BEGIN/END lines. All three are mandatory: a missing or non-numeric id makes the backend fail to load the datasource.",
						Value: map[string]any{
							"jsonData": map[string]any{
								"githubPlan":       "github-basic",
								"selectedAuthType": string(models.AuthTypeGithubApp),
								"appId":            "123456",
								"installationId":   "12345678",
							},
							"secureJsonData": map[string]any{
								"privateKey": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
							},
						},
					},
				},
				"enterpriseCloud": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "GitHub Enterprise Cloud",
						Description: "Enterprise Cloud is still served from github.com, so only jsonData.githubPlan changes. Leave jsonData.githubUrl unset — setting it would route the backend at an on-prem Enterprise Server instead.",
						Value: map[string]any{
							"jsonData": map[string]any{
								"githubPlan":       "github-enterprise-cloud",
								"selectedAuthType": string(models.AuthTypePAT),
							},
							"secureJsonData": map[string]any{
								"accessToken": "REPLACE_WITH_ACCESS_TOKEN",
							},
						},
					},
				},
				"enterpriseServer": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "GitHub Enterprise Server (on-prem)",
						Description: "Point the datasource at a self-hosted GitHub Enterprise Server. jsonData.githubUrl is required when jsonData.githubPlan is 'github-enterprise-server'. Omit the trailing slash: the backend concatenates '/api/v3' and '/api/graphql' onto the value.",
						Value: map[string]any{
							"jsonData": map[string]any{
								"githubPlan":       "github-enterprise-server",
								"githubUrl":        "https://github.example.com",
								"selectedAuthType": string(models.AuthTypePAT),
							},
							"secureJsonData": map[string]any{
								"accessToken": "REPLACE_WITH_ACCESS_TOKEN",
							},
						},
					},
				},
			},
		},
	})
}
