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
	"github.com/labstack/echo/v5"
)

func PlayerEndpoints(service service.PlayerService) rest.Routes {
	var out rest.Routes

	postHandler := fromPlayerServiceAwareHttpHandler(createPlayer, service)
	post := rest.NewRoute(http.MethodPost, "/players", postHandler)
	out = append(out, post)

	getHandler := fromPlayerServiceAwareHttpHandler(getPlayer, service)
	get := rest.NewRoute(http.MethodGet, "/players/:id", getHandler)
	out = append(out, get)

	listHandler := fromPlayerServiceAwareHttpHandler(listPlayers, service)
	list := rest.NewRoute(http.MethodGet, "/players", listHandler)
	out = append(out, list)

	deleteHandler := fromPlayerServiceAwareHttpHandler(deletePlayer, service)
	delete := rest.NewRoute(http.MethodDelete, "/players/:id", deleteHandler)
	out = append(out, delete)

	return out
}

// createPlayer godoc
//
//	@Summary		Create player
//	@Description	Creates a player and its homeworld.
//	@Tags			players
//	@Accept			json
//	@Produce		json
//	@Param			request	body		PlayerRequestDoc	true	"Player payload"
//	@Success		201		{object}	PlayerResponseDoc
//	@Failure		400		{string}	string	"Invalid player syntax"
//	@Failure		409		{string}	string	"Name already used"
//	@Failure		500		{object}	ToolkitErrorDoc
//	@Router			/players [post]
func createPlayer(c *echo.Context, s service.PlayerService) error {
	var playerDtoRequest communication.PlayerDtoRequest
	err := c.Bind(&playerDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid player syntax")
	}

	out, err := s.Create(c.Request().Context(), playerDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation) {
			return c.JSON(http.StatusConflict, "Name already used")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

// getPlayer godoc
//
//	@Summary		Get player
//	@Description	Returns a player by id.
//	@Tags			players
//	@Produce		json
//	@Param			id	path		string	true	"Player id (UUID)"
//	@Success		200	{object}	PlayerResponseDoc
//	@Failure		400	{string}	string	"Invalid id syntax"
//	@Failure		404	{string}	string	"No such player"
//	@Failure		500	{object}	ToolkitErrorDoc
//	@Router			/players/{id} [get]
func getPlayer(c *echo.Context, s service.PlayerService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := s.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such player")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

// listPlayers godoc
//
//	@Summary		List players
//	@Description	Returns players, optionally filtered by api_user.
//	@Tags			players
//	@Produce		json
//	@Param			api_user	query		string	false	"API user id (UUID)"
//	@Success		200			{array}		PlayerResponseDoc
//	@Failure		400			{string}	string	"Invalid id syntax"
//	@Failure		500			{object}	ToolkitErrorDoc
//	@Router			/players [get]
func listPlayers(c *echo.Context, s service.PlayerService) error {
	exists, apiUser, err := fetchIdFromQueryParam("api_user", c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	var players []communication.PlayerDtoResponse

	if exists {
		players, err = s.ListForApiUser(c.Request().Context(), apiUser)
	} else {
		players, err = s.List(c.Request().Context())
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	out, err := marshalNilToEmptySlice(players)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, out)
}

// deletePlayer godoc
//
//	@Summary		Delete player
//	@Description	Deletes a player by id.
//	@Tags			players
//	@Produce		json
//	@Param			id	path		string	true	"Player id (UUID)"
//	@Success		204	{string}	string
//	@Failure		400	{string}	string	"Invalid id syntax"
//	@Failure		404	{string}	string	"No such player"
//	@Failure		500	{object}	ToolkitErrorDoc
//	@Router			/players/{id} [delete]
func deletePlayer(c *echo.Context, s service.PlayerService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such player")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
