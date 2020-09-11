package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleRepositoriesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.RepositoriesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleRepositoriesQuery(ctx, query, q))
}

// HandleRepositories handles the plugin query for github tags
func (s *Server) HandleRepositories(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleRepositoriesQuery),
	}, nil
}
