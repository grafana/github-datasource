package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleCodeSearchQuery(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	model := &models.CodeSearchQuery{}
	if err := UnmarshalQuery(query.JSON, model); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCodeSearchQuery(ctx, model, query))
}

// HandleCodeSearch handles code-search queries.
func (s *QueryHandler) HandleCodeSearch(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{Responses: processQueries(ctx, req, s.handleCodeSearchQuery)}, nil
}
