package githubclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWorkflowRunsEndpoints(t *testing.T) {
	const base = "/api/v3/repos/acme/widget/actions/"
	for _, tc := range []struct {
		name, workflow, path string
	}{
		{"all workflows", "", base + "runs"},
		{"blank workflow", "  ", base + "runs"},
		{"by id", "42", base + "workflows/42/runs"},
		{"by filename", "ci.yml", base + "workflows/ci.yml/runs"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.path {
					t.Errorf("request path = %q, want %q", r.URL.Path, tc.path)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				params := r.URL.Query()
				assert.Equal(t, "main", params.Get("branch"))
				assert.Equal(t, "100", params.Get("per_page"))
				assert.Equal(t, "2026-11-16T00:00:00Z..2026-11-17T00:00:00Z", params.Get("created"))
				w.Header().Set("Content-Type", "application/json")
				switch params.Get("page") {
				case "1":
					w.Header().Set("Link", fmt.Sprintf("<http://%s%s?page=2>; rel=\"next\"", r.Host, r.URL.Path))
					fmt.Fprint(w, `{"total_count":2,"workflow_runs":[{"id":11,"workflow_id":42}]}`)
				case "2":
					fmt.Fprint(w, `{"total_count":2,"workflow_runs":[{"id":12,"workflow_id":43}]}`)
				default:
					t.Errorf("unexpected page %q", params.Get("page"))
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()
			rest, err := googlegithub.NewClient(server.Client()).WithEnterpriseURLs(server.URL, server.URL)
			require.NoError(t, err)
			client := &Client{restClient: rest}
			from := time.Date(2026, 11, 16, 0, 0, 0, 0, time.UTC)
			runs, err := client.GetWorkflowRuns(t.Context(), "acme", "widget", tc.workflow, "main", backend.TimeRange{From: from, To: from.Add(24 * time.Hour)})
			require.NoError(t, err)
			require.Len(t, runs, 2)
			assert.Equal(t, int64(11), runs[0].GetID())
			assert.Equal(t, int64(43), runs[1].GetWorkflowID())
		})
	}
}
