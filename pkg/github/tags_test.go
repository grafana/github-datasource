package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestGetAllTags(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListTagsOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListTags{}),
	)

	_, err := GetAllTags(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListTags(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListTagsOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
		from = time.Now().Add(-30 * 24 * time.Hour)
		to   = time.Now()
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListTags{}),
	)

	_, err := GetTagsInRange(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTagsDataFrames(t *testing.T) {
	committedAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	createdAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	user := GitActor{
		Name:  "firstCommitter",
		Email: "first@example.com",
		User: User{
			ID:      "1",
			Login:   "firstCommitter",
			Name:    "First Committer",
			Company: "ACME Corp",
			Email:   "first@example.com",
		},
	}

	commit1 := Commit{
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
		Author:  user,
	}

	commit2 := Commit{
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
		Author:  user,
	}

	tags := Tags{
		Tag{
			Name: "v1.0.0",
			Tagger: struct {
				Date githubv4.DateTime
				User User
			}{
				Date: githubv4.DateTime{
					Time: createdAt,
				},
				User: user.User,
			},
			Target: struct {
				OID    string
				Commit Commit "graphql:\"... on Commit\""
			}{
				OID:    "",
				Commit: commit1,
			},
		},
		Tag{
			Name: "v1.1.0",
			Tagger: struct {
				Date githubv4.DateTime
				User User
			}{
				Date: githubv4.DateTime{
					Time: createdAt,
				},
				User: user.User,
			},
			Target: struct {
				OID    string
				Commit Commit "graphql:\"... on Commit\""
			}{
				OID:    "",
				Commit: commit2,
			},
		},
	}

	if err := testutil.CheckGoldenFramer("tags", tags); err != nil {
		t.Fatal(err)
	}
}
