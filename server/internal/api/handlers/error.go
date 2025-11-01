package handlers

import "fmt"

type APIError struct {
	Message string
	Status  int
}

func NewAPIError(message string, status int) *APIError {
	return &APIError{message, status}
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %v", e.Message)
}
