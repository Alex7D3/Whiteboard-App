package main

import (
	"drawing-api/internal/db"
	"fmt"
)

func main() {
	err := db.InitDB()
	db := db.DB
	if err != nil {
		fmt.Println(err)
		return	
	}
	defer db.Close()

	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec(`CREATE TABLE IF NOT EXIST users (    
    	user_id serial PRIMARY KEY,
		email varchar(255) UNIQUE NOT NULL,
		username varchar(40) UNIQUE NOT NULL,
		password_hash varchar(255) NOT NULL
		created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`)
}
