package github

import (
	"context"
	"log"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// Issue represents a GitHub issue in a repository
type Issue struct {
	Title     string
	ClosedAt  githubv4.DateTime
	CreatedAt githubv4.DateTime
	Closed    bool
	Author    struct {
		User `graphql:"... on User"`
	}
}

// Issues is a slice of GitHub issues
type Issues []Issue

// Frame converts the list of issues to a Grafana DataFrame
func (c Issues) Frame() data.Frames {
	frame := data.NewFrame(
		"issues",
		data.NewField("title", nil, []string{}),
		data.NewField("author", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("closed", nil, []bool{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("closed_at", nil, []time.Time{}),
	)

	for _, v := range c {
		frame.AppendRow(
			v.Title,
			v.Author.User.Login,
			v.Author.User.Company,
			v.Closed,
			v.CreatedAt.Time,
			v.ClosedAt.Time,
		)
	}

	return data.Frames{frame}
}

// QueryGetIssues is the object representation of the graphql query for retrieving a paginated list of issues for a project
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
type QueryGetIssues struct {
	Repository struct {
		Issues struct {
			Nodes    Issues
			PageInfo PageInfo
		} `graphql:"issues(first: 100, after: $cursor, filterBy: $filters)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// GetAllIssues lists issues in a project. This function is slow and very prone to rate limiting.
func GetAllIssues(ctx context.Context, client Client, opts models.ListIssuesOptions) (Issues, error) {
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
		q := &QueryGetIssues{}
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

// GetIssuesInRange lists issues in a project given a time range.
func GetIssuesInRange(ctx context.Context, client Client, opts models.ListIssuesOptions, from time.Time, to time.Time) (Issues, error) {
	if opts.Filters == nil {
		opts.Filters = &githubv4.IssueFilters{}
	}

	opts.Filters.Since = &githubv4.DateTime{Time: from}

	issues, err := GetAllIssues(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	filtered := []Issue{}

	for i, v := range issues {
		if v.CreatedAt.After(from) && v.CreatedAt.Before(to) {
			filtered = append(filtered, issues[i])
		}
	}

	log.Printf("got %d issues, filtered to %d", len(issues), len(filtered))

	return filtered, nil
}
