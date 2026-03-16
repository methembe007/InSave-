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

	// Business metrics for budget service
	budgetsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "budget_budgets_created_total",
			Help: "Total number of budgets created",
		},
	)

	spendingTransactionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "budget_spending_transactions_total",
			Help: "Total number of spending transactions recorded",
		},
	)

	budgetAlertsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "budget_alerts_total",
			Help: "Total number of budget alerts generated",
		},
		[]string{"alert_type"},
	)

	activeBudgetsGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "budget_active_budgets",
			Help: "Number of currently active budgets",
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

// IncrementBudgetsCreated increments the budgets created counter
func IncrementBudgetsCreated() {
	budgetsCreatedTotal.Inc()
}

// IncrementSpendingTransactions increments the spending transactions counter
func IncrementSpendingTransactions() {
	spendingTransactionsTotal.Inc()
}

// IncrementBudgetAlerts increments the budget alerts counter
func IncrementBudgetAlerts(alertType string) {
	budgetAlertsTotal.WithLabelValues(alertType).Inc()
}

// SetActiveBudgets sets the active budgets gauge
func SetActiveBudgets(count float64) {
	activeBudgetsGauge.Set(count)
}
