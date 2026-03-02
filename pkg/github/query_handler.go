package github

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/pkg/errors"

	"github.com/grafana/github-datasource/pkg/models"
)

// QueryHandler is the main handler for datasource queries.
type QueryHandler struct {
	Datasource Datasource
}

// QueryHandlerFunc is the function signature used for mux.HandleFunc
type QueryHandlerFunc func(context.Context, backend.DataQuery) backend.DataResponse

func processQueries(ctx context.Context, req *backend.QueryDataRequest, handler QueryHandlerFunc) backend.Responses {
	res := backend.Responses{}
	for _, v := range req.Queries {
		res[v.RefID] = handler(ctx, v)
	}

	return res
}

// UnmarshalQuery attempts to unmarshal a query from JSON
func UnmarshalQuery(b []byte, v interface{}) *backend.DataResponse {
	if err := json.Unmarshal(b, v); err != nil {
		return &backend.DataResponse{
			Error:       errors.Wrap(err, "failed to unmarshal JSON request into query"),
			ErrorSource: backend.ErrorSourceDownstream,
		}
	}
	return nil
}

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
		models.QueryTypeProjectItems,
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

// applyFilters maps SQL filter predicates to GitHub API query options.
// It modifies the options map in-place and returns a list of GitHub search
// qualifiers for query types that use the search API.
func applyFilters(queryType string, options map[string]interface{}, filters []map[string]interface{}) []string {
	var searchQualifiers []string

	for _, f := range filters {
		column, _ := f["column"].(string)
		op, _ := f["op"].(string)
		value, _ := f["value"].(string)
		if column == "" || value == "" {
			continue
		}

		switch queryType {
		case models.QueryTypeIssues:
			switch column {
			case "state":
				if op == "==" || op == "=" {
					searchQualifiers = append(searchQualifiers, "state:"+value)
				}
			case "author":
				if op == "==" || op == "=" {
					searchQualifiers = append(searchQualifiers, "author:"+value)
				}
			case "labels":
				if op == "==" || op == "=" {
					for _, v := range parseJSONStringValues(value) {
						searchQualifiers = append(searchQualifiers, "label:"+v)
					}
				}
			case "assignees":
				if op == "==" || op == "=" {
					for _, v := range parseJSONStringValues(value) {
						searchQualifiers = append(searchQualifiers, "assignee:"+v)
					}
				}
			case "milestone":
				if op == "==" || op == "=" {
					searchQualifiers = append(searchQualifiers, "milestone:"+value)
				}
			}

		case models.QueryTypePullRequests, models.QueryTypePullRequestReviews:
			switch column {
			case "state":
				if op == "==" || op == "=" {
					searchQualifiers = append(searchQualifiers, "state:"+value)
				}
			case "author_login":
				if op == "==" || op == "=" {
					searchQualifiers = append(searchQualifiers, "author:"+value)
				}
			case "labels":
				if op == "==" || op == "=" {
					for _, v := range parseJSONStringValues(value) {
						searchQualifiers = append(searchQualifiers, "label:"+v)
					}
				}
			case "is_draft":
				if op == "==" || op == "=" {
					if value == "true" {
						searchQualifiers = append(searchQualifiers, "draft:true")
					} else {
						searchQualifiers = append(searchQualifiers, "draft:false")
					}
				}
			}

		case models.QueryTypeCodeScanning:
			switch column {
			case "state":
				if op == "==" || op == "=" {
					opts, _ := options["options"].(map[string]interface{})
					if opts == nil {
						opts = make(map[string]interface{})
						options["options"] = opts
					}
					opts["state"] = value
				}
			case "rule_severity":
				if op == "==" || op == "=" {
					opts, _ := options["options"].(map[string]interface{})
					if opts == nil {
						opts = make(map[string]interface{})
						options["options"] = opts
					}
					opts["severity"] = value
				}
			case "tool_name":
				if op == "==" || op == "=" {
					opts, _ := options["options"].(map[string]interface{})
					if opts == nil {
						opts = make(map[string]interface{})
						options["options"] = opts
					}
					opts["toolName"] = value
				}
			}

		case models.QueryTypeWorkflowRuns:
			switch column {
			case "head_branch":
				if op == "==" || op == "=" {
					opts, _ := options["options"].(map[string]interface{})
					if opts == nil {
						opts = make(map[string]interface{})
						options["options"] = opts
					}
					opts["branch"] = value
				}
			case "status":
				if op == "==" || op == "=" {
					opts, _ := options["options"].(map[string]interface{})
					if opts == nil {
						opts = make(map[string]interface{})
						options["options"] = opts
					}
					opts["status"] = value
				}
			case "event":
				if op == "==" || op == "=" {
					opts, _ := options["options"].(map[string]interface{})
					if opts == nil {
						opts = make(map[string]interface{})
						options["options"] = opts
					}
					opts["event"] = value
				}
			}

		case models.QueryTypeContributors:
			if column == "name" && (op == "like" || op == "==" || op == "=") {
				opts, _ := options["options"].(map[string]interface{})
				if opts == nil {
					opts = make(map[string]interface{})
					options["options"] = opts
				}
				opts["query"] = value
			}

		case models.QueryTypeLabels:
			if column == "name" && (op == "like" || op == "==" || op == "=") {
				opts, _ := options["options"].(map[string]interface{})
				if opts == nil {
					opts = make(map[string]interface{})
					options["options"] = opts
				}
				opts["query"] = value
			}

		case models.QueryTypeMilestones:
			if column == "title" && (op == "like" || op == "==" || op == "=") {
				opts, _ := options["options"].(map[string]interface{})
				if opts == nil {
					opts = make(map[string]interface{})
					options["options"] = opts
				}
				opts["query"] = value
			}

		case models.QueryTypePackages:
			switch column {
			case "name":
				if op == "==" || op == "=" {
					opts, _ := options["options"].(map[string]interface{})
					if opts == nil {
						opts = make(map[string]interface{})
						options["options"] = opts
					}
					existing, _ := opts["names"].(string)
					if existing != "" {
						opts["names"] = existing + "," + value
					} else {
						opts["names"] = value
					}
				}
			case "type":
				if op == "==" || op == "=" {
					opts, _ := options["options"].(map[string]interface{})
					if opts == nil {
						opts = make(map[string]interface{})
						options["options"] = opts
					}
					opts["packageType"] = value
				}
			}

		case models.QueryTypeRepositories:
			if column == "name" && (op == "like" || op == "==" || op == "=") {
				options["repository"] = value
			}
		}
	}

	if len(searchQualifiers) > 0 {
		opts, _ := options["options"].(map[string]interface{})
		if opts == nil {
			opts = make(map[string]interface{})
			options["options"] = opts
		}
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
	queries := make([]backend.DataQuery, 0, len(req.Queries))
	for _, q := range req.Queries {
		var raw map[string]interface{}
		if err := json.Unmarshal(q.JSON, &raw); err != nil {
			queries = append(queries, q)
			continue
		}
		grafanaSql, _ := raw["grafanaSql"].(bool)
		table, _ := raw["table"].(string)
		if !grafanaSql || table == "" {
			queries = append(queries, q)
			continue
		}

		// Table names use hyphens only (never underscores), so the first
		// underscore unambiguously separates the table name from the
		// owner/repo suffix: "issues_grafana_grafana" -> "issues" + "grafana_grafana".
		parts := strings.SplitN(table, "_", 2)
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

		normalized := map[string]interface{}{
			"refId":      raw["refId"],
			"datasource": raw["datasource"],
			"queryType":  queryType,
			"owner":      owner,
			"repository": repo,
			"options":    map[string]interface{}{},
		}

		if filters, ok := raw["filters"].([]interface{}); ok && len(filters) > 0 {
			filterMaps := make([]map[string]interface{}, 0, len(filters))
			for _, f := range filters {
				if fm, ok := f.(map[string]interface{}); ok {
					filterMaps = append(filterMaps, fm)
				}
			}
			applyFilters(queryType, normalized, filterMaps)
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

// GetQueryHandlers creates the QueryTypeMux type for handling queries
func GetQueryHandlers(s *QueryHandler) *datasource.QueryTypeMux {
	mux := datasource.NewQueryTypeMux()

	mux.HandleFunc(models.QueryTypeCommits, s.HandleCommits)
	mux.HandleFunc(models.QueryTypeIssues, s.HandleIssues)
	mux.HandleFunc(models.QueryTypeContributors, s.HandleContributors)
	mux.HandleFunc(models.QueryTypeLabels, s.HandleLabels)
	mux.HandleFunc(models.QueryTypePullRequests, s.HandlePullRequests)
	mux.HandleFunc(models.QueryTypePullRequestReviews, s.HandlePullRequestReviews)
	mux.HandleFunc(models.QueryTypeReleases, s.HandleReleases)
	mux.HandleFunc(models.QueryTypeTags, s.HandleTags)
	mux.HandleFunc(models.QueryTypePackages, s.HandlePackages)
	mux.HandleFunc(models.QueryTypeMilestones, s.HandleMilestones)
	mux.HandleFunc(models.QueryTypeRepositories, s.HandleRepositories)
	mux.HandleFunc(models.QueryTypeVulnerabilities, s.HandleVulnerabilities)
	mux.HandleFunc(models.QueryTypeProjects, s.HandleProjects)
	mux.HandleFunc(models.QueryTypeStargazers, s.HandleStargazers)
	mux.HandleFunc(models.QueryTypeWorkflows, s.HandleWorkflows)
	mux.HandleFunc(models.QueryTypeWorkflowUsage, s.HandleWorkflowUsage)
	mux.HandleFunc(models.QueryTypeWorkflowRuns, s.HandleWorkflowRuns)
	mux.HandleFunc(models.QueryTypeCodeScanning, s.HandleCodeScanning)
	mux.HandleFunc(models.QueryTypeDeployments, s.HandleDeployments)
	mux.HandleFunc(models.QueryTypeOrganizations, s.HandleOrganizations)

	return mux
}
