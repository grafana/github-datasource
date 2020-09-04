package plugin

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handlers struct {
	Labels     http.HandlerFunc
	Milestones http.HandlerFunc
}

func MustGetRouter(h Handlers) *mux.Router {
	router := mux.NewRouter()
	router.Path("/labels").Methods("GET").HandlerFunc(h.Labels)
	router.Path("/milestones").Methods("GET").HandlerFunc(h.Milestones)

	return router
}
