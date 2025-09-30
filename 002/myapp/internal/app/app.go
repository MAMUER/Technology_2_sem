package app

import (
	/*"encoding/json"
	"time"*/
	"fmt"
	"net/http"

	"github.com/MAMUER/myapp/internal/app/handlers"
	"github.com/MAMUER/myapp/utils"
)

/*type pingResp struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}*/

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = utils.NewID16()
		}
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handlers.Ping)

	// Корневой маршрут
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Hello, Go project structure!")
	})

	// Пример JSON-ручки: /ping
	/*mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(r)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(pingResp{
			Status: "ok",
			Time:   time.Now().UTC().Format(time.RFC3339),
		})
	})*/
	handler := withRequestID(mux)

	/*utils.LogInfo("Server is starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		utils.LogError("server error: " + err.Error())
	}*/

	utils.LogInfo("Server is starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		utils.LogError("server error: " + err.Error())
	}
}
