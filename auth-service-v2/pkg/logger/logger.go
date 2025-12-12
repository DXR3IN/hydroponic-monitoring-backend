package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func ErrorLogger(err error) {
	if err != nil {
		log.Error().Err(err).Msg("Error occurred")
	}
}

func InfoLogger(message string) {
	log.Info().Msg(message)
}

func DebugLogger(message string) {
	log.Debug().Msg(message)
}

func WarnLogger(message string) {
	log.Warn().Msg(message)
}
