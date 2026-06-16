package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleDeploymentsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.DeploymentsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleDeploymentsQuery(ctx, query, q))
}

// HandleDeployments handles the plugin query for github Deployments
func (s *QueryHandler) HandleDeployments(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleDeploymentsQuery),
	}, nil
}
