package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/plugin"
)

func main() {
	if err := datasource.Manage(models.PluginID, plugin.NewDataSourceInstance, datasource.ManageOpts{}); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
