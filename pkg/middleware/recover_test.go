package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRecoverMiddleware_CallsNext(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := Recover()
	callable := em(next)
	callable(&m)

	assert.True(*called)
}

func TestRecoverMiddleware_WhenNoErrorReturnsNoError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next := createHandlerFuncReturning(nil)

	em := Recover()
	callable := em(next)
	actual := callable(&m)

	assert.Nil(actual)
}

func TestRecoverMiddleware_WhenNoErrorDoesNotCallContextError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next := createHandlerFuncReturning(nil)

	em := Recover()
	callable := em(next)
	callable(&m)

	assert.Nil(m.reportedError)
}

var errDefault = fmt.Errorf("some error")

func TestRecoverMiddleware_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	m := mockEchoContext{}
	next := createHandlerFuncReturning(errDefault)

	em := Recover()
	callable := em(next)
	actual := callable(&m)

	assert.Equal(actual, errDefault)
}

func createPanickingHandlerFunc(err error) echo.HandlerFunc {
	return func(c echo.Context) error {
		panic(err)
	}
}

func TestRecoverMiddleware_SetsContextErrorOnPanic(t *testing.T) {
	assert := assert.New(t)

	req := &http.Request{
		Method: "GET",
	}
	res := &echo.Response{
		Status: http.StatusConflict,
	}
	m := newMockEchoContext(req, res)
	next := createPanickingHandlerFunc(errDefault)

	em := Recover()
	callable := em(next)
	callable(m)

	assert.Equal(errDefault, m.reportedError)
}

// func Recover() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			defer func() {
// 				maybeErr := recover()
// 				if maybeErr == nil {
// 					return
// 				}

// 				err, ok := maybeErr.(error)
// 				if !ok {
// 					err = fmt.Errorf("%v", maybeErr)
// 				}

// 				req := c.Request()
// 				res := c.Response()

// 				stack := debug.Stack()

// 				c.Error(err)
// 				c.Logger().Errorf(createErrorLog(req, res, string(stack), err))
// 			}()

// 			return next(c)
// 		}
// 	}
// }
