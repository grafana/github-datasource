package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handlePackagesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.PackagesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandlePackagesQuery(ctx, query, q))
}

// HandlePackages handles the plugin query for github Packages
func (s *Server) HandlePackages(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handlePackagesQuery),
	}, nil
}
