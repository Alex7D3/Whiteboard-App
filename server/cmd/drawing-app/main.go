package main

import (
	"drawing-api/internal/handlers"
	"drawing-api/internal/storage"
	"drawing-api/internal/service"
	"drawing-api/internal/router"
	"drawing-api/internal/db"
	"time"
	"os"
)
func main() {
	var tokenTimeout time.Duration = 5
	var dbTimeout time.Duration = time.Second * 10
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	err := db.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	h := handlers.NewAuthHandler(
		storage.NewPGUserStorage(db.DB),
		service.NewTokenService(jwtSecret, tokenTimeout),
		service.NewCookieService("auth_token", tokenTimeout),
		dbTimeout,
	)
	router.InitRouter(h)
}
