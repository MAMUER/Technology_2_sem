package main

import (
	"context"
	"log"
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
	"tech-ip-sem2/shared/middleware"

	"google.golang.org/grpc"
)

func main() {
	httpPortStr := os.Getenv("AUTH_PORT")
	if httpPortStr == "" {
		httpPortStr = "8081"
	}
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		log.Fatalf("Invalid HTTP port: %v", err)
	}

	// gRPC порт
	grpcPortStr := os.Getenv("AUTH_GRPC_PORT")
	if grpcPortStr == "" {
		grpcPortStr = "50051"
	}
	grpcPort, err := strconv.Atoi(grpcPortStr)
	if err != nil {
		log.Fatalf("Invalid gRPC port: %v", err)
	}

	authService := service.NewAuthService()

	httpHandlers := authhttp.NewHandlers(authService)
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("POST /v1/auth/login", httpHandlers.Login)
	httpMux.HandleFunc("GET /v1/auth/verify", httpHandlers.Verify)

	httpHandler := middleware.RequestID(httpMux)
	httpHandler = middleware.Logging(httpHandler)

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(httpPort),
		Handler: httpHandler,
	}

	grpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	authgrpc.RegisterAuthServiceServer(grpcServer, authService)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Printf("Auth HTTP service starting on port %d", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		log.Printf("Auth gRPC service starting on port %d", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down servers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	grpcServer.GracefulStop()

	wg.Wait()
	log.Println("Servers stopped")
}
