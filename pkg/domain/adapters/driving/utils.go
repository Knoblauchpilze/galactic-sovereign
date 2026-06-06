package drivingadapters

import "github.com/labstack/echo/v5"

type drivingAdapter[T any] = func(*echo.Context, T) error

func generateHandler[T any](handler drivingAdapter[T], usecase T) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return handler(c, usecase)
	}
}
