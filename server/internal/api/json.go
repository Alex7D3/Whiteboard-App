package api

import (
	"net/http"
	"encoding/json"
	"fmt"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return NewAPIError("Missing request body", http.StatusBadRequest)
	}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return NewAPIError("Invalid request body", http.StatusBadRequest)
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, output any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(output)
}

func WriteJSONMessage(w http.ResponseWriter, status int, output string) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"%s"}`, output)
	return nil
}

func WriteJSONError(w http.ResponseWriter, apiErr *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Status)
	fmt.Fprintf(w, `{"error":"%s"}`, apiErr)
}
