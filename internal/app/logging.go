package app

import (
	"log/slog"
	"os"
)

// logger is the package-level structured logger used by the wiring layer
// for startup and runtime diagnostics. It defaults to a text handler
// writing to stderr; InitLogger can reconfigure it (e.g. to JSON).
var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

// InitLogger configures the package logger. When jsonOutput is true records
// are emitted as JSON; otherwise a human-readable text handler writes to
// stderr. The configured logger is also installed as the default for any
// code that uses the log/slog package-level functions.
func InitLogger(jsonOutput bool) {
	var h slog.Handler
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if jsonOutput {
		h = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		h = slog.NewTextHandler(os.Stderr, opts)
	}
	logger = slog.New(h)
	slog.SetDefault(logger)
}
