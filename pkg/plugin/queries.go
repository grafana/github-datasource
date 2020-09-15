package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleCommitsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.CommitsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleCommitsQuery(ctx, query, q))
}

// HandleCommits handles the plugin query for github Commits
func (s *Server) HandleCommits(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleCommitsQuery),
	}, nil
}

func (s *Server) handleContributorsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.ContributorsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleContributorsQuery(ctx, query, q))
}

// HandleContributors handles the plugin query for github Contributors
func (s *Server) HandleContributors(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleContributorsQuery),
	}, nil
}

func (s *Server) handleIssuesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.IssuesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleIssuesQuery(ctx, query, q))
}

// HandleIssues handles the plugin query for github Issues
func (s *Server) HandleIssues(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleIssuesQuery),
	}, nil
}

func (s *Server) handleLabelsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.LabelsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleLabelsQuery(ctx, query, q))
}

// HandleLabels handles the plugin query for github Labels
func (s *Server) HandleLabels(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleLabelsQuery),
	}, nil
}

func (s *Server) handleMilestonesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.MilestonesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleMilestonesQuery(ctx, query, q))
}

// HandleMilestones handles the plugin query for github Milestones
func (s *Server) HandleMilestones(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleMilestonesQuery),
	}, nil
}

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

func (s *Server) handlePullRequestsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.PullRequestsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandlePullRequestsQuery(ctx, query, q))
}

// HandlePullRequests handles the plugin query for github PullRequests
func (s *Server) HandlePullRequests(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handlePullRequestsQuery),
	}, nil
}

func (s *Server) handleReleasesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.ReleasesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleReleasesQuery(ctx, query, q))
}

// HandleReleases handles the plugin query for github Releases
func (s *Server) HandleReleases(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleReleasesQuery),
	}, nil
}

func (s *Server) handleRepositoriesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.RepositoriesQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleRepositoriesQuery(ctx, query, q))
}

// HandleRepositories handles the plugin query for github tags
func (s *Server) HandleRepositories(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleRepositoriesQuery),
	}, nil
}

func (s *Server) handleTagsQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.TagsQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleTagsQuery(ctx, query, q))
}

// HandleTags handles the plugin query for github tags
func (s *Server) HandleTags(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleTagsQuery),
	}, nil
}
