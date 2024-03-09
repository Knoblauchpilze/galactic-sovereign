package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/labstack/echo"
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
				logger.Infof(createTimingLog(req, res, elapsed))
			} else {
				c.Error(err)
				logger.Warnf(createTimingLog(req, res, elapsed))
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

func pathFromRequest(req *http.Request) string {
	out := req.Host
	if req.URL.Path != "" {
		out += req.URL.Path
	}
	return out
}

func formatHttpStatusCode(status int) string {
	switch {
	case status >= 500:
		return logger.FormatWithColor(status, logger.Red)
	case status >= 400:
		return logger.FormatWithColor(status, logger.Yellow)
	case status >= 300:
		return logger.FormatWithColor(status, logger.Cyan)
	default:
		return logger.FormatWithColor(status, logger.Green)
	}
}
