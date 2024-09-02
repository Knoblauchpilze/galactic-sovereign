package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/game"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func PlanetEndpoints(service service.PlanetService, actions game.ActionService) rest.Routes {
	var out rest.Routes

	postHandler := fromPlanetServiceAwareHttpHandler(createPlanet, service, actions)
	post := rest.NewRoute(http.MethodPost, false, "/planets", postHandler)
	out = append(out, post)

	getHandler := fromPlanetServiceAwareHttpHandler(getPlanet, service, actions)
	get := rest.NewResourceRoute(http.MethodGet, false, "/planets", getHandler)
	out = append(out, get)

	listHandler := fromPlanetServiceAwareHttpHandler(listPlanets, service, actions)
	list := rest.NewRoute(http.MethodGet, false, "/planets", listHandler)
	out = append(out, list)

	deleteHandler := fromPlanetServiceAwareHttpHandler(deletePlanet, service, actions)
	delete := rest.NewResourceRoute(http.MethodDelete, false, "/planets", deleteHandler)
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
	exists, playerId, err := fetchIdFromQueryParam("player", c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	var planets []communication.PlanetDtoResponse

	if exists {
		planets, err = s.ListForPlayer(c.Request().Context(), playerId)
	} else {
		planets, err = s.List(c.Request().Context())
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	out, err := marshalNilToEmptySlice(planets)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, out)
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
