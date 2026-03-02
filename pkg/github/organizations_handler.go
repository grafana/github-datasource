package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleOrganizationsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.OrganizationsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleOrganizationsQuery(ctx, query, q))
}

// HandleOrganizations handles the plugin query for GitHub Organizations
func (s *QueryHandler) HandleOrganizations(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleOrganizationsQuery),
	}, nil
}
