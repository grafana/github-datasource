package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// NewGitHubInstance creates a new GitHubInstance using the settings to determine if things like the Caching Wrapper should be enabled
func NewGitHubInstance(ctx context.Context, settings models.Settings) instancemgmt.Instance {
	var (
		gh = github.NewDatasource(ctx, settings)
	)

	var d Datasource = gh

	if settings.CachingEnabled {
		d = WithCaching(d)
	}

	return d
}

// NewDataSourceInstance creates a new instance
func NewDataSourceInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	datasourceSettings, err := models.LoadSettings(settings)
	if err != nil {
		return nil, err
	}

	datasourceSettings.CachingEnabled = true

	return NewGitHubInstance(context.Background(), datasourceSettings), nil
}
