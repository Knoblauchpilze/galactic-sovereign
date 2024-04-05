package logger

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockEchoLogger struct {
	echo.Logger
}

func TestGetRequestLogger_ReturnsSetLogger(t *testing.T) {
	assert := assert.New(t)

	log := &mockEchoLogger{}
	ctx := context.WithValue(context.Background(), LogKey, log)

	actual := GetRequestLogger(ctx)
	assert.Equal(log, actual)
}

func TestGetRequestLogger_WhenLoggerButWithDifferentType_ReturnsLoggerWithNoRequestId(t *testing.T) {
	assert := assert.New(t)

	ctx := context.WithValue(context.Background(), LogKey, "not-a-logger")

	log := GetRequestLogger(ctx)
	assert.Equal("", log.Prefix())
}

func TestGetRequestLogger_WhenNoLogger_ReturnsLoggerWithNoRequestId(t *testing.T) {
	assert := assert.New(t)

	log := GetRequestLogger(context.Background())
	assert.Equal("", log.Prefix())
}
