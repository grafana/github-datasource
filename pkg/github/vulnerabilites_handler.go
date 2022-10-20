package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleVulnerabilitesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.VulnerabilityQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleVulnerabilitiesQuery(ctx, query, q))
}

// HandleVulnerabilitites handles the plugin query for github Vulnerabilitites
func (s *QueryHandler) HandleVulnerabilitites(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleVulnerabilitesQuery),
	}, nil
}
