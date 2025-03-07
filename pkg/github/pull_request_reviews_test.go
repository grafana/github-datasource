package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestListPullRequestReviews(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListPullRequestsOptions{
			Repository: "grafana",
			Owner:      "grafana",
			TimeField:  models.PullRequestClosedAt,
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "prCursor", "reviewCursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListPullRequestReviews{}),
	)

	_, err := GetPullRequestReviewsInRange(ctx, client, opts, time.Now().Add(-30*24*time.Hour), time.Now())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPullRequestReviewsDataFrame(t *testing.T) {
	openedAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	firstUser := models.User{
		ID:      "1",
		Login:   "testUser",
		Name:    "Test User",
		Company: "ACME corp",
		Email:   "user@example.com",
	}
	secondUser := models.User{
		ID:      "2",
		Login:   "testUser2",
		Name:    "Second User",
		Company: "ACME corp",
		Email:   "user2@example.com",
	}
	thirdUser := models.User{
		ID:      "3",
		Login:   "testUser3",
		Name:    "Third User",
		Company: "ACME corp",
		Email:   "user3@example.com",
	}

	pullRequestReviews := PullRequestReviews{
		{
			Number: 1,
			Title:  "PullRequest #1",
			URL:    "https://github.com/grafana/github-datasource/pulls/1",
			State:  githubv4.PullRequestStateOpen,
			Author: Author{
				User: firstUser,
			},
			Repository: Repository{
				NameWithOwner: "grafana/github-datasource",
			},
			Reviews: []Review{
				{
					URL:   "https://github.com/grafana/github-datasource/pull/1#pullrequestreview-2461579074",
					State: githubv4.PullRequestReviewStateApproved,
					Author: Author{
						User: secondUser,
					},
					CreatedAt: githubv4.DateTime{
						Time: openedAt,
					},
					UpdatedAt: githubv4.DateTime{
						Time: openedAt,
					},
					CommentsCount: 10,
				},
				{
					URL:   "https://github.com/grafana/github-datasource/pull/1#pullrequestreview-2461579074",
					State: githubv4.PullRequestReviewStateApproved,
					Author: Author{
						User: thirdUser,
					},
					CreatedAt: githubv4.DateTime{
						Time: openedAt,
					},
					UpdatedAt: githubv4.DateTime{
						Time: openedAt,
					},
					CommentsCount: 1,
				},
			},
		},
		{
			Number: 2,
			Title:  "PullRequest #2",
			URL:    "https://github.com/grafana/github-datasource/pulls/2",
			State:  githubv4.PullRequestStateOpen,
			Author: Author{
				User: secondUser,
			},
			Repository: Repository{
				NameWithOwner: "grafana/github-datasource",
			},
			Reviews: []Review{
				{
					URL:   "https://github.com/grafana/github-datasource/pull/1#pullrequestreview-2461579074",
					State: githubv4.PullRequestReviewStateApproved,
					Author: Author{
						User: firstUser,
					},
					CreatedAt: githubv4.DateTime{
						Time: openedAt,
					},
					UpdatedAt: githubv4.DateTime{
						Time: openedAt,
					},
					CommentsCount: 19,
				},
			},
		},
		{
			Number: 3,
			Title:  "PullRequest #2",
			URL:    "https://github.com/grafana/github-datasource/pulls/3",
			State:  githubv4.PullRequestStateOpen,
			Author: Author{
				User: secondUser,
			},
			Repository: Repository{
				NameWithOwner: "grafana/github-datasource",
			},
			Reviews: []Review{
				{
					URL:   "https://github.com/grafana/github-datasource/pull/1#pullrequestreview-2461579074",
					State: githubv4.PullRequestReviewStateApproved,
					Author: Author{
						User: firstUser,
					},
					CreatedAt: githubv4.DateTime{
						Time: openedAt,
					},
					UpdatedAt: githubv4.DateTime{
						Time: openedAt,
					},
					CommentsCount: 1,
				},
			},
		},
	}

	testutil.CheckGoldenFramer(t, "pull_request_reviews", pullRequestReviews)
}
