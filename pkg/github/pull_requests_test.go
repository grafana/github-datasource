package github

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-github-datasource/pkg/testutil"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/shurcooL/githubv4"
)

func TestListPullRequests(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListPullRequestsInRangeOptions{
			Repository: "grafana",
			Owner:      "grafana",
			TimeField:  models.PullRequestClosedAt,
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListPullRequests{}),
	)

	_, err := GetPullRequestsInRange(ctx, client, opts, time.Now().Add(-30*24*time.Hour), time.Now())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPullRequestsDataFrame(t *testing.T) {
	openedAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	firstUser := User{
		ID:      "1",
		Login:   "testUser",
		Name:    "Test User",
		Company: "ACME corp",
		Email:   "user@example.com",
	}
	secondUser := User{
		ID:      "2",
		Login:   "testUser2",
		Name:    "Second User",
		Company: "ACME corp",
		Email:   "user2@example.com",
	}

	pullRequests := PullRequests{
		PullRequest{
			Title: "PullRequest #1",
			State: githubv4.PullRequestStateOpen,
			Author: PullRequestAuthor{
				User: firstUser,
			},
			Closed:   false,
			IsDraft:  false,
			Locked:   false,
			Merged:   false,
			ClosedAt: githubv4.DateTime{},
			CreatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			MergedAt: githubv4.DateTime{
				Time: time.Time{},
			},
			Mergeable: githubv4.MergeableStateMergeable,
			MergedBy:  nil,
		},
		PullRequest{
			Title: "PullRequest #2",
			State: githubv4.PullRequestStateOpen,
			Author: PullRequestAuthor{
				User: secondUser,
			},
			Closed:   true,
			IsDraft:  false,
			Locked:   false,
			Merged:   true,
			ClosedAt: githubv4.DateTime{},
			CreatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: openedAt.Add(time.Hour * 2),
			},
			MergedAt: githubv4.DateTime{
				Time: openedAt.Add(time.Hour * 2),
			},
			Mergeable: githubv4.MergeableStateMergeable,
			MergedBy: &PullRequestAuthor{
				User: firstUser,
			},
		},
	}

	dr := backend.DataResponse{
		Frames: pullRequests.Frame(),
	}

	if err := experimental.CheckGoldenDataResponse(filepath.Join("testdata", "pullRequests.golden.txt"), &dr, UpdateGoldenFiles); err != nil {
		t.Fatal(err)
	}
}
