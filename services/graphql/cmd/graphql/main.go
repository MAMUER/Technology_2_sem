package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.uber.org/zap"

	"tech-ip-sem2/services/graphql/graph/generated"
	"tech-ip-sem2/services/graphql/graph/resolvers"
	"tech-ip-sem2/services/graphql/internal/middleware"
	"tech-ip-sem2/services/graphql/internal/repository"
	"tech-ip-sem2/services/graphql/internal/service"
	"tech-ip-sem2/shared/logger"
	sharedmw "tech-ip-sem2/shared/middleware"
)

func main() {
	log := logger.New("graphql")
	defer log.Sync()

	portStr := os.Getenv("GRAPHQL_PORT")
	if portStr == "" {
		portStr = "8090"
	}
	port, _ := strconv.Atoi(portStr)

	// Подключение к БД
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbHost == "" {
		dbHost = "localhost"
		dbPort = "5432"
		dbUser = "postgres"
		dbPass = "postgres"
		dbName = "postgres"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	repo, err := repository.NewPostgresTaskRepository(connStr)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer repo.Close()

	// Сервисы
	taskService := service.NewTaskService(repo, log)
	resolver := resolvers.NewResolver(taskService, log)

	// GraphQL сервер
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	)

	// Маршруты
	mux := http.NewServeMux()
	mux.Handle("/query", srv)
	mux.Handle("/", playground.Handler("GraphQL", "/query"))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"service": "graphql",
			"version": "1.0.0",
		})
	})

	// Middleware
	handler := sharedmw.RequestID(mux)
	handler = middleware.AuthMiddleware(log)(handler)
	handler = sharedmw.AccessLog(log)(handler)

	log.Info("GraphQL running", zap.Int("port", port))
	if err := http.ListenAndServe(":"+strconv.Itoa(port), handler); err != nil {
		log.Fatal("Server failed", zap.Error(err))
	}
}
