package middleware

import (
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestResponseEnvelope_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next, called := createHandlerFuncWithCalledBoolean()

	em := ResponseEnvelope()
	callable := em(next)
	callable(m)

	assert.True(*called)
}

func TestResponseEnvelope_AssignsNewLogger(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(nil)

	em := ResponseEnvelope()
	callable := em(next)
	callable(m)

	assert.True(m.loggerChanged)
}

func TestResponseEnvelope_AddLoggerToRequestContext(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(nil)

	em := ResponseEnvelope()
	callable := em(next)
	callable(m)

	assert.True(m.requestChanged)
	ctx := m.request.Context()
	actual := ctx.Value(logger.LogKey)
	assert.NotNil(actual)
	_, ok := actual.(echo.Logger)
	assert.True(ok)
}

func TestResponseEnvelope_OverridesResponseWriter(t *testing.T) {
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

func TestResponseEnvelope_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(errDefault)

	em := ResponseEnvelope()
	callable := em(next)
	actual := callable(m)

	assert.Equal(errDefault, actual)
}
