package githubclient

import (
	"context"
	"fmt"

	googlegithub "github.com/google/go-github/v53/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Client is a wrapper of GitHub clients that can access the GraphQL and rest API.
type Client struct {
	restClient    *googlegithub.Client
	graphqlClient *githubv4.Client
}

// New instantiates a new GitHub API client.
func New(ctx context.Context, settings models.Settings) (*Client, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: settings.AccessToken},
	)

	httpClient := oauth2.NewClient(ctx, src)

	if settings.GithubURL == "" {
		return &Client{
			restClient:    googlegithub.NewClient(httpClient),
			graphqlClient: githubv4.NewClient(httpClient),
		}, nil
	}

	restClient, err := googlegithub.NewEnterpriseClient(settings.GithubURL, settings.GithubURL, httpClient)
	if err != nil {
		return nil, fmt.Errorf("instantiating enterprise rest client: %w", err)
	}

	return &Client{
		restClient:    restClient,
		graphqlClient: githubv4.NewEnterpriseClient(fmt.Sprintf("%s/api/graphql", settings.GithubURL), httpClient),
	}, nil
}

// Query sends a query to the GitHub GraphQL API.
func (client *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return client.graphqlClient.Query(ctx, q, variables)
}

// ListWorkflows sends a request to the GitHub rest API to list the workflows in a specific repository.
func (client *Client) ListWorkflows(ctx context.Context, owner, repo string, opts *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error) {
	return client.restClient.Actions.ListWorkflows(ctx, owner, repo, opts)
}
