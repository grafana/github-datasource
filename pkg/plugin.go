package main

import (
	"net/http"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"context"
)

func main() {
	// Setup the plugin environment
	backend.SetupPluginEnvironment("grafanan-github-datasource")

	mux := http.NewServeMux()
	ds := Init(mux)
	// httpResourceHandler := httpadapter.New(mux)

	err := backend.Serve(backend.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	})

	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}

}

// QueryData .....
func (plugin *GithubStatsDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return nil, nil
}

// CheckHealth .....
func (plugin *GithubStatsDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}

// Init creates the google sheets datasource and sets up all the routes
func Init(mux *http.ServeMux) *GithubStatsDatasource {
	// fix me
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	ds := &GithubStatsDatasource{client: githubv4.NewClient(httpClient)}
	return ds
}
