package monitoring

import (
	"context"
	"fmt"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
)

var (
	nrApp *newrelic.Application
)

// InitNewRelic initializes New Relic monitoring
func InitNewRelic(licenseKey, appName, environment string) error {
	if licenseKey == "" {
		logger.Info("New Relic license key not provided, skipping initialization")
		return nil
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(config *newrelic.Config) {
			config.Labels = map[string]string{
				"environment": environment,
				"service":     "bocchi-api",
			}
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create New Relic application: %w", err)
	}

	nrApp = app
	logger.Info("New Relic monitoring initialized")
	return nil
}

// GetNewRelicApp returns the New Relic application instance
func GetNewRelicApp() *newrelic.Application {
	return nrApp
}

// NewRelicMiddleware returns HTTP middleware for New Relic monitoring
func NewRelicMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if nrApp == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Create New Relic transaction
			txn := nrApp.StartTransaction(r.URL.Path)
			defer txn.End()

			// Add request attributes
			txn.AddAttribute("http.method", r.Method)
			txn.AddAttribute("http.url", r.URL.String())
			txn.AddAttribute("http.user_agent", r.UserAgent())

			// Set web request and response writer
			txn.SetWebRequestHTTP(r)
			w = txn.SetWebResponse(w)

			// Add transaction to request context
			ctx := newrelic.NewContext(r.Context(), txn)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// RecordCustomEvent records a custom event to New Relic
func RecordCustomEvent(eventType string, params map[string]interface{}) {
	if nrApp == nil {
		return
	}

	nrApp.RecordCustomEvent(eventType, params)
}

// RecordCustomMetric records a custom metric to New Relic
func RecordCustomMetric(name string, value float64) {
	if nrApp == nil {
		return
	}

	nrApp.RecordCustomMetric(name, value)
}

// NoticeError reports an error to New Relic
func NoticeError(ctx context.Context, err error) {
	if nrApp == nil {
		return
	}

	// Try to get transaction from context
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.NoticeError(err)
	} else {
		// If no transaction context, create a background transaction for orphaned errors
		txn := nrApp.StartTransaction("error-without-transaction-context")
		defer txn.End()
		txn.NoticeError(err)
	}
}

// StartBackgroundTransaction starts a background transaction for non-web operations
func StartBackgroundTransaction(name string) *newrelic.Transaction {
	if nrApp == nil {
		return nil
	}

	return nrApp.StartTransaction(name)
}

// Shutdown gracefully shuts down New Relic monitoring
func Shutdown() {
	if nrApp != nil {
		nrApp.Shutdown(10) // 10 second timeout
		logger.Info("New Relic monitoring shut down")
	}
}