package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func MustGetRouter(h *Datasource) *mux.Router {
	if h == nil {
		panic("datasource is nil")
	}

	router := mux.NewRouter()
	router.Path("/auth").Methods("GET").HandlerFunc(HandleAuth(h))
	router.Path("/auth/callback").Methods("GET").HandlerFunc(HandleAuthCallback(h))

	return router
}

func HandleAuth(ds *Datasource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := ds.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)

		http.Redirect(w, r, url, http.StatusFound)
	}
}

func HandleAuthCallback(ds *Datasource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := ds.oauthConfig.Exchange(r.Context(), r.URL.Query().Get("code"), oauth2.AccessTypeOffline)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(fmt.Sprintf("%+v | %s : %s", token, token.RefreshToken, token.AccessToken)))
	}
}
