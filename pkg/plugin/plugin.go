package plugin

import (
	"context"
	"net/http"

	dserrors "github.com/grafana/github-datasource/pkg/errors"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

// Handler is the plugin entrypoint and implements all of the necessary handler functions for dataqueries, healthchecks, and resources.
type Handler struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirement
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (cr *Handler) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	h, err := cr.im.Get(req.PluginContext)
	if err != nil {
		return nil, err
	}

	if val, ok := h.(*Instance); ok {
		return HandleQueryData(ctx, val, req)
	}
	return nil, dserrors.ErrorBadDatasource
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (cr *Handler) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	h, err := cr.im.Get(req.PluginContext)
	if err != nil {
		return nil, err
	}

	if val, ok := h.(*Instance); ok {
		return CheckHealth(ctx, val, req)
	}

	return &backend.CheckHealthResult{
		Status:      backend.HealthStatusError,
		Message:     dserrors.ErrorBadDatasource.Error(),
		JSONDetails: nil,
	}, dserrors.ErrorBadDatasource
}

// ServeHTTP is the main HTTP handler for serving resource calls
func (cr *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pluginCtx := httpadapter.PluginConfigFromContext(r.Context())

	h, err := cr.im.Get(pluginCtx)
	if err != nil {
		panic(err)
	}

	if ds, ok := h.(*Instance); ok {
		GetRouter(ds.Handlers).ServeHTTP(w, r)
		return
	}

	panic(dserrors.ErrorBadDatasource)
}

// GetDatasourceOpts returns the datasource.ServeOpts for this plugin
func GetDatasourceOpts() datasource.ServeOpts {
	handler := &Handler{
		// creates a instance manager for your plugin. The function passed
		// into `NewInstanceManger` is called when the instance is created
		// for the first time or when a datasource configuration changed.
		im: datasource.NewInstanceManager(newDataSourceInstance),
	}

	return datasource.ServeOpts{
		QueryDataHandler:    handler,
		CheckHealthHandler:  handler,
		CallResourceHandler: httpadapter.New(handler),
	}
}
