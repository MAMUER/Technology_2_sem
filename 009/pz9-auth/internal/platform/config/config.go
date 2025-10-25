package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN     string
	BcryptCost int // например, 12
	Addr       string
}

func Load() Config {
	// Загружаем .env файл
	_ = godotenv.Load()

	cost := 12
	if v := os.Getenv("BCRYPT_COST"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			cost = parsed
		} else {
			log.Printf("Invalid BCRYPT_COST '%s', using default: %d", v, cost)
		}
	}

	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	log.Printf("Config loaded: Addr=%s, BcryptCost=%d", addr, cost)

	return Config{
		DB_DSN:     dsn,
		BcryptCost: cost,
		Addr:       addr,
	}
}
