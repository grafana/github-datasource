package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleFileContributorsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.FileContributorsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleFileContributorsQuery(ctx, query, q))
}

// HandleFileContributors handles the plugin query for github file contributors
func (s *QueryHandler) HandleFileContributors(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleFileContributorsQuery),
	}, nil
}
