package middleware

import (
	"net/http"

	"tech-ip-sem2/shared/metrics"
)

func Metrics(metrics *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return metrics.Middleware(next)
	}
}
