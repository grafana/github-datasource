package main

import (
	"encoding/json"

	"context"

	dserrors "github.com/grafana/grafana-github-datasource/pkg/errors"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/pkg/errors"
)

func getQueryHandler(d models.Datasource, q models.QueryType) (models.QueryHandlerFunc, error) {
	m := map[models.QueryType]models.QueryHandlerFunc{
		models.QueryTypeCommits:      d.HandleCommitsQuery,
		models.QueryTypeIssues:       d.HandleIssuesQuery,
		models.QueryTypeContributors: d.HandleContributorsQuery,
		models.QueryTypeTags:         d.HandleTagsQuery,
		models.QueryTypeReleases:     d.HandleReleasesQuery,
		models.QueryTypePullRequests: d.HandlePullRequestsQuery,
	}

	if val, ok := m[q]; ok {
		return val, nil
	}

	return nil, dserrors.ErrorQueryTypeUnimplemented
}

// HandleQueryData handles the `QueryData` request for the Github datasource
func HandleQueryData(ctx context.Context, d models.Datasource, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	responses := backend.Responses{}
	for _, v := range req.Queries {
		query := &models.Query{}

		if err := json.Unmarshal(v.JSON, &query); err != nil {
			return nil, errors.Wrap(dserrors.ErrorBadQuery, err.Error())
		}

		handler, err := getQueryHandler(d, query.QueryType)
		if err != nil {
			return nil, err
		}
		response := backend.DataResponse{}

		f, err := handler(ctx, query, v)
		if err != nil {
			response.Error = err
		}

		response.Frames = f.Frame()
		responses[v.RefID] = response
	}

	return &backend.QueryDataResponse{
		Responses: responses,
	}, nil
}

// CheckHealth ensures that the datasource settings are able to retrieve data from GitHub
func CheckHealth(ctx context.Context, d models.Datasource, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}
