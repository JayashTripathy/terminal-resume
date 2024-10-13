package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}


type ApiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

func 	MakeHTTPHandler(fn ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			apiErr := ApiError{
				Error: err.Error(),
			}

			WriteJSON(w, http.StatusInternalServerError, apiErr)
		}

	}
}