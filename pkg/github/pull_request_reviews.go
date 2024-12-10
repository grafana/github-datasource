package github

import (
	"context"
	"fmt"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryListPullRequests lists all pull requests in a repository
//
//		{
//		  search(query: "is:pr repo:grafana/grafana merged:2020-08-19..*", type: ISSUE, first: 100) {
//		    nodes {
//		      ... on PullRequest {
//	         reviews(first: 100) {
//
//
//		      }
//		  }
//		}
type QueryListPullRequestReviews struct {
	Search struct {
		Nodes []struct {
			Reviews struct {
				Nodes []struct {
					Review Review `graphql:"... on PullRequestReview"`
				}
			} `graphql:"reviews(first: 100)"`
		}
		PageInfo models.PageInfo
	} `graphql:"search(query: $query, type: ISSUE, first: 100, after: $cursor)"`
}

type ReviewComments struct {
	TotalCount int64
}

type ReviewAuthor struct {
	User models.User `graphql:"... on User"`
}

type Review struct {
	Author   ReviewAuthor
	State    githubv4.PullRequestReviewState
	Comments ReviewComments `graphsql:"comments(first: 0)"`
}

// PullRequestReviews is a list of GitHub Pull Request Reviews
type PullRequestReviews []Review

// Frames coverts the list of Pull Request Reviews to a Grafana DataFrame
func (r PullRequestReviews) Frames() data.Frames {
	frame := data.NewFrame(
		"pull_request_reviews",
		data.NewField("state", nil, []string{}),
		data.NewField("author_name", nil, []string{}),
		data.NewField("author_login", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("comment_count", nil, []int64{}),
	)

	for _, v := range r {
		frame.AppendRow(
			string(v.State),
			v.Author.User.Name,
			v.Author.User.Login,
			v.Author.User.Company,
			v.Comments.TotalCount,
		)
	}

	return data.Frames{frame}
}

// GetAllPullRequestReviews uses the graphql search endpoint API to search all pull requests in the repository
// and all reviews for those pull requests.
func GetAllPullRequestReviews(ctx context.Context, client models.Client, opts models.ListPullRequestsOptions) (PullRequestReviews, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"query":  githubv4.String(buildQuery(opts)),
		}

		pullRequestReviews = PullRequestReviews{}
	)

	for {
		q := &QueryListPullRequestReviews{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		reviews := make(PullRequestReviews, 0)
		for _, pr := range q.Search.Nodes {
			for _, v := range pr.Reviews.Nodes {
				reviews = append(reviews, v.Review)
			}
		}

		pullRequestReviews = append(pullRequestReviews, reviews...)

		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Search.PageInfo.EndCursor
	}

	return pullRequestReviews, nil
}

// GetPullRequestReviewsInRange uses the graphql search endpoint API to find pull request reviews in the given time range.
func GetPullRequestReviewsInRange(ctx context.Context, client models.Client, opts models.ListPullRequestsOptions, from time.Time, to time.Time) (PullRequestReviews, error) {
	var q string

	if opts.TimeField != models.PullRequestNone {
		q = fmt.Sprintf("%s:%s..%s", opts.TimeField.String(), from.Format(time.RFC3339), to.Format(time.RFC3339))
	}

	if opts.Query != nil {
		q = fmt.Sprintf("%s %s", *opts.Query, q)
	}

	return GetAllPullRequestReviews(ctx, client, models.ListPullRequestsOptions{
		Repository: opts.Repository,
		Owner:      opts.Owner,
		TimeField:  opts.TimeField,
		Query:      &q,
	})
}
