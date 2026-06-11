package github

import (
	"context"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
)

func TestListBranches(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListBranchesOptions{
			Repository: "grafana",
			Owner:      "grafana",
			Query:      "release/",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "name", "owner", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListBranches{}),
	)

	_, err := GetAllBranches(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}
