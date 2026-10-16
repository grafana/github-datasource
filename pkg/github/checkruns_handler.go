package github

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
)

func (s *QueryHandler) handleCheckRunsRequests(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.CheckRunsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCheckRunsQuery(ctx, query, q))
}

// HandleCheckRuns handles the plugin query for github check runs
func (s *QueryHandler) HandleCheckRuns(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleCheckRunsRequests),
	}, nil
}
