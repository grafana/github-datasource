package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-github-datasource/pkg/testutil"
)

func TestListPullRequests(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListPullRequestsInRangeOptions{
			Repository: "grafana",
			Owner:      "grafana",
			TimeField:  models.PullRequestClosedAt,
		}
	)

	testVariables := testutil.GetTestVariablesFunction("query", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListPullRequests{}),
	)

	_, err := GetPullRequestsInRange(ctx, client, opts, time.Now().Add(-30*24*time.Hour), time.Now())
	if err != nil {
		t.Fatal(err)
	}
}
