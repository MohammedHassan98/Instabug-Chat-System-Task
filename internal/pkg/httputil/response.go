package httputil

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error  string      `json:"error,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, err string) {
	WriteJSON(w, status, ErrorResponse{Error: err})
}

func WriteValidationErrors(w http.ResponseWriter, errors interface{}) {
	WriteJSON(w, http.StatusBadRequest, ErrorResponse{Errors: errors})
}
