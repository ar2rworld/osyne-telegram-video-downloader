package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(production bool) *Logger {
	var logger zerolog.Logger

	if production {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	return &Logger{logger}
}
