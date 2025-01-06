package writer

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrEmptyHTTPStatus   = errors.New("HTTP status must be set")
	ErrEmptyErrorMessage = errors.New("error message cannot be empty")
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// WriteJSON writes the provided data as JSON with a given status code.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	if status == 0 {
		return ErrEmptyHTTPStatus
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteErrJSON writes an error response in JSON format.
func WriteErrJSON(w http.ResponseWriter, status int, title, detail string) error {
	if status == 0 {
		return ErrEmptyHTTPStatus
	}
	if title == "" {
		return ErrEmptyErrorMessage
	}

	errResp := ErrorResponse{
		Status: status,
		Title:  title,
		Detail: detail,
	}

	return WriteJSON(w, status, errResp)
}
