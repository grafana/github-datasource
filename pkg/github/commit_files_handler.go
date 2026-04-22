package github

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
)

func (s *QueryHandler) handleCommitFilesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.CommitFilesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCommitFilesQuery(ctx, query, q))
}

// HandleCommitFiles handles the plugin query for files changed in a GitHub commit
func (s *QueryHandler) HandleCommitFiles(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleCommitFilesQuery),
	}, nil
}

func (s *QueryHandler) handlePullRequestFilesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.PullRequestFilesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandlePullRequestFilesQuery(ctx, query, q))
}

// HandlePullRequestFiles handles the plugin query for files changed in a GitHub pull request
func (s *QueryHandler) HandlePullRequestFiles(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handlePullRequestFilesQuery),
	}, nil
}
