package github

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// An Organization is a single GitHub organization
type Organization struct {
	Name string
}

// Organizations is a slice of GitHub Organizations
type Organizations []Organization

// Frames converts the list of Organizations to a Grafana DataFrame
func (c Organizations) Frames() data.Frames {
	return data.Frames{}
}

// QueryListOrganizations is the GraphQL query for listing organizations
type QueryListOrganizations struct {
	Viewer struct {
		Organizations struct {
			Nodes    []Organization
			PageInfo PageInfo
		} `graphql:"organizations(first: 100, after: $cursor)"`
	}
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
