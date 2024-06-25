package logger

import (
	"github.com/rs/zerolog"
	"os"
)

type ZeroLogger struct {
	zerologger *zerolog.Logger
}

var _ Logger = (*ZeroLogger)(nil)

func NewZeroLogger() Logger {
	logger := zerolog.New(os.Stdout)

	return &ZeroLogger{
		zerologger: &logger,
	}
}

func (l *ZeroLogger) Info(msg string) {
	l.zerologger.Info().Msg(msg)
}

func (l *ZeroLogger) Error(msg string) {
	l.zerologger.Error().Msg(msg)
}

func (l *ZeroLogger) Fatal(msg string) {
	l.zerologger.Fatal().Msg(msg)
}
