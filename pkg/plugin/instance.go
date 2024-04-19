package plugin

import (
	"context"
	"fmt"

	"github.com/grafana/github-datasource/pkg/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// NewGitHubInstance creates a new GitHubInstance using the settings to determine if things like the Caching Wrapper should be enabled
func NewGitHubInstance(ctx context.Context, settings models.Settings) (instancemgmt.Instance, error) {
	gh, err := github.NewDatasource(ctx, settings)
	if err != nil {
		return nil, fmt.Errorf("instantiating github datasource: %w", err)
	}
	return WithCaching(gh), nil
}

// NewDataSourceInstance creates a new instance
func NewDataSourceInstance(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	datasourceSettings, err := models.LoadSettings(settings)
	if err != nil {
		return nil, err
	}

	instance, err := NewGitHubInstance(context.Background(), datasourceSettings)
	if err != nil {
		return instance, fmt.Errorf("instantiating github instance")
	}

	return instance, nil
}
