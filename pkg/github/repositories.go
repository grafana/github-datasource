package github

import (
	"context"

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

// ListRepositoriesOptions is the options for listing repositories
type ListRepositoriesOptions struct {
	Organization string
}

// GetAllRepositories retrieves all available repositories for an organization
func GetAllRepositories(ctx context.Context, client Client, opts ListRepositoriesOptions) ([]Repository, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"name":   githubv4.String(opts.Organization),
		}

		repos = []Repository{}
	)
	for {
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
