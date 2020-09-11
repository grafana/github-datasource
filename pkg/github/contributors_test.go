package github

import (
	"context"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
)

func TestGetAllContributors(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListContributorsOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListContributors{}),
	)

	_, err := GetAllContributors(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetContributorsWithQuery(t *testing.T) {
	var (
		ctx  = context.Background()
		q    = "test query"
		opts = models.ListContributorsOptions{
			Repository: "grafana",
			Owner:      "grafana",
			Query:      &q,
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "query")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListContributors{}),
	)

	_, err := GetAllContributors(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestContributorsDataFrame(t *testing.T) {
	contributors := GitActors{
		GitActor{
			Name:  "Example User",
			Email: "user1@example.com",
			User: User{
				Login:   "exuser1",
				Name:    "Example User",
				Company: "ACME Corp",
				Email:   "user1@example.com",
				URL:     "https://github.com/user1",
			},
		},
		GitActor{
			Name:  "Example User2",
			Email: "user2@example.com",
			User: User{
				Login:   "exuser2",
				Name:    "Example User2",
				Company: "ACME Corp",
				Email:   "user2@example.com",
				URL:     "https://github.com/user2",
			},
		},
	}

	if err := testutil.CheckGoldenFramer("contributors", contributors); err != nil {
		t.Fatal(err)
	}
}
