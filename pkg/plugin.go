package main

import (
	"net/http"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/prometheus/client_golang/prometheus"

	"context"
)

const metricNamespace = "sheets_datasource"

func main() {
	// Setup the plugin environment
	backend.SetupPluginEnvironment("grafanan-github-datasource")

	mux := http.NewServeMux()
	ds := Init(mux)
	//httpResourceHandler := httpadapter.New(mux)

	err := backend.Serve(backend.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	})

	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}

type GithubStatsDatasource struct {
}

func (plugin *GithubStatsDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return nil, nil
}

func (plugin *GithubStatsDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}

// Init creates the google sheets datasource and sets up all the routes
func Init(mux *http.ServeMux) *GithubStatsDatasource {
	queriesTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "data_query_total",
			Help:      "data query counter",
			Namespace: metricNamespace,
		},
		[]string{"scenario"},
	)
	prometheus.MustRegister(queriesTotal)

	ds := &GithubStatsDatasource{}
	return ds
}
