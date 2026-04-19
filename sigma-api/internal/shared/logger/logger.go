package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	Log  *slog.Logger
	once sync.Once
)

// InitLogger initializes a structured logger.
func InitLogger() {
	once.Do(func() {
		// Use JSON handler for structured, zero-allocation-friendly logging
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		Log = slog.New(handler)
		slog.SetDefault(Log)
	})
}
