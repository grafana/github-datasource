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

// QueryListPullRequests lists all pull requests in a repository
//
//	{
//	  search(query: "is:pr repo:grafana/grafana merged:2020-08-19..*", type: ISSUE, first: 100) {
//	    nodes {
//	      ... on PullRequest {
//	        id
//	        title
//	      }
//	  }
//	}
type QueryListPullRequests struct {
	Search struct {
		Nodes []struct {
			PullRequest PullRequest `graphql:"... on PullRequest"`
		}
		PageInfo models.PageInfo
	} `graphql:"search(query: $query, type: ISSUE, first: 100, after: $cursor)"`
}

// PullRequestAuthor is the structure of the Author object in a Pull Request (which requires a graphQL object expansion on `User`)
type PullRequestAuthor struct {
	User models.User `graphql:"... on User"`
}

// PullRequest is a GitHub pull request
type PullRequest struct {
	Number     int64
	Title      string
	URL        string
	Additions  int64
	Deletions  int64
	State      githubv4.PullRequestState
	Author     PullRequestAuthor
	Closed     bool
	IsDraft    bool
	Locked     bool
	Merged     bool
	ClosedAt   githubv4.DateTime
	CreatedAt  githubv4.DateTime
	UpdatedAt  githubv4.DateTime
	MergedAt   githubv4.DateTime
	Mergeable  githubv4.MergeableState
	MergedBy   *PullRequestAuthor
	Repository Repository
}

// PullRequests is a list of GitHub Pull Requests
type PullRequests []PullRequest

// Frames converts the list of Pull Requests to a Grafana DataFrame
func (p PullRequests) Frames() data.Frames {
	openTime := data.NewField("open_time", nil, []float64{})
	openTime.Config = &data.FieldConfig{
		Unit: "s", // The values are in seconds
	}

	frame := data.NewFrame(
		"pull_requests",
		data.NewField("number", nil, []int64{}),
		data.NewField("title", nil, []string{}),
		data.NewField("url", nil, []string{}),
		data.NewField("additions", nil, []int64{}),
		data.NewField("deletions", nil, []int64{}),
		data.NewField("repository", nil, []string{}),
		data.NewField("state", nil, []string{}),
		data.NewField("author_name", nil, []string{}),
		data.NewField("author_login", nil, []string{}),
		data.NewField("author_email", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("closed", nil, []bool{}),
		data.NewField("is_draft", nil, []bool{}),
		data.NewField("locked", nil, []bool{}),
		data.NewField("merged", nil, []bool{}),
		data.NewField("mergeable", nil, []string{}),
		data.NewField("closed_at", nil, []*time.Time{}),
		data.NewField("merged_at", nil, []*time.Time{}),
		data.NewField("merged_by_name", nil, []*string{}),
		data.NewField("merged_by_login", nil, []*string{}),
		data.NewField("merged_by_email", nil, []*string{}),
		data.NewField("merged_by_company", nil, []*string{}),
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("created_at", nil, []time.Time{}),
		openTime,
	)

	for _, v := range p {
		var (
			closedAt        *time.Time
			mergedAt        *time.Time
			mergedByName    *string
			mergedByLogin   *string
			mergedByEmail   *string
			mergedByCompany *string
			secondsOpen     = time.Now().UTC().Sub(v.CreatedAt.UTC()).Round(time.Second).Seconds()
		)

		if !v.ClosedAt.IsZero() {
			t := v.ClosedAt.Time
			closedAt = &t
		}

		if !v.MergedAt.IsZero() {
			t := v.MergedAt.Time
			mergedAt = &t
		}

		if closedAt != nil {
			secondsOpen = v.ClosedAt.UTC().Sub(v.CreatedAt.UTC()).Seconds()
		}

		if mergedAt != nil {
			secondsOpen = v.MergedAt.UTC().Sub(v.CreatedAt.UTC()).Seconds()
		}

		mergedBy := v.MergedBy
		if mergedBy != nil {
			mergedByNameT := v.MergedBy.User.Name
			if len(mergedByNameT) != 0 {
				mergedByName = &mergedByNameT
			}

			mergedByLoginT := v.MergedBy.User.Login
			if len(mergedByLoginT) != 0 {
				mergedByLogin = &mergedByLoginT
			}

			mergedByEmailT := v.MergedBy.User.Email
			if len(mergedByEmailT) != 0 {
				mergedByEmail = &mergedByEmailT
			}

			mergedByCompanyT := v.MergedBy.User.Company
			if len(mergedByCompanyT) != 0 {
				mergedByCompany = &mergedByCompanyT
			}
		}

		frame.AppendRow(
			v.Number,
			v.Title,
			v.URL,
			v.Additions,
			v.Deletions,
			v.Repository.NameWithOwner,
			string(v.State),
			v.Author.User.Name,
			v.Author.User.Login,
			v.Author.User.Email,
			v.Author.User.Company,
			v.Closed,
			v.IsDraft,
			v.Locked,
			v.Merged,
			string(v.Mergeable),
			closedAt,
			mergedAt,
			mergedByName,
			mergedByLogin,
			mergedByEmail,
			mergedByCompany,
			v.UpdatedAt.Time,
			v.CreatedAt.Time,
			secondsOpen,
		)
	}

	return data.Frames{frame}
}

// GetAllPullRequests uses the graphql search endpoint API to search all pull requests in the repository
func GetAllPullRequests(ctx context.Context, client models.Client, opts models.ListPullRequestsOptions) (PullRequests, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"query":  githubv4.String(buildQuery(opts)),
		}

		pullRequests = []PullRequest{}
	)

	for {
		q := &QueryListPullRequests{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}
		prs := make([]PullRequest, len(q.Search.Nodes))

		for i, v := range q.Search.Nodes {
			prs[i] = v.PullRequest
		}

		pullRequests = append(pullRequests, prs...)

		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Search.PageInfo.EndCursor
	}

	return pullRequests, nil
}

// GetPullRequestsInRange uses the graphql search endpoint API to find pull requests in the given time range.
func GetPullRequestsInRange(ctx context.Context, client models.Client, opts models.ListPullRequestsOptions, from time.Time, to time.Time) (PullRequests, error) {
	var q string

	if opts.TimeField != models.PullRequestNone {
		q = fmt.Sprintf("%s:%s..%s", opts.TimeField.String(), from.Format(time.RFC3339), to.Format(time.RFC3339))
	}

	if opts.Query != nil {
		q = fmt.Sprintf("%s %s", *opts.Query, q)
	}

	return GetAllPullRequests(ctx, client, models.ListPullRequestsOptions{
		Repository: opts.Repository,
		Owner:      opts.Owner,
		TimeField:  opts.TimeField,
		Query:      &q,
	})
}

// buildQuery builds the "query" field for Pull Request searches
func buildQuery(opts models.ListPullRequestsOptions) string {
	search := []string{
		"is:pr",
	}

	if opts.Repository == "" {
		search = append(search, fmt.Sprintf("org:%s", opts.Owner))
	} else {
		search = append(search, fmt.Sprintf("repo:%s/%s", opts.Owner, opts.Repository))
	}

	if opts.Query != nil {
		queryString, err := InterPolateMacros(*opts.Query)
		if err != nil {
			return strings.Join(search, " ")
		}
		search = append(search, queryString)
	}

	return strings.Join(search, " ")
}
