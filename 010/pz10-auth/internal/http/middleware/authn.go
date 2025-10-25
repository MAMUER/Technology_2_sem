package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"example.com/pz10-auth/internal/platform/jwt"
)

// AuthN middleware для аутентификации JWT токенов
func AuthN(v jwt.Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error": "authorization header required"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := v.Parse(tokenStr)
			if err != nil {
				fmt.Printf("JWT Validation Error: %v\n", err)
				http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
				return
			}

			fmt.Printf("JWT Validated: user %v, role %v\n", claims["sub"], claims["role"])

			// Используем общую функцию WithClaims
			ctx := WithClaims(r.Context(), claims)
			newReq := r.WithContext(ctx)
			next.ServeHTTP(w, newReq)
		})
	}
}
