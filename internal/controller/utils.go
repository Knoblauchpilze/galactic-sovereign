package controller

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func marshalNilToEmptySlice[T any](in []T) ([]byte, error) {
	toMarshal := make([]T, 0)
	if in != nil {
		toMarshal = in
	}

	return json.Marshal(toMarshal)
}

func fetchIdFromQueryParam(key string, c echo.Context) (exists bool, id uuid.UUID, err error) {
	maybeId := c.QueryParam(key)
	exists = (maybeId != "")
	if maybeId == "" {
		return exists, uuid.UUID{}, nil
	}

	id, err = uuid.Parse(maybeId)
	return exists, id, err
}
