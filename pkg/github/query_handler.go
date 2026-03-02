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

// pullRequestsTablePrefix is the table name prefix used by dsabstraction for pull-requests.
// Table format: "pull-requests_<owner>_<repo>" e.g. "pull-requests_grafana_grafana".
const pullRequestsTablePrefix = "pull-requests_"

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
		if !grafanaSql || table == "" || !strings.HasPrefix(table, pullRequestsTablePrefix) {
			queries = append(queries, q)
			continue
		}
		// Table format: "pull-requests_<owner>_<repo>" (repo may contain underscores)
		parts := strings.SplitN(table, "_", 3)
		if len(parts) < 3 {
			queries = append(queries, q)
			continue
		}
		owner, repo := parts[1], parts[2]
		normalized := map[string]interface{}{
			"refId":       raw["refId"],
			"datasource":  raw["datasource"],
			"queryType":   models.QueryTypePullRequests,
			"owner":       owner,
			"repository":  repo,
			"options":     map[string]interface{}{},
		}
		jsonBytes, err := json.Marshal(normalized)
		if err != nil {
			queries = append(queries, q)
			continue
		}
		queries = append(queries, backend.DataQuery{
			RefID:         q.RefID,
			QueryType:     models.QueryTypePullRequests,
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

	// This could be a map[models.QueryType]datasource.QueryHandlerFunc and then a loop to handle all of them.
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
