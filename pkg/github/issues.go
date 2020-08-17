package github

import (
	"context"

	"github.com/shurcooL/githubv4"
)

// QueryListIssues is the object representation of the graphql query for retrieving a paginated list of issues for a project
// {
//   repository(name: "grafana", owner: "grafana") {
//     issues(first: 100, filters: {}) {
//       nodes {
//         closedAt
//         title
//         closed
//         author {
//           ... on User {
//             id
//             email
//             login
//           }
//         }
//       }
//     }
//   }
// }
type QueryListIssues struct {
	Repository struct {
		Issues struct {
			Nodes    []Issue
			PageInfo PageInfo
		} `graphql:"issues(first: 100, after: $cursor, filterBy: $filters)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// Issue represents a GitHub issue in a repository
type Issue struct {
	Title    string
	ClosedAt githubv4.DateTime
	Closed   bool
	Author   struct {
		User `graphql:"... on User"`
	}
}

// ListIssuesOptions provides options when retrieving issues
type ListIssuesOptions struct {
	Repository string
	Owner      string
	Filters    *githubv4.IssueFilters
}

// ListIssues lists issues in a project. This function is slow and very prone to rate limiting.
func ListIssues(ctx context.Context, client Client, opts ListIssuesOptions) ([]Issue, error) {
	var (
		variables = map[string]interface{}{
			"cursor":  (*githubv4.String)(nil),
			"name":    githubv4.String(opts.Repository),
			"owner":   githubv4.String(opts.Owner),
			"filters": opts.Filters,
		}

		issues = []Issue{}
	)

	for {
		q := &QueryListIssues{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}
		issues = append(issues, q.Repository.Issues.Nodes...)
		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Issues.PageInfo.EndCursor
	}

	return issues, nil
}
