package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	taskshttp "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/services/tasks/internal/service"
	"tech-ip-sem2/shared/middleware"
)

func main() {
	portStr := os.Getenv("TASKS_PORT")
	if portStr == "" {
		portStr = "8082"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	authBaseURL := os.Getenv("AUTH_BASE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://localhost:8081"
	}

	authClient := authclient.NewClient(authBaseURL, 3*time.Second)

	tasksService := service.NewTasksService()

	handlers := taskshttp.NewHandlers(tasksService, authClient)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/tasks", handlers.AuthMiddleware(handlers.CreateTask))
	mux.HandleFunc("GET /v1/tasks", handlers.AuthMiddleware(handlers.ListTasks))
	mux.HandleFunc("GET /v1/tasks/{id}", handlers.AuthMiddleware(handlers.GetTask))
	mux.HandleFunc("PATCH /v1/tasks/{id}", handlers.AuthMiddleware(handlers.UpdateTask))
	mux.HandleFunc("DELETE /v1/tasks/{id}", handlers.AuthMiddleware(handlers.DeleteTask))

	handler := middleware.RequestID(mux)
	handler = middleware.Logging(handler)

	addr := ":" + strconv.Itoa(port)
	log.Printf("Tasks service starting on port %d", port)
	log.Printf("Auth service URL: %s", authBaseURL)
	log.Fatal(http.ListenAndServe(addr, handler))
}
