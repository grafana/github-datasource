package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestGetAllRepositories(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListRepositoriesOptions{
			Owner: "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListRepositories{}),
	)

	_, err := GetAllRepositories(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRepositoriesDataFrame(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2023-08-25T16:21:56+00:00")

	repositories := Repositories{
		Repository{
			Name: "grafana",
			Owner: struct{ Login string }{
				Login: "grafana",
			},
			NameWithOwner:  "grafana/grafana",
			URL:            "github.com/grafana/grafana",
			ForkCount:      10,
			StargazerCount: 123,
			IsFork:         true,
			IsMirror:       true,
			IsPrivate:      false,
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			UpdatedAt: githubv4.DateTime{
				Time: updatedAt,
			},
		},
		Repository{
			Name: "loki",
			Owner: struct{ Login string }{
				Login: "grafana",
			},
			NameWithOwner: "grafana/loki",
			URL:           "github.com/grafana/loki",
			ForkCount:     12,
			IsFork:        true,
			IsMirror:      true,
			IsPrivate:     false,
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
		},
	}

	testutil.CheckGoldenFramer(t, "repositories", repositories)
}
