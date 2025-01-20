package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UniverseEndpoints(service service.UniverseService) rest.Routes {
	var out rest.Routes

	postHandler := fromUniverseServiceAwareHttpHandler(createUniverse, service)
	post := rest.NewRoute(http.MethodPost, "/universes", postHandler)
	out = append(out, post)

	getHandler := fromUniverseServiceAwareHttpHandler(getUniverse, service)
	get := rest.NewRoute(http.MethodGet, "/universes/:id", getHandler)
	out = append(out, get)

	listHandler := fromUniverseServiceAwareHttpHandler(listUniverses, service)
	list := rest.NewRoute(http.MethodGet, "/universes", listHandler)
	out = append(out, list)

	deleteHandler := fromUniverseServiceAwareHttpHandler(deleteUniverse, service)
	delete := rest.NewRoute(http.MethodDelete, "/universes/:id", deleteHandler)
	out = append(out, delete)

	return out
}

func createUniverse(c echo.Context, s service.UniverseService) error {
	var universeDtoRequest communication.UniverseDtoRequest
	err := c.Bind(&universeDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid universe syntax")
	}

	out, err := s.Create(c.Request().Context(), universeDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation) {
			return c.JSON(http.StatusConflict, "Name already used")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func getUniverse(c echo.Context, s service.UniverseService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := s.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such universe")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func listUniverses(c echo.Context, s service.UniverseService) error {
	universes, err := s.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	out, err := marshalNilToEmptySlice(universes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, out)
}

func deleteUniverse(c echo.Context, s service.UniverseService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such universe")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
