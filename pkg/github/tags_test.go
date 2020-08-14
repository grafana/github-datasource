package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/testutil"
)

func TestGetAllTags(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = ListTagsOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListTags{}),
	)

	_, err := GetAllTags(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListTags(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = ListTagsOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
		from = time.Now().Add(-30 * 24 * time.Hour)
		to   = time.Now()
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListTags{}),
	)

	_, err := GetTagsInRange(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}
}
