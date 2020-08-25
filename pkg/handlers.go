package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type OAuth2Client interface {
	HandleAuth(http.ResponseWriter, *http.Request)
	HandleAuthCallback(http.ResponseWriter, *http.Request)
}

type LabelsHandler interface {
	HandleGetLabels(http.ResponseWriter, *http.Request)
}

type MilestonesHandler interface {
	HandleGetMilestones(http.ResponseWriter, *http.Request)
}

type Handlers struct {
	OAuth2     OAuth2Client
	Labels     LabelsHandler
	Milestones MilestonesHandler
}

func MustGetRouter(h Handlers) *mux.Router {
	router := mux.NewRouter()
	router.Path("/labels").Methods("GET").HandlerFunc(h.Labels.HandleGetLabels)
	router.Path("/milestones").Methods("GET").HandlerFunc(h.Milestones.HandleGetMilestones)
	router.Path("/auth").Methods("GET").HandlerFunc(h.OAuth2.HandleAuth)
	router.Path("/auth/callback").Methods("GET").HandlerFunc(h.OAuth2.HandleAuthCallback)

	return router
}
