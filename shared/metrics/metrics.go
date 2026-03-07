package metrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	requestsTotal    *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	requestsInFlight prometheus.Gauge
}

func New(service string) *Metrics {
	requestsTotal := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "http_requests_total",
			Help:        "Total number of HTTP requests",
			ConstLabels: prometheus.Labels{"service": service},
		},
		[]string{"method", "route", "status"},
	)

	requestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "http_request_duration_seconds",
			Help:        "Duration of HTTP requests in seconds",
			Buckets:     []float64{0.01, 0.05, 0.1, 0.3, 1, 3, 5, 10},
			ConstLabels: prometheus.Labels{"service": service},
		},
		[]string{"method", "route"},
	)

	requestsInFlight := promauto.NewGauge(
		prometheus.GaugeOpts{
			Name:        "http_requests_in_flight",
			Help:        "Current number of in-flight HTTP requests",
			ConstLabels: prometheus.Labels{"service": service},
		},
	)

	return &Metrics{
		requestsTotal:    requestsTotal,
		requestDuration:  requestDuration,
		requestsInFlight: requestsInFlight,
	}
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.requestsInFlight.Inc()
		defer m.requestsInFlight.Dec()

		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		route := normalizeRoute(r.URL.Path)

		start := time.Now()

		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start).Seconds()

		m.requestsTotal.WithLabelValues(
			r.Method,
			route,
			strconv.Itoa(wrapper.statusCode),
		).Inc()

		m.requestDuration.WithLabelValues(
			r.Method,
			route,
		).Observe(duration)
	})
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func normalizeRoute(path string) string {
	if path == "/metrics" {
		return path
	}

	if strings.Contains(path, "/v1/auth/") {
		switch {
		case strings.Contains(path, "/v1/auth/login"):
			return "/v1/auth/login"
		case strings.Contains(path, "/v1/auth/verify"):
			return "/v1/auth/verify"
		}
	}

	if strings.Contains(path, "/v1/tasks/") {
		parts := strings.Split(path, "/")
		if len(parts) >= 4 && parts[3] != "tasks" {
			return "/v1/tasks/:id"
		}
	}

	if path == "/v1/tasks" {
		return "/v1/tasks"
	}

	return path
}
