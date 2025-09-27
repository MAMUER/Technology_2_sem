package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Текстовый ответ
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
		w.Write([]byte("Hello, world!"))
	})

	// Пока временный JSON (без UUID — добавим на шаге 4)
	/*mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user{
			ID:   "temp",
			Name: "Gopher",
		})
	})*/

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user{
			ID:   uuid.NewString(), // теперь реальный UUID
			Name: "Gopher",
		})
	})

	// Новый эндпоинт /health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := healthResponse{
			Status: "ok",
			Time:   time.Now().Format(time.RFC3339),
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	//addr := ":8080"
	addr := ":" + port
	log.Printf("Starting on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}