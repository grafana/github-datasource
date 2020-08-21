package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	d := map[string]string{
		"error": err.Error(),
	}

	b, marshalError := json.Marshal(d)
	if marshalError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.Wrapf(marshalError, "error when marshalling error '%s'", err.Error()).Error()))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(b)
}

func WriteResponse(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
