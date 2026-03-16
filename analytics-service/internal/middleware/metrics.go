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

	// Business metrics for analytics service
	analysisRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_analysis_requests_total",
			Help: "Total number of analysis requests",
		},
		[]string{"analysis_type"},
	)

	recommendationsGeneratedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_recommendations_generated_total",
			Help: "Total number of recommendations generated",
		},
	)

	financialHealthScoreGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "analytics_financial_health_score",
			Help: "Financial health score per user",
		},
		[]string{"user_id"},
	)

	cacheHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	cacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_cache_misses_total",
			Help: "Total number of cache misses",
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

// IncrementAnalysisRequests increments the analysis requests counter
func IncrementAnalysisRequests(analysisType string) {
	analysisRequestsTotal.WithLabelValues(analysisType).Inc()
}

// IncrementRecommendationsGenerated increments the recommendations generated counter
func IncrementRecommendationsGenerated() {
	recommendationsGeneratedTotal.Inc()
}

// SetFinancialHealthScore sets the financial health score for a user
func SetFinancialHealthScore(userID string, score float64) {
	financialHealthScoreGauge.WithLabelValues(userID).Set(score)
}

// IncrementCacheHits increments the cache hits counter
func IncrementCacheHits() {
	cacheHitsTotal.Inc()
}

// IncrementCacheMisses increments the cache misses counter
func IncrementCacheMisses() {
	cacheMissesTotal.Inc()
}
