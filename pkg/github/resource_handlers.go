package github

import (
	"context"
	"net/http"

	"github.com/grafana/grafana-github-datasource/pkg/httputil"
	"github.com/grafana/grafana-github-datasource/pkg/models"
)

func (d *Datasource) HandleAuth(http.ResponseWriter, *http.Request)         {}
func (d *Datasource) HandleAuthCallback(http.ResponseWriter, *http.Request) {}

func handleGetLabels(ctx context.Context, client Client, r *http.Request) (Labels, error) {
	q := r.URL.Query()
	opts := models.ListLabelsOptions{
		Repository: q.Get("repository"),
		Owner:      q.Get("owner"),
		Query:      q.Get("query"),
	}

	labels, err := GetAllLabels(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

func (d *Datasource) HandleGetLabels(w http.ResponseWriter, r *http.Request) {
	labels, err := handleGetLabels(r.Context(), d.client, r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httputil.WriteResponse(w, labels)
}

func handleGetMilestones(ctx context.Context, client Client, r *http.Request) (Milestones, error) {
	q := r.URL.Query()
	opts := models.ListMilestonesOptions{
		Repository: q.Get("repository"),
		Owner:      q.Get("owner"),
		Query:      q.Get("query"),
	}

	milestones, err := GetAllMilestones(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	return milestones, nil
}

func (d *Datasource) HandleGetMilestones(w http.ResponseWriter, r *http.Request) {
	milestones, err := handleGetMilestones(r.Context(), d.client, r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httputil.WriteResponse(w, milestones)
}
