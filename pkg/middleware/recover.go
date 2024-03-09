package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/labstack/echo"
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
				logger.Errorf(createErrorLog(req, res, string(stack), err))

				// stack := make([]byte, config.StackSize)
				// length := runtime.Stack(stack, !config.DisableStackAll)
				// if !config.DisablePrintStack {
				// 	c.Logger().Printf("[PANIC RECOVER] %v %s\n", err, stack[:length])
				// }
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
