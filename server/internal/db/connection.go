package db

import (
	"fmt"
	"log"
	"os"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	print(connStr)
	DB, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL database (sqlx)")
	return nil
}
