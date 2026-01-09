package main

import (
	"drawing-api/internal/handlers"
	"drawing-api/internal/storage"
	"drawing-api/internal/service"
	"drawing-api/internal/router"
	"drawing-api/internal/ws"
	"drawing-api/internal/db"
	"time"
	"os"
)

const (
	accessExpiry = time.Minute * 5
	refreshExpiry = time.Hour * 24 * 7
	dbTimeout = time.Second * 10
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func main() {
	err := db.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	authHandler := handlers.NewAuthHandler(
		storage.NewPGUserStorage(db.DB),
		storage.NewPGSessionStorage(db.DB),
		service.NewTokenService(jwtSecret),
		service.NewCookieService("refresh_token"),
		dbTimeout,
		accessExpiry,
		refreshExpiry,
	)

	wsHandler := handlers.NewWsHandler(ws.NewHub())
			
	router.InitRouter(authHandler, wsHandler)
}
