package models

import (
	"context"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type QueryHandlerFunc func(context.Context, *Query, backend.DataQuery) (dfutil.Framer, error)

// The Datasource type handles the requests sent to the datasource backend
type Datasource interface {
	HandleIssuesQuery(context.Context, *Query, backend.DataQuery) (dfutil.Framer, error)
	HandleCommitsQuery(context.Context, *Query, backend.DataQuery) (dfutil.Framer, error)
	HandleTagsQuery(context.Context, *Query, backend.DataQuery) (dfutil.Framer, error)
	HandleReleasesQuery(context.Context, *Query, backend.DataQuery) (dfutil.Framer, error)
	HandleContributorsQuery(context.Context, *Query, backend.DataQuery) (dfutil.Framer, error)
	HandlePullRequestsQuery(context.Context, *Query, backend.DataQuery) (dfutil.Framer, error)
}
