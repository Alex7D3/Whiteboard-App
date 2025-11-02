package api

import (
	"net/http"
	"log"
	"errors"
	"drawing-api/internal/api/handlers"

)

type AppHandler func(http.ResponseWriter, *http.Request) error

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		var apiErr APIError
		if errors.As(err, &apiErr) {
			WriteJSONError(w, apiErr)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func InitRouter(authHandler *handlers.AuthHandler) {
	mux := http.NewServeMux()

	// Public
	mux.Handle("POST /register", AppHandler(authHandler.Register))
	mux.Handle("POST /login",    AppHandler(authHandler.Login))
	mux.Handle("POST /refresh", authHandler.Refresh)
	mux.Handle("POST /logout", authHandler.Logout)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
