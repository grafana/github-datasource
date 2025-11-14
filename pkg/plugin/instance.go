package plugin

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"

	"github.com/grafana/github-datasource/pkg/github"
	"github.com/grafana/github-datasource/pkg/models"
)

// NewGitHubInstance creates a new GitHubInstance using the settings to determine if things like the Caching Wrapper should be enabled
func NewGitHubInstance(ctx context.Context, settings models.Settings) (instancemgmt.Instance, error) {
	gh, err := github.NewDatasource(ctx, settings)
	if err != nil {
		return nil, err
	}

	var d Datasource = gh

	if settings.CachingEnabled {
		d = WithCaching(d)
	}

	return d, nil
}

// NewDataSourceInstance creates a new instance
func NewDataSourceInstance(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	datasourceSettings, err := models.LoadSettings(settings)
	if err != nil {
		return nil, err
	}
	
	datasourceSettings.CachingEnabled = true

	instance, err := NewGitHubInstance(ctx, datasourceSettings)
	if err != nil {
		return instance, fmt.Errorf("instantiating github instance: %w", err)
	}

	return instance, nil
}
