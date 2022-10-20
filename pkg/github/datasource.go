package github

import (
	"context"
	"fmt"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Make sure Datasource implements required interfaces.
var (
	_ backend.QueryDataHandler   = (*Datasource)(nil)
	_ backend.CheckHealthHandler = (*Datasource)(nil)
)

// Datasource handles requests to GitHub
type Datasource struct {
	client *githubv4.Client
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
		Query:      query.Options.Query,
	}

	return GetAllVulnerabilities(ctx, d.client, opt)
}

// HandleProjectsQuery is the query handler for listing GitHub Projects
func (d *Datasource) HandleProjectsQuery(ctx context.Context, query *models.ProjectsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListProjectsOptions{
		Organization: query.Options.Organization,
	}

	return GetAllProjects(ctx, d.client, opt)
}

// CheckHealth is the health check for GitHub
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	_, err := GetAllRepositories(ctx, d.client, models.ListRepositoriesOptions{
		Owner: "grafana",
	})

	if err != nil {
		return newHealthResult(backend.HealthStatusError, "Health check failed")
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
func NewDatasource(ctx context.Context, settings models.Settings) *Datasource {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: settings.AccessToken},
	)

	httpClient := oauth2.NewClient(ctx, src)

	if settings.GithubURL == "" {
		return &Datasource{
			client: githubv4.NewClient(httpClient),
		}
	}

	return &Datasource{
		client: githubv4.NewEnterpriseClient(fmt.Sprintf("%s/api/graphql", settings.GithubURL), httpClient),
	}
}

func newHealthResult(status backend.HealthStatus, message string) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
