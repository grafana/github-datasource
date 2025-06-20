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
			TimeField:  models.IssueClosedAt,
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
			Number: 1,
			Title:  "Issue #1",
			ClosedAt: githubv4.DateTime{
				Time: time.Time{},
			},
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			Closed: false,
			Labels: struct { Nodes []struct{ Name string } }{
				Nodes: []struct{ Name string }{
					{ Name: "bug" },
					{ Name: "help wanted" },
				},
			},
			Author: struct { models.User "graphql:\"... on User\"" }{
				User: models.User{
					ID:      "1",
					Login:   "firstUser",
					Name:    "First User",
					Company: "ACME Corp",
					Email:   "first@example.com",
					URL:     "",
				},
			},
			Assignees: struct { Nodes []struct { models.User } }{
				Nodes: []struct { models.User }{
					{ User: models.User{ Login: "firstUser" } },
					{ User: models.User{ Login: "secondUser" } },
				},
			},
			Repository: Repository{
				Name: "grafana",
				Owner: struct{ Login string }{
					Login: "grafana",
				},
				NameWithOwner: "grafana/grafana",
				URL:           "github.com/grafana/grafana",
				ForkCount:     10,
				IsFork:        true,
				IsMirror:      true,
				IsPrivate:     false,
				CreatedAt: githubv4.DateTime{
					Time: createdAt,
				},
			},
		},
		Issue{
			Number: 2,
			Title:  "Issue #2",
			ClosedAt: githubv4.DateTime{
				Time: createdAt.Add(time.Hour * 6),
			},
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: createdAt.Add(time.Hour * 6),
			},
			Closed: true,
			Labels: struct { Nodes []struct{ Name string } }{
				Nodes: []struct{ Name string }{
					{ Name: "enhancement" },
				},
			},
			Author: struct { models.User "graphql:\"... on User\"" }{
				User: models.User{
					ID:      "2",
					Login:   "secondUser",
					Name:    "Second User",
					Company: "ACME Corp",
					Email:   "second@example.com",
					URL:     "",
				},
			},
			Assignees: struct { Nodes []struct { models.User } }{
				Nodes: []struct { models.User }{
					{ User: models.User{ Login: "firstUser" } },
				},
			},
			Repository: Repository{
				Name: "grafana",
				Owner: struct{ Login string }{
					Login: "grafana",
				},
				NameWithOwner: "grafana/grafana",
				URL:           "github.com/grafana/grafana",
				ForkCount:     10,
				IsFork:        true,
				IsMirror:      true,
				IsPrivate:     false,
				CreatedAt: githubv4.DateTime{
					Time: createdAt,
				},
			},
		},
		Issue{
			Number: 3,
			Title:  "Issue #3",
			ClosedAt: githubv4.DateTime{
				Time: time.Time{},
			},
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			Closed: false,
			Labels: struct { Nodes []struct{ Name string } }{
				Nodes: []struct{ Name string }{},
			},
			Author: struct { models.User "graphql:\"... on User\"" }{
				User: models.User{
					ID:      "3",
					Login:   "firstUser",
					Name:    "First User",
					Company: "ACME Corp",
					Email:   "first@example.com",
					URL:     "",
				},
			},
			Assignees: struct { Nodes []struct { models.User }}{
				Nodes: []struct { models.User }{},
			},
			Repository: Repository{
				Name: "grafana",
				Owner: struct{ Login string }{
					Login: "grafana",
				},
				NameWithOwner: "grafana/grafana",
				URL:           "github.com/grafana/grafana",
				ForkCount:     10,
				IsFork:        true,
				IsMirror:      true,
				IsPrivate:     false,
				CreatedAt: githubv4.DateTime{
					Time: createdAt,
				},
			},
		},
	}

	testutil.CheckGoldenFramer(t, "issues", issues)
}
