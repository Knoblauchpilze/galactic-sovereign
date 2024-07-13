package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UniverseEndpoints(service service.UniverseService) rest.Routes {
	var out rest.Routes

	postHandler := fromUniverseServiceAwareHttpHandler(createUniverse, service)
	post := rest.NewRoute(http.MethodPost, false, "/universes", postHandler)
	out = append(out, post)

	getHandler := fromUniverseServiceAwareHttpHandler(getUniverse, service)
	get := rest.NewResourceRoute(http.MethodGet, false, "/universes", getHandler)
	out = append(out, get)

	listHandler := fromUniverseServiceAwareHttpHandler(listUniverses, service)
	list := rest.NewRoute(http.MethodGet, false, "/universes", listHandler)
	out = append(out, list)

	deleteHandler := fromUniverseServiceAwareHttpHandler(deleteUniverse, service)
	delete := rest.NewResourceRoute(http.MethodDelete, false, "/universes", deleteHandler)
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
		if errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey) {
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
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
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
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such universe")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
