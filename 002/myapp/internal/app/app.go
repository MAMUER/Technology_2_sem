package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MAMUER/myapp/utils"
)

type pingResp struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func Run() {
	mux := http.NewServeMux()

	// Корневой маршрут
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Hello, Go project structure!")
	})

	// Пример JSON-ручки: /ping
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(pingResp{
			Status: "ok",
			Time:   time.Now().UTC().Format(time.RFC3339),
		})
	})

	utils.LogInfo("Server is starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		utils.LogError("server error: " + err.Error())
	}
}
utils/logger.go
package utils

import (
	"fmt"
	"net/http"
	"time"
)

func LogRequest(r *http.Request) {
	fmt.Printf("[%s] %s %s %s\n",
		time.Now().Format(time.RFC3339),
		r.RemoteAddr,
		r.Method,
		r.URL.Path,
	)
}

func LogInfo(msg string) {
	fmt.Printf("[INFO] %s %s\n", time.Now().Format(time.RFC3339), msg)
}

func LogError(msg string) {
	fmt.Printf("[ERROR] %s %s\n", time.Now().Format(time.RFC3339), msg)
}