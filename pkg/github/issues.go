package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
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

// Frames converts the list of issues to a Grafana DataFrame
func (c Issues) Frames() data.Frames {
	frame := data.NewFrame(
		"issues",
		data.NewField("title", nil, []string{}),
		data.NewField("author", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("closed", nil, []bool{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("closed_at", nil, []*time.Time{}),
	)

	for _, v := range c {
		var closedAt *time.Time
		if !v.ClosedAt.Time.IsZero() {
			closedAt = &v.ClosedAt.Time
		}

		frame.AppendRow(
			v.Title,
			v.Author.User.Login,
			v.Author.User.Company,
			v.Closed,
			v.CreatedAt.Time,
			closedAt,
		)
	}

	return data.Frames{frame}
}

// QuerySearchIssues is the object representation of the graphql query for retrieving a paginated list of issues using the search query
// {
//   search(query: "is:issue repo:grafana/grafana opened:2020-08-19..*", type: ISSUE, first: 100) {
//     nodes {
//       ... on PullRequest {
//         id
//         title
//       }
//   }
// }
type QuerySearchIssues struct {
	Search struct {
		Nodes []struct {
			Issue Issue `graphql:"... on Issue"`
		}
		PageInfo PageInfo
	} `graphql:"search(query: $query, type: ISSUE, first: 100, after: $cursor)"`
}

// GetIssuesInRange lists issues in a project given a time range.
func GetIssuesInRange(ctx context.Context, client Client, opts models.ListIssuesOptions, from time.Time, to time.Time) (Issues, error) {
	search := []string{
		"is:issue",
		fmt.Sprintf("repo:%s/%s", opts.Owner, opts.Repository),
		fmt.Sprintf("%s:%s..%s", opts.TimeField.String(), from.Format(time.RFC3339), to.Format(time.RFC3339)),
	}

	if opts.Query != nil {
		search = append(search, *opts.Query)
	}

	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"query":  githubv4.String(strings.Join(search, " ")),
		}

		issues = []Issue{}
	)

	for {
		q := &QuerySearchIssues{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}
		is := make([]Issue, len(q.Search.Nodes))

		for i, v := range q.Search.Nodes {
			is[i] = v.Issue
		}

		issues = append(issues, is...)

		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Search.PageInfo.EndCursor
	}

	return issues, nil
}
