package github

import (
	"context"
	"testing"

	"github.com/grafana/github-datasource/pkg/testutil"
)

func TestGetAllOrganizations(t *testing.T) {
	var (
		ctx = context.Background()
	)

	testVariables := func(t *testing.T, variables map[string]interface{}) {
	}

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListOrganizations{}),
	)

	_, err := GetAllOrganizations(ctx, client)
	if err != nil {
		t.Fatal(err)
	}
}
