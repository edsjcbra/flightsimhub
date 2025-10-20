package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/edsjcbra/flightsimhub/config"
)

var DB *pgxpool.Pool

func Connect() {
	dbURL := config.AppConfig.DBUrl

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ Error creating DB pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	DB = pool
	log.Println("✅ Connected to PostgreSQL successfully")
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("🔌 Database connection closed")
	}
}
