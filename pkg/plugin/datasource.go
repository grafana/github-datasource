package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// The Datasource type handles the requests sent to the datasource backend
type Datasource interface {
	HandleRepositoriesQuery(context.Context, *models.RepositoriesQuery, backend.DataQuery) (github.Repositories, error)
	HandleIssuesQuery(context.Context, *models.IssuesQuery, backend.DataQuery) (github.Issues, error)
	HandleCommitsQuery(context.Context, *models.CommitsQuery, backend.DataQuery) (github.Commits, error)
	HandleTagsQuery(context.Context, *models.TagsQuery, backend.DataQuery) (github.Tags, error)
	HandleReleasesQuery(context.Context, *models.ReleasesQuery, backend.DataQuery) (github.Releases, error)
	HandleContributorsQuery(context.Context, *models.ContributorsQuery, backend.DataQuery) (github.Users, error)
	HandlePullRequestsQuery(context.Context, *models.PullRequestsQuery, backend.DataQuery) (github.PullRequests, error)
	HandleLabelsQuery(context.Context, *models.LabelsQuery, backend.DataQuery) (github.Labels, error)
	HandlePackagesQuery(context.Context, *models.PackagesQuery, backend.DataQuery) (github.Packages, error)
	HandleMilestonesQuery(context.Context, *models.MilestonesQuery, backend.DataQuery) (github.Milestones, error)
	CheckHealth(context.Context) error
}

// HandleQueryData handles the `QueryData` request for the Github datasource
func HandleQueryData(ctx context.Context, d Datasource, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	m := GetQueryHandlers(&Server{
		Datasource: d,
	})

	return m.QueryData(ctx, req)
}

// CheckHealth ensures that the datasource settings are able to retrieve data from GitHub
func CheckHealth(ctx context.Context, d Datasource, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if err := d.CheckHealth(ctx); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}
