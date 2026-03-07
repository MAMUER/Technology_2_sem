package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"tech-ip-sem2/shared/logger"
)

func DebugRequestID(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerID := r.Header.Get("X-Request-ID")
			ctxID := GetRequestID(r.Context())

			log.WithRequestID(ctxID).Debug("DEBUG: Request ID check",
				zap.String("header_x_request_id", headerID),
				zap.String("context_request_id", ctxID),
				zap.String("url", r.URL.String()),
			)

			next.ServeHTTP(w, r)
		})
	}
}
