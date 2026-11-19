package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// HandleMergeQueueTimings handles plugin queries for GitHub merge queue attempts.
func (s *QueryHandler) HandleMergeQueueTimings(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	responses := processQueries(ctx, req, func(ctx context.Context, q backend.DataQuery) backend.DataResponse {
		var query models.Query
		if err := UnmarshalQuery(q.JSON, &query); err != nil {
			return *err
		}
		return dfutil.FrameResponseWithError(GetMergeQueueTimings(ctx, s.Datasource.client, query.Owner, query.Repository, q.TimeRange.From, q.TimeRange.To))
	})
	return &backend.QueryDataResponse{Responses: responses}, nil
}
