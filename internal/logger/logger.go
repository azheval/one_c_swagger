package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

func New(level string, logPath string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "DEBUG":
		lvl = slog.LevelDebug
	case "INFO":
		lvl = slog.LevelInfo
	case "WARN":
		lvl = slog.LevelWarn
	case "ERROR":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err := os.MkdirAll(logPath, 0755); err != nil {
			slog.Error("Error creating log directory", "path", logPath, "error", err)
			return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
		}
	}

	logFile, err := os.OpenFile(filepath.Join(logPath, "one_c_swagger.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Error opening log file", "error", err)
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	}

	return slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: lvl}))
}
