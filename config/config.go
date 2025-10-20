package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

var AppConfig Config

func LoadConfig() {
	_ = godotenv.Load() // ok se não existir

	AppConfig.Port = getEnv("PORT", "8080")
	AppConfig.DatabaseURL = getEnv("DATABASE_URL", "")
	AppConfig.JWTSecret = getEnv("JWT_SECRET", "supersecret")

	if AppConfig.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL not set in environment")
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
