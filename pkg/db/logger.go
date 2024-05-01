package db

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type pgxLoggerImpl struct {
	ignoreUnknownMessages bool
	logger                echo.Logger
}

const pgxPrefixString = "pgx"
const pgxPrepareMessage = "Prepare"
const pgxQueryMessage = "Query"

func new(ignoreUnknownMessages bool) tracelog.Logger {
	return &pgxLoggerImpl{
		ignoreUnknownMessages: ignoreUnknownMessages,
		logger:                logger.New(pgxPrefixString),
	}
}

func (l *pgxLoggerImpl) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	var outMsg string
	knownMessage := true

	switch msg {
	case pgxPrepareMessage:
		outMsg = prepareSqlMessage(msg, data)
	case pgxQueryMessage:
		outMsg = prepareSqlMessage(msg, data)
	default:
		knownMessage = false
	}

	if !knownMessage && !l.ignoreUnknownMessages {
		outMsg = fmt.Sprintf("%s %s", msg, flattenMap(data))
	}

	if outMsg == "" {
		return
	}

	switch level {
	case tracelog.LogLevelTrace:
		l.logger.Debugf(outMsg)
	case tracelog.LogLevelDebug:
		l.logger.Debugf(outMsg)
	case tracelog.LogLevelInfo:
		l.logger.Infof(outMsg)
	case tracelog.LogLevelWarn:
		l.logger.Warnf(outMsg)
	case tracelog.LogLevelError:
		l.logger.Errorf(outMsg)
	}
}

func toTracelogLevel(level log.Lvl) tracelog.LogLevel {
	switch level {
	case log.DEBUG:
		return tracelog.LogLevelDebug
	case log.INFO:
		return tracelog.LogLevelInfo
	case log.WARN:
		return tracelog.LogLevelWarn
	case log.ERROR:
		return tracelog.LogLevelError
	default:
		return tracelog.LogLevelNone
	}
}

func flattenMap(data map[string]interface{}) string {
	// Order of maps is not deterministic
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var values []string
	for _, key := range keys {
		value := data[key]
		values = append(values, fmt.Sprintf("%v: %v", key, value))
	}

	return strings.Join(values, " ")
}

func prepareSqlMessage(message string, data map[string]interface{}) string {
	out := message

	sql, ok := data["sql"]
	if ok {
		out += fmt.Sprintf(" %v", sql)
	}

	args, ok := data["args"]
	if ok {
		out += fmt.Sprintf(" args=%v", args)
	}

	duration, ok := data["time"]
	if ok {
		out += fmt.Sprintf(", time=%v", duration)
	}

	return out
}
