package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestListMilestones(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListMilestonesOptions{
			Repository: "grafana",
			Owner:      "grafana",
			Query:      "test",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "name", "owner", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListMilestones{}),
	)

	_, err := GetAllMilestones(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMilestonesDataframe(t *testing.T) {
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

	milestones := Milestones{
		{
			Closed: false,
			Creator: struct {
				User models.User "graphql:\"... on User\""
			}{
				User: firstUser,
			},
			DueOn: githubv4.DateTime{
				Time: openedAt.Add(100 * time.Hour),
			},
			ClosedAt: githubv4.DateTime{},
			CreatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			State: githubv4.MilestoneStateOpen,
			Title: "first milestone",
		},
		{
			Closed: true,
			Creator: struct {
				User models.User "graphql:\"... on User\""
			}{
				User: secondUser,
			},
			DueOn: githubv4.DateTime{
				Time: openedAt.Add(100 * time.Hour),
			},
			ClosedAt: githubv4.DateTime{
				Time: openedAt.Add(10 * time.Hour),
			},
			CreatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			State: githubv4.MilestoneStateClosed,
			Title: "second milestone",
		},
		{
			Closed: false,
			Creator: struct {
				User models.User "graphql:\"... on User\""
			}{
				User: secondUser,
			},
			DueOn: githubv4.DateTime{
				Time: openedAt.Add(120 * time.Hour),
			},
			ClosedAt: githubv4.DateTime{},
			CreatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			State: githubv4.MilestoneStateOpen,
			Title: "third milestone",
		},
	}

	testutil.CheckGoldenFramer(t, "milestones", milestones)
}
