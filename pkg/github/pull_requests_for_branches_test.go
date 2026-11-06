package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

// graphqlClient sends Query through a real githubv4 client so the runtime-built query struct is encoded and decoded for real.
type graphqlClient struct {
	*testutil.TestClient
	gql *githubv4.Client
}

func (c *graphqlClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.gql.Query(ctx, q, variables)
}

type graphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func newGraphQLServer(t *testing.T, response string, got *graphqlRequest) models.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(got); err != nil {
			t.Errorf("decoding request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	return &graphqlClient{
		TestClient: testutil.NewTestClient(t, nil, nil),
		gql:        githubv4.NewEnterpriseClient(server.URL, server.Client()),
	}
}

const pullRequestsForBranchesResponse = `{"data":{"repository":{
  "nameWithOwner":"grafana/sigil",
  "defaultBranchRef":{"name":"main"},
  "b0":{"nodes":[
    {"number":12,"state":"MERGED","headRefName":"feat/a","baseRefName":"main","createdAt":"2026-10-01T10:00:00Z","mergedAt":"2026-10-02T10:00:00Z","closedAt":"2026-10-02T10:00:00Z","headRepository":{"nameWithOwner":"grafana/sigil"}},
    {"number":7,"state":"CLOSED","headRefName":"feat/a","baseRefName":"main","createdAt":"2026-09-01T10:00:00Z","mergedAt":null,"closedAt":"2026-09-03T10:00:00Z","headRepository":null}
  ]},
  "b1":{"nodes":[]},
  "b2":{"nodes":[
    {"number":30,"state":"OPEN","headRefName":"fix/\"quoted\" $x","baseRefName":"release-1.0","createdAt":"2026-10-05T10:00:00Z","mergedAt":null,"closedAt":null,"headRepository":{"nameWithOwner":"someone/sigil"}}
  ]}
}}}`

func TestGetPullRequestsForBranches(t *testing.T) {
	var got graphqlRequest
	client := newGraphQLServer(t, pullRequestsForBranchesResponse, &got)

	branches := []string{"feat/a", "no-pr", `fix/"quoted" $x`, "feat/a", ""}
	prs, err := GetPullRequestsForBranches(context.Background(), client, models.ListPullRequestsForBranchesOptions{
		Owner:      "grafana",
		Repository: "sigil",
		Branches:   branches,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("branch names are variables, never query text", func(t *testing.T) {
		for _, b := range []string{"feat/a", "no-pr", "quoted"} {
			if strings.Contains(got.Query, b) {
				t.Errorf("query text contains branch %q: %s", b, got.Query)
			}
		}
		for i, b := range []string{"feat/a", "no-pr", `fix/"quoted" $x`} {
			if v := got.Variables[fmt.Sprintf("h%d", i)]; v != b {
				t.Errorf("variable h%d = %v, want %q", i, v, b)
			}
		}
		if _, ok := got.Variables["h3"]; ok {
			t.Error("duplicate and empty branches must be dropped before building the query")
		}
		for _, want := range []string{"$h0:String!", "b0: pullRequests(headRefName: $h0, first: 5", "repository(owner: $owner, name: $name)", "defaultBranchRef{name}"} {
			if !strings.Contains(got.Query, want) {
				t.Errorf("query missing %q: %s", want, got.Query)
			}
		}
	})

	t.Run("rows", func(t *testing.T) {
		if len(prs) != 3 {
			t.Fatalf("got %d pull requests, want 3", len(prs))
		}
		if prs[0].Repository != "grafana/sigil" || prs[0].DefaultBranch == nil || *prs[0].DefaultBranch != "main" {
			t.Errorf("repository/default branch = %q/%v", prs[0].Repository, prs[0].DefaultBranch)
		}
		if prs[0].PullRequest.MergedAt == nil || prs[1].PullRequest.MergedAt != nil {
			t.Error("merged_at must be set only on the merged pull request")
		}
		if prs[1].PullRequest.HeadRepository != nil {
			t.Error("a deleted head repository must decode as nil")
		}
		if prs[2].PullRequest.HeadRepository.NameWithOwner != "someone/sigil" {
			t.Error("fork head repository not decoded")
		}
	})

	testutil.CheckGoldenFramer(t, "pull_requests_for_branches", prs)
}

func TestGetPullRequestsForBranchesNullDefaultBranch(t *testing.T) {
	var got graphqlRequest
	client := newGraphQLServer(t, `{"data":{"repository":{"nameWithOwner":"grafana/empty","defaultBranchRef":null,"b0":{"nodes":[]}}}}`, &got)

	prs, err := GetPullRequestsForBranches(context.Background(), client, models.ListPullRequestsForBranchesOptions{
		Owner: "grafana", Repository: "empty", Branches: []string{"x"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(prs) != 0 {
		t.Fatalf("got %d pull requests, want 0", len(prs))
	}
}

func TestGetPullRequestsForBranchesLimits(t *testing.T) {
	client := testutil.NewTestClient(t, nil, func(t *testing.T, _ interface{}) {
		t.Error("no GitHub request may be sent")
	})

	t.Run("missing input returns nothing", func(t *testing.T) {
		for _, opts := range []models.ListPullRequestsForBranchesOptions{
			{Repository: "sigil", Branches: []string{"a"}},
			{Owner: "grafana", Branches: []string{"a"}},
			{Owner: "grafana", Repository: "sigil", Branches: []string{""}},
		} {
			prs, err := GetPullRequestsForBranches(context.Background(), client, opts)
			if err != nil || len(prs) != 0 {
				t.Errorf("%+v: got %d rows, err %v", opts, len(prs), err)
			}
		}
	})

	t.Run("too many branches", func(t *testing.T) {
		branches := make([]string, models.MaxPullRequestsForBranches+1)
		for i := range branches {
			branches[i] = fmt.Sprintf("b%d", i)
		}
		_, err := GetPullRequestsForBranches(context.Background(), client, models.ListPullRequestsForBranchesOptions{
			Owner: "grafana", Repository: "sigil", Branches: branches,
		})
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}
