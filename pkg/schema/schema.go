package schema

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/grafana/dsconfig/dsconfig"
	dsschema "github.com/grafana/dsconfig/schema"
	sdkSchema "github.com/grafana/grafana-plugin-sdk-go/experimental/pluginschema"
)

const TargetAPIVersion = "v0alpha1"

//go:embed dsconfig.json
var configSchemaJSON []byte

func DSConfigSchema() (*dsconfig.Schema, error) {
	var s dsconfig.Schema
	if err := json.Unmarshal(configSchemaJSON, &s); err != nil {
		return nil, fmt.Errorf("parse dsconfig.json: %w", err)
	}
	return &s, nil
}

func NewSDKSchema() *sdkSchema.PluginSchema {
	cfg, err := DSConfigSchema()
	if err != nil {
		panic(err)
	}
	settings, err := cfg.ToPluginSchemaSettings()
	if err != nil {
		panic(err)
	}
	var schema *sdkSchema.PluginSchema
	schema, err = dsschema.NewPluginSchema(TargetAPIVersion, settings, newSettingsExamples())
	if err != nil {
		panic(err)
	}
	return schema
}

func newSettingsExamples() *sdkSchema.SettingsExamples {
	return &sdkSchema.SettingsExamples{}
}
