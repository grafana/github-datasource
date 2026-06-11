package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleBranchesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.BranchesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleBranchesQuery(ctx, query, q))
}

// HandleBranches handles the plugin query for github branches
func (s *QueryHandler) HandleBranches(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleBranchesQuery),
	}, nil
}
