package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pkg/errors"
)

// WriteError writes an error in JSON format to the ResponseWriter
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	d := map[string]string{
		"error": err.Error(),
	}

	b, marshalError := json.Marshal(d)
	if marshalError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(errors.Wrapf(marshalError, "error when marshalling error '%s'", err.Error()).Error())); err != nil {
			log.DefaultLogger.Error(err.Error())
		}
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		log.DefaultLogger.Error(err.Error())
	}
}

// WriteResponse writes a standard HTTP response to the ResponseWriter in JSON format
func WriteResponse(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		log.DefaultLogger.Error(err.Error())
	}
}
