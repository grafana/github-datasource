package plugin

import (
	"context"
	"encoding/json"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	dserrors "github.com/grafana/grafana-github-datasource/pkg/errors"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/pkg/errors"
)

type QueryHandlerFunc func(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error)

// The Datasource type handles the requests sent to the datasource backend
type Datasource interface {
	HandleIssuesQuery(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error)
	HandleCommitsQuery(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error)
	HandleTagsQuery(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error)
	HandleReleasesQuery(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error)
	HandleContributorsQuery(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error)
	HandlePullRequestsQuery(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error)
	HandleLabelsQuery(context.Context, *models.Query, backend.DataQuery) (dfutil.Framer, error)
	CheckHealth(context.Context) error
}

func getQueryHandler(d Datasource, q string) (QueryHandlerFunc, error) {
	m := map[string]QueryHandlerFunc{
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

func processQuery(ctx context.Context, d Datasource, v backend.DataQuery) (dfutil.Framer, error) {
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
func HandleQueryData(ctx context.Context, d Datasource, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	responses := backend.Responses{}
	for _, v := range req.Queries {
		responses[v.RefID] = dfutil.FrameResponseWithError(processQuery(ctx, d, v))
	}

	return &backend.QueryDataResponse{
		Responses: responses,
	}, nil
}

// CheckHealth ensures that the datasource settings are able to retrieve data from GitHub
func CheckHealth(ctx context.Context, d Datasource, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if err := d.CheckHealth(ctx); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{Status: backend.HealthStatusOk}, nil
}
