package logger

import (
	"os"
	"sync"

	"github.com/GlebKirsan/go-final-project/internal/config"
	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

var (
	logger Logger
	once   sync.Once
)

func Get() *Logger {
	once.Do(func() {
		zeroLogger := zerolog.New(os.Stderr).With().
			Timestamp().
			Logger().
			Output(zerolog.ConsoleWriter{Out: os.Stderr})
		cfg := config.Get()
		switch cfg.LogLevel {
		case "debug":
			zeroLogger.Level(zerolog.DebugLevel)
		case "info":
			zeroLogger.Level(zerolog.InfoLevel)
		case "warn":
			zeroLogger.Level(zerolog.WarnLevel)
		case "err":
			zeroLogger.Level(zerolog.ErrorLevel)
		case "fatal":
			zeroLogger.Level(zerolog.FatalLevel)
		default:
			zeroLogger.Level(zerolog.InfoLevel)
		}
		logger = Logger{Logger: &zeroLogger}
	})

	return &logger
}
