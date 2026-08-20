package githubclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/stretchr/testify/require"
)

func TestListCheckRunsForRef(t *testing.T) {
	var gotPath, gotQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"total_count": 1,
			"check_runs": [
				{
					"id": 4,
					"name": "build",
					"head_sha": "abc123",
					"status": "completed",
					"conclusion": "success",
					"started_at": "2024-01-01T10:00:00Z",
					"completed_at": "2024-01-01T10:05:00Z",
					"app": {"name": "GitHub Actions", "slug": "github-actions"},
					"check_suite": {"id": 42},
					"output": {"title": "Build succeeded", "annotations_count": 0}
				}
			]
		}`))
	}))
	defer server.Close()

	restClient, err := googlegithub.NewClient(nil).WithEnterpriseURLs(server.URL, server.URL)
	require.NoError(t, err)

	client := &Client{restClient: restClient}

	results, resp, err := client.ListCheckRunsForRef(context.Background(), "grafana", "github-datasource", "abc123", &googlegithub.ListCheckRunsOptions{
		CheckName:   googlegithub.Ptr("build"),
		Status:      googlegithub.Ptr("completed"),
		Filter:      googlegithub.Ptr("all"),
		ListOptions: googlegithub.ListOptions{PerPage: 100},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, results)

	require.Equal(t, "/api/v3/repos/grafana/github-datasource/commits/abc123/check-runs", gotPath)
	require.Contains(t, gotQuery, "check_name=build")
	require.Contains(t, gotQuery, "status=completed")
	require.Contains(t, gotQuery, "filter=all")
	require.Contains(t, gotQuery, "per_page=100")

	require.Len(t, results.CheckRuns, 1)
	require.Equal(t, int64(4), results.CheckRuns[0].GetID())
	require.Equal(t, "build", results.CheckRuns[0].GetName())
	require.Equal(t, "success", results.CheckRuns[0].GetConclusion())
	require.Equal(t, "GitHub Actions", results.CheckRuns[0].GetApp().GetName())
	require.Equal(t, int64(42), results.CheckRuns[0].GetCheckSuite().GetID())
}

func TestListCheckRunsForRefError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "Not Found"}`))
	}))
	defer server.Close()

	restClient, err := googlegithub.NewClient(nil).WithEnterpriseURLs(server.URL, server.URL)
	require.NoError(t, err)

	client := &Client{restClient: restClient}

	results, resp, err := client.ListCheckRunsForRef(context.Background(), "grafana", "github-datasource", "missing", nil)
	require.Error(t, err)
	require.Nil(t, results)
	require.Nil(t, resp)
}
