package middleware

import "context"

type key int

const (
	claimsKey key = iota
)

// GetClaims извлекает claims из context
func GetClaims(ctx context.Context) (map[string]interface{}, bool) {
	claims, ok := ctx.Value(claimsKey).(map[string]interface{})
	return claims, ok
}

// WithClaims добавляет claims в context
func WithClaims(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}
