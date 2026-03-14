package github

import (
	"encoding/json"
	"strings"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	schemas "github.com/grafana/schemads"
)

// tableToQueryType maps normalized table names to their QueryType constants.
// Built from the query type constants via normalizeTableNames so the schema
// table definitions and this map stay in sync automatically.
var tableToQueryType = func() map[string]string {
	qts := []string{
		models.QueryTypeCommits,
		models.QueryTypeIssues,
		models.QueryTypePullRequests,
		models.QueryTypePullRequestReviews,
		models.QueryTypeRepositories,
		models.QueryTypeContributors,
		models.QueryTypeTags,
		models.QueryTypeReleases,
		models.QueryTypeLabels,
		models.QueryTypeMilestones,
		models.QueryTypePackages,
		models.QueryTypeVulnerabilities,
		models.QueryTypeProjects,
		models.QueryTypeStargazers,
		models.QueryTypeWorkflows,
		models.QueryTypeWorkflowUsage,
		models.QueryTypeWorkflowRuns,
		models.QueryTypeCodeScanning,
		models.QueryTypeOrganizations,
		models.QueryTypeGraphQL,
	}
	m := make(map[string]string, len(qts))
	for _, qt := range qts {
		m[normalizeTableNames(qt)] = qt
	}
	return m
}()

// parseJSONStringValues tries to parse s as a JSON array of strings.
// If successful it returns the individual elements; otherwise it returns
// the original string as a single-element slice so callers can always
// range over the result.
func parseJSONStringValues(s string) []string {
	var arr []string
	if err := json.Unmarshal([]byte(s), &arr); err == nil && len(arr) > 0 {
		return arr
	}
	return []string{s}
}

func extractFilterValues(condition schemas.FilterCondition) []string {
	out := make([]string, 0, len(condition.Values)+1)
	for _, v := range condition.Values {
		if v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 && condition.Value != "" {
		out = append(out, condition.Value)
	}
	return out
}

// applyFilters maps SQL filter predicates to GitHub API query options.
// It modifies the options map in-place and returns a list of GitHub search
// qualifiers for query types that use the search API.
func applyFilters(queryType string, options map[string]interface{}, filters []schemas.ColumnFilter) []string {
	var searchQualifiers []string

	opts, _ := options["options"].(map[string]interface{})
	if opts == nil {
		opts = make(map[string]interface{})
		options["options"] = opts
	}

	appendEqualitySearchQualifier := func(name, operator string, values []string, isJSON bool) {
		if operator == "==" || operator == "=" || operator == "in" {
			for _, value := range values {
				if isJSON {
					for _, v := range parseJSONStringValues(value) {
						searchQualifiers = append(searchQualifiers, name+":"+v)
					}
				} else {
					searchQualifiers = append(searchQualifiers, name+":"+value)
				}
			}
		}
	}
	setOption := func(name, op, value string) {
		if op == "==" || op == "=" || op == "in" {
			opts[name] = value
		}
	}

	for _, f := range filters {
		if f.Name == "" || len(f.Conditions) == 0 {
			continue
		}

		for _, condition := range f.Conditions {
			op := strings.ToLower(strings.TrimSpace(condition.Operator))
			values := extractFilterValues(condition)
			if len(values) == 0 {
				continue
			}

			switch queryType {
			case models.QueryTypeIssues:
				switch f.Name {
				case "state":
					appendEqualitySearchQualifier(f.Name, op, values, false)
				case "author":
					appendEqualitySearchQualifier(f.Name, op, values, false)
				case "labels":
					appendEqualitySearchQualifier("label", op, values, true)
				case "assignees":
					appendEqualitySearchQualifier("assignee", op, values, true)
				case "milestone":
					appendEqualitySearchQualifier("milestone", op, values, false)
				}
			case models.QueryTypePullRequests, models.QueryTypePullRequestReviews:
				switch f.Name {
				case "state":
					appendEqualitySearchQualifier("state", op, values, false)
				case "author_login":
					appendEqualitySearchQualifier("author", op, values, false)
				case "labels":
					appendEqualitySearchQualifier("label", op, values, true)
				case "is_draft":
					if op == "==" || op == "=" || op == "in" {
						for _, value := range values {
							if value == "true" {
								searchQualifiers = append(searchQualifiers, "draft:true")
							} else {
								searchQualifiers = append(searchQualifiers, "draft:false")
							}
						}
					}
				}
			case models.QueryTypeCodeScanning:
				switch f.Name {
				case "state":
					setOption("state", op, values[0])
				case "rule_severity":
					setOption("severity", op, values[0])
				case "tool_name":
					setOption("toolName", op, values[0])
				}
			case models.QueryTypeWorkflowRuns:
				switch f.Name {
				case "head_branch":
					setOption("branch", op, values[0])
				case "status":
					setOption("status", op, values[0])
				case "event":
					setOption("event", op, values[0])
				}
			case models.QueryTypeContributors:
				if f.Name == "name" && (op == "like" || op == "==" || op == "=" || op == "in") {
					opts["query"] = values[0]
				}
			case models.QueryTypeLabels:
				if f.Name == "name" && (op == "like" || op == "==" || op == "=" || op == "in") {
					opts["query"] = values[0]
				}
			case models.QueryTypeMilestones:
				if f.Name == "title" && (op == "like" || op == "==" || op == "=" || op == "in") {
					opts["query"] = values[0]
				}
			case models.QueryTypePackages:
				switch f.Name {
				case "name":
					if op == "==" || op == "=" || op == "in" {
						for _, value := range values {
							existing, _ := opts["names"].(string)
							if existing != "" {
								opts["names"] = existing + "," + value
							} else {
								opts["names"] = value
							}
						}
					}
				case "type":
					if op == "==" || op == "=" || op == "in" {
						opts["packageType"] = values[0]
					}
				}
			case models.QueryTypeRepositories:
				if f.Name == "name" && (op == "like" || op == "==" || op == "=" || op == "in") {
					options["repository"] = values[0]
				}
			}
		}
	}

	if len(searchQualifiers) > 0 {
		existing, _ := opts["query"].(string)
		combined := strings.Join(searchQualifiers, " ")
		if existing != "" {
			opts["query"] = existing + " " + combined
		} else {
			opts["query"] = combined
		}
	}

	return searchQualifiers
}

func normalizeGrafanaSQLRequest(req *backend.QueryDataRequest) *backend.QueryDataRequest {
	if req == nil || len(req.Queries) == 0 {
		return req
	}

	grafanaConfig := req.PluginContext.GrafanaConfig
	queries := make([]backend.DataQuery, 0, len(req.Queries))
	for _, q := range req.Queries {
		var query schemas.GenericQuery
		if err := json.Unmarshal(q.JSON, &query); err != nil {
			queries = append(queries, q)
			continue
		}
		if !query.GrafanaSql || query.Table == "" {
			queries = append(queries, q)
			continue
		}

		if query.GrafanaSql {
			if grafanaConfig == nil {
				backend.Logger.Warn("grafanaConfig is not set, skipping query")
				continue
			}
			if !grafanaConfig.FeatureToggles().IsEnabled("dsAbstractionApp") {
				backend.Logger.Warn("dsAbstractionApp is not enabled, skipping query")
				continue
			}
		}

		// Table names use hyphens only (never underscores), so the first
		// underscore unambiguously separates the table name from the
		// owner/repo suffix: "issues_grafana_grafana" -> "issues" + "grafana_grafana".
		parts := strings.SplitN(query.Table, "_", 2)
		queryType, ok := tableToQueryType[parts[0]]
		if !ok {
			queries = append(queries, q)
			continue
		}

		var owner, repo string
		if len(parts) == 2 {
			ownerRepo := strings.SplitN(parts[1], "_", 2)
			owner = ownerRepo[0]
			if len(ownerRepo) == 2 {
				repo = ownerRepo[1]
			}
		}
		if v := strings.TrimSpace(query.TableParameterValues["organization"]); v != "" {
			owner = v
		}
		if v := strings.TrimSpace(query.TableParameterValues["repository"]); v != "" {
			repo = v
		}

		normalized := map[string]interface{}{
			"refId":      query.RefID,
			"datasource": query.Datasource,
			"queryType":  queryType,
			"owner":      owner,
			"repository": repo,
			"options":    map[string]interface{}{},
		}
		if queryType == models.QueryTypeProjects && owner != "" {
			opts, _ := normalized["options"].(map[string]interface{})
			opts["organization"] = owner
		}
		if v := strings.TrimSpace(query.TableParameterValues["workflow"]); v != "" {
			opts, _ := normalized["options"].(map[string]interface{})
			opts["workflow"] = v
		}

		if len(query.Filters) > 0 {
			applyFilters(queryType, normalized, query.Filters)
		}

		jsonBytes, err := json.Marshal(normalized)
		if err != nil {
			queries = append(queries, q)
			continue
		}
		queries = append(queries, backend.DataQuery{
			RefID:         q.RefID,
			QueryType:     queryType,
			MaxDataPoints: q.MaxDataPoints,
			Interval:      q.Interval,
			TimeRange:     q.TimeRange,
			JSON:          jsonBytes,
		})
	}
	return &backend.QueryDataRequest{
		PluginContext: req.PluginContext,
		Headers:       req.Headers,
		Queries:       queries,
	}
}
