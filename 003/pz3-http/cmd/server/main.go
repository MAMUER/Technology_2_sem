package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/pz3-http/internal/api"
	"example.com/pz3-http/internal/storage"
)

func main() {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	mux := http.NewServeMux()
	/*mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})*/

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	/*
		// Коллекция
		mux.HandleFunc("GET /tasks", h.ListTasks)
		mux.HandleFunc("POST /tasks", h.CreateTask)
		// Элемент
		mux.HandleFunc("GET /tasks/", h.GetTask)
	*/

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListTasks(w, r)
		case http.MethodPost:
			h.CreateTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTask(w, r)
		case http.MethodPatch:
			h.PatchTask(w, r)
		case http.MethodDelete:
			h.DeleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	/*
		// Подключаем логирование
		handler := api.Logging(mux)
	*/
	handler := api.Logging(api.CORS(mux))

	addr := ":8080"

	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Println("listening on", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server exited gracefully")
	/*
		log.Println("listening on", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatal(err)
		}
	*/
}
