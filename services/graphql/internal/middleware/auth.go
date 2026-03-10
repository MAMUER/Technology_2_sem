package middleware

import (
	"context"
	"net/http"
	"strings"

	"tech-ip-sem2/shared/logger"
)

type contextKey string

const (
	SubjectKey contextKey = "subject"
)

func AuthMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					ctx = context.WithValue(ctx, SubjectKey, "student")
					log.Debug("authenticated via token")
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// cookie
			cookie, err := r.Cookie("session_id")
			if err == nil && cookie.Value != "" {
				ctx = context.WithValue(ctx, SubjectKey, "student")
				log.Debug("authenticated via cookie")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// По умолчанию - анонимный пользователь
			ctx = context.WithValue(ctx, SubjectKey, "anonymous")
			log.Debug("no authentication, using anonymous")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSubject(ctx context.Context) string {
	if subject, ok := ctx.Value(SubjectKey).(string); ok {
		return subject
	}
	return "anonymous"
}
