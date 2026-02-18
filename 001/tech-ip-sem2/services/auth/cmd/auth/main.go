// services/auth/cmd/auth/main.go (исправленный)
package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	authhttp "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/middleware"
)

func main() {
	// Загрузка конфигурации
	portStr := os.Getenv("AUTH_PORT")
	if portStr == "" {
		portStr = "8081"
	}

	// Проверяем, что порт - это число
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	// Создание сервиса
	authService := service.NewAuthService()

	// Создание обработчиков
	handlers := authhttp.NewHandlers(authService)

	// Настройка роутера
	mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("POST /v1/auth/login", handlers.Login)
	mux.HandleFunc("GET /v1/auth/verify", handlers.Verify)

	// Применяем middleware
	handler := middleware.RequestID(mux)
	handler = middleware.Logging(handler)

	addr := ":" + strconv.Itoa(port)
	log.Printf("Auth service starting on port %d", port)
	log.Fatal(http.ListenAndServe(addr, handler))
}
