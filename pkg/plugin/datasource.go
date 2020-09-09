package plugin

import (
	"context"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// The Datasource type handles the requests sent to the datasource backend
type Datasource interface {
	HandleIssuesQuery(context.Context, *models.IssuesQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleCommitsQuery(context.Context, *models.CommitsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleTagsQuery(context.Context, *models.TagsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleReleasesQuery(context.Context, *models.ReleasesQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleContributorsQuery(context.Context, *models.ContributorsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandlePullRequestsQuery(context.Context, *models.PullRequestsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleLabelsQuery(context.Context, *models.LabelsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandlePackagesQuery(context.Context, *models.PackagesQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleMilestonesQuery(context.Context, *models.MilestonesQuery, backend.DataQuery) (dfutil.Framer, error)
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
