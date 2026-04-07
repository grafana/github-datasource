package github

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
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

// GetQueryHandlers creates the QueryTypeMux type for handling queries
func GetQueryHandlers(s *QueryHandler) *datasource.QueryTypeMux {
	mux := datasource.NewQueryTypeMux()

	mux.HandleFunc(string(models.QueryTypeCommits), s.HandleCommits)
	mux.HandleFunc(string(models.QueryTypeIssues), s.HandleIssues)
	mux.HandleFunc(string(models.QueryTypeContributors), s.HandleContributors)
	mux.HandleFunc(string(models.QueryTypeLabels), s.HandleLabels)
	mux.HandleFunc(string(models.QueryTypePullRequests), s.HandlePullRequests)
	mux.HandleFunc(string(models.QueryTypePullRequestReviews), s.HandlePullRequestReviews)
	mux.HandleFunc(string(models.QueryTypeReleases), s.HandleReleases)
	mux.HandleFunc(string(models.QueryTypeTags), s.HandleTags)
	mux.HandleFunc(string(models.QueryTypePackages), s.HandlePackages)
	mux.HandleFunc(string(models.QueryTypeMilestones), s.HandleMilestones)
	mux.HandleFunc(string(models.QueryTypeRepositories), s.HandleRepositories)
	mux.HandleFunc(string(models.QueryTypeVulnerabilities), s.HandleVulnerabilities)
	mux.HandleFunc(string(models.QueryTypeProjects), s.HandleProjects)
	mux.HandleFunc(string(models.QueryTypeStargazers), s.HandleStargazers)
	mux.HandleFunc(string(models.QueryTypeWorkflows), s.HandleWorkflows)
	mux.HandleFunc(string(models.QueryTypeWorkflowUsage), s.HandleWorkflowUsage)
	mux.HandleFunc(string(models.QueryTypeWorkflowRuns), s.HandleWorkflowRuns)
	mux.HandleFunc(string(models.QueryTypeCodeScanning), s.HandleCodeScanning)
	mux.HandleFunc(string(models.QueryTypeDeployments), s.HandleDeployments)
	mux.HandleFunc(string(models.QueryTypeOrganizations), s.HandleOrganizations)
	mux.HandleFunc(string(models.QueryTypeCommitFiles), s.HandleCommitFiles)
	mux.HandleFunc(string(models.QueryTypePullRequestFiles), s.HandlePullRequestFiles)

	return mux
}
