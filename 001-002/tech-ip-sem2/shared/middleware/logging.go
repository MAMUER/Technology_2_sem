package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := GetRequestID(r.Context())
		if requestID == "" {
			requestID = "unknown"
		}

		log.Printf("[%s] %s %s started", requestID, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s %s completed in %v", requestID, r.Method, r.URL.Path, time.Since(start))
	})
}
