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
	createdAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	repositories := Repositories{
		Repository{
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

	if err := testutil.CheckGoldenFramer("repositories", repositories); err != nil {
		t.Fatal(err)
	}
}
