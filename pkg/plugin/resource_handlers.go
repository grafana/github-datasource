package plugin

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handlers stores the list of http.HandlerFunc functions for the different resource calls
type Handlers struct {
	Labels     http.HandlerFunc
	Milestones http.HandlerFunc
}

// GetRouter creates the gorilla/mux router for the HTTP handlers
func GetRouter(h Handlers) *mux.Router {
	router := mux.NewRouter()
	router.Path("/labels").Methods("GET").HandlerFunc(h.Labels)
	router.Path("/milestones").Methods("GET").HandlerFunc(h.Milestones)

	return router
}
