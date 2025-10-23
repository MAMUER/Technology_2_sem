package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"example.com/pz10-auth/internal/platform/jwt"
)

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
			fmt.Printf("All claims: %+v\n", claims)
			
			// Добавляем claims в контекст
			ctx := context.WithValue(r.Context(), "claims", claims)
			
			// Создаем новый запрос с контекстом
			newReq := r.WithContext(ctx)
			
			// Передаем управление следующему обработчику
			next.ServeHTTP(w, newReq)
		})
	}
}