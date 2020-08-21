package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-github-datasource/pkg/testutil"
)

func TestGetAllCommits(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListCommitsOptions{
			Repository: "test",
			Ref:        "master",
			Owner:      "kminehart-test",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "ref")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListCommits{}),
	)

	_, err := GetAllCommits(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListCommits(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListCommitsOptions{
			Repository: "grafana",
			Ref:        "master",
			Owner:      "grafana",
		}
		from = time.Now().Add(-7 * 24 * time.Hour)
		to   = time.Now()
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "ref", "cursor", "since", "until")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListCommitsInRange{}),
	)

	_, err := GetCommitsInRange(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}
}
