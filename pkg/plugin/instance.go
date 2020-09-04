package plugin

import (
	"context"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-github-datasource/pkg/github"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// Instance is the root Datasource implementation that wraps a Datasource
type Instance struct {
	Datasource Datasource
	Handlers   Handlers
}

func (i *Instance) HandleIssuesQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleIssuesQuery(ctx, q, req)
}

func (i *Instance) HandleCommitsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleCommitsQuery(ctx, q, req)
}

func (i *Instance) HandleTagsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleTagsQuery(ctx, q, req)
}

func (i *Instance) HandleReleasesQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleReleasesQuery(ctx, q, req)
}

func (i *Instance) HandleContributorsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleContributorsQuery(ctx, q, req)
}

func (i *Instance) HandlePullRequestsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandlePullRequestsQuery(ctx, q, req)
}

func (i *Instance) HandleLabelsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleLabelsQuery(ctx, q, req)
}

func (i *Instance) CheckHealth(ctx context.Context) error {
	return i.Datasource.CheckHealth(ctx)
}

func NewGitHubInstance(ctx context.Context, settings models.Settings) *Instance {
	var (
		gh = github.NewDatasource(ctx, settings)
	)

	var d Datasource = gh

	if settings.CachingEnabled {
		d = WithCaching(d)
	}

	// TODO: wrap these HTTP handlers with a caching wrapper
	return &Instance{
		Datasource: d,
		Handlers: Handlers{
			Labels:     gh.HandleGetLabels,
			Milestones: gh.HandleGetMilestones,
		},
	}
}

func newDataSourceInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	datasourceSettings, err := models.LoadSettings(settings)
	if err != nil {
		return nil, err
	}

	datasourceSettings.CachingEnabled = true

	return NewGitHubInstance(context.Background(), datasourceSettings), nil
}
