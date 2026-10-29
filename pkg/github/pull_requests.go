package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/github-datasource/pkg/dfutil"
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

// PullRequest is a GitHub pull request.
// The fields guarded by @include are optional field groups (see models.PullRequestField) and are only fetched when selected.
type PullRequest struct {
	Number    int64
	Title     string
	URL       string
	Additions int64 `graphql:"additions @include(if: $withDiffstat)"`
	Deletions int64 `graphql:"deletions @include(if: $withDiffstat)"`
	State     githubv4.PullRequestState
	Author    PullRequestAuthor `graphql:"author @include(if: $withAuthor)"`
	Closed    bool
	IsDraft   bool
	Locked    bool
	Merged    bool
	ClosedAt  githubv4.DateTime
	CreatedAt githubv4.DateTime
	UpdatedAt githubv4.DateTime
	MergedAt  githubv4.DateTime
	Mergeable githubv4.MergeableState `graphql:"mergeable @include(if: $withMergeable)"`
	MergedBy  *PullRequestAuthor      `graphql:"mergedBy @include(if: $withMergedBy)"`
	Labels    struct {
		Nodes []struct {
			Name string
		}
	} `graphql:"labels(first: 100) @include(if: $withLabels)"`
	Repository Repository
}

// PullRequests is a list of GitHub Pull Requests
type PullRequests []PullRequest

// PullRequestsWithFields is a list of GitHub Pull Requests fetched with only some of the optional field groups
type PullRequestsWithFields struct {
	PullRequests PullRequests
	Options      models.ListPullRequestsOptions
}

// pullRequestsFramer keeps the original frame for queries that don't select fields, and only the selected columns otherwise
func pullRequestsFramer(p PullRequests, opts models.ListPullRequestsOptions) dfutil.Framer {
	if opts.Fields == nil {
		return p
	}
	return PullRequestsWithFields{PullRequests: p, Options: opts}
}

// Frames converts the list of Pull Requests to a Grafana DataFrame
func (p PullRequests) Frames() data.Frames {
	return p.frames(models.ListPullRequestsOptions{})
}

// Frames converts the list of Pull Requests to a Grafana DataFrame with only the columns of the selected field groups
func (p PullRequestsWithFields) Frames() data.Frames {
	return p.PullRequests.frames(p.Options)
}

func (p PullRequests) frames(opts models.ListPullRequestsOptions) data.Frames {
	var (
		withDiffstat  = opts.Includes(models.PullRequestFieldDiffstat)
		withAuthor    = opts.Includes(models.PullRequestFieldAuthor)
		withMergeable = opts.Includes(models.PullRequestFieldMergeable)
		withMergedBy  = opts.Includes(models.PullRequestFieldMergedBy)
		withLabels    = opts.Includes(models.PullRequestFieldLabels)
	)

	openTime := data.NewField("open_time", nil, []float64{})
	openTime.Config = &data.FieldConfig{
		Unit: "s", // The values are in seconds
	}

	frame := data.NewFrame(
		"pull_requests",
		data.NewField("number", nil, []int64{}),
		data.NewField("title", nil, []string{}),
		data.NewField("url", nil, []string{}),
	)
	if withDiffstat {
		frame.Fields = append(frame.Fields,
			data.NewField("additions", nil, []int64{}),
			data.NewField("deletions", nil, []int64{}),
		)
	}
	frame.Fields = append(frame.Fields,
		data.NewField("repository", nil, []string{}),
		data.NewField("state", nil, []string{}),
	)
	if withAuthor {
		frame.Fields = append(frame.Fields,
			data.NewField("author_name", nil, []string{}),
			data.NewField("author_login", nil, []string{}),
			data.NewField("author_email", nil, []string{}),
			data.NewField("author_company", nil, []string{}),
		)
	}
	frame.Fields = append(frame.Fields,
		data.NewField("closed", nil, []bool{}),
		data.NewField("is_draft", nil, []bool{}),
		data.NewField("locked", nil, []bool{}),
		data.NewField("merged", nil, []bool{}),
	)
	if withMergeable {
		frame.Fields = append(frame.Fields, data.NewField("mergeable", nil, []string{}))
	}
	frame.Fields = append(frame.Fields,
		data.NewField("closed_at", nil, []*time.Time{}),
		data.NewField("merged_at", nil, []*time.Time{}),
	)
	if withMergedBy {
		frame.Fields = append(frame.Fields,
			data.NewField("merged_by_name", nil, []*string{}),
			data.NewField("merged_by_login", nil, []*string{}),
			data.NewField("merged_by_email", nil, []*string{}),
			data.NewField("merged_by_company", nil, []*string{}),
		)
	}
	frame.Fields = append(frame.Fields,
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("created_at", nil, []time.Time{}),
		openTime,
	)
	if withLabels {
		frame.Fields = append(frame.Fields, data.NewField("labels", nil, []json.RawMessage{}))
	}

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

		row := []interface{}{v.Number, v.Title, v.URL}
		if withDiffstat {
			row = append(row, v.Additions, v.Deletions)
		}
		row = append(row, v.Repository.NameWithOwner, string(v.State))
		if withAuthor {
			row = append(row, v.Author.User.Name, v.Author.User.Login, v.Author.User.Email, v.Author.User.Company)
		}
		row = append(row, v.Closed, v.IsDraft, v.Locked, v.Merged)
		if withMergeable {
			row = append(row, string(v.Mergeable))
		}
		row = append(row, closedAt, mergedAt)
		if withMergedBy {
			row = append(row, mergedByName, mergedByLogin, mergedByEmail, mergedByCompany)
		}
		row = append(row, v.UpdatedAt.Time, v.CreatedAt.Time, secondsOpen)
		if withLabels {
			labels := make([]string, len(v.Labels.Nodes))
			for i, label := range v.Labels.Nodes {
				labels[i] = label.Name
			}

			labelsBytes, _ := json.Marshal(labels)
			row = append(row, json.RawMessage(labelsBytes))
		}
		frame.AppendRow(row...)
	}

	return data.Frames{frame}
}

// GetAllPullRequests uses the graphql search endpoint API to search all pull requests in the repository
func GetAllPullRequests(ctx context.Context, client models.Client, opts models.ListPullRequestsOptions) (PullRequests, error) {
	var (
		variables = map[string]interface{}{
			"cursor":        (*githubv4.String)(nil),
			"query":         githubv4.String(buildQuery(opts)),
			"withDiffstat":  githubv4.Boolean(opts.Includes(models.PullRequestFieldDiffstat)),
			"withMergeable": githubv4.Boolean(opts.Includes(models.PullRequestFieldMergeable)),
			"withAuthor":    githubv4.Boolean(opts.Includes(models.PullRequestFieldAuthor)),
			"withMergedBy":  githubv4.Boolean(opts.Includes(models.PullRequestFieldMergedBy)),
			"withLabels":    githubv4.Boolean(opts.Includes(models.PullRequestFieldLabels)),
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
		Fields:     opts.Fields,
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
