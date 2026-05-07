package main

import (
	"embed"
	"io/fs"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/pluginschema"
)

//go:embed schema/v0alpha1/*.json
var schemaFS embed.FS

// loadSchema returns the embedded v0alpha1 PluginSchema.
func loadSchema() (*pluginschema.PluginSchema, error) {
	// the composite provider expects files under {apiVersion}/...
	// our embed roots at schema/, so strip that prefix.
	sub, err := fs.Sub(schemaFS, "schema")
	if err != nil {
		return nil, err
	}
	return pluginschema.NewCompositeFileSchemaProvider(sub).Get("v0alpha1")
}
