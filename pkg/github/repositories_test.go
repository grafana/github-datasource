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

	user := User{
		ID:      "1",
		Login:   "exampleUser",
		Name:    "Example User",
		Company: "ACME Corp",
		Email:   "user@example.com",
	}

	releases := Repositories{
		Repository{
			ID:           "1",
			Name:         "Repository #1",
			Author:       user,
			IsDraft:      true,
			IsPrerelease: false,
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			PublishedAt: githubv4.DateTime{},
			TagName:     "v1.0.0",
			URL:         "https://example.com/v1.0.0",
		},
		Repository{
			ID:           "1",
			Name:         "Repository #2",
			Author:       user,
			IsDraft:      true,
			IsPrerelease: false,
			CreatedAt: githubv4.DateTime{
				Time: createdAt,
			},
			PublishedAt: githubv4.DateTime{
				Time: createdAt.Add(time.Hour),
			},
			TagName: "v1.1.0",
			URL:     "https://example.com/v1.1.0",
		},
	}

	if err := testutil.CheckGoldenFramer("releases", releases); err != nil {
		t.Fatal(err)
	}
}
