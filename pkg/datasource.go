package main

import (
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"context"
)

// The Datasource type handles the requests sent to the datasource backend
type Datasource struct {
	client      *githubv4.Client
	oauthConfig *oauth2.Config
}

// HandHandleQueryData handles the `QueryData` request for the Github datasource
func (d *Datasource) HandleQueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return nil, nil
}

// CheckHealth ensures that the datasource settings are able to retrieve data from GitHub
func (plugin *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}

func NewDatasource(ctx context.Context, settings models.Settings) *Datasource {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: settings.AccessToken},
	)

	httpClient := oauth2.NewClient(ctx, src)

	return &Datasource{
		client: githubv4.NewClient(httpClient),
		oauthConfig: &oauth2.Config{
			ClientID:     "dd372c64898d6e07b3d4",
			ClientSecret: "f650d0e2c735668444156bc4fb44ab913b21461c",
			Scopes:       []string{"repo", "read:org"},
			Endpoint:     github.Endpoint,
		},
	}
}
