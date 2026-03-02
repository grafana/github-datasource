package github

import (
	"encoding/json"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestTableToQueryTypeCoversAllTypes(t *testing.T) {
	allQueryTypes := []string{
		models.QueryTypeCommits, models.QueryTypeIssues,
		models.QueryTypePullRequests, models.QueryTypePullRequestReviews,
		models.QueryTypeRepositories, models.QueryTypeContributors,
		models.QueryTypeTags, models.QueryTypeReleases,
		models.QueryTypeLabels, models.QueryTypeMilestones,
		models.QueryTypePackages, models.QueryTypeVulnerabilities,
		models.QueryTypeProjects, models.QueryTypeProjectItems,
		models.QueryTypeStargazers, models.QueryTypeWorkflows,
		models.QueryTypeWorkflowUsage, models.QueryTypeWorkflowRuns,
		models.QueryTypeCodeScanning, models.QueryTypeOrganizations,
		models.QueryTypeGraphQL,
	}
	for _, qt := range allQueryTypes {
		tableName := normalizeTableNames(qt)
		if _, ok := tableToQueryType[tableName]; !ok {
			t.Errorf("tableToQueryType missing entry for %q (table name %q)", qt, tableName)
		}
	}
}

func TestNormalizeAllTableTypes(t *testing.T) {
	tests := []struct {
		name      string
		table     string
		wantType  string
		wantOwner string
		wantRepo  string
		unchanged bool
	}{
		{name: "issues", table: "issues_grafana_grafana", wantType: models.QueryTypeIssues, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "pull-requests", table: "pull-requests_grafana_grafana", wantType: models.QueryTypePullRequests, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "pull-request-reviews", table: "pull-request-reviews_grafana_grafana", wantType: models.QueryTypePullRequestReviews, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "commits", table: "commits_grafana_grafana", wantType: models.QueryTypeCommits, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "code-scanning", table: "code-scanning_grafana_grafana", wantType: models.QueryTypeCodeScanning, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "workflow-runs", table: "workflow-runs_grafana_grafana", wantType: models.QueryTypeWorkflowRuns, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "workflow-usage", table: "workflow-usage_grafana_grafana", wantType: models.QueryTypeWorkflowUsage, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "stargazers", table: "stargazers_grafana_grafana", wantType: models.QueryTypeStargazers, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "labels", table: "labels_grafana_grafana", wantType: models.QueryTypeLabels, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "milestones", table: "milestones_grafana_grafana", wantType: models.QueryTypeMilestones, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "packages", table: "packages_grafana_grafana", wantType: models.QueryTypePackages, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "vulnerabilities", table: "vulnerabilities_grafana_grafana", wantType: models.QueryTypeVulnerabilities, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "contributors", table: "contributors_grafana_grafana", wantType: models.QueryTypeContributors, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "tags", table: "tags_grafana_grafana", wantType: models.QueryTypeTags, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "releases", table: "releases_grafana_grafana", wantType: models.QueryTypeReleases, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "workflows", table: "workflows_grafana_grafana", wantType: models.QueryTypeWorkflows, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "repositories owner only", table: "repositories_grafana", wantType: models.QueryTypeRepositories, wantOwner: "grafana", wantRepo: ""},
		{name: "projects owner only", table: "projects_grafana", wantType: models.QueryTypeProjects, wantOwner: "grafana", wantRepo: ""},
		{name: "organizations no sub-tables", table: "organizations", wantType: models.QueryTypeOrganizations, wantOwner: "", wantRepo: ""},
		{name: "repo with underscores", table: "pull-requests_foo_my_repo", wantType: models.QueryTypePullRequests, wantOwner: "foo", wantRepo: "my_repo"},
		{name: "unknown table unchanged", table: "unknown_grafana_grafana", unchanged: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryJSON, _ := json.Marshal(map[string]interface{}{"refId": "A", "grafanaSql": true, "table": tt.table})
			req := &backend.QueryDataRequest{
				Queries: []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
			}
			out := normalizeGrafanaSQLRequest(req)
			if out == nil || len(out.Queries) != 1 {
				t.Fatalf("expected one query")
			}
			q := out.Queries[0]
			if tt.unchanged {
				if string(q.JSON) != string(queryJSON) {
					t.Errorf("expected query to be unchanged")
				}
				return
			}
			if q.QueryType != tt.wantType {
				t.Errorf("queryType: got %q, want %q", q.QueryType, tt.wantType)
			}
			var raw map[string]interface{}
			if err := json.Unmarshal(q.JSON, &raw); err != nil {
				t.Fatal(err)
			}
			if raw["owner"] != tt.wantOwner {
				t.Errorf("owner: got %v, want %q", raw["owner"], tt.wantOwner)
			}
			if raw["repository"] != tt.wantRepo {
				t.Errorf("repository: got %v, want %q", raw["repository"], tt.wantRepo)
			}
		})
	}
}

func TestNormalizeGrafanaSQLRequest(t *testing.T) {
	t.Run("rewrites grafanaSql pull-requests query with table owner_repo", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","datasource":{"type":"grafana-github-datasource","uid":"PE7186B25D2FDE54B"},"filters":null,"grafanaSql":true,"table":"pull-requests_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query, got %v", out)
		}
		q := out.Queries[0]
		if q.QueryType != models.QueryTypePullRequests {
			t.Errorf("queryType: got %q, want %q", q.QueryType, models.QueryTypePullRequests)
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(q.JSON, &raw); err != nil {
			t.Fatal(err)
		}
		if raw["queryType"] != models.QueryTypePullRequests {
			t.Errorf("JSON queryType: got %v", raw["queryType"])
		}
		if raw["owner"] != "grafana" || raw["repository"] != "grafana" {
			t.Errorf("owner/repository: got %v / %v", raw["owner"], raw["repository"])
		}
	})

	t.Run("rewrites grafanaSql issues query", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		q := out.Queries[0]
		if q.QueryType != models.QueryTypeIssues {
			t.Errorf("queryType: got %q, want %q", q.QueryType, models.QueryTypeIssues)
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(q.JSON, &raw); err != nil {
			t.Fatal(err)
		}
		if raw["owner"] != "grafana" || raw["repository"] != "grafana" {
			t.Errorf("owner/repository: got %v / %v", raw["owner"], raw["repository"])
		}
	})

	t.Run("rewrites grafanaSql commits query", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"commits_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if out.Queries[0].QueryType != models.QueryTypeCommits {
			t.Errorf("queryType: got %q, want %q", out.Queries[0].QueryType, models.QueryTypeCommits)
		}
	})

	t.Run("rewrites grafanaSql code-scanning query", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"code-scanning_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if out.Queries[0].QueryType != models.QueryTypeCodeScanning {
			t.Errorf("queryType: got %q, want %q", out.Queries[0].QueryType, models.QueryTypeCodeScanning)
		}
	})

	t.Run("rewrites grafanaSql organizations query (no sub-tables)", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"organizations"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		q := out.Queries[0]
		if q.QueryType != models.QueryTypeOrganizations {
			t.Errorf("queryType: got %q, want %q", q.QueryType, models.QueryTypeOrganizations)
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(q.JSON, &raw); err != nil {
			t.Fatal(err)
		}
		if raw["owner"] != "" {
			t.Errorf("expected empty owner, got %v", raw["owner"])
		}
		if raw["repository"] != "" {
			t.Errorf("expected empty repository, got %v", raw["repository"])
		}
	})

	t.Run("rewrites grafanaSql workflow-runs query", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"workflow-runs_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if out.Queries[0].QueryType != models.QueryTypeWorkflowRuns {
			t.Errorf("queryType: got %q, want %q", out.Queries[0].QueryType, models.QueryTypeWorkflowRuns)
		}
	})

	t.Run("leaves non-grafanaSql query unchanged", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","queryType":"Pull_Requests","owner":"grafana","repository":"grafana"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if string(out.Queries[0].JSON) != string(queryJSON) {
			t.Errorf("query should be unchanged")
		}
	})

	t.Run("leaves grafanaSql query with unknown table unchanged", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"unknown_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if string(out.Queries[0].JSON) != string(queryJSON) {
			t.Errorf("query should be unchanged")
		}
	})

	t.Run("parses repo name with underscores", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"B","grafanaSql":true,"table":"pull-requests_foo_my_repo"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "B", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		if raw["owner"] != "foo" || raw["repository"] != "my_repo" {
			t.Errorf("owner/repository: got %v / %v", raw["owner"], raw["repository"])
		}
	})

	t.Run("handles nil request", func(t *testing.T) {
		out := normalizeGrafanaSQLRequest(nil)
		if out != nil {
			t.Errorf("expected nil, got %v", out)
		}
	})

	t.Run("handles empty queries", func(t *testing.T) {
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{},
		}
		out := normalizeGrafanaSQLRequest(req)
		if len(out.Queries) != 0 {
			t.Errorf("expected empty queries")
		}
	})
}

func TestNormalizeGrafanaSQLRequestWithFilters(t *testing.T) {
	t.Run("pushes down state filter for issues", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"column":"state","op":"==","value":"open"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts, _ := raw["options"].(map[string]interface{})
		if opts == nil {
			t.Fatal("expected options to be set")
		}
		query, _ := opts["query"].(string)
		if query != "state:open" {
			t.Errorf("expected query 'state:open', got %q", query)
		}
	})

	t.Run("pushes down author filter for issues", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"column":"author","op":"==","value":"octocat"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		query := opts["query"].(string)
		if query != "author:octocat" {
			t.Errorf("expected query 'author:octocat', got %q", query)
		}
	})

	t.Run("pushes down multiple filters for issues", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"column":"state","op":"==","value":"open"},{"column":"labels","op":"==","value":"bug"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		query := opts["query"].(string)
		if query != "state:open label:bug" {
			t.Errorf("expected query 'state:open label:bug', got %q", query)
		}
	})

	t.Run("pushes down state filter for code-scanning", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"code-scanning_grafana_grafana","filters":[{"column":"state","op":"==","value":"open"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		state, _ := opts["state"].(string)
		if state != "open" {
			t.Errorf("expected state 'open', got %q", state)
		}
	})

	t.Run("pushes down branch filter for workflow-runs", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"workflow-runs_grafana_grafana","filters":[{"column":"head_branch","op":"==","value":"main"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		branch, _ := opts["branch"].(string)
		if branch != "main" {
			t.Errorf("expected branch 'main', got %q", branch)
		}
	})

	t.Run("pushes down status filter for workflow-runs", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"workflow-runs_grafana_grafana","filters":[{"column":"status","op":"==","value":"completed"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		status, _ := opts["status"].(string)
		if status != "completed" {
			t.Errorf("expected status 'completed', got %q", status)
		}
	})

	t.Run("pushes down draft filter for pull-requests", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"pull-requests_grafana_grafana","filters":[{"column":"is_draft","op":"==","value":"true"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		query := opts["query"].(string)
		if query != "draft:true" {
			t.Errorf("expected query 'draft:true', got %q", query)
		}
	})

	t.Run("pushes down name filter for labels (like)", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"labels_grafana_grafana","filters":[{"column":"name","op":"like","value":"bug"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		query, _ := opts["query"].(string)
		if query != "bug" {
			t.Errorf("expected query 'bug', got %q", query)
		}
	})

	t.Run("pushes down package type filter", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"packages_grafana_grafana","filters":[{"column":"type","op":"==","value":"DOCKER"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		pkgType, _ := opts["packageType"].(string)
		if pkgType != "DOCKER" {
			t.Errorf("expected packageType 'DOCKER', got %q", pkgType)
		}
	})

	t.Run("pushes down name filter for repositories", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"repositories_grafana","filters":[{"column":"name","op":"like","value":"grafana"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		if raw["repository"] != "grafana" {
			t.Errorf("expected repository 'grafana', got %v", raw["repository"])
		}
	})

	t.Run("pushes down tool_name filter for code-scanning", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"code-scanning_grafana_grafana","filters":[{"column":"tool_name","op":"==","value":"CodeQL"}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		toolName, _ := opts["toolName"].(string)
		if toolName != "CodeQL" {
			t.Errorf("expected toolName 'CodeQL', got %q", toolName)
		}
	})

	t.Run("ignores filters with empty values", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"column":"state","op":"==","value":""}]}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		if _, exists := opts["query"]; exists {
			t.Errorf("expected no query set for empty filter value")
		}
	})

	t.Run("handles null filters", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":null}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if out.Queries[0].QueryType != models.QueryTypeIssues {
			t.Errorf("queryType: got %q, want %q", out.Queries[0].QueryType, models.QueryTypeIssues)
		}
	})
}
