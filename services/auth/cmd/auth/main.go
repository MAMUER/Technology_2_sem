package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	authgrpc "tech-ip-sem2/services/auth/internal/grpc"
	authhttp "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/logger"
	"tech-ip-sem2/shared/metrics"
	"tech-ip-sem2/shared/middleware"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	log := logger.New("auth")
	defer log.Sync()

	metrics := metrics.New("auth")

	httpPortStr := os.Getenv("AUTH_PORT")
	if httpPortStr == "" {
		httpPortStr = "8081"
	}
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		log.Fatal("Invalid HTTP port", zap.Error(err))
	}

	grpcPortStr := os.Getenv("AUTH_GRPC_PORT")
	if grpcPortStr == "" {
		grpcPortStr = "50051"
	}
	grpcPort, err := strconv.Atoi(grpcPortStr)
	if err != nil {
		log.Fatal("Invalid gRPC port", zap.Error(err))
	}

	authService := service.NewAuthService(log)
	sessionService := service.NewSessionService(log)

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			sessionService.CleanupExpired()
		}
	}()

	httpHandlers := authhttp.NewHandlers(authService, sessionService, log)
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("POST /v1/auth/login", httpHandlers.Login)
	httpMux.HandleFunc("POST /v1/auth/logout", httpHandlers.Logout)
	httpMux.HandleFunc("GET /v1/auth/verify", httpHandlers.Verify)
	httpMux.HandleFunc("GET /v1/auth/csrf", httpHandlers.GetCSRFToken)

	httpMux.Handle("GET /metrics", metrics.Handler())

	handler := middleware.RequestID(httpMux)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.DebugRequestID(log)(handler)
	handler = middleware.Metrics(metrics)(handler)
	handler = middleware.AccessLog(log)(handler)

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(httpPort),
		Handler: handler,
	}

	grpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(grpcPort))
	if err != nil {
		log.Fatal("Failed to listen gRPC", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	authgrpc.RegisterAuthServiceServer(grpcServer, authService, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Info("Auth HTTP service starting",
			zap.Int("port", httpPort),
			zap.String("metrics", "/metrics"),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		log.Info("Auth gRPC service starting", zap.Int("port", grpcPort))
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Error("gRPC server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down servers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP shutdown error", zap.Error(err))
	}

	grpcServer.GracefulStop()

	wg.Wait()
	log.Info("Servers stopped")
}
