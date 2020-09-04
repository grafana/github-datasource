package github

import (
	"context"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Datasource handles requests to GitHub
type Datasource struct {
	client *githubv4.Client
}

func (d *Datasource) HandleIssuesQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.IssueOptionsWithRepo(query.IssuesOptions, query.Owner, query.Repository)
	return GetIssuesInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandleCommitsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.CommitsOptionsWithRepo(query.CommitsOptions, query.Owner, query.Repository)
	return GetCommitsInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandleTagsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListTagsOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
	}

	if req.TimeRange.From.Unix() <= 0 && req.TimeRange.To.Unix() <= 0 {
		return GetAllTags(ctx, d.client, opt)
	}

	return GetTagsInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandleReleasesQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListReleasesOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
	}

	if req.TimeRange.From.Unix() <= 0 && req.TimeRange.To.Unix() <= 0 {
		return GetAllReleases(ctx, d.client, opt)
	}
	return GetReleasesInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandlePullRequestsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.PullRequestOptionsWithRepo(query.PullRequestsOptions, query.Owner, query.Repository)

	if req.TimeRange.From.Unix() <= 0 && req.TimeRange.To.Unix() <= 0 {
		return GetAllPullRequests(ctx, d.client, opt)
	}
	return GetPullRequestsInRange(ctx, d.client, opt, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandleContributorsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListContributorsOptions{
		Owner:      query.Owner,
		Repository: query.Repository,
		Query:      query.ContributorsOptions.Query,
	}

	return GetAllContributors(ctx, d.client, opt)
}

func (d *Datasource) HandleLabelsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	opt := models.ListLabelsOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
		Query:      query.LabelsOptions.Query,
	}

	return GetAllLabels(ctx, d.client, opt)
}

// CheckHealth calls frequently used endpoints to determine if the client has sufficient privileges
func (d *Datasource) CheckHealth(ctx context.Context) error {
	_, err := GetAllRepositories(ctx, d.client, models.ListRepositoriesOptions{
		Organization: "grafana",
	})

	if err != nil {
		return errors.Wrap(err, "failed to list repositories in the Grafana organization")
	}

	return nil
}

func NewDatasource(ctx context.Context, settings models.Settings) *Datasource {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: settings.AccessToken},
	)

	httpClient := oauth2.NewClient(ctx, src)

	return &Datasource{
		client: githubv4.NewClient(httpClient),
	}
}
