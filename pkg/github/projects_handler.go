package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleProjectsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.ProjectsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleProjectsQuery(ctx, query, q))
}

// HandleProjects handles the plugin query for github Projects
func (s *QueryHandler) HandleProjects(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleProjectsQuery),
	}, nil
}
