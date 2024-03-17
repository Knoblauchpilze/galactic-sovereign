package logger

import (
	"regexp"
	"testing"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogger_UsesLoggerImpl(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := zerolog.Logger{}
	prettyLogger = m

	l := New("")
	_, ok := l.(*loggerImpl)
	assert.True(ok)
}

func TestLogger_DefaultLevelIsTrace(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(zerolog.TraceLevel, prettyLogger.GetLevel())
}

func TestLogger_Prefix(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := zerolog.Logger{}
	prettyLogger = m

	l := New("prefix")
	assert.Equal("prefix", l.Prefix())

	l.SetPrefix("anotherPrefix")
	assert.Equal("anotherPrefix", l.Prefix())
}

func TestLogger_Level(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := zerolog.Logger{}
	prettyLogger = m

	l := New("prefix")
	assert.Equal(log.DEBUG, l.Level())

	l.SetLevel(log.ERROR)
	assert.Equal(log.ERROR, l.Level())
}

func TestLogger_UsesConsoleWriter(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = &m

	l := New("prefix")
	actual := l.Output()

	scw, ok := actual.(*safeConsoleWriter)
	assert.True(ok)
	assert.Equal(&m, scw.writer)
}

func TestLogger_AllowsReplacingOutput(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = &m

	m2 := mockIoWriter{}

	l := New("prefix")
	l.SetOutput(&m2)

	actual := l.Output()
	assert.Equal(&m2, actual)
}

func TestLogger_Debugf(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")

	l.Debugf("%s", "hello")

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m DBG \[prefix\] hello\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Infof(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")

	l.Infof("%s", "hello")

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m \x1b\[32mINF\x1b\[0m \x1b\[1m\[prefix\] hello\x1b\[0m\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Warnf(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")

	l.Warnf("%s", "hello")

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m \x1b\[33mWRN\x1b\[0m \x1b\[1m\[prefix\] hello\x1b\[0m\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Errorf(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")

	l.Errorf("%s", "hello")

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m \x1b\[31mERR\x1b\[0m \x1b\[1m\[prefix\] hello\x1b\[0m\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Panicf(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")
	f := func() {
		l.Panicf("%s", "hello")
	}

	assert.PanicsWithValuef("[prefix] hello", f, "")

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m \x1b\[31mPNC\x1b\[0m \x1b\[1m\[prefix\] hello\x1b\[0m\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_WithHeader(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")
	l.SetHeader("header")

	l.Debugf("%s", "hello")

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m DBG \[header\] \[prefix\] hello\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Debugj(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")

	l.Debugj(map[string]interface{}{"key": "value"})

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m DBG \x1b\[36mid=\x1b\[0mprefix \x1b\[36mkey=\x1b\[0mvalue\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Infoj(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")

	l.Infoj(map[string]interface{}{"key": "value"})

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m \x1b\[32mINF\x1b\[0m \x1b\[36mid=\x1b\[0mprefix \x1b\[36mkey=\x1b\[0mvalue\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Warnj(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")

	l.Warnj(map[string]interface{}{"key": "value"})

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m \x1b\[33mWRN\x1b\[0m \x1b\[36mid=\x1b\[0mprefix \x1b\[36mkey=\x1b\[0mvalue\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Errorj(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")

	l.Errorj(map[string]interface{}{"key": "value"})

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m \x1b\[31mERR\x1b\[0m \x1b\[36mid=\x1b\[0mprefix \x1b\[36mkey=\x1b\[0mvalue\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_Panicj(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")
	f := func() {
		l.Panicj(map[string]interface{}{"key": "value"})
	}

	assert.PanicsWithValuef("", f, "")

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m \x1b\[31mPNC\x1b\[0m \x1b\[36mid=\x1b\[0mprefix \x1b\[36mkey=\x1b\[0mvalue\n`)
	assert.True(matcher.MatchString(actual))
}

func TestLogger_JsonWithHeader(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetDefaultLogger)

	m := mockIoWriter{}
	consoleWriter = zerolog.ConsoleWriter{Out: &m, TimeFormat: time.DateTime}

	l := New("prefix")
	l.SetHeader("header")

	l.Debugj(map[string]interface{}{"key": "value"})

	assert.Equal(1, m.called)
	actual := string(m.data[0])
	matcher := regexp.MustCompile(`\x1b\[90m[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\x1b\[0m DBG \x1b\[36mheader=\x1b\[0mheader \x1b\[36mid=\x1b\[0mprefix \x1b\[36mkey=\x1b\[0mvalue\n`)
	assert.True(matcher.MatchString(actual))
}

func resetDefaultLogger() {
	prettyLogger = zerolog.New(nil).Output(newSafeConsoleWriter())
}
