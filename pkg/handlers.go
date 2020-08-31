package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ChoicesHandler http.HandlerFunc

type Handlers struct {
	Choices    http.HandlerFunc
	Labels     http.HandlerFunc
	Milestones http.HandlerFunc
}

func MustGetRouter(h Handlers) *mux.Router {
	router := mux.NewRouter()
	// The "choices" endpoint is analogous to the "HandleQueryData" gRPC handler
	// but is used for possible values for dashboard variables
	router.Path("/choices/{queryType}").Methods("POST").HandlerFunc(h.Choices)
	router.Path("/labels").Methods("GET").HandlerFunc(h.Labels)
	router.Path("/milestones").Methods("GET").HandlerFunc(h.Milestones)

	return router
}
