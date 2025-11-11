package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/pz7-redis/internal/cache"
)

func main() {
	// Получаем адрес Redis из переменной окружения или используем по умолчанию
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379" // значение по умолчанию для Docker Compose
	}

	log.Printf("Connecting to Redis at: %s", redisAddr)
	c := cache.New(redisAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		if key == "" || value == "" {
			http.Error(w, "key and value required", http.StatusBadRequest)
			return
		}
		err := c.Set(key, value, 10*time.Second) // TTL = 10 сек
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "OK: %s=%s (TTL 10s)", key, value)
	})

	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key required", http.StatusBadRequest)
			return
		}
		val, err := c.Get(key)
		if err != nil {
			// Возвращаем 200 OK с сообщением, что ключ не найден
			fmt.Fprintf(w, "redis: nil")
			return
		}
		fmt.Fprintf(w, "VALUE: %s=%s", key, val)
	})

	mux.HandleFunc("/ttl", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key required", http.StatusBadRequest)
			return
		}
		ttl, err := c.TTL(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "TTL for %s: %v", key, ttl)
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}