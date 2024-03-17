package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorMiddleware_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ErrorMiddleware()
	callable := em(next)
	callable(&m)

	assert.True(*called)
}

func TestErrorMiddleware_WhenNoErrorReturnsNoError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next := createHandlerFuncReturning(nil)

	em := ErrorMiddleware()
	callable := em(next)
	actual := callable(&m)

	assert.Nil(actual)
}

func TestErrorMiddleware_ConvertsErrorWithCodeToHttpError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	err := errors.NewCode(36)
	next := createHandlerFuncReturning(err)

	em := ErrorMiddleware()
	callable := em(next)
	actual := callable(&m)

	httpErr, ok := actual.(*echo.HTTPError)
	assert.True(ok)
	assert.Equal(httpErr.Code, http.StatusInternalServerError)
	assert.Equal(httpErr.Message, err)
}

func TestErrorMiddleware_PropagatesUnknownError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	err := fmt.Errorf("some error")
	next := createHandlerFuncReturning(err)

	em := ErrorMiddleware()
	callable := em(next)
	actual := callable(&m)

	assert.Equal(actual, err)
}
