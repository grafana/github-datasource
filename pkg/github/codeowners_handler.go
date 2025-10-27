package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleCodeownersQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.CodeownersQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCodeownersQuery(ctx, query, q))
}

// HandleCodeowners handles the plugin query for github codeowners
func (s *QueryHandler) HandleCodeowners(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleCodeownersQuery),
	}, nil
}
