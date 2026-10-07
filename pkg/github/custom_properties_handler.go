package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleCustomPropertiesQuery(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	model := &models.CustomPropertiesQuery{}
	if err := UnmarshalQuery(query.JSON, model); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCustomPropertiesQuery(ctx, model, query))
}

// HandleCustomProperties handles organization custom-property queries.
func (s *QueryHandler) HandleCustomProperties(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{Responses: processQueries(ctx, req, s.handleCustomPropertiesQuery)}, nil
}
