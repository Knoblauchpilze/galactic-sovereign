package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRecover_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := Recover()
	callable := em(next)
	callable(&m)

	assert.True(*called)
}

func TestRecover_WhenNoErrorReturnsNoError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next := createHandlerFuncReturning(nil)

	em := Recover()
	callable := em(next)
	actual := callable(&m)

	assert.Nil(actual)
}

func TestRecover_WhenNoErrorDoesNotCallContextError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next := createHandlerFuncReturning(nil)

	em := Recover()
	callable := em(next)
	callable(&m)

	assert.Nil(m.reportedError)
}

var errDefault = fmt.Errorf("some error")

func TestRecover_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next := createHandlerFuncReturning(errDefault)

	em := Recover()
	callable := em(next)
	actual := callable(&m)

	assert.Equal(actual, errDefault)
}

func createPanickingHandlerFunc(err interface{}) echo.HandlerFunc {
	return func(c echo.Context) error {
		panic(err)
	}
}

func TestRecover_SetsContextErrorOnPanic(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusConflict)
	next := createPanickingHandlerFunc(errDefault)

	em := Recover()
	callable := em(next)
	callable(m)

	assert.Equal(errDefault, m.reportedError)
}

func TestRecover_ConvertsToErrorUnknownPanic(t *testing.T) {
	assert := assert.New(t)
	m := newMockEchoContext(http.StatusConflict)
	next := createPanickingHandlerFunc(36)

	em := Recover()
	callable := em(next)
	callable(m)

	expected := fmt.Errorf("%v", 36)
	assert.Equal(expected, m.reportedError)
}
