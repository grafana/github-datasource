package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleReleasesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.ReleasesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleReleasesQuery(ctx, query, q))
}

// HandleReleases handles the plugin query for github Releases
func (s *Server) HandleReleases(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleReleasesQuery),
	}, nil
}
