package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/testutil"
)

func TestGetAllReleases(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = ListReleasesOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListReleases{}),
	)

	_, err := GetAllReleases(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListReleases(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = ListReleasesOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListReleases{}),
	)

	_, err := GetReleasesInRange(ctx, client, opts, time.Now().Add(-30*24*time.Hour), time.Now())
	if err != nil {
		t.Fatal(err)
	}
}
