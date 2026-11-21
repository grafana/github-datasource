package githubclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gh "github.com/google/go-github/v84/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSecurityAdvisories(t *testing.T) {
	for _, tc := range []struct {
		name, repo, path, state string
	}{
		{"repository", "plugin", "/api/v3/repos/grafana/plugin/security-advisories", "draft"},
		{"organization", "", "/api/v3/orgs/grafana/security-advisories", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, tc.path, r.URL.Path)
				assert.Equal(t, tc.state, r.URL.Query().Get("state"))
				assert.Equal(t, tc.state != "", r.URL.Query().Has("state"))
				assert.Equal(t, "next", r.URL.Query().Get("after"))
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[{"ghsa_id":"GHSA-example"}]`))
			}))
			defer server.Close()

			restClient, err := gh.NewClient(nil).WithEnterpriseURLs(server.URL, server.URL)
			require.NoError(t, err)
			client := &Client{restClient: restClient}
			advisories, _, err := client.ListSecurityAdvisories(t.Context(), "grafana", tc.repo, &gh.ListRepositorySecurityAdvisoriesOptions{
				ListCursorOptions: gh.ListCursorOptions{After: "next"},
				State:             tc.state,
			})
			require.NoError(t, err)
			require.Len(t, advisories, 1)
			require.Equal(t, "GHSA-example", advisories[0].GetGHSAID())
		})
	}
}
