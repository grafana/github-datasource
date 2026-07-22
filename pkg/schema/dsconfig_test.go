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
			Examples: map[string]*spec3.Example{},
		},
	})
}
