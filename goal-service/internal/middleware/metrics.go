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

	// Business metrics for goal service
	goalsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "goal_goals_created_total",
			Help: "Total number of goals created",
		},
	)

	goalsCompletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "goal_goals_completed_total",
			Help: "Total number of goals completed",
		},
	)

	contributionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "goal_contributions_total",
			Help: "Total number of contributions made to goals",
		},
	)

	activeGoalsGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goal_active_goals",
			Help: "Number of currently active goals",
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

// IncrementGoalsCreated increments the goals created counter
func IncrementGoalsCreated() {
	goalsCreatedTotal.Inc()
}

// IncrementGoalsCompleted increments the goals completed counter
func IncrementGoalsCompleted() {
	goalsCompletedTotal.Inc()
}

// IncrementContributions increments the contributions counter
func IncrementContributions() {
	contributionsTotal.Inc()
}

// SetActiveGoals sets the active goals gauge
func SetActiveGoals(count float64) {
	activeGoalsGauge.Set(count)
}
