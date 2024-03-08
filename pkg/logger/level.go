package logger

import (
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

func fromLogLevel(level log.Lvl) zerolog.Level {
	switch level {
	case log.DEBUG:
		return zerolog.DebugLevel
	case log.INFO:
		return zerolog.InfoLevel
	case log.WARN:
		return zerolog.WarnLevel
	case log.ERROR:
		return zerolog.ErrorLevel
	default:
		return zerolog.NoLevel
	}
}

func fromZerologLevel(level zerolog.Level) log.Lvl {
	switch level {
	case zerolog.DebugLevel:
		return log.DEBUG
	case zerolog.InfoLevel:
		return log.INFO
	case zerolog.WarnLevel:
		return log.WARN
	case zerolog.ErrorLevel:
		return log.ERROR
	default:
		return log.OFF
	}
}
