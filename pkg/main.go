package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
)

func main() {
	err := datasource.Serve(newDatasource())

	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
