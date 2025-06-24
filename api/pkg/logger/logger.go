package logger

import (
	"context"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Level represents the log level
type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
)

// Init initializes the logger with JSON format
func Init(level Level) {
	// Set time format to RFC3339
	zerolog.TimeFieldFormat = time.RFC3339

	// Configure zerolog
	zerolog.SetGlobalLevel(parseLevel(level))

	// Set up console writer for development
	if os.Getenv("ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		// Production: JSON format
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

// parseLevel converts our Level type to zerolog.Level
func parseLevel(level Level) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

// WithContext adds context fields to the logger
func WithContext(fields map[string]interface{}) zerolog.Logger {
	ctx := log.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return ctx.Logger()
}

// Debug logs a debug message
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Info logs an info message
func Info(msg string) {
	log.Info().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	log.Warn().Msg(msg)
}

// Error logs an error message and sends to Sentry
func Error(msg string, err error) {
	log.Error().Err(err).Msg(msg)
	
	// Also send to Sentry if available
	if err != nil {
		sentry.CaptureException(err)
	}
}

// Fatal logs a fatal message, sends to Sentry, and exits
func Fatal(msg string, err error) {
	log.Fatal().Err(err).Msg(msg)
	
	// Send to Sentry before exiting
	if err != nil {
		sentry.CaptureException(err)
		sentry.Flush(2 * time.Second) // Wait for Sentry to send the error
	}
}

// InfoWithFields logs an info message with structured fields
func InfoWithFields(msg string, fields map[string]interface{}) {
	event := log.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// ErrorWithFields logs an error message with structured fields and sends to Sentry
func ErrorWithFields(msg string, err error, fields map[string]interface{}) {
	event := log.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
	
	// Send to Sentry with additional context
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope) {
			for k, v := range fields {
				scope.SetExtra(k, v)
			}
			scope.SetTag("component", "logger")
			sentry.CaptureException(err)
		})
	}
}

// ErrorWithContext logs an error message with context and sends to Sentry
func ErrorWithContext(ctx context.Context, msg string, err error) {
	log.Error().Err(err).Msg(msg)
	
	// Send to Sentry with context
	if err != nil {
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.CaptureException(err)
		} else {
			sentry.CaptureException(err)
		}
	}
}

// ErrorWithContextAndFields logs an error with context and fields, sends to Sentry
func ErrorWithContextAndFields(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	event := log.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
	
	// Send to Sentry with context and fields
	if err != nil {
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				for k, v := range fields {
					scope.SetExtra(k, v)
				}
				scope.SetTag("component", "logger")
				hub.CaptureException(err)
			})
		} else {
			sentry.WithScope(func(scope *sentry.Scope) {
				for k, v := range fields {
					scope.SetExtra(k, v)
				}
				scope.SetTag("component", "logger")
				sentry.CaptureException(err)
			})
		}
	}
}