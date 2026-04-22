package github

import (
	"encoding/json"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/config"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/featuretoggles"
)

func pluginCtxWithFeatureToggle() backend.PluginContext {
	return backend.PluginContext{
		GrafanaConfig: config.NewGrafanaCfg(map[string]string{
			featuretoggles.EnabledFeatures: "dsAbstractionApp",
		}),
	}
}

func TestTableToQueryTypeCoversAllTypes(t *testing.T) {
	allQueryTypes := []models.QueryType{
		models.QueryTypeCommits, models.QueryTypeIssues,
		models.QueryTypePullRequests, models.QueryTypePullRequestReviews,
		models.QueryTypeRepositories, models.QueryTypeContributors,
		models.QueryTypeTags, models.QueryTypeReleases,
		models.QueryTypeLabels, models.QueryTypeMilestones,
		models.QueryTypePackages, models.QueryTypeVulnerabilities,
		models.QueryTypeProjects,
		models.QueryTypeStargazers, models.QueryTypeWorkflows,
		models.QueryTypeWorkflowUsage, models.QueryTypeWorkflowRuns,
		models.QueryTypeCodeScanning, models.QueryTypeDeployments,
		models.QueryTypeOrganizations, models.QueryTypeGraphQL,
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
		wantType  models.QueryType
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
		{name: "deployments", table: "deployments_grafana_grafana", wantType: models.QueryTypeDeployments, wantOwner: "grafana", wantRepo: "grafana"},
		{name: "repositories owner only", table: "repositories_grafana", wantType: models.QueryTypeRepositories, wantOwner: "grafana", wantRepo: ""},
		{name: "projects owner only", table: "projects_grafana", wantType: models.QueryTypeProjects, wantOwner: "grafana", wantRepo: ""},
		{name: "organizations no table parameters", table: "organizations", wantType: models.QueryTypeOrganizations, wantOwner: "", wantRepo: ""},
		{name: "repo with underscores", table: "pull-requests_foo_my_repo", wantType: models.QueryTypePullRequests, wantOwner: "foo", wantRepo: "my_repo"},
		{name: "unknown table unchanged", table: "unknown_grafana_grafana", unchanged: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryJSON, _ := json.Marshal(map[string]interface{}{"refId": "A", "grafanaSql": true, "table": tt.table})
			req := &backend.QueryDataRequest{
				PluginContext: pluginCtxWithFeatureToggle(),
				Queries:       []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
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
			if q.QueryType != string(tt.wantType) {
				t.Errorf("queryType: got %q, want %q", q.QueryType, string(tt.wantType))
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
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query, got %v", out)
		}
		q := out.Queries[0]
		if q.QueryType != string(models.QueryTypePullRequests) {
			t.Errorf("queryType: got %q, want %q", q.QueryType, models.QueryTypePullRequests)
		}
		var raw map[string]any
		if err := json.Unmarshal(q.JSON, &raw); err != nil {
			t.Fatal(err)
		}
		if raw["queryType"] != string(models.QueryTypePullRequests) {
			t.Errorf("JSON queryType: got %v", raw["queryType"])
		}
		if raw["owner"] != "grafana" || raw["repository"] != "grafana" {
			t.Errorf("owner/repository: got %v / %v", raw["owner"], raw["repository"])
		}
	})

	t.Run("rewrites grafanaSql issues query", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		q := out.Queries[0]
		if q.QueryType != string(models.QueryTypeIssues) {
			t.Errorf("queryType: got %q, want %q", q.QueryType, string(models.QueryTypeIssues))
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
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if out.Queries[0].QueryType != string(models.QueryTypeCommits) {
			t.Errorf("queryType: got %q, want %q", out.Queries[0].QueryType, string(models.QueryTypeCommits))
		}
	})

	t.Run("rewrites grafanaSql code-scanning query", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"code-scanning_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if out.Queries[0].QueryType != string(models.QueryTypeCodeScanning) {
			t.Errorf("queryType: got %q, want %q", out.Queries[0].QueryType, string(models.QueryTypeCodeScanning))
		}
	})

	t.Run("rewrites grafanaSql organizations query (no table parameters)", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"organizations"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		q := out.Queries[0]
		if q.QueryType != string(models.QueryTypeOrganizations) {
			t.Errorf("queryType: got %q, want %q", q.QueryType, string(models.QueryTypeOrganizations))
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
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if out.Queries[0].QueryType != string(models.QueryTypeWorkflowRuns) {
			t.Errorf("queryType: got %q, want %q", out.Queries[0].QueryType, string(models.QueryTypeWorkflowRuns))
		}
	})

	t.Run("maps projects owner to options.organization", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"projects_grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		org, _ := opts["organization"].(string)
		if org != "grafana" {
			t.Errorf("expected options.organization 'grafana', got %q", org)
		}
	})

	t.Run("uses tableParameterValues for owner and repository", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues","tableParameterValues":{"organization":"grafana","repository":"grafana"}}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		if raw["owner"] != "grafana" || raw["repository"] != "grafana" {
			t.Errorf("owner/repository: got %v / %v", raw["owner"], raw["repository"])
		}
	})

	t.Run("tableParameterValues override owner and repository from table suffix", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_oldorg_oldrepo","tableParameterValues":{"organization":"grafana","repository":"grafana"}}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		if raw["owner"] != "grafana" || raw["repository"] != "grafana" {
			t.Errorf("owner/repository: got %v / %v", raw["owner"], raw["repository"])
		}
	})

	t.Run("maps workflow table parameter to options.workflow", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"workflow-usage","tableParameterValues":{"organization":"grafana","repository":"grafana","workflow":"build.yml"}}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		if raw["owner"] != "grafana" || raw["repository"] != "grafana" {
			t.Errorf("owner/repository: got %v / %v", raw["owner"], raw["repository"])
		}
		opts, _ := raw["options"].(map[string]interface{})
		workflow, _ := opts["workflow"].(string)
		if workflow != "build.yml" {
			t.Errorf("expected options.workflow 'build.yml', got %q", workflow)
		}
	})

	t.Run("maps projects organization from tableParameterValues", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"projects","tableParameterValues":{"organization":"grafana"}}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		org, _ := opts["organization"].(string)
		if org != "grafana" {
			t.Errorf("expected options.organization 'grafana', got %q", org)
		}
	})

	t.Run("leaves non-grafanaSql query unchanged", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","queryType":"Pull_Requests","owner":"grafana","repository":"grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
			PluginContext: pluginCtxWithFeatureToggle(),
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
			PluginContext: pluginCtxWithFeatureToggle(),
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
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries:       []backend.DataQuery{},
		}
		out := normalizeGrafanaSQLRequest(req)
		if len(out.Queries) != 0 {
			t.Errorf("expected empty queries")
		}
	})
}

func TestNormalizeGrafanaSQLRequestWithFilters(t *testing.T) {
	t.Run("pushes down state filter for issues", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"name":"state","conditions":[{"operator":"=","value":"open"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"name":"author","conditions":[{"operator":"=","value":"octocat"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"name":"state","conditions":[{"operator":"=","value":"open"}]},{"name":"labels","conditions":[{"operator":"=","value":"bug"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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

	t.Run("pushes down JSON array value for issues labels", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"name":"labels","conditions":[{"operator":"in","values":["bug","triage"]}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		if query != "label:bug label:triage" {
			t.Errorf("expected query 'label:bug label:triage', got %q", query)
		}
	})

	t.Run("pushes down IN values for assignees", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"name":"assignees","conditions":[{"operator":"in","values":["alice","bob"]}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		if query != "assignee:alice assignee:bob" {
			t.Errorf("expected query 'assignee:alice assignee:bob', got %q", query)
		}
	})

	t.Run("pushes down state filter for code-scanning", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"code-scanning_grafana_grafana","filters":[{"name":"state","conditions":[{"operator":"=","value":"open"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"workflow-runs_grafana_grafana","filters":[{"name":"head_branch","conditions":[{"operator":"=","value":"main"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"workflow-runs_grafana_grafana","filters":[{"name":"status","conditions":[{"operator":"=","value":"completed"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"pull-requests_grafana_grafana","filters":[{"name":"is_draft","conditions":[{"operator":"=","value":"true"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"labels_grafana_grafana","filters":[{"name":"name","conditions":[{"operator":"like","value":"bug"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"packages_grafana_grafana","filters":[{"name":"type","conditions":[{"operator":"=","value":"DOCKER"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"repositories_grafana","filters":[{"name":"name","conditions":[{"operator":"like","value":"grafana"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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

	t.Run("pushes down is_fork filter for repositories", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"repositories_grafana","filters":[{"name":"is_fork","conditions":[{"operator":"=","value":"true"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		repo, _ := raw["repository"].(string)
		if repo != "fork:only" {
			t.Errorf("expected repository 'fork:only', got %q", repo)
		}
	})

	t.Run("pushes down is_private filter for repositories", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"repositories_grafana","filters":[{"name":"is_private","conditions":[{"operator":"=","value":"true"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		repo, _ := raw["repository"].(string)
		if repo != "is:private" {
			t.Errorf("expected repository 'is:private', got %q", repo)
		}
	})

	t.Run("pushes down is_private false as is:public", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"repositories_grafana","filters":[{"name":"is_private","conditions":[{"operator":"=","value":"false"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		repo, _ := raw["repository"].(string)
		if repo != "is:public" {
			t.Errorf("expected repository 'is:public', got %q", repo)
		}
	})

	t.Run("pushes down tool_name filter for code-scanning", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"code-scanning_grafana_grafana","filters":[{"name":"tool_name","conditions":[{"operator":"=","value":"CodeQL"}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana","filters":[{"name":"state","conditions":[{"operator":"=","value":""}]}]}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
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
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query")
		}
		if out.Queries[0].QueryType != string(models.QueryTypeIssues) {
			t.Errorf("queryType: got %q, want %q", out.Queries[0].QueryType, string(models.QueryTypeIssues))
		}
	})
}

func TestNormalizeGrafanaSQLRequestWithoutFeatureToggle(t *testing.T) {
	t.Run("drops grafanaSql query when GrafanaConfig is nil", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil {
			t.Fatal("expected non-nil response")
		}
		if len(out.Queries) != 0 {
			t.Errorf("expected grafanaSql query to be dropped when GrafanaConfig is nil, got %d queries", len(out.Queries))
		}
	})

	t.Run("drops grafanaSql query when feature toggle is not set", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{
				GrafanaConfig: config.NewGrafanaCfg(map[string]string{}),
			},
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil {
			t.Fatal("expected non-nil response")
		}
		if len(out.Queries) != 0 {
			t.Errorf("expected grafanaSql query to be dropped when feature toggle is not set, got %d queries", len(out.Queries))
		}
	})

	t.Run("preserves non-grafanaSql queries when toggle is off", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","queryType":"Pull_Requests","owner":"grafana","repository":"grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{
				GrafanaConfig: config.NewGrafanaCfg(map[string]string{}),
			},
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: queryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query, got %v", out)
		}
		if string(out.Queries[0].JSON) != string(queryJSON) {
			t.Error("non-grafanaSql query should be preserved unchanged")
		}
	})

	t.Run("drops only grafanaSql queries from mixed request", func(t *testing.T) {
		sqlQueryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana"}`)
		normalQueryJSON := []byte(`{"refId":"B","queryType":"Pull_Requests","owner":"grafana","repository":"grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{
				GrafanaConfig: config.NewGrafanaCfg(map[string]string{}),
			},
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: sqlQueryJSON},
				{RefID: "B", JSON: normalQueryJSON},
			},
		}
		out := normalizeGrafanaSQLRequest(req)
		if out == nil || len(out.Queries) != 1 {
			t.Fatalf("expected one query (non-SQL preserved), got %d", len(out.Queries))
		}
		if out.Queries[0].RefID != "B" {
			t.Errorf("expected preserved query to be refId B, got %s", out.Queries[0].RefID)
		}
	})
}

func TestTimeFieldDefaults(t *testing.T) {
	t.Run("issues default timeField is created (0)", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries:       []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		tf, ok := opts["timeField"].(float64)
		if !ok {
			t.Fatal("expected timeField to be set")
		}
		if int(tf) != int(models.IssueCreatedAt) {
			t.Errorf("expected timeField %d (IssueCreatedAt), got %d", models.IssueCreatedAt, int(tf))
		}
	})

	t.Run("pull-requests default timeField is created (1)", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"pull-requests_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries:       []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		tf := int(opts["timeField"].(float64))
		if tf != int(models.PullRequestCreatedAt) {
			t.Errorf("expected timeField %d (PullRequestCreatedAt), got %d", models.PullRequestCreatedAt, tf)
		}
	})

	t.Run("workflows default timeField is none (0)", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"workflows_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries:       []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		tf := int(opts["timeField"].(float64))
		if tf != int(models.WorkflowTimeFieldNone) {
			t.Errorf("expected timeField %d (WorkflowTimeFieldNone), got %d", models.WorkflowTimeFieldNone, tf)
		}
	})

	t.Run("timeField table parameter overrides default for PRs", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"pull-requests","tableParameterValues":{"organization":"grafana","repository":"grafana","timeField":"merged"}}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries:       []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		tf := int(opts["timeField"].(float64))
		if tf != int(models.PullRequestMergedAt) {
			t.Errorf("expected timeField %d (PullRequestMergedAt), got %d", models.PullRequestMergedAt, tf)
		}
	})

	t.Run("timeField table parameter sets closed for issues", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"issues","tableParameterValues":{"organization":"grafana","repository":"grafana","timeField":"closed"}}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries:       []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		tf := int(opts["timeField"].(float64))
		if tf != int(models.IssueClosedAt) {
			t.Errorf("expected timeField %d (IssueClosedAt), got %d", models.IssueClosedAt, tf)
		}
	})

	t.Run("timeField table parameter sets updated for workflows", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"workflows","tableParameterValues":{"organization":"grafana","repository":"grafana","timeField":"updated"}}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries:       []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		tf := int(opts["timeField"].(float64))
		if tf != int(models.WorkflowUpdatedAt) {
			t.Errorf("expected timeField %d (WorkflowUpdatedAt), got %d", models.WorkflowUpdatedAt, tf)
		}
	})

	t.Run("commits does not get timeField set", func(t *testing.T) {
		queryJSON := []byte(`{"refId":"A","grafanaSql":true,"table":"commits_grafana_grafana"}`)
		req := &backend.QueryDataRequest{
			PluginContext: pluginCtxWithFeatureToggle(),
			Queries:       []backend.DataQuery{{RefID: "A", JSON: queryJSON}},
		}
		out := normalizeGrafanaSQLRequest(req)
		var raw map[string]interface{}
		if err := json.Unmarshal(out.Queries[0].JSON, &raw); err != nil {
			t.Fatal(err)
		}
		opts := raw["options"].(map[string]interface{})
		if _, exists := opts["timeField"]; exists {
			t.Error("expected no timeField for commits")
		}
	})
}
