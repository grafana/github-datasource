package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestGetAllCommits(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListCommitsOptions{
			Repository: "test",
			Ref:        "master",
			Owner:      "kminehart-test",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "ref")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListCommits{}),
	)

	_, err := GetAllCommits(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListCommits(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListCommitsOptions{
			Repository: "grafana",
			Ref:        "master",
			Owner:      "grafana",
		}
		from = time.Now().Add(-7 * 24 * time.Hour)
		to   = time.Now()
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "ref", "cursor", "since", "until")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListCommitsInRange{}),
	)

	_, err := GetCommitsInRange(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommitsDataframe(t *testing.T) {
	committedAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	commits := Commits{
		Commit{
			OID: "",
			PushedDate: githubv4.DateTime{
				Time: committedAt.Add(time.Minute * 2),
			},
			AuthoredDate: githubv4.DateTime{
				Time: committedAt,
			},
			CommittedDate: githubv4.DateTime{
				Time: committedAt,
			},
			Message: "commit #1",
			Author: GitActor{
				Name:  "firstCommitter",
				Email: "first@example.com",
				User: User{
					ID:      "1",
					Login:   "firstCommitter",
					Name:    "First Committer",
					Company: "ACME Corp",
					Email:   "first@example.com",
				},
			},
		},
		Commit{
			OID: "",
			PushedDate: githubv4.DateTime{
				Time: committedAt.Add(time.Hour * 2),
			},
			AuthoredDate: githubv4.DateTime{
				Time: committedAt.Add(time.Hour),
			},
			CommittedDate: githubv4.DateTime{
				Time: committedAt.Add(time.Hour),
			},
			Message: "commit #2",
			Author: GitActor{
				Name:  "secondCommitter",
				Email: "second@example.com",
				User: User{
					ID:      "1",
					Login:   "secondCommitter",
					Name:    "Second Committer",
					Company: "ACME Corp",
					Email:   "second@example.com",
				},
			},
		},
	}

	if err := testutil.CheckGoldenFramer("commits", commits); err != nil {
		t.Fatal(err)
	}
}
