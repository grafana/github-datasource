package github

import (
	"context"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
)

func TestListMilestones(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListMilestonesOptions{
			Repository: "grafana",
			Owner:      "grafana",
			Query:      "test",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "name", "owner", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListMilestones{}),
	)

	_, err := GetAllMilestones(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}
