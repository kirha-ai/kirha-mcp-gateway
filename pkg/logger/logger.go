// Package logger provides a structured logging implementation using slog.
// It formats logs for Google Cloud Platform (GCP) compatibility and supports
// configurable log levels through environment variables.
package logger

import (
	"log/slog"
	"os"
	"sync"
	"time"
)

// lock is a mutex to ensure thread-safe access to the global logger.
var lock = &sync.Mutex{}

// globalLogger is the singleton instance of the logger.
var globalLogger *slog.Logger

// New creates a new logger instance with the specified name.
// It returns a singleton logger configured for GCP-compatible JSON output.
// If logging is disabled via configuration, it returns a silent logger.
// The logger is created once and reused across the application with different names.
//
// Parameters:
//   - name: The name to identify this logger instance (e.g., "kirha_client", "mcp_server")
//
// Returns:
//   - *slog.Logger: A structured logger instance with the specified name
//
// Example:
//
//	logger := logger.New("my_component")
//	logger.Info("Starting component")
//	logger.Error("Operation failed", slog.String("error", err.Error()))
func New(name string) *slog.Logger {
	lock.Lock()
	defer lock.Unlock()

	if globalLogger != nil {
		return globalLogger.With(slog.String("logger", name))
	}

	cfg := LoadConfig()
	if !cfg.EnableLogs {
		return slog.New(&SilentHandler{})
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, gcpHandlerOpts))
	logger = logger.With(slog.String("application", "kirha-mcp-gateway"))
	globalLogger = logger
	logger = logger.With(slog.String("logger", name))
	return logger
}

// gcpHandlerOpts defines the options for the slog handler to format log attributes.
var gcpHandlerOpts = &slog.HandlerOptions{
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.LevelKey:
			return slog.Attr{
				Key:   "severity",
				Value: slog.StringValue(levelToGCP(a.Value.String())),
			}

		case slog.MessageKey:
			return slog.String("message", a.Value.String())

		case slog.TimeKey:
			return slog.String("time", a.Value.Time().Format(time.RFC3339Nano))

		default:
			return a
		}
	},
}

// levelToGCP converts slog levels to GCP logging levels.
func levelToGCP(level string) string {
	switch level {
	case "DEBUG":
		return "DEBUG"
	case "INFO":
		return "INFO"
	case "WARN":
		return "WARNING"
	case "ERROR":
		return "ERROR"
	default:
		return level
	}
}
