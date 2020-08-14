package github

import (
	"context"

	"github.com/shurcooL/githubv4"
)

// QueryListOrganizations is the GraphQL query for listing organizations
type QueryListOrganizations struct {
	Viewer struct {
		Organizations struct {
			Nodes    []Organization
			PageInfo PageInfo
		} `graphql:"organizations(first: 100, after: $cursor)"`
	}
}

// An Organization is a single GitHub organization
type Organization struct {
	Name string
}

// GetAllOrganizations lists the available organizations for the client
func GetAllOrganizations(ctx context.Context, client Client) ([]Organization, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
		}

		organizations = []Organization{}
	)

	for {
		q := &QueryListOrganizations{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}
		organizations = append(organizations, q.Viewer.Organizations.Nodes...)
		if !q.Viewer.Organizations.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Viewer.Organizations.PageInfo.EndCursor
	}

	return organizations, nil
}
