package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func New(level, format string) (*slog.Logger, error) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}

	opts := &slog.HandlerOptions{Level: lvl}

	var h slog.Handler
	switch format {
	case "json":
		h = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		h = slog.NewTextHandler(os.Stdout, opts)
	default:
		return nil, fmt.Errorf("invalid log format %q (use text or json)", format)
	}

	return slog.New(h), nil
}
