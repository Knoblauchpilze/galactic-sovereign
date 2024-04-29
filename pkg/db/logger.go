package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/jackc/pgx/v5/tracelog"
)

type pgxLoggerImpl struct{}

func (l *pgxLoggerImpl) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	outMsg := fmt.Sprintf("pgx %s %s", msg, flattenMap(data))

	switch level {
	case tracelog.LogLevelTrace:
		logger.Tracef(outMsg)
	case tracelog.LogLevelDebug:
		logger.Debugf(outMsg)
	case tracelog.LogLevelInfo:
		logger.Infof(outMsg)
	case tracelog.LogLevelWarn:
		logger.Warnf(outMsg)
	case tracelog.LogLevelError:
		logger.Errorf(outMsg)
	case tracelog.LogLevelNone:
		logger.Tracef(outMsg)
	}
}

func flattenMap(data map[string]interface{}) string {
	var values []string
	for key, value := range data {
		values = append(values, fmt.Sprintf("%v: %v", key, value))
	}

	return strings.Join(values, " ")
}
