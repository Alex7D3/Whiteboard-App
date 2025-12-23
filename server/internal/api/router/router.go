package router

import (
	"net/http"
	"log"
	"drawing-api/internal/api/handlers"
	"drawing-api/internal/api"

)


func InitRouter(authHandler *handlers.AuthHandler) {
	mux := http.NewServeMux()

	// Public
	mux.Handle("POST /register", api.AppHandler(authHandler.Register))
	mux.Handle("POST /login",    api.AppHandler(authHandler.Login))

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
