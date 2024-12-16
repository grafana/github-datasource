package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/grafana/github-datasource/pkg/dfutil"
	githubclient "github.com/grafana/github-datasource/pkg/github/client"
	"github.com/grafana/github-datasource/pkg/github/projects"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// Make sure Datasource implements required interfaces.
var (
	_ backend.QueryDataHandler   = (*Datasource)(nil)
	_ backend.CheckHealthHandler = (*Datasource)(nil)
)

// Datasource handles requests to GitHub
type Datasource struct {
	client *githubclient.Client
}

// HandleRepositoriesQuery is the query handler for listing GitHub Repositories
func (d *Datasource) HandleRepositoriesQuery(ctx context.Context, query *models.RepositoriesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListRepositoriesOptions{
		Owner:      query.Owner,
		Repository: query.Repository,
	}

	return GetAllRepositories(ctx, d.client, opt)
}

// HandleIssuesQuery is the query handler for listing GitHub Issues
func (d *Datasource) HandleIssuesQuery(ctx context.Context, query *models.IssuesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.IssueOptionsWithRepo(query.Options, query.Owner, query.Repository)
	return GetIssuesInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

// HandleCommitsQuery is the query handler for listing GitHub Commits
func (d *Datasource) HandleCommitsQuery(ctx context.Context, query *models.CommitsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.CommitsOptionsWithRepo(query.Options, query.Owner, query.Repository)
	return GetCommitsInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

// HandleTagsQuery is the query handler for listing GitHub Tags
func (d *Datasource) HandleTagsQuery(ctx context.Context, query *models.TagsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListTagsOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
	}

	if req.TimeRange.From.Unix() <= 0 && req.TimeRange.To.Unix() <= 0 {
		return GetAllTags(ctx, d.client, opt)
	}

	return GetTagsInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

// HandleReleasesQuery is the query handler for listing GitHub Releases
func (d *Datasource) HandleReleasesQuery(ctx context.Context, query *models.ReleasesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListReleasesOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
	}

	if req.TimeRange.From.Unix() <= 0 && req.TimeRange.To.Unix() <= 0 {
		return GetAllReleases(ctx, d.client, opt)
	}
	return GetReleasesInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

// HandlePullRequestsQuery is the query handler for listing GitHub PullRequests
func (d *Datasource) HandlePullRequestsQuery(ctx context.Context, query *models.PullRequestsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.PullRequestOptionsWithRepo(query.Options, query.Owner, query.Repository)

	if req.TimeRange.From.Unix() <= 0 && req.TimeRange.To.Unix() <= 0 {
		return GetAllPullRequests(ctx, d.client, opt)
	}
	return GetPullRequestsInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

// HandleContributorsQuery is the query handler for listing GitHub Contributors
func (d *Datasource) HandleContributorsQuery(ctx context.Context, query *models.ContributorsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListContributorsOptions{
		Owner:      query.Owner,
		Repository: query.Repository,
		Query:      query.Options.Query,
	}

	return GetAllContributors(ctx, d.client, opt)
}

// HandleLabelsQuery is the query handler for listing GitHub Labels
func (d *Datasource) HandleLabelsQuery(ctx context.Context, query *models.LabelsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListLabelsOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
		Query:      query.Options.Query,
	}

	return GetAllLabels(ctx, d.client, opt)
}

// HandleMilestonesQuery is the query handler for listing GitHub Milestones
func (d *Datasource) HandleMilestonesQuery(ctx context.Context, query *models.MilestonesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListMilestonesOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
		Query:      query.Options.Query,
	}

	return GetAllMilestones(ctx, d.client, opt)
}

// HandlePackagesQuery is the query handler for listing GitHub Packages
func (d *Datasource) HandlePackagesQuery(ctx context.Context, query *models.PackagesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.PackagesOptionsWithRepo(query.Options, query.Owner, query.Repository)

	return GetAllPackages(ctx, d.client, opt)
}

// HandleVulnerabilitiesQuery is the query handler for listing GitHub Packages
func (d *Datasource) HandleVulnerabilitiesQuery(ctx context.Context, query *models.VulnerabilityQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListVulnerabilitiesOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
	}

	return GetAllVulnerabilities(ctx, d.client, opt)
}

// HandleProjectsQuery is the query handler for listing GitHub Projects
func (d *Datasource) HandleProjectsQuery(ctx context.Context, query *models.ProjectsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ProjectOptions{
		Organization: query.Options.Organization,
		Number:       query.Options.Number,
		User:         query.Options.User,
		Kind:         query.Options.Kind,
		Filters:      query.Options.Filters,
	}

	if projects.ProjectNumber(query.Options.Number) > 0 {
		return projects.GetAllProjectItems(ctx, d.client, opt)
	}
	return projects.GetAllProjects(ctx, d.client, opt)
}

// HandleStargazersQuery is the query handler for listing stargazers of a GitHub repository
func (d *Datasource) HandleStargazersQuery(ctx context.Context, query *models.StargazersQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListStargazersOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
	}

	return GetStargazers(ctx, d.client, opt, req.TimeRange)
}

// HandleWorkflowsQuery is the query handler for listing workflows of a GitHub repository
func (d *Datasource) HandleWorkflowsQuery(ctx context.Context, query *models.WorkflowsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListWorkflowsOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
		TimeField:  query.Options.TimeField,
	}

	return GetWorkflows(ctx, d.client, opt, req.TimeRange)
}

// HandleWorkflowUsageQuery is the query handler for getting the usage information of a specific workflow
func (d *Datasource) HandleWorkflowUsageQuery(ctx context.Context, query *models.WorkflowUsageQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.WorkflowUsageOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
		Workflow:   query.Options.Workflow,
	}

	return GetWorkflowUsage(ctx, d.client, opt, req.TimeRange)
}

// HandleWorkflowRunsQuery is the query handler for listing workflow runs of a GitHub repository
func (d *Datasource) HandleWorkflowRunsQuery(ctx context.Context, query *models.WorkflowRunsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.WorkflowRunsOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
		Workflow:   query.Options.Workflow,
		Branch:     query.Options.Branch,
	}

	return GetWorkflowRuns(ctx, d.client, opt, req.TimeRange)
}

// CheckHealth is the health check for GitHub
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	_, err := GetAllRepositories(ctx, d.client, models.ListRepositoriesOptions{
		Owner:      "grafana",
		Repository: "github-datasource",
	})
	if err != nil {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			return newHealthResult(backend.HealthStatusError, "401 Unauthorized. Check your API key/Access token")
		}
		if strings.Contains(err.Error(), "404 Not Found") {
			return newHealthResult(backend.HealthStatusError, "404 Not Found. Check the Github Enterprise Server URL")
		}
		if strings.HasSuffix(err.Error(), "no such host") {
			return newHealthResult(backend.HealthStatusError, "Unable to reach the Github Enterprise Server URL from the Grafana server. Check the Github Enterprise Server URL and/or proxy settings")
		}
		return newHealthResult(backend.HealthStatusError, fmt.Sprintf("Health check failed. %s", err.Error()))
	}

	return newHealthResult(backend.HealthStatusOk, "Data source is working")
}

// QueryData runs the query
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	m := GetQueryHandlers(&QueryHandler{
		Datasource: *d,
	})

	return m.QueryData(ctx, req)
}

// NewDatasource creates a new datasource for handling queries
func NewDatasource(ctx context.Context, settings models.Settings) (*Datasource, error) {
	client, err := githubclient.New(ctx, settings)
	if err != nil {
		return nil, err
	}
	return &Datasource{client: client}, nil
}

func newHealthResult(status backend.HealthStatus, message string) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
