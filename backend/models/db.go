package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", "postgres://postgres:123456@localhost:5432/hdwallet_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}
