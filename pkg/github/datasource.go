package github

import (
	"context"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Datasource handles requests to GitHub
type Datasource struct {
	client      *githubv4.Client
	oauthConfig *oauth2.Config
}

func (d *Datasource) HandleIssuesQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return GetIssuesInRange(ctx, d.client, query.IssuesOptions, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandleCommitsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return GetCommitsInRange(ctx, d.client, query.CommitsOptions, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandleTagsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return GetTagsInRange(ctx, d.client, query.TagsOptions, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandleReleasesQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return GetReleasesInRange(ctx, d.client, query.ReleasesOptions, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandlePullRequestsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return GetPullRequestsInRange(ctx, d.client, query.PullRequestsOptions, req.TimeRange.From, req.TimeRange.To)
}

func (d *Datasource) HandleContributorsQuery(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error) {
	return nil, nil
}

func NewDatasource(ctx context.Context, settings models.Settings) *Datasource {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: settings.AccessToken},
	)

	httpClient := oauth2.NewClient(ctx, src)

	return &Datasource{
		client:      githubv4.NewClient(httpClient),
		oauthConfig: GetOAuthConfig(settings.ClientID, settings.ClientSecret),
	}
}
