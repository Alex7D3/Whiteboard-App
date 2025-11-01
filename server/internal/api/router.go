package api

import (
	"net/http"
	"log"
	"errors"
	"drawing-api/internal/api/handlers"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		var apiErr handlers.APIError
		if errors.As(err, &apiErr) {
			handlers.WriteJSONError(w, apiErr)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func InitRouter(authHandler *handlers.AuthHandler) {
	mux := http.NewServeMux()

	// Public
	mux.Handle("POST /register", appHandler(authHandler.Register))
	mux.Handle("POST /login",    appHandler(authHandler.Login))
	mux.HandleFunc("POST /refresh", handlers.Refresh)
	mux.HandleFunc("POST /logout", handlers.Logout)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
