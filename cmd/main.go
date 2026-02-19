package main

import (
	"log"
	"os"
	"subscriptions-tracker/internal/db"
	"subscriptions-tracker/internal/handlers"
	"subscriptions-tracker/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file not found, using system env")
	}

	cfg := db.Config{
		DSN: os.Getenv("DB_DSN"),
	}

	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("Error connecting to db", err)
	}
	defer db.Close()

	handler := &handlers.Handler{
		DB: db,
	}

	server.RunServer(handler)
}
