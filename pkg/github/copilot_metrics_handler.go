package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleCopilotMetricsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.CopilotMetricsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCopilotMetricsQuery(ctx, query, q))
}

// HandleCopilotMetrics handles the plugin query for GitHub Copilot metrics
func (s *QueryHandler) HandleCopilotMetrics(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleCopilotMetricsQuery),
	}, nil
}

func (s *QueryHandler) handleCopilotMetricsTeamQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.CopilotMetricsTeamQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCopilotMetricsTeamQuery(ctx, query, q))
}

// HandleCopilotMetricsTeam handles the plugin query for GitHub Copilot metrics for a team
func (s *QueryHandler) HandleCopilotMetricsTeam(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleCopilotMetricsTeamQuery),
	}, nil
}
