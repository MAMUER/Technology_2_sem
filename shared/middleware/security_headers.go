package middleware

import (
	"net/http"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Защита от MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Защита от clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Включает XSS фильтр в браузерах
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Content Security Policy (базовая)
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self'; "+
				"img-src 'self' data:; "+
				"font-src 'self'; "+
				"connect-src 'self';")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// HSTS (включаем только если HTTPS)
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		next.ServeHTTP(w, r)
	})
}
