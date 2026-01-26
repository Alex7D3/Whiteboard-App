package router

import (
	"net/http"
	"log"
	"drawing-api/internal/api"
	"drawing-api/internal/handlers"
	"os"
)


func InitRouter(authHandler *handlers.AuthHandler, wsHandler *handlers.WsHandler) {
	mux := http.NewServeMux()
	
	// Public
	mux.Handle("POST /api/register", api.AppHandler(authHandler.Register))
	mux.Handle("POST /api/login",    api.AppHandler(authHandler.Login))
	mux.Handle("POST /api/logout",   api.AppHandler(authHandler.Logout))
	mux.Handle("POST /api/refresh",  api.AppHandler(authHandler.Refresh))

	// Authenticated
	mux.Handle("POST /ws/create-room",       api.AppHandler(authHandler.Authorize(wsHandler.CreateRoom)))
	mux.Handle("GET /ws/join-room/{roomID}", api.AppHandler(authHandler.AuthorizeWS(wsHandler.JoinRoom)))
	
	corsMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENT_URI"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return 
		}
		mux.ServeHTTP(w, r)
	})

	go wsHandler.Hub.Run()

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", corsMux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
