package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		log.Fatal("POSTGRES_URL is not set")
	}

	var err error
	for i := range 5 {
		DB, err = pgxpool.New(context.Background(), url)
		if err == nil {
			if err = DB.Ping(context.Background()); err == nil {
				log.Println("Connected to PostgreSQL database")
				return
			}
			DB.Close()
		}

		log.Printf("Failed to connect to PostgreSQL, retrying in 2 seconds... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Unable to connect to database after retries: %v\n", err)
}
