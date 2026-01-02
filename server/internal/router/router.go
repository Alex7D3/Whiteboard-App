package router

import (
	"net/http"
	"log"
	"drawing-api/internal/api"
	"drawing-api/internal/handlers"
)


func InitRouter(authHandler *handlers.AuthHandler, wsHandler *handlers.WsHandler) {
	mux := http.NewServeMux()

	// Public
	mux.Handle("POST /register", api.AppHandler(authHandler.Register))
	mux.Handle("POST /login",    api.AppHandler(authHandler.Login))
	mux.Handle("POST /logout",   api.AppHandler(authHandler.Logout))
	mux.Handle("POST /refresh",  api.AppHandler(authHandler.Refresh))

	mux.Handle("POST /ws/create-room",       api.AppHandler(authHandler.Authorize(wsHandler.CreateRoom)))
	mux.Handle("GET /ws/join-room/{roomID}", api.AppHandler(authHandler.Authorize(wsHandler.JoinRoom)))

	go wsHandler.Hub.Run()

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
