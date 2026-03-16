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

	// Business metrics for savings service
	savingsTransactionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "savings_transactions_total",
			Help: "Total number of savings transactions created",
		},
	)

	savingsAmountTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "savings_amount_total",
			Help: "Total amount saved across all transactions",
		},
	)

	currentStreakGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "savings_current_streak",
			Help: "Current savings streak per user",
		},
		[]string{"user_id"},
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

// IncrementSavingsTransactions increments the savings transactions counter
func IncrementSavingsTransactions() {
	savingsTransactionsTotal.Inc()
}

// AddSavingsAmount adds to the total savings amount
func AddSavingsAmount(amount float64) {
	savingsAmountTotal.Add(amount)
}

// SetCurrentStreak sets the current streak for a user
func SetCurrentStreak(userID string, streak float64) {
	currentStreakGauge.WithLabelValues(userID).Set(streak)
}
