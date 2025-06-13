package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleWorkflowsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.WorkflowsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}

	return dfutil.FrameResponseWithError(s.Datasource.HandleWorkflowsQuery(ctx, query, q))
}

// HandleWorkflows handles the plugin query for GitHub workflows
func (s *QueryHandler) HandleWorkflows(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleWorkflowsQuery),
	}, nil
}

func (s *QueryHandler) handleWorkflowUsageQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.WorkflowUsageQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}

	return dfutil.FrameResponseWithError(s.Datasource.HandleWorkflowUsageQuery(ctx, query, q))
}

// HandleWorkflowUsage handles the plugin query for GitHub workflows
func (s *QueryHandler) HandleWorkflowUsage(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleWorkflowUsageQuery),
	}, nil
}

func (s *QueryHandler) handleWorkflowRunsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.WorkflowRunsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}

	return dfutil.FrameResponseWithError(s.Datasource.HandleWorkflowRunsQuery(ctx, query, q))
}

// HandleWorkflowRuns handles the plugin query for GitHub workflows
func (s *QueryHandler) HandleWorkflowRuns(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleWorkflowRunsQuery),
	}, nil
}
