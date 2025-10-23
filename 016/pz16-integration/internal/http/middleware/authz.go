package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func AuthZRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool)
	for _, role := range allowedRoles {
		allowed[role] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claimsValue := r.Context().Value("claims")
			if claimsValue == nil {
				fmt.Println("AuthZ: No claims in context")
				http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// ПРИНИМАЕМ ЛЮБОЙ ТИП CLAIMS
			var role string
			
			// Пробуем разные типы claims
			switch claims := claimsValue.(type) {
			case map[string]interface{}:
				if r, ok := claims["role"].(string); ok {
					role = r
				}
			case jwt.MapClaims:
				if r, ok := claims["role"].(string); ok {
					role = r
				}
			default:
				fmt.Printf("AuthZ: Unknown claims type: %T\n", claimsValue)
				http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
				return
			}

			if role == "" {
				fmt.Println("AuthZ: No role found in claims")
				http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
				return
			}

			fmt.Printf("AuthZ: User role = %s, allowed roles = %v\n", role, allowedRoles)

			if !allowed[role] {
				http.Error(w, `{"error": "insufficient permissions"}`, http.StatusForbidden)
				return
			}

			fmt.Println("AuthZ: Authorization successful")
			next.ServeHTTP(w, r)
		})
	}
}