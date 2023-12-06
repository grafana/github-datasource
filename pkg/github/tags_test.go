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
	createdAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	user := GitActor{
		Name:  "firstCommitter",
		Email: "first@example.com",
		User: models.User{
			ID:      "1",
			Login:   "firstCommitter",
			Name:    "First Committer",
			Company: "ACME Corp",
			Email:   "first@example.com",
		},
	}

	tags := Tags{
		Tag{
			Name: "v1.0.0",
			OID:  "",
			Tagger: GitActor{
				Email: user.Email,
				Date: githubv4.GitTimestamp{
					Time: createdAt,
				},
				User: user.User,
			},
		},
		Tag{
			Name: "v1.1.0",
			Tagger: GitActor{
				Email: user.Email,
				Date: githubv4.GitTimestamp{
					Time: createdAt,
				},
				User: user.User,
			},
			OID: "",
		},
	}

	testutil.CheckGoldenFramer(t, "tags", tags)
}
