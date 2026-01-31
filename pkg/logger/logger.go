package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	SystemLogger zerolog.Logger
	AuthLogger   zerolog.Logger
	AuditLogger  zerolog.Logger
)

type Config struct {
	LogDir      string
	Environment string // "development" or "production"
}

func Init(cfg Config) {
	// Ensure log directory exists
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		panic(err)
	}

	SystemLogger = newLogger(cfg, "system.log")
	AuthLogger = newLogger(cfg, "auth.log")
	AuditLogger = newLogger(cfg, "audit.log")
}

func newLogger(cfg Config, filename string) zerolog.Logger {
	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.LogDir, filename),
		MaxSize:    10,   // megabytes
		MaxBackups: 3,    // number of backups
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}

	var writers []io.Writer
	writers = append(writers, fileLogger)

	// If development, also log to console (colored)
	if cfg.Environment == "development" {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	}

	multi := io.MultiWriter(writers...)

	return zerolog.New(multi).With().Timestamp().Logger()
}
