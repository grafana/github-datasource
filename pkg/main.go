package main

import (
	"os"

	"github.com/grafana/github-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
)

const dsID = "grafana-github-datasource"

func main() {
	if err := datasource.Manage(dsID, plugin.NewDataSourceInstance, datasource.ManageOpts{}); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
