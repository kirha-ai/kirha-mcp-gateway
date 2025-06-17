package logger

import (
	"context"
	"log/slog"
)

// SilentHandler implements slog.Handler but discards all logs.
type SilentHandler struct{}

// Ensure interface compliance
var _ slog.Handler = (*SilentHandler)(nil)

// Enabled always returns false, disabling all logging at any level.
func (h *SilentHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

// Handle does nothing and returns nil.
func (h *SilentHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the same handler (no-op).
func (h *SilentHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the same handler (no-op).
func (h *SilentHandler) WithGroup(_ string) slog.Handler {
	return h
}
