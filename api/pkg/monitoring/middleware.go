package monitoring

import (
	"net/http"
	"time"

	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
)

// InitMonitoring initializes both New Relic and Sentry monitoring
func InitMonitoring(newRelicKey, sentryDSN, appName, environment, release string) error {
	// Initialize New Relic
	if err := InitNewRelic(newRelicKey, appName, environment); err != nil {
		logger.Error("Failed to initialize New Relic", err)
		return err
	}

	// Initialize Sentry
	if err := InitSentry(sentryDSN, environment, release); err != nil {
		logger.Error("Failed to initialize Sentry", err)
		return err
	}

	logger.Info("Monitoring services initialized successfully")
	return nil
}

// MonitoringMiddleware returns a combined middleware for both New Relic and Sentry
func MonitoringMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Chain the middlewares: Sentry -> New Relic -> next handler
		return SentryMiddleware()(NewRelicMiddleware()(next))
	}
}

// RequestIDMiddleware adds a request ID to the context for better tracing
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate request ID if not present
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}

			// Set request ID in response header
			w.Header().Set("X-Request-ID", requestID)

			// Add request ID to context for logging
			ctx := r.Context()
			
			// Add breadcrumb to Sentry
			AddBreadcrumb(ctx, &struct {
				Message string            `json:"message"`
				Level   string            `json:"level"`
				Data    map[string]string `json:"data"`
			}{
				Message: "HTTP Request",
				Level:   "info",
				Data: map[string]string{
					"request_id": requestID,
					"method":     r.Method,
					"url":        r.URL.String(),
				},
			})

			// Set tag in Sentry
			SetTag(ctx, "request_id", requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// PerformanceMiddleware tracks request performance metrics
func PerformanceMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrappedWriter := &performanceResponseWriter{ResponseWriter: w}

			next.ServeHTTP(wrappedWriter, r)

			// Calculate duration
			duration := time.Since(start)

			// Record custom metrics to New Relic
			RecordCustomMetric("Custom/HTTP/ResponseTime", duration.Seconds())
			RecordCustomMetric("Custom/HTTP/RequestCount", 1)
			
			if wrappedWriter.statusCode >= 400 {
				RecordCustomMetric("Custom/HTTP/ErrorCount", 1)
			}

			// Log performance data
			logger.InfoWithFields("HTTP request completed", map[string]interface{}{
				"method":      r.Method,
				"url":         r.URL.String(),
				"status_code": wrappedWriter.statusCode,
				"duration_ms": duration.Milliseconds(),
				"duration":    duration.String(),
			})
		})
	}
}

// performanceResponseWriter wraps http.ResponseWriter to capture status code
type performanceResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *performanceResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *performanceResponseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = 200
	}
	return rw.ResponseWriter.Write(b)
}

// ShutdownMonitoring gracefully shuts down all monitoring services
func ShutdownMonitoring() {
	logger.Info("Shutting down monitoring services...")
	
	// Shutdown New Relic
	Shutdown()
	
	// Shutdown Sentry
	CloseSentry()
	
	logger.Info("All monitoring services shut down successfully")
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	// Simple implementation - in production, you might want to use UUID
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}