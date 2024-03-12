package logger

import (
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var echoLogLevels = map[log.Lvl]zerolog.Level{
	log.DEBUG: zerolog.DebugLevel,
	log.INFO:  zerolog.InfoLevel,
	log.WARN:  zerolog.WarnLevel,
	log.ERROR: zerolog.ErrorLevel,
	log.OFF:   zerolog.NoLevel,
}

func TestLevel_fromEchoLogLevel(t *testing.T) {
	for in, expected := range echoLogLevels {
		t.Run("", func(t *testing.T) {
			actual := fromEchoLogLevel(in)
			assert.Equal(t, expected, actual)
		})
	}
}

var zerologLevels = map[zerolog.Level]log.Lvl{
	zerolog.DebugLevel: log.DEBUG,
	zerolog.InfoLevel:  log.INFO,
	zerolog.WarnLevel:  log.WARN,
	zerolog.ErrorLevel: log.ERROR,
	zerolog.FatalLevel: log.ERROR,
	zerolog.PanicLevel: log.ERROR,
	zerolog.NoLevel:    log.OFF,
	zerolog.Disabled:   log.OFF,
}

func TestLevel_fromZerologLevel(t *testing.T) {
	for in, expected := range zerologLevels {
		t.Run("", func(t *testing.T) {
			actual := fromZerologLevel(in)
			assert.Equal(t, expected, actual)
		})
	}
}
