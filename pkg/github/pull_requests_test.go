package github

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/shurcooL/githubv4"
)

func TestListPullRequests(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListPullRequestsOptions{
			Repository: "grafana",
			Owner:      "grafana",
			TimeField:  models.PullRequestClosedAt,
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "cursor", "withDiffstat", "withMergeable", "withAuthor", "withMergedBy", "withLabels")

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

	pullRequests := PullRequests{
		{
			Number: 1,
			Title:  "PullRequest #1",
			URL:    "https://github.com/grafana/github-datasource/pulls/1",
			State:  githubv4.PullRequestStateOpen,
			Author: PullRequestAuthor{
				User: firstUser,
			},
			Repository: Repository{
				NameWithOwner: "grafana/github-datasource",
			},
			Closed:    true,
			IsDraft:   false,
			Locked:    false,
			Merged:    true,
			Additions: 5,
			Deletions: 1,
			CreatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			MergedAt: githubv4.DateTime{
				Time: openedAt.Add(100 * time.Minute),
			},
			ClosedAt: githubv4.DateTime{
				Time: openedAt.Add(100 * time.Minute),
			},
			Mergeable: githubv4.MergeableStateMergeable,
			MergedBy:  nil,
			Labels: struct {
				Nodes []struct{ Name string }
			}{
				Nodes: []struct{ Name string }{
					{Name: "bug"},
					{Name: "enhancement"},
				},
			},
		},
		{
			Number: 2,
			Title:  "PullRequest #2",
			URL:    "https://github.com/grafana/github-datasource/pulls/2",
			State:  githubv4.PullRequestStateOpen,
			Author: PullRequestAuthor{
				User: secondUser,
			},
			Repository: Repository{
				NameWithOwner: "grafana/github-datasource",
			},
			Closed:  true,
			IsDraft: false,
			Locked:  false,
			Merged:  true,
			MergedAt: githubv4.DateTime{
				Time: openedAt.Add(100 * time.Minute),
			},
			ClosedAt: githubv4.DateTime{
				Time: openedAt.Add(100 * time.Minute),
			},
			CreatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: openedAt.Add(time.Hour * 2),
			},
			Mergeable: githubv4.MergeableStateMergeable,
			MergedBy: &PullRequestAuthor{
				User: firstUser,
			},
			Labels: struct {
				Nodes []struct{ Name string }
			}{
				Nodes: []struct{ Name string }{
					{Name: "documentation"},
				},
			},
		},
		{
			Number: 3,
			Title:  "PullRequest #2",
			URL:    "https://github.com/grafana/github-datasource/pulls/3",
			State:  githubv4.PullRequestStateOpen,
			Author: PullRequestAuthor{
				User: secondUser,
			},
			Repository: Repository{
				NameWithOwner: "grafana/github-datasource",
			},
			Closed:  false,
			IsDraft: false,
			Locked:  false,
			Merged:  false,
			MergedAt: githubv4.DateTime{
				Time: openedAt.Add(100 * time.Minute),
			},
			ClosedAt: githubv4.DateTime{},
			CreatedAt: githubv4.DateTime{
				Time: openedAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: openedAt.Add(time.Hour * 2),
			},
			Mergeable: githubv4.MergeableStateMergeable,
			Labels: struct {
				Nodes []struct{ Name string }
			}{
				Nodes: []struct{ Name string }{},
			},
		},
	}

	testutil.CheckGoldenFramer(t, "pull_requests", pullRequests)
	testutil.CheckGoldenFramer(t, "pull_requests_selected_fields", PullRequestsWithFields{
		PullRequests: pullRequests,
		Options:      models.ListPullRequestsOptions{Fields: []models.PullRequestField{models.PullRequestFieldDiffstat, models.PullRequestFieldLabels}},
	})
	testutil.CheckGoldenFramer(t, "pull_requests_no_additional_fields", PullRequestsWithFields{
		PullRequests: pullRequests,
		Options:      models.ListPullRequestsOptions{Fields: []models.PullRequestField{}},
	})
}

func TestBuildQuery(t *testing.T) {
	t.Run("Searching pull requests with a Repository and organization should use the repo field", func(t *testing.T) {
		opts := models.ListPullRequestsOptions{
			Owner:      "grafana",
			Repository: "github-datasource",
		}

		var (
			result = buildQuery(opts)
			expect = "is:pr repo:grafana/github-datasource"
		)

		if result != expect {
			t.Fatalf("Unexpected result from buildQuery. Expected '%s', received '%s'", expect, result)
		}
	})

	t.Run("Issue #61 - Searching pull requests without a Repository should search the entire org", func(t *testing.T) {
		opts := models.ListPullRequestsOptions{
			Owner:      "test",
			Repository: "",
		}

		var (
			result = buildQuery(opts)
			expect = "is:pr org:test"
		)
		if result != expect {
			t.Fatalf("Unexpected result from buildQuery. Expected '%s', received '%s'", expect, result)
		}
	})
}

func TestListPullRequestsFieldVariables(t *testing.T) {
	tests := []struct {
		name   string
		fields []models.PullRequestField
		want   map[string]bool
	}{
		{
			name:   "fields unset fetches every group",
			fields: nil,
			want:   map[string]bool{"withDiffstat": true, "withMergeable": true, "withAuthor": true, "withMergedBy": true, "withLabels": true},
		},
		{
			name:   "selected fields fetches only those groups",
			fields: []models.PullRequestField{models.PullRequestFieldLabels},
			want:   map[string]bool{"withDiffstat": false, "withMergeable": false, "withAuthor": false, "withMergedBy": false, "withLabels": true},
		},
		{
			name:   "empty fields fetches no optional group",
			fields: []models.PullRequestField{},
			want:   map[string]bool{"withDiffstat": false, "withMergeable": false, "withAuthor": false, "withMergedBy": false, "withLabels": false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := models.ListPullRequestsOptions{
				Repository: "grafana",
				Owner:      "grafana",
				TimeField:  models.PullRequestClosedAt,
				Fields:     tt.fields,
			}
			client := testutil.NewTestClient(t,
				func(t *testing.T, variables map[string]interface{}) {
					for name, want := range tt.want {
						if got := variables[name]; got != githubv4.Boolean(want) {
							t.Errorf("variable %s = %v, want %v", name, got, want)
						}
					}
				},
				testutil.GetTestQueryFunction(&QueryListPullRequests{}),
			)

			// GetPullRequestsInRange rebuilds the options, so this also checks that Fields is passed on
			if _, err := GetPullRequestsInRange(context.Background(), client, opts, time.Now().Add(-24*time.Hour), time.Now()); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPullRequestsFramer(t *testing.T) {
	prs := PullRequests{{Number: 1}}

	if _, ok := pullRequestsFramer(prs, models.ListPullRequestsOptions{}).(PullRequests); !ok {
		t.Error("fields unset should keep the original PullRequests frame")
	}

	fields := []models.PullRequestField{models.PullRequestFieldLabels}
	framer, ok := pullRequestsFramer(prs, models.ListPullRequestsOptions{Fields: fields}).(PullRequestsWithFields)
	if !ok {
		t.Fatal("selected fields should use PullRequestsWithFields")
	}
	if got := len(framer.Frames()[0].Fields); got != 15 {
		t.Errorf("labels only frame has %d columns, want 15", got)
	}
}

func TestHandlePullRequestsQueryInvalidFields(t *testing.T) {
	d := &Datasource{}
	query := &models.PullRequestsQuery{
		Options: models.ListPullRequestsOptions{Fields: []models.PullRequestField{"Labels"}},
	}

	// The fields are validated before the client is used, so the Datasource needs no client
	_, err := d.HandlePullRequestsQuery(context.Background(), query, backend.DataQuery{})
	if err == nil || !strings.Contains(err.Error(), `invalid pull request field "Labels"`) {
		t.Fatalf("got error %v, want an invalid field error", err)
	}
}
