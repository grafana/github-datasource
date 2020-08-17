package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/testutil"
)

func TestGetAllContributors(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = ListContributorsOptions{
			Repository: "grafana",
			Ref:        "master",
			Owner:      "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "ref")

	// Listing contributors just list commits :)
	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListCommits{}),
	)

	_, err := GetAllContributors(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetContributorsInRange(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = ListContributorsOptions{
			Repository: "grafana",
			Ref:        "master",
			Owner:      "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "ref", "cursor", "until", "since")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListCommitsInRange{}),
	)

	_, err := GetContributorsInRange(ctx, client, opts, time.Now().Add(-7*24*time.Hour), time.Now())
	if err != nil {
		t.Fatal(err)
	}
}
