package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool() (*pgxpool.Pool, error) {
	// Получаем URL из переменной окружения или используем дефолтный
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://teacher_app:secure_password_123@localhost:5433/notes?sslmode=disable"
	}

	// Парсим конфигурацию
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Настройка пула соединений
	cfg.MaxConns = 20                    // Максимальное количество соединений
	cfg.MinConns = 5                     // Минимальное количество соединений
	cfg.MaxConnLifetime = time.Hour      // Максимальное время жизни соединения
	cfg.MaxConnIdleTime = 5 * time.Minute // Максимальное время простоя соединения
	cfg.HealthCheckPeriod = time.Minute  // Период проверки здоровья

	// Создаем пул
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL. Pool stats: MaxConns=%d, MinConns=%d", 
		cfg.MaxConns, cfg.MinConns)
	return pool, nil
}