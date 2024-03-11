package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				maybeErr := recover()
				if maybeErr == nil {
					return
				}

				err, ok := maybeErr.(error)
				if !ok {
					err = fmt.Errorf("%v", maybeErr)
				}

				req := c.Request()
				res := c.Response()

				stack := debug.Stack()

				c.Error(err)
				c.Logger().Errorf(createErrorLog(req, res, string(stack), err))
			}()

			return next(c)

		}
	}
}

func createErrorLog(req *http.Request, res *echo.Response, stack string, err error) string {
	var out string

	out += fmt.Sprintf("%v", req.Method)
	out += fmt.Sprintf(" %v", pathFromRequest(req))
	out += fmt.Sprintf(" generated panic %v, stack: %v", err, stack)
	out += fmt.Sprintf(" -> %s", formatHttpStatusCode(res.Status))

	return out
}
