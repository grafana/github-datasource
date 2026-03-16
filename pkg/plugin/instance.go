package plugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	schemas "github.com/grafana/schemads"

	"github.com/grafana/github-datasource/pkg/configschema"
	"github.com/grafana/github-datasource/pkg/github"
	"github.com/grafana/github-datasource/pkg/models"
)

// GitHubInstanceWithSchema wraps the GitHub datasource with schema support.
type GitHubInstanceWithSchema struct {
	Datasource
	*schemas.SchemaDatasource
}

func (g *GitHubInstanceWithSchema) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req.Path == "config-schema" {
		return sender.Send(&backend.CallResourceResponse{
			Status:  http.StatusOK,
			Headers: map[string][]string{"Content-Type": {"application/json"}},
			Body:    configschema.ConfigSchemaJSON,
		})
	}
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
		backend.Logger.Error("instance does not implement Datasource interface")
		return instance, nil
	}

	var ghDs *github.Datasource
	if datasourceSettings.CachingEnabled {
		cachedDs, ok := instance.(*CachedDatasource)
		if !ok {
			backend.Logger.Error("instance is not a cached datasource")
			return instance, nil
		}
		ghDs, ok = cachedDs.datasource.(*github.Datasource)
		if !ok {
			backend.Logger.Error("cached datasource is not a github datasource")
			return instance, nil
		}
	} else {
		ghDs, ok = instance.(*github.Datasource)
		if !ok {
			backend.Logger.Error("instance is not a github datasource")
			return instance, nil
		}
	}
	// Add schema support
	schemaProvider := github.NewSchemaProvider(ghDs)
	schemaDs := schemas.NewSchemaDatasource(
		schemaProvider,
		schemaProvider,
		schemaProvider,
		schemaProvider,
		nil,
		nil,
	)

	return &GitHubInstanceWithSchema{
		Datasource:       ds,
		SchemaDatasource: schemaDs,
	}, nil
}
