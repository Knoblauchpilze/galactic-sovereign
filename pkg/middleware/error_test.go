package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUnit_Error_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next, called := createHandlerFuncWithCalledBoolean()

	em := Error()
	callable := em(next)
	callable(ctx)

	assert.True(*called)
}

func TestUnit_Error_WhenNoErrorReturnsNoError(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(nil)

	em := Error()
	callable := em(next)
	actual := callable(ctx)

	assert.Nil(actual)
}

func TestUnit_Error_ConvertsErrorWithCodeToHttpError(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	err := errors.NewCode(36)
	next := createHandlerFuncReturning(err)

	em := Error()
	callable := em(next)
	actual := callable(ctx)

	httpErr, ok := actual.(*echo.HTTPError)
	assert.True(ok)
	assert.Equal(httpErr.Code, http.StatusInternalServerError)
	assert.Equal(httpErr.Message, err)
}

func TestUnit_Error_PropagatesUnknownError(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	err := fmt.Errorf("some error")
	next := createHandlerFuncReturning(err)

	em := Error()
	callable := em(next)
	actual := callable(ctx)

	assert.Equal(actual, err)
}
