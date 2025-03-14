package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *QueryHandler) handlePullRequestReviewsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.PullRequestsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleReviewsQuery(ctx, query, q))
}

// HandlePullRequests handles the plugin query for github PullRequests
func (s *QueryHandler) HandlePullRequestReviews(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handlePullRequestReviewsQuery),
	}, nil
}
