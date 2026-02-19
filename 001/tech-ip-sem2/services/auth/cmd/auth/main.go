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
	portStr := os.Getenv("AUTH_PORT")
	if portStr == "" {
		portStr = "8081"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	authService := service.NewAuthService()

	handlers := authhttp.NewHandlers(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/auth/login", handlers.Login)
	mux.HandleFunc("GET /v1/auth/verify", handlers.Verify)

	handler := middleware.RequestID(mux)
	handler = middleware.Logging(handler)

	addr := ":" + strconv.Itoa(port)
	log.Printf("Auth service starting on port %d", port)
	log.Fatal(http.ListenAndServe(addr, handler))
}
