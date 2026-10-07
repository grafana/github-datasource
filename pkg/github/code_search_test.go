package github

import (
	"context"
	"testing"

	googlegithub "github.com/google/go-github/v84/github"
)

type codeSearchMockClient struct {
	queries    []string
	incomplete bool
}

func (client *codeSearchMockClient) SearchCode(_ context.Context, query string, options *googlegithub.SearchOptions) (*googlegithub.CodeSearchResult, *googlegithub.Response, error) {
	client.queries = append(client.queries, query)
	if options.Page == 1 {
		return &googlegithub.CodeSearchResult{
			Total:             googlegithub.Ptr(2),
			IncompleteResults: googlegithub.Ptr(client.incomplete),
			CodeResults: []*googlegithub.CodeResult{
				{Path: googlegithub.Ptr(".github/settings.yaml"), Repository: &googlegithub.Repository{FullName: googlegithub.Ptr("rwest/managed")}},
			},
		}, &googlegithub.Response{NextPage: 2}, nil
	}
	return &googlegithub.CodeSearchResult{
		Total:             googlegithub.Ptr(2),
		IncompleteResults: googlegithub.Ptr(client.incomplete),
		CodeResults: []*googlegithub.CodeResult{
			{Path: googlegithub.Ptr(".github/examples/settings.yaml"), Repository: &googlegithub.Repository{FullName: googlegithub.Ptr("rwest/example")}},
		},
	}, &googlegithub.Response{}, nil
}

func TestGetCodeSearchRequiresCompleteResults(t *testing.T) {
	_, err := getCodeSearch(context.Background(), &codeSearchMockClient{incomplete: true}, "org:rwest settings", "", false, true)
	if err == nil {
		t.Fatal("expected incomplete code search to fail")
	}
}

func TestGetCodeSearchPaginatesAndFiltersExactPath(t *testing.T) {
	client := &codeSearchMockClient{}
	search, err := getCodeSearch(context.Background(), client, "org:rwest filename:settings.yaml path:.github", ".github/settings.yaml", true, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(client.queries) != 2 || len(search.Items) != 1 {
		t.Fatalf("expected two pages and one exact-path result, got pages=%d items=%d", len(client.queries), len(search.Items))
	}
	if !search.Complete || search.Total != 2 {
		t.Fatalf("expected complete two-result search, got complete=%v total=%d", search.Complete, search.Total)
	}
	frame := search.Frames()[0]
	assertFrameStringPointer(t, frame, "full_name", 0, "rwest/managed")
}
