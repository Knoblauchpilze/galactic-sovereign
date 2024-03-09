package logger

import (
	"github.com/rs/zerolog"
)

var prettyLogger = zerolog.New(nil).Output(newSafeConsoleWriter())

func Tracef(format string, v ...interface{}) {
	prettyLogger.Trace().Timestamp().Msgf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	prettyLogger.Debug().Timestamp().Msgf(format, v...)
}

func Infof(format string, v ...interface{}) {
	prettyLogger.Info().Timestamp().Msgf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	prettyLogger.Warn().Timestamp().Msgf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	prettyLogger.Error().Timestamp().Msgf(format, v...)
}
