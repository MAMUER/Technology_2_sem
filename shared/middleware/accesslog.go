package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"tech-ip-sem2/shared/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func AccessLog(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := newResponseWriter(w)

			requestID := GetRequestID(r.Context())

			if requestID != "" {
				log.WithRequestID(requestID).Debug("request started",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
				)
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			loggerWithID := log.Logger
			if requestID != "" {
				loggerWithID = log.WithRequestID(requestID)
			}

			loggerWithID.Info("request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.status),
				zap.Float64("duration_ms", float64(duration.Milliseconds())),
				zap.String("remote_ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}
