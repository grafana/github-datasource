package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestSearchIssues(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListIssuesOptions{
			Repository: "grafana",
			Owner:      "grafana",
			TimeField:  models.IssuetClosedAt,
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QuerySearchIssues{}),
	)

	_, err := GetIssuesInRange(ctx, client, opts, time.Now().Add(-30*24*time.Hour), time.Now())
	if err != nil {
		t.Fatal(err)
	}
}

func TestIssuesDataframe(t *testing.T) {
	createdAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	issues := Issues{
		Issue{
			Title:    "Issue #1",
			ClosedAt: githubv4.DateTime{},
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			Closed: false,
			Author: struct {
				User "graphql:\"... on User\""
			}{
				User: User{
					ID:      "1",
					Login:   "firstUser",
					Name:    "First User",
					Company: "ACME Corp",
					Email:   "first@example.com",
					URL:     "",
				},
			},
		},
		Issue{
			Title: "Issue #2",
			ClosedAt: githubv4.DateTime{
				Time: createdAt.Add(time.Hour * 6),
			},
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			Closed: true,
			Author: struct {
				User "graphql:\"... on User\""
			}{
				User: User{
					ID:      "2",
					Login:   "secondUser",
					Name:    "Second User",
					Company: "ACME Corp",
					Email:   "second@example.com",
					URL:     "",
				},
			},
		},
	}

	if err := testutil.CheckGoldenFramer("issues", issues); err != nil {
		t.Fatal(err)
	}
}
