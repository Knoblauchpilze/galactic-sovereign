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

func PlayerEndpoints(service service.PlayerService) rest.Routes {
	var out rest.Routes

	postHandler := fromPlayerServiceAwareHttpHandler(createPlayer, service)
	post := rest.NewRoute(http.MethodPost, false, "/players", postHandler)
	out = append(out, post)

	getHandler := fromPlayerServiceAwareHttpHandler(getPlayer, service)
	get := rest.NewResourceRoute(http.MethodGet, true, "/players", getHandler)
	out = append(out, get)

	listHandler := fromPlayerServiceAwareHttpHandler(listPlayers, service)
	list := rest.NewRoute(http.MethodGet, true, "/players", listHandler)
	out = append(out, list)

	deleteHandler := fromPlayerServiceAwareHttpHandler(deletePlayer, service)
	delete := rest.NewResourceRoute(http.MethodDelete, true, "/players", deleteHandler)
	out = append(out, delete)

	return out
}

func createPlayer(c echo.Context, s service.PlayerService) error {
	var playerDtoRequest communication.PlayerDtoRequest
	err := c.Bind(&playerDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid player syntax")
	}

	out, err := s.Create(c.Request().Context(), playerDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey) {
			return c.JSON(http.StatusConflict, "Name already used")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func getPlayer(c echo.Context, s service.PlayerService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := s.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such player")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func listPlayers(c echo.Context, s service.PlayerService) error {
	out, err := s.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func deletePlayer(c echo.Context, s service.PlayerService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such player")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
