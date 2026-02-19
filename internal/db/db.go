package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	DSN string
}

func Connect(config Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to DB successfully")
	return db, nil
}
