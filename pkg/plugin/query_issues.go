package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleIssuesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.IssuesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleIssuesQuery(ctx, query, q))
}

// HandleIssues handles the plugin query for github Issues
func (s *Server) HandleIssues(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleIssuesQuery),
	}, nil
}
