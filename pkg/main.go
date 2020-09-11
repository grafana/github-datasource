package main

import (
	"os"

	"github.com/grafana/github-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
)

func main() {
	err := datasource.Serve(plugin.GetDatasourceOpts())

	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
