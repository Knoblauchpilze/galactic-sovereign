package logger

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type KeyType string

const LogKey KeyType = "loggerKey"

func RegisterRequestLogger(requestId uuid.UUID, ctx context.Context) (echo.Logger, context.Context) {
	log := New(requestId.String())
	return log, context.WithValue(ctx, LogKey, log)
}

func GetRequestLogger(ctx context.Context) echo.Logger {
	log, ok := ctx.Value(LogKey).(echo.Logger)
	if ok && log != nil {
		return log
	}

	return New("")
}
