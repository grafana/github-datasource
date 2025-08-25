package github

import (
	"context"
	"encoding/json"
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
	Number    int64
	Title     string
	ClosedAt  githubv4.DateTime
	CreatedAt githubv4.DateTime
	UpdatedAt githubv4.DateTime
	Closed    bool
	Labels    struct {
		Nodes []struct {
			Name string
		}
	} `graphql:"labels(first: 100)"`
	Assignees struct {
		Nodes []struct {
			models.User
		}
	} `graphql:"assignees(first: 10)"`
	Author struct {
		models.User `graphql:"... on User"`
	}
	Repository Repository
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
		data.NewField("repo", nil, []string{}),
		data.NewField("number", nil, []int64{}),
		data.NewField("closed", nil, []bool{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("closed_at", nil, []*time.Time{}),
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("labels", nil, []json.RawMessage{}),
		data.NewField("assignees", nil, []json.RawMessage{}),
	)

	for _, v := range c {
		var closedAt *time.Time
		if !v.ClosedAt.IsZero() {
			t := v.ClosedAt.Time
			closedAt = &t
		}

		labels := make([]string, len(v.Labels.Nodes))
		for i, label := range v.Labels.Nodes {
			labels[i] = label.Name
		}

		assignees := make([]string, len(v.Assignees.Nodes))
		for i, assignee := range v.Assignees.Nodes {
			assignees[i] = assignee.Login
		}

		labelsBytes, _ := json.Marshal(labels)
		rawLabelArray := json.RawMessage(labelsBytes)

		assigneesBytes, _ := json.Marshal(assignees)
		rawAssigneesArray := json.RawMessage(assigneesBytes)

		frame.AppendRow(
			v.Title,
			v.Author.Login,
			v.Author.Company,
			v.Repository.NameWithOwner,
			v.Number,
			v.Closed,
			v.CreatedAt.Time,
			closedAt,
			v.UpdatedAt.Time,
			rawLabelArray,
			rawAssigneesArray,
		)
	}

	return data.Frames{frame}
}

// QuerySearchIssues is the object representation of the graphql query for retrieving a paginated list of issues using the search query
//
//	{
//	  search(query: "is:issue repo:grafana/grafana opened:2020-08-19..*", type: ISSUE, first: 100) {
//	    nodes {
//	      ... on PullRequest {
//	        id
//	        title
//	      }
//	  }
//	}
type QuerySearchIssues struct {
	Search struct {
		Nodes []struct {
			Issue Issue `graphql:"... on Issue"`
		}
		PageInfo models.PageInfo
	} `graphql:"search(query: $query, type: ISSUE, first: 100, after: $cursor)"`
}

// GetIssuesInRange lists issues in a project given a time range.
func GetIssuesInRange(ctx context.Context, client models.Client, opts models.ListIssuesOptions, from time.Time, to time.Time) (Issues, error) {

	filter := fmt.Sprintf("repo:%s/%s", opts.Owner, opts.Repository)
	if opts.Repository == "" {
		filter = fmt.Sprintf("owner:%s", opts.Owner)
	}

	search := []string{
		"is:issue",
		filter,
		fmt.Sprintf("%s:%s..%s", opts.TimeField.String(), from.Format(time.RFC3339), to.Format(time.RFC3339)),
	}

	if opts.Query != nil {
		queryString, err := InterPolateMacros(*opts.Query)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		search = append(search, queryString)
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
