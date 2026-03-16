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

	// Business metrics for notification service
	notificationsSentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notification_notifications_sent_total",
			Help: "Total number of notifications sent",
		},
		[]string{"type"},
	)

	notificationFailuresTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notification_failures_total",
			Help: "Total number of notification delivery failures",
		},
		[]string{"type"},
	)

	unreadNotificationsGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "notification_unread_notifications",
			Help: "Number of unread notifications",
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

// IncrementNotificationsSent increments the notifications sent counter
func IncrementNotificationsSent(notificationType string) {
	notificationsSentTotal.WithLabelValues(notificationType).Inc()
}

// IncrementNotificationFailures increments the notification failures counter
func IncrementNotificationFailures(notificationType string) {
	notificationFailuresTotal.WithLabelValues(notificationType).Inc()
}

// SetUnreadNotifications sets the unread notifications gauge
func SetUnreadNotifications(count float64) {
	unreadNotificationsGauge.Set(count)
}
