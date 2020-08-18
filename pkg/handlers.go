package main

import (
	"github.com/gorilla/mux"
)

func MustGetRouter(h OAuth2Client) *mux.Router {
	if h == nil {
		panic("datasource is nil")
	}

	router := mux.NewRouter()
	router.Path("/auth").Methods("GET").HandlerFunc(h.HandleAuth)
	router.Path("/auth/callback").Methods("GET").HandlerFunc(h.HandleAuthCallback)

	return router
}
