package github

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
)

func (s *QueryHandler) handleCodeScanningRequests(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.CodeScanningQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCodeScanningQuery(ctx, query, q))
}

// HandleCodeScanning handles the plugin query for github code scanning
func (s *QueryHandler) HandleCodeScanning(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleCodeScanningRequests),
	}, nil
}
