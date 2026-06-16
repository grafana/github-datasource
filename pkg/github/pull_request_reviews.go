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

// QueryListPullRequestReviews lists all pull request reviews in a repository
//
//			{
//			  search(query: "is:pr repo:grafana/grafana merged:2020-08-19..*", type: ISSUE, first: 100) {
//			    nodes {
//			      ... on PullRequest {
//	              Number
//	              Title
//	              URL
//	              State
//	              Author
//	              Repository
//				     reviews(first: 100) {
//			           createdAt
//			           updatedAt
//			           state
//			           url
//			           author {
//			             id
//			             login
//			             name
//			             company
//			             email
//			             url
//			           }
//			           comments(first: 0) {
//			             totalCount
//		               }
//				     }
//				  }
//				}
//			  }
//			}
type QueryListPullRequestReviews struct {
	Search struct {
		Nodes []struct {
			PullRequest struct {
				Number     int64
				Title      string
				URL        string
				State      githubv4.PullRequestState
				Author     Author
				Repository Repository
				Reviews    struct {
					Nodes []struct {
						Review struct {
							CreatedAt githubv4.DateTime
							UpdatedAt githubv4.DateTime
							URL       string
							Author    Author
							State     githubv4.PullRequestReviewState
							Comments  struct {
								TotalCount int64
							} `graphql:"comments(first: 0)"`
						} `graphql:"... on PullRequestReview"`
					}
					PageInfo models.PageInfo
				} `graphql:"reviews(first: 100, after: $reviewCursor)"`
			} `graphql:"... on PullRequest"`
		}
		PageInfo models.PageInfo
	} `graphql:"search(query: $query, type: ISSUE, first: 100, after: $prCursor)"`
}

type Author struct {
	User models.User `graphql:"... on User"`
}

type Review struct {
	CreatedAt     githubv4.DateTime
	UpdatedAt     githubv4.DateTime
	URL           string
	Author        Author
	State         githubv4.PullRequestReviewState
	CommentsCount int64
}

type PullRequestWithReviews struct {
	Number     int64
	Title      string
	State      githubv4.PullRequestState
	URL        string
	Author     Author
	Repository Repository
	Reviews    []Review
}

// PullRequestReviews is a list of GitHub Pull Request Reviews
type PullRequestReviews []PullRequestWithReviews

// Frames coverts the list of Pull Request Reviews to a Grafana DataFrame
func (prs PullRequestReviews) Frames() data.Frames {
	frame := data.NewFrame(
		"pull_request_reviews",
		data.NewField("pull_request_number", nil, []int64{}),
		data.NewField("pull_request_title", nil, []string{}),
		data.NewField("pull_request_state", nil, []string{}),
		data.NewField("pull_request_url", nil, []string{}),
		data.NewField("pull_request_author_name", nil, []string{}),
		data.NewField("pull_request_author_login", nil, []string{}),
		data.NewField("pull_request_author_email", nil, []string{}),
		data.NewField("pull_request_author_company", nil, []string{}),
		data.NewField("repository", nil, []string{}),
		data.NewField("review_author_name", nil, []string{}),
		data.NewField("review_author_login", nil, []string{}),
		data.NewField("review_author_email", nil, []string{}),
		data.NewField("review_author_company", nil, []string{}),
		data.NewField("review_url", nil, []string{}),
		data.NewField("review_state", nil, []string{}),
		data.NewField("review_comment_count", nil, []int64{}),
		data.NewField("review_updated_at", nil, []time.Time{}),
		data.NewField("review_created_at", nil, []time.Time{}),
	)

	for _, pr := range prs {
		for _, review := range pr.Reviews {
			frame.AppendRow(
				pr.Number,
				pr.Title,
				string(pr.State),
				pr.URL,
				pr.Author.User.Name,
				pr.Author.User.Login,
				pr.Author.User.Email,
				pr.Author.User.Company,
				pr.Repository.NameWithOwner,
				review.Author.User.Name,
				review.Author.User.Login,
				review.Author.User.Email,
				review.Author.User.Company,
				review.URL,
				string(review.State),
				review.CommentsCount,
				review.UpdatedAt.Time,
				review.CreatedAt.Time,
			)
		}
	}

	return data.Frames{frame}
}

// GetAllPullRequestReviews uses the graphql search endpoint API to search all pull requests in the repository
// and all reviews for those pull requests.
func GetAllPullRequestReviews(ctx context.Context, client models.Client, opts models.ListPullRequestsOptions) (PullRequestReviews, error) {
	var (
		variables = map[string]interface{}{
			"prCursor":     (*githubv4.String)(nil),
			"reviewCursor": (*githubv4.String)(nil),
			"query":        githubv4.String(buildQuery(opts)),
		}

		pullRequestReviews = PullRequestReviews{}
	)

	for {
		q := &QueryListPullRequestReviews{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		prs := make([]PullRequestWithReviews, len(q.Search.Nodes))

		for i, prNode := range q.Search.Nodes {
			pr := prNode.PullRequest

			prs[i] = PullRequestWithReviews{
				Number:     pr.Number,
				Title:      pr.Title,
				State:      pr.State,
				URL:        pr.URL,
				Author:     pr.Author,
				Repository: pr.Repository,
			}

			for {
				for _, reviewNode := range pr.Reviews.Nodes {
					review := reviewNode.Review

					prs[i].Reviews = append(prs[i].Reviews, Review{
						CreatedAt:     review.CreatedAt,
						UpdatedAt:     review.UpdatedAt,
						URL:           review.URL,
						Author:        review.Author,
						State:         review.State,
						CommentsCount: review.Comments.TotalCount,
					})
				}

				if !pr.Reviews.PageInfo.HasNextPage {
					variables["reviewCursor"] = (*githubv4.String)(nil)
					break
				}

				variables["reviewCursor"] = pr.Reviews.PageInfo.EndCursor
				if err := client.Query(ctx, q, variables); err != nil {
					return nil, errors.WithStack(err)
				}
			}
		}

		pullRequestReviews = append(pullRequestReviews, prs...)

		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables["prCursor"] = q.Search.PageInfo.EndCursor
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
