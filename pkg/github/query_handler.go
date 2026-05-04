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
	register := func(qt models.QueryType, handler backend.QueryDataHandlerFunc) {
		mux.HandleFunc(string(qt), handler)
	}

	register(models.QueryTypeCommits, s.HandleCommits)
	register(models.QueryTypeIssues, s.HandleIssues)
	register(models.QueryTypeContributors, s.HandleContributors)
	register(models.QueryTypeLabels, s.HandleLabels)
	register(models.QueryTypePullRequests, s.HandlePullRequests)
	register(models.QueryTypePullRequestReviews, s.HandlePullRequestReviews)
	register(models.QueryTypeReleases, s.HandleReleases)
	register(models.QueryTypeTags, s.HandleTags)
	register(models.QueryTypePackages, s.HandlePackages)
	register(models.QueryTypeMilestones, s.HandleMilestones)
	register(models.QueryTypeRepositories, s.HandleRepositories)
	register(models.QueryTypeVulnerabilities, s.HandleVulnerabilities)
	register(models.QueryTypeProjects, s.HandleProjects)
	register(models.QueryTypeStargazers, s.HandleStargazers)
	register(models.QueryTypeWorkflows, s.HandleWorkflows)
	register(models.QueryTypeWorkflowUsage, s.HandleWorkflowUsage)
	register(models.QueryTypeWorkflowRuns, s.HandleWorkflowRuns)
	register(models.QueryTypeCodeScanning, s.HandleCodeScanning)
	register(models.QueryTypeDeployments, s.HandleDeployments)
	register(models.QueryTypeOrganizations, s.HandleOrganizations)
	register(models.QueryTypeCommitFiles, s.HandleCommitFiles)
	register(models.QueryTypePullRequestFiles, s.HandlePullRequestFiles)

	return mux
}
