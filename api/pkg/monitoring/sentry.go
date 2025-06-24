package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
)

// InitSentry initializes Sentry error monitoring
func InitSentry(dsn, environment, release string) error {
	if dsn == "" {
		logger.Info("Sentry DSN not provided, skipping initialization")
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		Release:          release,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
		ProfilesSampleRate: 1.0,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Filter out sensitive information
			if event.Request != nil {
				// Remove sensitive headers
				if event.Request.Headers != nil {
					delete(event.Request.Headers, "Authorization")
					delete(event.Request.Headers, "Cookie")
					delete(event.Request.Headers, "X-Api-Key")
				}
			}
			return event
		},
	})

	if err != nil {
		return fmt.Errorf("failed to initialize Sentry: %w", err)
	}

	logger.Info("Sentry error monitoring initialized")
	return nil
}

// SentryMiddleware returns HTTP middleware for Sentry error capture
func SentryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a new hub for this request
			hub := sentry.GetHubFromContext(r.Context())
			if hub == nil {
				hub = sentry.CurrentHub().Clone()
			}

			// Set request context
			hub.Scope().SetRequest(r)
			hub.Scope().SetTag("component", "http-handler")
			hub.Scope().SetTag("method", r.Method)
			hub.Scope().SetTag("url", r.URL.Path)

			// Add user context if available (you can customize this)
			if userID := r.Header.Get("X-User-ID"); userID != "" {
				hub.Scope().SetUser(sentry.User{
					ID: userID,
				})
			}

			// Create transaction for performance monitoring
			transaction := sentry.StartTransaction(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
			transaction.SetTag("http.method", r.Method)
			transaction.SetTag("http.url", r.URL.String())
			
			// Add transaction to context
			ctx := transaction.Context()
			r = r.WithContext(sentry.SetHubOnContext(ctx, hub))

			defer transaction.Finish()

			// Wrap response writer to capture status code
			wrappedWriter := &responseWriter{ResponseWriter: w}

			// Recover from panics and send to Sentry
			defer func() {
				if err := recover(); err != nil {
					hub.RecoverWithContext(r.Context(), err)
					transaction.Status = sentry.SpanStatusInternalError
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(wrappedWriter, r)

			// Set transaction status based on response
			if wrappedWriter.statusCode >= 400 {
				if wrappedWriter.statusCode >= 500 {
					transaction.Status = sentry.SpanStatusInternalError
				} else {
					transaction.Status = sentry.SpanStatusInvalidArgument
				}
			} else {
				transaction.Status = sentry.SpanStatusOK
			}
		})
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

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = 200
	}
	return rw.ResponseWriter.Write(b)
}

// CaptureError captures an error and sends it to Sentry
func CaptureError(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
	} else {
		sentry.CaptureException(err)
	}
}

// CaptureMessage captures a message and sends it to Sentry
func CaptureMessage(ctx context.Context, message string, level sentry.Level) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureMessage(message)
	} else {
		sentry.CaptureMessage(message)
	}
}

// AddBreadcrumb adds a breadcrumb to the current scope
func AddBreadcrumb(ctx context.Context, breadcrumb *sentry.Breadcrumb) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.AddBreadcrumb(breadcrumb, nil)
	} else {
		sentry.AddBreadcrumb(breadcrumb)
	}
}

// SetTag sets a tag in the current scope
func SetTag(ctx context.Context, key, value string) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetTag(key, value)
	} else {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag(key, value)
		})
	}
}

// SetUser sets user information in the current scope
func SetUser(ctx context.Context, user sentry.User) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetUser(user)
	} else {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(user)
		})
	}
}

// StartTransaction starts a new Sentry transaction for performance monitoring
func StartTransaction(ctx context.Context, name string) *sentry.Span {
	return sentry.StartTransaction(ctx, name)
}

// Flush flushes the Sentry client buffer
func FlushSentry(timeout time.Duration) bool {
	return sentry.Flush(timeout)
}

// CloseSentry closes the Sentry client
func CloseSentry() {
	FlushSentry(5 * time.Second)
	logger.Info("Sentry monitoring shut down")
}