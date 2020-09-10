package plugin

import (
	"context"
	"encoding/json"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/pkg/errors"
)

// Server is the
type Server struct {
	Datasource Datasource
}

type QueryHandlerFunc func(context.Context, backend.DataQuery) backend.DataResponse

func UnmarshalQuery(b []byte, v interface{}) *backend.DataResponse {
	if err := json.Unmarshal(b, v); err != nil {
		return &backend.DataResponse{
			Error: errors.Wrap(err, "failed to unmarshal JSON request into query"),
		}
	}
	return nil
}

func processQueries(ctx context.Context, req *backend.QueryDataRequest, handler QueryHandlerFunc) backend.Responses {
	res := backend.Responses{}
	for _, v := range req.Queries {
		res[v.RefID] = handler(ctx, v)
	}

	return res
}

func GetQueryHandlers(s *Server) *datasource.QueryTypeMux {
	mux := datasource.NewQueryTypeMux()

	// This could be a map[models.QueryType]datasource.QueryHandlerFunc and then a loop to handle all of them.
	mux.HandleFunc(models.QueryTypeCommits, s.HandleCommits)
	mux.HandleFunc(models.QueryTypeIssues, s.HandleIssues)
	mux.HandleFunc(models.QueryTypeContributors, s.HandleContributors)
	mux.HandleFunc(models.QueryTypeLabels, s.HandleLabels)
	mux.HandleFunc(models.QueryTypePullRequests, s.HandlePullRequests)
	mux.HandleFunc(models.QueryTypeReleases, s.HandleReleases)
	mux.HandleFunc(models.QueryTypeTags, s.HandleTags)
	mux.HandleFunc(models.QueryTypePackages, s.HandlePackages)
	mux.HandleFunc(models.QueryTypeMilestones, s.HandleMilestones)

	return mux
}
