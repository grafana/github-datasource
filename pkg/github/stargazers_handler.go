package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handleStargazersQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.StargazersQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleStargazersQuery(ctx, query, q))
}

// HandleStargazers handles the plugin query for GitHub stargazers
func (s *QueryHandler) HandleStargazers(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleStargazersQuery),
	}, nil
}
