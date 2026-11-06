package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handlePullRequestsForBranchesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.PullRequestsForBranchesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandlePullRequestsForBranchesQuery(ctx, query, q))
}

// HandlePullRequestsForBranches handles the plugin query for looking up github pull requests by head branch
func (s *QueryHandler) HandlePullRequestsForBranches(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handlePullRequestsForBranchesQuery),
	}, nil
}
