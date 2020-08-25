package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestGetIssues(t *testing.T) {
	var (
		ctx = context.Background()
	)

	t.Run("With no filters", func(t *testing.T) {
		opts := models.ListIssuesOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}

		testIssueVariables := func(t *testing.T, variables map[string]interface{}) {
			if _, ok := variables["filters"]; !ok {
				t.Fatal("the 'filters' key must always be included in the variables map")
			}
			if variables["filters"] != (*githubv4.IssueFilters)(nil) {
				t.Fatal("if filters is not provided in the arguments, it should be nil in the variables")
			}

			testutil.EnsureKeysAreSet(t, variables, "name", "owner", "filters")
		}

		client := testutil.NewTestClient(t,
			testIssueVariables,
			testutil.GetTestQueryFunction(&QueryGetIssues{}),
		)

		_, err := GetAllIssues(ctx, client, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("With a 'since' filter", func(t *testing.T) {
		opts := models.ListIssuesOptions{
			Repository: "grafana",
			Owner:      "grafana",
			Filters: &githubv4.IssueFilters{
				Since: &githubv4.DateTime{
					Time: time.Now().Add(-30 * 24 * time.Hour),
				},
			},
		}

		testIssueVariables := func(t *testing.T, variables map[string]interface{}) {
			filters, ok := variables["filters"]
			if !ok {
				t.Fatal("Filters are included in the options, but not in the variables given to the query")
			}

			switch filters.(type) {
			case *githubv4.IssueFilters:
				break
			default:
				t.Fatal("Unexpected type of variables['filters']. Expected '*githubv4.IssueFilters'")
			}

			testutil.EnsureKeysAreSet(t, variables, "name", "owner", "filters")
		}

		client := testutil.NewTestClient(t,
			testIssueVariables,
			testutil.GetTestQueryFunction(&QueryGetIssues{}),
		)

		_, err := GetAllIssues(ctx, client, opts)
		if err != nil {
			t.Fatal(err)
		}
	})
}

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

	_, err := SearchIssues(ctx, client, opts, time.Now().Add(-30*24*time.Hour), time.Now())
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
