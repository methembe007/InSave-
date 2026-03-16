package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics for auth service
	registrationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	loginsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_logins_total",
			Help: "Total number of successful logins",
		},
	)

	loginFailuresTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_login_failures_total",
			Help: "Total number of failed login attempts",
		},
	)

	activeTokensGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "auth_active_tokens",
			Help: "Number of currently active tokens",
		},
	)

	rateLimitHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
	)
)

// MetricsMiddleware wraps HTTP handlers to collect metrics
func MetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.statusCode)

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// IncrementRegistrations increments the registration counter
func IncrementRegistrations() {
	registrationsTotal.Inc()
}

// IncrementLogins increments the successful login counter
func IncrementLogins() {
	loginsTotal.Inc()
}

// IncrementLoginFailures increments the failed login counter
func IncrementLoginFailures() {
	loginFailuresTotal.Inc()
}

// SetActiveTokens sets the active tokens gauge
func SetActiveTokens(count float64) {
	activeTokensGauge.Set(count)
}

// IncrementRateLimitHits increments the rate limit hits counter
func IncrementRateLimitHits() {
	rateLimitHitsTotal.Inc()
}
