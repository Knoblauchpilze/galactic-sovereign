package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestResponseEnvelopeMiddleware_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next, called := createHandlerFuncWithCalledBoolean()

	em := ResponseEnvelope()
	callable := em(next)
	callable(m)

	assert.True(*called)
}

func TestResponseEnvelopeMiddleware_SetsRequestId(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(nil)

	em := ResponseEnvelope()
	callable := em(next)
	callable(m)

	actual, ok := m.values[string(requestIdKey)]
	assert.True(ok)
	assert.IsType(uuid.UUID{}, actual)
}

func TestResponseEnvelopeMiddleware_AssignsNewLogger(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(nil)

	em := ResponseEnvelope()
	callable := em(next)
	callable(m)

	assert.True(m.loggerChanged)
}

func TestResponseEnvelopeMiddleware_AddLogKeyToRequestContext(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(nil)

	em := ResponseEnvelope()
	callable := em(next)
	callable(m)

	assert.True(m.requestChanged)
	ctx := m.request.Context()
	actual := ctx.Value(logKey)
	assert.NotNil(actual)
	_, ok := actual.(echo.Logger)
	assert.True(ok)
}

func TestResponseEnvelopeMiddleware_OverridesResponseWriter(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(nil)

	w := m.response.Writer

	em := ResponseEnvelope()
	callable := em(next)
	callable(m)

	actual := m.response.Writer
	assert.NotEqual(w, actual)

	assert.IsType(&envelopeResponseWriter{}, actual)
	actualW := actual.(*envelopeResponseWriter).writer
	assert.Equal(w, actualW)
}

func TestResponseEnvelopeMiddleware_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(errDefault)

	em := ResponseEnvelope()
	callable := em(next)
	actual := callable(m)

	assert.Equal(errDefault, actual)
}

func TestGetLoggerFromContext_ReturnsSetLogger(t *testing.T) {
	assert := assert.New(t)

	log := &mockEchoLogger{}
	ctx := context.WithValue(context.Background(), logKey, log)

	actual := GetLoggerFromContext(ctx)
	assert.Equal(log, actual)
}

func TestGetLoggerFromContext_WhenLoggerButWithDifferentType_UsesKey(t *testing.T) {
	assert := assert.New(t)

	ctx := context.WithValue(context.Background(), logKey, "not-a-logger")
	ctx = context.WithValue(ctx, requestIdKey, defaultUuid)

	log := GetLoggerFromContext(ctx)
	assert.Equal(defaultUuid.String(), log.Prefix())
}

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")

func TestGetLoggerFromContext_WhenNoLoggerButARequestKey_UsesIt(t *testing.T) {
	assert := assert.New(t)

	ctx := context.WithValue(context.Background(), requestIdKey, defaultUuid)

	log := GetLoggerFromContext(ctx)
	assert.Equal(defaultUuid.String(), log.Prefix())
}

func TestGetLoggerFromContext_WhenNoLoggerAndNoKey_ReturnsDefaultLogger(t *testing.T) {
	assert := assert.New(t)

	log := GetLoggerFromContext(context.Background())
	assert.Equal("", log.Prefix())
}
