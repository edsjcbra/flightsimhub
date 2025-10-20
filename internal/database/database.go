package database

import (
	"context"
	"log"
	"time"

	"github.com/edsjcbra/flightsimhub/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	Pool, err = pgxpool.New(ctx, config.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// ping
	if err = Pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully")
}

func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("🔌 Database connection closed")
	}
}
