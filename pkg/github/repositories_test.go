package github

import (
	"context"
	"testing"

	"github.com/grafana/grafana-github-datasource/pkg/testutil"
)

func TestGetAllRepositories(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = ListRepositoriesOptions{
			Organization: "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListRepositories{}),
	)

	_, err := GetAllRepositories(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}
