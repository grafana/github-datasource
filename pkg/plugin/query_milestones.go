package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleMilestonesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.MilestonesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleMilestonesQuery(ctx, query, q))
}

// HandleMilestones handles the plugin query for github Milestones
func (s *Server) HandleMilestones(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleMilestonesQuery),
	}, nil
}
