package github

import (
	"context"
	"net/http"

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

func (d *Datasource) HandleAuth(http.ResponseWriter, *http.Request)         {}
func (d *Datasource) HandleAuthCallback(http.ResponseWriter, *http.Request) {}

func (d *Datasource) HandleIssuesQuery(context context.Context, query *models.Query, req backend.DataQuery) (backend.DataResponse, error) {
	return backend.DataResponse{}, nil
}

func (d *Datasource) HandleCommitsQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (backend.DataResponse, error) {
	commits, err := GetCommitsInRange(ctx, d.client, ListCommitsOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
		Ref:        query.Ref,
	}, req.TimeRange.From, req.TimeRange.To)

	if err != nil {
		return backend.DataResponse{}, err
	}

	frames, err := commits.ToDataFrame()
	if err != nil {
		return backend.DataResponse{}, err
	}

	return backend.DataResponse{
		Frames: frames,
	}, nil
}

func (d *Datasource) HandleTagsQuery(context.Context, *models.Query, backend.DataQuery) (backend.DataResponse, error) {
	return backend.DataResponse{}, nil
}

func (d *Datasource) HandleReleasesQuery(ctx context.Context, query *models.Query, req backend.DataQuery) (backend.DataResponse, error) {
	releases, err := GetReleasesInRange(ctx, d.client, ListReleasesOptions{
		Repository: query.Repository,
		Owner:      query.Owner,
	}, req.TimeRange.From, req.TimeRange.To)

	if err != nil {
		return backend.DataResponse{}, err
	}

	frames, err := releases.ToDataFrame()
	if err != nil {
		return backend.DataResponse{}, err
	}

	return backend.DataResponse{
		Frames: frames,
	}, nil

}

func (d *Datasource) HandleContributorsQuery(context.Context, *models.Query, backend.DataQuery) (backend.DataResponse, error) {
	return backend.DataResponse{}, nil
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
