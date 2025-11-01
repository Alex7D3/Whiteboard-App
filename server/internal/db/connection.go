package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_NAME"),
	)

	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return err
	}

	// Optional: ping again to verify connection
	if err := DB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL database (sqlx)")
	return nil
}
