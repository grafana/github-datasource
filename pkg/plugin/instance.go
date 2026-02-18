package plugin

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	schemas "github.com/grafana/schemads"

	"github.com/grafana/github-datasource/pkg/github"
	"github.com/grafana/github-datasource/pkg/models"
)

// GitHubInstanceWithSchema wraps the GitHub datasource with schema support.
type GitHubInstanceWithSchema struct {
	Datasource
	*schemas.SchemaDatasource
}

func (g *GitHubInstanceWithSchema) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return g.SchemaDatasource.CallResource(ctx, req, sender)
}

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

	ds, ok := instance.(Datasource)
	if !ok {
		return instance, fmt.Errorf("instance does not implement Datasource interface")
	}

	var ghDs *github.Datasource
	if datasourceSettings.CachingEnabled {
		ghDs = instance.(*CachedDatasource).datasource.(*github.Datasource)
	} else {
		ghDs = instance.(*github.Datasource)
	}
	// Add schema support
	schemaHandler := github.NewSchemaHandler(ghDs)

	schemaDs := schemas.NewSchemaDatasource(schemaHandler, nil)
	_, err = schemaDs.NewDatasource(ctx, settings)
	if err != nil {
		return nil, fmt.Errorf("creating schema datasource: %w", err)
	}

	return &GitHubInstanceWithSchema{
		Datasource:       ds,
		SchemaDatasource: schemaDs,
	}, nil
}
