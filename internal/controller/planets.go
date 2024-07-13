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

func PlanetEndpoints(service service.PlanetService) rest.Routes {
	var out rest.Routes

	postHandler := fromPlanetServiceAwareHttpHandler(createPlanet, service)
	post := rest.NewRoute(http.MethodPost, false, "/planets", postHandler)
	out = append(out, post)

	getHandler := fromPlanetServiceAwareHttpHandler(getPlanet, service)
	get := rest.NewResourceRoute(http.MethodGet, true, "/planets", getHandler)
	out = append(out, get)

	listHandler := fromPlanetServiceAwareHttpHandler(listPlanets, service)
	list := rest.NewRoute(http.MethodGet, true, "/planets", listHandler)
	out = append(out, list)

	deleteHandler := fromPlanetServiceAwareHttpHandler(deletePlanet, service)
	delete := rest.NewResourceRoute(http.MethodDelete, true, "/planets", deleteHandler)
	out = append(out, delete)

	return out
}

func createPlanet(c echo.Context, s service.PlanetService) error {
	var planetDtoRequest communication.PlanetDtoRequest
	err := c.Bind(&planetDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid planet syntax")
	}

	out, err := s.Create(c.Request().Context(), planetDtoRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func getPlanet(c echo.Context, s service.PlanetService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := s.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such planet")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func listPlanets(c echo.Context, s service.PlanetService) error {
	out, err := s.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func deletePlanet(c echo.Context, s service.PlanetService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such planet")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
