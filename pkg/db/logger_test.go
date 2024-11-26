package db

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/logger"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	logger.Logger

	debugCalled int
	infoCalled  int
	warnCalled  int
	errorCalled int

	format string
	args   []interface{}
}

func TestUnit_ToTraceLogLevel(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		in       log.Lvl
		expected tracelog.LogLevel
	}

	testCases := map[string]testCase{
		"debug": {in: log.DEBUG, expected: tracelog.LogLevelDebug},
		"info":  {in: log.INFO, expected: tracelog.LogLevelInfo},
		"warn":  {in: log.WARN, expected: tracelog.LogLevelWarn},
		"error": {in: log.ERROR, expected: tracelog.LogLevelError},
		"off":   {in: log.OFF, expected: tracelog.LogLevelNone},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := toTracelogLevel(testCase.in)

			assert.Equal(testCase.expected, actual)
		})
	}
}

func TestUnit_FlattenMap(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		in       map[string]interface{}
		expected string
	}

	testCases := []testCase{
		{in: nil, expected: ""},
		{in: map[string]interface{}{"key": 1}, expected: "key: 1"},
		{in: map[string]interface{}{"key": []float32{1.2, -4.5}}, expected: "key: [1.2 -4.5]"},
		{in: map[string]interface{}{"key": "value", "key2": 36}, expected: "key: value key2: 36"},
		{in: map[string]interface{}{"key": -59.9, "key2": "haha", "key3": 72}, expected: "key: -59.9 key2: haha key3: 72"},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			actual := flattenMap(testCase.in)

			assert.Equal(testCase.expected, actual)
		})
	}
}

var date = time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC)

func TestUnit_PrepareSqlMessage(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		msg      string
		args     map[string]interface{}
		expected string
	}

	testCases := []testCase{
		{
			msg:      "dummy",
			args:     map[string]interface{}{},
			expected: "dummy",
		},
		{
			msg: "Query",
			args: map[string]interface{}{
				"sql": "select * from table",
			},
			expected: "Query select * from table",
		},
		{
			msg: "Query",
			args: map[string]interface{}{
				"args": []interface{}{"aa-ee", 36},
			},
			expected: "Query args=[aa-ee 36]",
		},
		{
			msg: "Prepare",
			args: map[string]interface{}{
				"time": date,
			},
			expected: "Prepare, time=2024-04-01 11:08:47.651387237 +0000 UTC",
		},
		{
			msg: "Prepare",
			args: map[string]interface{}{
				"sql":  "select * from table where id = $1",
				"args": []interface{}{27, "aa-ee"},
			},
			expected: "Prepare select * from table where id = $1 args=[27 aa-ee]",
		},
		{
			msg: "Query",
			args: map[string]interface{}{
				"sql":  "select * from table where id = $1",
				"time": date,
			},
			expected: "Query select * from table where id = $1, time=2024-04-01 11:08:47.651387237 +0000 UTC",
		},
		{
			msg: "Query",
			args: map[string]interface{}{
				"args": []interface{}{27, "aa-ee"},
				"time": date,
			},
			expected: "Query args=[27 aa-ee], time=2024-04-01 11:08:47.651387237 +0000 UTC",
		},
		{
			msg: "Query",
			args: map[string]interface{}{
				"sql":  "select * from table where id = $1 and name = $2",
				"args": []interface{}{27, "aa-ee"},
				"time": date,
			},
			expected: "Query select * from table where id = $1 and name = $2 args=[27 aa-ee], time=2024-04-01 11:08:47.651387237 +0000 UTC",
		},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			actual := prepareSqlMessage(testCase.msg, testCase.args)
			assert.Equal(testCase.expected, actual)
		})
	}
}

func TestUnit_PgxLoggerImpl_Log_Trace(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		log: m,
	}

	l.Log(context.Background(), tracelog.LogLevelTrace, "message", map[string]any{"toto": 1})

	assert.Equal(1, m.debugCalled)
	assert.Equal("message toto: 1", m.format)
	assert.Equal(0, len(m.args))
}

func TestUnit_PgxLoggerImpl_Log_Debug(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		log: m,
	}

	l.Log(context.Background(), tracelog.LogLevelDebug, "message", map[string]any{"toto": 1})

	assert.Equal(1, m.debugCalled)
	assert.Equal("message toto: 1", m.format)
	assert.Equal(0, len(m.args))
}

func TestUnit_PgxLoggerImpl_Log_Info(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		log: m,
	}

	l.Log(context.Background(), tracelog.LogLevelInfo, "message", map[string]any{"toto": 1})

	assert.Equal(1, m.infoCalled)
	assert.Equal("message toto: 1", m.format)
	assert.Equal(0, len(m.args))
}

func TestUnit_PgxLoggerImpl_Log_Warn(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		log: m,
	}

	l.Log(context.Background(), tracelog.LogLevelWarn, "message", map[string]any{"toto": 1})

	assert.Equal(1, m.warnCalled)
	assert.Equal("message toto: 1", m.format)
	assert.Equal(0, len(m.args))
}

func TestUnit_PgxLoggerImpl_Log_Error(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		log: m,
	}

	l.Log(context.Background(), tracelog.LogLevelError, "message", map[string]any{"toto": 1})

	assert.Equal(1, m.errorCalled)
	assert.Equal("message toto: 1", m.format)
	assert.Equal(0, len(m.args))
}

func TestUnit_PgxLoggerImpl_Log_None(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		log: m,
	}

	l.Log(context.Background(), tracelog.LogLevelNone, "message", map[string]any{"toto": 1})

	assert.Equal(0, m.debugCalled)
	assert.Equal(0, m.infoCalled)
	assert.Equal(0, m.warnCalled)
	assert.Equal(0, m.errorCalled)
}

func TestUnit_PgxLoggerImpl_WhenMessageIsUnknownAndSetToIgnore_ExpectNoLog(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		ignoreUnknownMessages: true,
		log:                   m,
	}

	l.Log(context.Background(), tracelog.LogLevelError, "message", map[string]any{"toto": 1})

	assert.Equal(0, m.debugCalled)
	assert.Equal(0, m.infoCalled)
	assert.Equal(0, m.warnCalled)
	assert.Equal(0, m.errorCalled)
}

func TestUnit_PgxLoggerImpl_WhenMessageIsKnownAndSetToIgnore_ExpectFormattedLog(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		ignoreUnknownMessages: true,
		log:                   m,
	}

	args := map[string]interface{}{
		"sql":  "select * from table where id = $1",
		"args": []interface{}{"aa-ee"},
	}
	l.Log(context.Background(), tracelog.LogLevelInfo, "Query", args)

	assert.Equal(1, m.infoCalled)
	assert.Equal("Query select * from table where id = $1 args=[aa-ee]", m.format)
	assert.Equal(0, len(m.args))
}

func TestUnit_PgxLoggerImpl_prepareMessage_ExpectFormattedLog(t *testing.T) {
	assert := assert.New(t)

	m := &mockLogger{}
	l := pgxLoggerImpl{
		ignoreUnknownMessages: true,
		log:                   m,
	}

	args := map[string]interface{}{
		"sql":  "select * from table where id = $1",
		"args": []interface{}{"aa-ee"},
	}
	l.Log(context.Background(), tracelog.LogLevelInfo, "Prepare", args)

	assert.Equal(1, m.infoCalled)
	assert.Equal("Prepare select * from table where id = $1 args=[aa-ee]", m.format)
	assert.Equal(0, len(m.args))
}

func (m *mockLogger) Output() io.Writer {
	return os.Stdout
}

func (m *mockLogger) Level() logger.Level {
	return logger.DEBUG
}

func (m *mockLogger) SetLevel(logger.Level) {}

func (m *mockLogger) Prefix() string {
	return ""
}

func (m *mockLogger) SetPrefix(string) {}

func (m *mockLogger) Header() string {
	return ""
}

func (m *mockLogger) SetHeader(string) {}

func (m *mockLogger) Debugf(format string, args ...interface{}) {
	m.debugCalled++
	m.format = format
	m.args = append(m.args, args...)
}

func (m *mockLogger) Infof(format string, args ...interface{}) {
	m.infoCalled++
	m.format = format
	m.args = append(m.args, args...)
}

func (m *mockLogger) Warnf(format string, args ...interface{}) {
	m.warnCalled++
	m.format = format
	m.args = append(m.args, args...)
}

func (m *mockLogger) Errorf(format string, args ...interface{}) {
	m.errorCalled++
	m.format = format
	m.args = append(m.args, args...)
}
