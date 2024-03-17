package middleware

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
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

	actual, ok := m.values[RequestIdKey]
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

func TestResponseEnvelopeMiddleware_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusOK)
	next := createHandlerFuncReturning(errDefault)

	em := ResponseEnvelope()
	callable := em(next)
	actual := callable(m)

	assert.Equal(errDefault, actual)
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

// func ResponseEnvelope() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			id := uuid.Must(uuid.NewRandom())
// 			log := logger.New(id.String())

// 			c.Set(RequestIdKey, id)
// 			c.SetLogger(log)

// 			w := &envelopeResponseWriter{
// 				response: responseEnvelope{
// 					RequestId: id,
// 				},
// 				writer: c.Response().Writer,
// 				logger: log,
// 			}
// 			c.Response().Writer = w

// 			return next(c)
// 		}
// 	}
// }
