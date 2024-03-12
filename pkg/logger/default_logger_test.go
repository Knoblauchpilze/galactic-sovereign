package logger

import (
	"regexp"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var matcherStr = `{"level":"${level}","time":"[0-9A-Z:+-]+","message":"26 test"}`
var pattern = "${level}"

func TestDefaultLogger_Tracef(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetPrettyLogger)

	m := &mockIoWriter{}
	prettyLogger = zerolog.New(m)

	Tracef("%d test", 26)

	assert.Equal(1, m.called)
	assert.Equal(1, len(m.data))

	actual := string(m.data[0])
	matcher := regexp.MustCompile(strings.ReplaceAll(matcherStr, pattern, "trace"))
	assert.True(matcher.MatchString(actual))
}

func TestDefaultLogger_Debugf(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetPrettyLogger)

	m := &mockIoWriter{}
	prettyLogger = zerolog.New(m)

	Debugf("%d test", 26)

	assert.Equal(1, m.called)
	assert.Equal(1, len(m.data))

	actual := string(m.data[0])
	matcher := regexp.MustCompile(strings.ReplaceAll(matcherStr, pattern, "debug"))
	assert.True(matcher.MatchString(actual))
}

func TestDefaultLogger_Infof(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetPrettyLogger)

	m := &mockIoWriter{}
	prettyLogger = zerolog.New(m)

	Infof("%d test", 26)

	assert.Equal(1, m.called)
	assert.Equal(1, len(m.data))

	actual := string(m.data[0])
	matcher := regexp.MustCompile(strings.ReplaceAll(matcherStr, pattern, "info"))
	assert.True(matcher.MatchString(actual))
}

func TestDefaultLogger_Warnf(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetPrettyLogger)

	m := &mockIoWriter{}
	prettyLogger = zerolog.New(m)

	Warnf("%d test", 26)

	assert.Equal(1, m.called)
	assert.Equal(1, len(m.data))

	actual := string(m.data[0])
	matcher := regexp.MustCompile(strings.ReplaceAll(matcherStr, pattern, "warn"))
	assert.True(matcher.MatchString(actual))
}

func TestDefaultLogger_Errorf(t *testing.T) {
	assert := assert.New(t)
	t.Cleanup(resetPrettyLogger)

	m := &mockIoWriter{}
	prettyLogger = zerolog.New(m)

	Errorf("%d test", 26)

	assert.Equal(1, m.called)
	assert.Equal(1, len(m.data))

	actual := string(m.data[0])
	matcher := regexp.MustCompile(strings.ReplaceAll(matcherStr, pattern, "error"))
	assert.True(matcher.MatchString(actual))
}

func resetPrettyLogger() {
	prettyLogger = zerolog.New(nil).Output(newSafeConsoleWriter())
}
