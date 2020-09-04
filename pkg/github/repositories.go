package github

import (
	"context"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/shurcooL/githubv4"
)

// QueryListRepositories is the GraphQL query for retrieving a list of repositories for an organization
type QueryListRepositories struct {
	Organization struct {
		Repositories struct {
			Nodes    []Repository
			PageInfo PageInfo
		} `graphql:"repositories(first: 100, after: $cursor)"`
	} `graphql:"organization(login: $name)"`
}

// GetAllRepositories retrieves all available repositories for an organization
func GetAllRepositories(ctx context.Context, client Client, opts models.ListRepositoriesOptions) ([]Repository, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"name":   githubv4.String(opts.Organization),
		}

		repos = []Repository{}
	)
	for i := 0; i < PageNumberLimit; i++ {
		q := &QueryListRepositories{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}

		repos = append(repos, q.Organization.Repositories.Nodes...)
		if !q.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Organization.Repositories.PageInfo.EndCursor
	}

	return repos, nil
}
