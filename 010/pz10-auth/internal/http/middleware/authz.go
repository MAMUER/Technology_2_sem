package middleware

import (
	"fmt"
	"net/http"
)

func AuthZRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool)
	for _, role := range allowedRoles {
		allowed[role] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok {
				fmt.Println("AuthZ: No claims in context")
				http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
				return
			}

			role, _ := claims["role"].(string)
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
