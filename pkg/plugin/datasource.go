package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// The Datasource type handles the requests sent to the datasource backend
type Datasource interface {
	HandleRepositoriesQuery(context.Context, *models.RepositoriesQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleIssuesQuery(context.Context, *models.IssuesQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleCommitsQuery(context.Context, *models.CommitsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleTagsQuery(context.Context, *models.TagsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleReleasesQuery(context.Context, *models.ReleasesQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleContributorsQuery(context.Context, *models.ContributorsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandlePullRequestsQuery(context.Context, *models.PullRequestsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleLabelsQuery(context.Context, *models.LabelsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandlePackagesQuery(context.Context, *models.PackagesQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleMilestonesQuery(context.Context, *models.MilestonesQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleVulnerabilitiesQuery(context.Context, *models.VulnerabilityQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleProjectsQuery(context.Context, *models.ProjectsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleStargazersQuery(context.Context, *models.StargazersQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleWorkflowsQuery(context.Context, *models.WorkflowsQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleWorkflowUsageQuery(context.Context, *models.WorkflowUsageQuery, backend.DataQuery) (dfutil.Framer, error)
	HandleWorkflowRunsQuery(context.Context, *models.WorkflowRunsQuery, backend.DataQuery) (dfutil.Framer, error)
	CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error)
	QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error)
}
