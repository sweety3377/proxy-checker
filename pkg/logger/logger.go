package logger

import (
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"os"
)

// New logger
func New() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorStackFieldName = "stacktrace"
	zerolog.InterfaceMarshalFunc = json.Marshal

	logger := zerolog.New(os.Stderr).
		With().Timestamp().
		Logger()

	logger.Info().Msg("logs is enabled")

	return &logger
}
