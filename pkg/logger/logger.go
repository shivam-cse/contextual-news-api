package logger

import (
    "log/slog"
    "os"
)


// In a real production app, you might check an environment variable
// to decide whether to use NewTextHandler (for development) or
// NewJSONHandler (for production).
func New() *slog.Logger {
    opts := &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }
    return slog.New(slog.NewTextHandler(os.Stdout, opts))
}