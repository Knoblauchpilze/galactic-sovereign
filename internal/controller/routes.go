package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving"
	ports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/labstack/echo/v5"
)

func UniverseEndpoints(usecase ports.ForManagingUniverse) rest.Routes {
	var out rest.Routes

	handler := generateHandler(driving.CreateUniverse, usecase)
	post := rest.NewRoute(http.MethodPost, "/universes", handler)
	out = append(out, post)

	handler = generateHandler(driving.GetUniverse, usecase)
	get := rest.NewRoute(http.MethodGet, "/universes/:id", handler)
	out = append(out, get)

	handler = generateHandler(driving.ListUniverses, usecase)
	list := rest.NewRoute(http.MethodGet, "/universes", handler)
	out = append(out, list)

	handler = generateHandler(driving.DeleteUniverse, usecase)
	delete := rest.NewRoute(http.MethodDelete, "/universes/:id", handler)
	out = append(out, delete)

	return out
}

type drivingAdapter[T any] = func(*echo.Context, T) error

func generateHandler[T any](handler drivingAdapter[T], usecase T) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return handler(c, usecase)
	}
}
