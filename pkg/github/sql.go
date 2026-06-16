package github

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	schemas "github.com/grafana/schemads"
)

// tableToQueryType maps normalized table names to their QueryType constants.
// Built from the query type constants via normalizeTableNames so the schema
// table definitions and this map stay in sync automatically.
var tableToQueryType = func() map[string]models.QueryType {
	qts := []models.QueryType{
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
		models.QueryTypeDeployments,
		models.QueryTypeOrganizations,
		models.QueryTypeGraphQL,
	}
	m := make(map[string]models.QueryType, len(qts))
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

func anyToString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func extractFilterValues(condition schemas.FilterCondition) []string {
	out := make([]string, 0, len(condition.Values)+1)
	for _, v := range condition.Values {
		if s := anyToString(v); s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		if s := anyToString(condition.Value); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// applyFilters maps SQL filter predicates to GitHub API query options.
// It modifies the options map in-place and returns a list of GitHub search
// qualifiers for query types that use the search API.
func applyFilters(queryType models.QueryType, options map[string]interface{}, filters []schemas.ColumnFilter) []string {
	var searchQualifiers []string

	opts, _ := options["options"].(map[string]interface{})
	if opts == nil {
		opts = make(map[string]interface{})
		options["options"] = opts
	}

	appendEqualitySearchQualifier := func(name string, operator schemas.Operator, values []string, isJSON bool) {
		if operator == schemas.OperatorEquals || operator == schemas.OperatorIn {
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
	setOption := func(name string, operator schemas.Operator, value string) {
		if operator == schemas.OperatorEquals || operator == schemas.OperatorIn {
			opts[name] = value
		}
	}

	for _, f := range filters {
		if f.Name == "" || len(f.Conditions) == 0 {
			continue
		}

		for _, condition := range f.Conditions {
			values := extractFilterValues(condition)
			if len(values) == 0 {
				continue
			}

			switch queryType {
			case models.QueryTypeIssues:
				switch f.Name {
				case "state":
					appendEqualitySearchQualifier(f.Name, condition.Operator, values, false)
				case "author":
					appendEqualitySearchQualifier(f.Name, condition.Operator, values, false)
				case "labels":
					appendEqualitySearchQualifier("label", condition.Operator, values, true)
				case "assignees":
					appendEqualitySearchQualifier("assignee", condition.Operator, values, true)
				case "milestone":
					appendEqualitySearchQualifier("milestone", condition.Operator, values, false)
				}
			case models.QueryTypePullRequests, models.QueryTypePullRequestReviews:
				switch f.Name {
				case "author_login":
					appendEqualitySearchQualifier("author", condition.Operator, values, false)
				case "labels":
					appendEqualitySearchQualifier("label", condition.Operator, values, true)
				case "is_draft":
					if condition.Operator == schemas.OperatorEquals || condition.Operator == schemas.OperatorIn {
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
					setOption("state", condition.Operator, values[0])
				case "rule_severity":
					setOption("severity", condition.Operator, values[0])
				case "tool_name":
					setOption("toolName", condition.Operator, values[0])
				}
			case models.QueryTypeWorkflowRuns:
				switch f.Name {
				case "head_branch":
					setOption("branch", condition.Operator, values[0])
				case "status":
					setOption("status", condition.Operator, values[0])
				case "event":
					setOption("event", condition.Operator, values[0])
				}
			case models.QueryTypeContributors:
				if f.Name == "name" && (condition.Operator == schemas.OperatorLike || condition.Operator == schemas.OperatorEquals || condition.Operator == schemas.OperatorIn) {
					opts["query"] = values[0]
				}
			case models.QueryTypeLabels:
				if f.Name == "name" && (condition.Operator == schemas.OperatorLike || condition.Operator == schemas.OperatorEquals || condition.Operator == schemas.OperatorIn) {
					opts["query"] = values[0]
				}
			case models.QueryTypeMilestones:
				if f.Name == "title" && (condition.Operator == schemas.OperatorLike || condition.Operator == schemas.OperatorEquals || condition.Operator == schemas.OperatorIn) {
					opts["query"] = values[0]
				}
			case models.QueryTypePackages:
				switch f.Name {
				case "name":
					if condition.Operator == schemas.OperatorEquals || condition.Operator == schemas.OperatorIn {
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
					if condition.Operator == schemas.OperatorEquals || condition.Operator == schemas.OperatorIn {
						opts["packageType"] = values[0]
					}
				}
			case models.QueryTypeRepositories:
				appendRepoSearchQualifier := func(qualifier string) {
					existing, _ := options["repository"].(string)
					options["repository"] = strings.TrimSpace(existing + " " + qualifier)
				}
				switch f.Name {
				case "name":
					if condition.Operator == schemas.OperatorLike || condition.Operator == schemas.OperatorEquals || condition.Operator == schemas.OperatorIn {
						appendRepoSearchQualifier(values[0])
					}
				case "is_fork":
					if condition.Operator == schemas.OperatorEquals && values[0] == "true" {
						appendRepoSearchQualifier("fork:only")
					}
				case "is_private":
					if condition.Operator == schemas.OperatorEquals {
						if values[0] == "true" {
							appendRepoSearchQualifier("is:private")
						} else {
							appendRepoSearchQualifier("is:public")
						}
					}
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

func resolveTimeField(queryType models.QueryType, value string) (any, bool) {
	switch queryType {
	case models.QueryTypeIssues:
		switch value {
		case "created":
			return models.IssueCreatedAt, true
		case "closed":
			return models.IssueClosedAt, true
		case "updated":
			return models.IssueUpdatedAt, true
		}
	case models.QueryTypePullRequests, models.QueryTypePullRequestReviews:
		switch value {
		case "closed":
			return models.PullRequestClosedAt, true
		case "created":
			return models.PullRequestCreatedAt, true
		case "merged":
			return models.PullRequestMergedAt, true
		case "updated":
			return models.PullRequestUpdatedAt, true
		}
	case models.QueryTypeWorkflows:
		switch value {
		case "created":
			return models.WorkflowCreatedAt, true
		case "updated":
			return models.WorkflowUpdatedAt, true
		}
	}
	return 0, false
}

func defaultTimeField(queryType models.QueryType) int {
	switch queryType {
	case models.QueryTypePullRequests, models.QueryTypePullRequestReviews:
		return int(models.PullRequestCreatedAt)
	default:
		return 0
	}
}

func normalizeGrafanaSQLRequest(req *backend.QueryDataRequest) *backend.QueryDataRequest {
	if req == nil || len(req.Queries) == 0 {
		return req
	}

	grafanaConfig := req.PluginContext.GrafanaConfig
	queries := make([]backend.DataQuery, 0, len(req.Queries))
	for _, q := range req.Queries {
		var query schemas.Query
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
		if v := strings.TrimSpace(anyToString(query.TableParameterValues["organization"])); v != "" {
			owner = v
		}
		if v := strings.TrimSpace(anyToString(query.TableParameterValues["repository"])); v != "" {
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
		if v := strings.TrimSpace(anyToString(query.TableParameterValues["workflow"])); v != "" {
			opts, _ := normalized["options"].(map[string]interface{})
			opts["workflow"] = v
		}

		switch queryType {
		case models.QueryTypeIssues, models.QueryTypePullRequests, models.QueryTypePullRequestReviews, models.QueryTypeWorkflows:
			opts, _ := normalized["options"].(map[string]interface{})
			if tfStr := strings.TrimSpace(anyToString(query.TableParameterValues["timeField"])); tfStr != "" {
				if tf, ok := resolveTimeField(queryType, tfStr); ok {
					opts["timeField"] = tf
				} else {
					opts["timeField"] = defaultTimeField(queryType)
				}
			} else {
				opts["timeField"] = defaultTimeField(queryType)
			}
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
			QueryType:     string(queryType),
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
