package github

import (
	"context"
	"encoding/json"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/pkg/errors"
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
			Error: errors.Wrap(err, "failed to unmarshal JSON request into query"),
		}
	}
	return nil
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

	return mux
}
