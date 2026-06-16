package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleVulnerabilitiesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.VulnerabilityQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleVulnerabilitiesQuery(ctx, query, q))
}

// HandleVulnerabilities handles the plugin query for github Vulnerabilities
func (s *QueryHandler) HandleVulnerabilities(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleVulnerabilitiesQuery),
	}, nil
}
