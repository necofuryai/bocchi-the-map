package monitoring

import (
	"context"
	"crypto/rand"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
)

// contextKey is a custom type for context keys to prevent collisions
type contextKey string

// requestIDKey is the typed key for request ID in context
const requestIDKey contextKey = "request_id"

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
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			
			// Add breadcrumb to Sentry
			AddBreadcrumb(ctx, &sentry.Breadcrumb{
				Message: "HTTP Request",
				Level:   sentry.LevelInfo,
				Data: map[string]interface{}{
					"request_id": requestID,
					"method":     r.Method,
					"url":        r.URL.String(),
				},
			})

			// Set tag in Sentry
			SetTag(ctx, "request_id", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
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
			
			// Record detailed error metrics
			if wrappedWriter.statusCode >= 400 && wrappedWriter.statusCode < 500 {
				RecordCustomMetric("Custom/HTTP/ClientErrorCount", 1)
			} else if wrappedWriter.statusCode >= 500 {
				RecordCustomMetric("Custom/HTTP/ServerErrorCount", 1)
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

// randomString generates a cryptographically secure random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	
	// Use crypto/rand for cryptographically secure random bytes
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		// If crypto/rand fails, this is a critical security issue
		// We should not continue with weak random generation
		panic("crypto/rand failure: cannot generate secure request ID")
	}
	
	for i := range b {
		b[i] = charset[randomBytes[i]%byte(len(charset))]
	}
	return string(b)
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// StartTrace starts a new trace span for monitoring
func StartTrace(ctx context.Context, name string) context.Context {
	// Add breadcrumb to Sentry for tracing
	AddBreadcrumb(ctx, &sentry.Breadcrumb{
		Message: "Trace Started: " + name,
		Level:   sentry.LevelInfo,
		Data: map[string]interface{}{
			"trace_name": name,
			"timestamp":  time.Now().Unix(),
		},
	})
	
	// Return context with trace information
	return context.WithValue(ctx, "trace_name", name)
}

// EndTrace ends a trace span
func EndTrace(ctx context.Context) {
	if traceName, ok := ctx.Value("trace_name").(string); ok {
		AddBreadcrumb(ctx, &sentry.Breadcrumb{
			Message: "Trace Ended: " + traceName,
			Level:   sentry.LevelInfo,
			Data: map[string]interface{}{
				"trace_name": traceName,
				"timestamp":  time.Now().Unix(),
			},
		})
	}
}

// AddUserContext adds user information to monitoring context
func AddUserContext(ctx context.Context, userID, email string) {
	// Set user context in Sentry
	SetUser(ctx, sentry.User{
		ID:    userID,
		Email: email,
	})
	
	// Set tags for better filtering
	SetTag(ctx, "user_id", userID)
	if email != "" {
		SetTag(ctx, "user_email", email)
	}
}

// RecordAuthFailure records authentication failure metrics
func RecordAuthFailure(ctx context.Context, errorType string) {
	// Record custom metric to New Relic
	RecordCustomMetric("Custom/Auth/FailureCount", 1)
	RecordCustomMetric("Custom/Auth/Failure/"+errorType, 1)
	
	// Add breadcrumb to Sentry
	AddBreadcrumb(ctx, &sentry.Breadcrumb{
		Message: "Authentication Failed",
		Level:   sentry.LevelWarning,
		Data: map[string]interface{}{
			"error_type": errorType,
			"timestamp":  time.Now().Unix(),
		},
	})
}

// RecordRateLimitExceeded records rate limit exceeded metrics
func RecordRateLimitExceeded(ctx context.Context, ip string) {
	// Record custom metric to New Relic
	RecordCustomMetric("Custom/RateLimit/ExceededCount", 1)
	
	// Add breadcrumb to Sentry
	AddBreadcrumb(ctx, &sentry.Breadcrumb{
		Message: "Rate Limit Exceeded",
		Level:   sentry.LevelWarning,
		Data: map[string]interface{}{
			"ip":        ip,
			"timestamp": time.Now().Unix(),
		},
	})
	
	// Set context for this incident
	SetTag(ctx, "rate_limit_ip", ip)
}