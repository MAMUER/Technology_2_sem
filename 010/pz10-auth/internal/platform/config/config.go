package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	JWTSecret  []byte
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func Load() Config {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	accessTTL := os.Getenv("JWT_ACCESS_TTL")
	if accessTTL == "" {
		accessTTL = "15m"
	}
	accessDur, err := time.ParseDuration(accessTTL)
	if err != nil {
		log.Fatal("bad JWT_ACCESS_TTL")
	}

	refreshTTL := os.Getenv("JWT_REFRESH_TTL")
	if refreshTTL == "" {
		refreshTTL = "168h" // 7 дней
	}
	refreshDur, err := time.ParseDuration(refreshTTL)
	if err != nil {
		log.Fatal("bad JWT_REFRESH_TTL")
	}

	return Config{
		Port:       ":" + port,
		JWTSecret:  []byte(secret),
		AccessTTL:  accessDur,
		RefreshTTL: refreshDur,
	}
}