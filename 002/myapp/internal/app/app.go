package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/MAMUER/myapp/internal/app/handlers"
	"github.com/MAMUER/myapp/utils"
	"github.com/joho/godotenv"
)

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = utils.NewID16()
		}
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

// getPort получает порт из .env файла или переменной окружения
func getPort() string {
	// Загружаем .env файл (игнорируем ошибку, если файла нет)
	_ = godotenv.Load()

	// Получаем порт из переменной окружения APP_PORT
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // порт по умолчанию
	}

	// Проверяем, что порт валидный
	if _, err := strconv.Atoi(port); err != nil {
		utils.LogError(fmt.Sprintf("Invalid port %s, using default 8080", port))
		port = "8080"
	}

	return port
}

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handlers.Ping)

	mux.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		utils.WriteErr(w, http.StatusBadRequest, "bad_request_example")
	})

	// Корневой маршрут
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Hello, Go project structure!")
	})

	handler := withRequestID(mux)

	port := getPort()
	addr := ":" + port

	utils.LogInfo(fmt.Sprintf("Server is starting on %s", addr))
	if err := http.ListenAndServe(addr, handler); err != nil {
		utils.LogError("server error: " + err.Error())
	}
}
