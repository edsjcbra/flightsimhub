package database

import (
	"context"
	"log"

	"github.com/edsjcbra/flightsimhub/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool exportado
var Pool *pgxpool.Pool

func Connect() {
	connStr := config.AppConfig.DatabaseURL
	var err error

	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL successfully")
}

func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("🔌 Database connection closed")
	}
}
