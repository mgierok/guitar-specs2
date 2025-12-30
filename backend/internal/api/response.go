package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ListResponse[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, code, message string, details interface{}) {
	WriteJSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}
