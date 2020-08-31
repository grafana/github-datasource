package main

import (
	"encoding/json"

	"context"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	dserrors "github.com/grafana/grafana-github-datasource/pkg/errors"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/pkg/errors"
)

func getQueryHandler(d models.Datasource, q string) (models.QueryHandlerFunc, error) {
	m := map[string]models.QueryHandlerFunc{
		models.QueryTypeCommits:      d.HandleCommitsQuery,
		models.QueryTypeIssues:       d.HandleIssuesQuery,
		models.QueryTypeContributors: d.HandleContributorsQuery,
		models.QueryTypeTags:         d.HandleTagsQuery,
		models.QueryTypeReleases:     d.HandleReleasesQuery,
		models.QueryTypePullRequests: d.HandlePullRequestsQuery,
		models.QueryTypeLabels:       d.HandleLabelsQuery,
	}

	if val, ok := m[q]; ok {
		return val, nil
	}

	return nil, dserrors.ErrorQueryTypeUnimplemented
}

func processQuery(ctx context.Context, d models.Datasource, v backend.DataQuery) (dfutil.Framer, error) {
	var (
		query = &models.Query{}
	)

	if err := json.Unmarshal(v.JSON, &query); err != nil {
		return nil, errors.Wrap(dserrors.ErrorBadQuery, err.Error())
	}

	handler, err := getQueryHandler(d, v.QueryType)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get query handler for querytype '%s'", v.QueryType)
	}

	return handler(ctx, query, v)
}

// HandleQueryData handles the `QueryData` request for the Github datasource
func HandleQueryData(ctx context.Context, d models.Datasource, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	responses := backend.Responses{}
	for _, v := range req.Queries {
		responses[v.RefID] = dfutil.FrameResponseWithError(processQuery(ctx, d, v))
	}

	return &backend.QueryDataResponse{
		Responses: responses,
	}, nil
}

// CheckHealth ensures that the datasource settings are able to retrieve data from GitHub
func CheckHealth(ctx context.Context, d models.Datasource, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}
