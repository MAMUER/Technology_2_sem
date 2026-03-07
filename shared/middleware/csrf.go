package middleware

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"tech-ip-sem2/shared/logger"
)

type csrfResponse struct {
	Error string `json:"error"`
}

// CSRFMiddleware проверяет CSRF токен только для запросов с cookies
func CSRFMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := GetRequestID(r.Context())
			log := log.WithRequestID(requestID)

			// Опасные методы, требующие CSRF защиты
			dangerousMethods := map[string]bool{
				http.MethodPost:   true,
				http.MethodPut:    true,
				http.MethodPatch:  true,
				http.MethodDelete: true,
			}

			// Проверяем только опасные методы
			if dangerousMethods[r.Method] {
				// Проверяем, есть ли session cookie
				_, sessionErr := r.Cookie("session_id")
				_, csrfCookieErr := r.Cookie("csrf_token")

				// Если есть session cookie и csrf cookie, значит это браузерный запрос - нужна CSRF проверка
				if sessionErr == nil && csrfCookieErr == nil {
					// Получаем CSRF токен из cookie
					csrfCookie, _ := r.Cookie("csrf_token")

					// Получаем CSRF токен из заголовка
					csrfHeader := r.Header.Get("X-CSRF-Token")
					if csrfHeader == "" {
						log.Warn("CSRF header missing for cookie-based request", zap.String("method", r.Method))
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusForbidden)
						json.NewEncoder(w).Encode(csrfResponse{Error: "CSRF token missing in header"})
						return
					}

					// Сравниваем токены
					if csrfCookie.Value != csrfHeader {
						log.Warn("CSRF token mismatch",
							zap.String("cookie", csrfCookie.Value[:8]+"..."),
							zap.String("header", csrfHeader[:8]+"..."),
						)
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusForbidden)
						json.NewEncoder(w).Encode(csrfResponse{Error: "CSRF token invalid"})
						return
					}

					log.Debug("CSRF check passed for cookie-based request", zap.String("method", r.Method))
				} else {
					// Для API запросов без cookies пропускаем CSRF проверку
					log.Debug("Skipping CSRF check for API request (no cookies)", zap.String("method", r.Method))
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
