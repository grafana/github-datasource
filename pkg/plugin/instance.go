package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// Instance is the root Datasource implementation that wraps a Datasource
type Instance struct {
	Datasource Datasource
	Handlers   Handlers
}

// HandleRepositoriesQuery ...
func (i *Instance) HandleRepositoriesQuery(ctx context.Context, q *models.RepositoriesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleRepositoriesQuery(ctx, q, req)
}

// HandleIssuesQuery ...
func (i *Instance) HandleIssuesQuery(ctx context.Context, q *models.IssuesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleIssuesQuery(ctx, q, req)
}

// HandleCommitsQuery ...
func (i *Instance) HandleCommitsQuery(ctx context.Context, q *models.CommitsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleCommitsQuery(ctx, q, req)
}

// HandleTagsQuery ...
func (i *Instance) HandleTagsQuery(ctx context.Context, q *models.TagsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleTagsQuery(ctx, q, req)
}

// HandleReleasesQuery ...
func (i *Instance) HandleReleasesQuery(ctx context.Context, q *models.ReleasesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleReleasesQuery(ctx, q, req)
}

// HandleContributorsQuery ...
func (i *Instance) HandleContributorsQuery(ctx context.Context, q *models.ContributorsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleContributorsQuery(ctx, q, req)
}

// HandlePullRequestsQuery ...
func (i *Instance) HandlePullRequestsQuery(ctx context.Context, q *models.PullRequestsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandlePullRequestsQuery(ctx, q, req)
}

// HandleLabelsQuery ...
func (i *Instance) HandleLabelsQuery(ctx context.Context, q *models.LabelsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleLabelsQuery(ctx, q, req)
}

// HandlePackagesQuery ...
func (i *Instance) HandlePackagesQuery(ctx context.Context, q *models.PackagesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandlePackagesQuery(ctx, q, req)
}

// HandleMilestonesQuery ...
func (i *Instance) HandleMilestonesQuery(ctx context.Context, q *models.MilestonesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleMilestonesQuery(ctx, q, req)
}

// HandleVulnerabilitiesQuery ...
func (i *Instance) HandleVulnerabilitiesQuery(ctx context.Context, q *models.VulnerabilityQuery, req backend.DataQuery) (dfutil.Framer, error) {
	return i.Datasource.HandleVulnerabilitiesQuery(ctx, q, req)
}

// CheckHealth ...
func (i *Instance) CheckHealth(ctx context.Context) error {
	return i.Datasource.CheckHealth(ctx)
}

// NewGitHubInstance creates a new GitHubInstance using the settings to determine if things like the Caching Wrapper should be enabled
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
