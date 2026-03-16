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

	// Business metrics for user service
	profileUpdatesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_profile_updates_total",
			Help: "Total number of profile updates",
		},
	)

	accountDeletionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_account_deletions_total",
			Help: "Total number of account deletions",
		},
	)

	activeUsersGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "user_active_users",
			Help: "Number of currently active users",
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

// IncrementProfileUpdates increments the profile updates counter
func IncrementProfileUpdates() {
	profileUpdatesTotal.Inc()
}

// IncrementAccountDeletions increments the account deletions counter
func IncrementAccountDeletions() {
	accountDeletionsTotal.Inc()
}

// SetActiveUsers sets the active users gauge
func SetActiveUsers(count float64) {
	activeUsersGauge.Set(count)
}
