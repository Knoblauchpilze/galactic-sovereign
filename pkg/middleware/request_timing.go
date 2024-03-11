package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestTiming() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			err := next(c)
			elapsed := time.Since(start)

			if err == nil {
				c.Logger().Infof(createTimingLog(req, res, elapsed))
			} else {
				c.Error(err)
				c.Logger().Warnf(createTimingLog(req, res, elapsed))
			}

			return nil
		}
	}
}

func createTimingLog(req *http.Request, res *echo.Response, elapsed time.Duration) string {
	var out string

	out += fmt.Sprintf("%v", req.Method)
	out += fmt.Sprintf(" %v", pathFromRequest(req))
	out += fmt.Sprintf(" processed in %v", elapsed)
	out += fmt.Sprintf(" -> %s", formatHttpStatusCode(res.Status))

	return out
}
