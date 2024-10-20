package middleware

import (
	"regexp"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestResponseEnvelope_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next, called := createHandlerFuncWithCalledBoolean()

	em := ResponseEnvelope()
	callable := em(next)
	callable(ctx)

	assert.True(*called)
}

func TestResponseEnvelope_AssignsNewLogger(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(nil)

	l := ctx.Logger()

	em := ResponseEnvelope()
	callable := em(next)
	callable(ctx)

	actual := ctx.Logger()
	assert.NotEqual(l, actual)
}

func TestResponseEnvelope_SetsUuidPrefixForRequestLogger(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(nil)

	em := ResponseEnvelope()
	callable := em(next)
	callable(ctx)

	actual := ctx.Logger()
	// https://stackoverflow.com/questions/136505/searching-for-uuids-in-text-with-regex
	pattern := regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	assert.True(pattern.MatchString(actual.Prefix()))
}

func TestResponseEnvelope_AddLoggerToRequestContext(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(nil)

	em := ResponseEnvelope()
	callable := em(next)
	callable(ctx)

	reqCtx := ctx.Request().Context()
	actual := reqCtx.Value(logger.LogKey)
	assert.NotNil(actual)
	_, ok := actual.(echo.Logger)
	assert.True(ok)
}

func TestResponseEnvelope_OverridesResponseWriter(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(nil)

	w := ctx.Response().Writer

	em := ResponseEnvelope()
	callable := em(next)
	callable(ctx)

	actual := ctx.Response().Writer
	assert.NotEqual(w, actual)

	assert.IsType(&envelopeResponseWriter{}, actual)
	actualW := actual.(*envelopeResponseWriter).writer
	assert.Equal(w, actualW)
}

func TestResponseEnvelope_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(errDefault)

	em := ResponseEnvelope()
	callable := em(next)
	actual := callable(ctx)

	assert.Equal(errDefault, actual)
}
