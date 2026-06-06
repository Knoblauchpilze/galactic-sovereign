package drivingadapters

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type drivingAdapter[T any] = func(*echo.Context, T) error

func generateHandler[T any](handler drivingAdapter[T], usecase T) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return handler(c, usecase)
	}
}

func fetchIdFromQueryParam(key string, c *echo.Context) (exists bool, id uuid.UUID, err error) {
	maybeId := c.QueryParam(key)
	exists = (maybeId != "")
	if maybeId == "" {
		return exists, uuid.UUID{}, nil
	}

	id, err = uuid.Parse(maybeId)
	return exists, id, err
}
