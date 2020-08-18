package main

import (
	"encoding/json"
	"net/http"

	"context"

	dserrors "github.com/grafana/grafana-github-datasource/pkg/errors"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/pkg/errors"
)

type QueryHandlerFunc func(context.Context, *models.Query, backend.DataQuery) (backend.DataResponse, error)

// The Datasource type handles the requests sent to the datasource backend
type Datasource interface {
	HandleIssuesQuery(context.Context, *models.Query, backend.DataQuery) (backend.DataResponse, error)
	HandleCommitsQuery(context.Context, *models.Query, backend.DataQuery) (backend.DataResponse, error)
	HandleTagsQuery(context.Context, *models.Query, backend.DataQuery) (backend.DataResponse, error)
	HandleReleasesQuery(context.Context, *models.Query, backend.DataQuery) (backend.DataResponse, error)
	HandleContributorsQuery(context.Context, *models.Query, backend.DataQuery) (backend.DataResponse, error)
}

type OAuth2Client interface {
	HandleAuth(http.ResponseWriter, *http.Request)
	HandleAuthCallback(http.ResponseWriter, *http.Request)
}

func getQueryHandler(d Datasource, q models.QueryType) (QueryHandlerFunc, error) {
	m := map[models.QueryType]QueryHandlerFunc{
		models.QueryTypeCommits:      d.HandleCommitsQuery,
		models.QueryTypeIssues:       d.HandleIssuesQuery,
		models.QueryTypeContributors: d.HandleContributorsQuery,
		models.QueryTypeTags:         d.HandleTagsQuery,
		models.QueryTypeReleases:     d.HandleReleasesQuery,
	}

	if val, ok := m[q]; ok {
		return val, nil
	}

	return nil, dserrors.ErrorQueryTypeUnimplemented
}

// HandleQueryData handles the `QueryData` request for the Github datasource
func HandleQueryData(ctx context.Context, d Datasource, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
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

		response, err := handler(ctx, query, v)
		if err != nil {
			response.Error = err
		}

		responses[v.RefID] = response
	}

	return &backend.QueryDataResponse{
		Responses: responses,
	}, nil
}

// CheckHealth ensures that the datasource settings are able to retrieve data from GitHub
func CheckHealth(ctx context.Context, d Datasource, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}
