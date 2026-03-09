package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"tech-ip-sem2/services/tasks/internal/cache"
	"tech-ip-sem2/services/tasks/internal/client/authclient"
	taskshttp "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/services/tasks/internal/rabbitmq"
	"tech-ip-sem2/services/tasks/internal/repository"
	"tech-ip-sem2/services/tasks/internal/service"
	"tech-ip-sem2/shared/logger"
	"tech-ip-sem2/shared/metrics"
	"tech-ip-sem2/shared/middleware"
)

func main() {
	// Логгер
	log := logger.New("tasks")
	defer log.Sync()

	instanceID := os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		instanceID = "unknown"
	}

	log.Info("Tasks service starting",
		zap.String("instance", instanceID),
		zap.String("version", "1.0.0"),
	)

	// Метрики
	metrics := metrics.New("tasks")

	portStr := os.Getenv("TASKS_PORT")
	if portStr == "" {
		portStr = "8082"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Invalid port number", zap.Error(err))
	}

	// Auth gRPC клиент
	authGRPCAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authGRPCAddr == "" {
		authGRPCAddr = "localhost:50051"
	}

	authClient, err := authclient.NewClient(authGRPCAddr, 3*time.Second, log)
	if err != nil {
		log.Fatal("Failed to create auth client", zap.Error(err))
	}
	defer authClient.Close()

	// Подключение к PostgreSQL
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	var taskRepo repository.TaskRepository
	if dbHost != "" && dbUser != "" {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode)

		repo, err := repository.NewPostgresTaskRepository(connStr)
		if err != nil {
			log.Warn("Failed to connect to database, falling back to in-memory storage",
				zap.Error(err))
			taskRepo = nil
		} else {
			taskRepo = repo
			log.Info("Connected to PostgreSQL database")
			defer repo.Close()
		}
	} else {
		log.Info("Database not configured, using in-memory storage")
	}

	// Подключение к Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	baseTTL, _ := strconv.Atoi(os.Getenv("CACHE_TTL_SECONDS"))
	if baseTTL == 0 {
		baseTTL = 120
	}

	jitterMax, _ := strconv.Atoi(os.Getenv("CACHE_TTL_JITTER_SECONDS"))
	if jitterMax == 0 {
		jitterMax = 30
	}

	redisCache := cache.NewRedisCache(cache.CacheConfig{
		Addr:      redisAddr,
		Password:  redisPassword,
		DB:        redisDB,
		BaseTTL:   baseTTL,
		JitterMax: jitterMax,
	}, log)

	if redisCache.IsEnabled() {
		log.Info("Redis cache enabled",
			zap.String("addr", redisAddr),
			zap.Int("base_ttl", baseTTL),
			zap.Int("jitter_max", jitterMax),
		)
		defer redisCache.Close()
	} else {
		log.Info("Redis cache disabled")
	}

	// Подключение к RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("RABBITMQ_QUEUE")

	var rabbitPublisher *rabbitmq.Publisher
	if rabbitURL != "" && queueName != "" {
		log.Info("Attempting to connect to RabbitMQ",
			zap.String("url", rabbitURL),
			zap.String("queue", queueName))

		pub, err := rabbitmq.NewPublisher(rabbitmq.PublisherConfig{
			URL:   rabbitURL,
			Queue: queueName,
		}, log)
		if err != nil {
			log.Warn("Failed to connect to RabbitMQ, events will be disabled",
				zap.Error(err))
			rabbitPublisher = nil
		} else {
			rabbitPublisher = pub
			log.Info("RabbitMQ publisher connected successfully",
				zap.String("queue", queueName))
			defer rabbitPublisher.Close()
		}
	} else {
		log.Info("RabbitMQ not configured (missing URL or queue), publisher disabled")
	}

	// Сервис задач с кэшем и RabbitMQ
	tasksService := service.NewTasksService(log, taskRepo, redisCache, rabbitPublisher)
	handlers := taskshttp.NewHandlers(tasksService, authClient, log)

	mux := http.NewServeMux()

	// Эндпоинты API
	mux.HandleFunc("POST /v1/tasks", handlers.AuthMiddleware(handlers.CreateTask))
	mux.HandleFunc("GET /v1/tasks", handlers.AuthMiddleware(handlers.ListTasks))
	mux.HandleFunc("GET /v1/tasks/search", handlers.AuthMiddleware(handlers.SearchTasks))
	mux.HandleFunc("GET /v1/tasks/{id}", handlers.AuthMiddleware(handlers.GetTask))
	mux.HandleFunc("PATCH /v1/tasks/{id}", handlers.AuthMiddleware(handlers.UpdateTask))
	mux.HandleFunc("DELETE /v1/tasks/{id}", handlers.AuthMiddleware(handlers.DeleteTask))

	// Метрики и health
	mux.Handle("GET /metrics", metrics.Handler())
	mux.HandleFunc("GET /health", handlers.Health)

	// Middleware
	handler := middleware.RequestID(mux)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.AccessLog(log)(handler)
	handler = middleware.Metrics(metrics)(handler)
	handler = middleware.CSRFMiddleware(log)(handler)

	addr := ":" + strconv.Itoa(port)
	log.Info("Tasks service starting",
		zap.Int("port", port),
		zap.String("auth_grpc_addr", authGRPCAddr),
		zap.Bool("database_enabled", taskRepo != nil),
		zap.Bool("cache_enabled", redisCache.IsEnabled()),
		zap.Bool("rabbitmq_enabled", rabbitPublisher != nil),
		zap.String("instance", instanceID),
	)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("Server failed", zap.Error(err))
	}
}
