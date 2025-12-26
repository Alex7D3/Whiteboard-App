package api

import (
	"fmt"
	"net/http"
	"errors"
)

type APIError struct {
	Message string
	Status  int
}

func NewAPIError(message string, status int) *APIError {
	return &APIError{message, status}
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error - %v", e.Message)
}

type AppHandler func(http.ResponseWriter, *http.Request) error

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			WriteJSONError(w, apiErr)
		} else {
			WriteJSONError(w, NewAPIError("Internal server error", http.StatusInternalServerError))
		}
	}
}
