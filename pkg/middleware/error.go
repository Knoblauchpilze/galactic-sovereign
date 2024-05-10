package middleware

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/labstack/echo/v4"
)

func Error() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			if impl, ok := err.(errors.ErrorWithCode); ok {
				httpCode := errorCodeToHttpErrorCode(impl.Code())
				return echo.NewHTTPError(httpCode, impl)
			}

			return err
		}
	}
}

func errorCodeToHttpErrorCode(code errors.ErrorCode) int {
	switch code {
	default:
		return http.StatusInternalServerError
	}
}
