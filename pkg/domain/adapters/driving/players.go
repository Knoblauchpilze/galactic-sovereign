package drivingadapters

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func PlayerEndpoints(usecase drivingports.ForManagingPlayer) rest.Routes {
	var out rest.Routes

	handler := generateHandler(createPlayer, usecase)
	post := rest.NewRoute(http.MethodPost, "/players", handler)
	out = append(out, post)

	handler = generateHandler(getPlayer, usecase)
	get := rest.NewRoute(http.MethodGet, "/players/:id", handler)
	out = append(out, get)

	handler = generateHandler(listPlayers, usecase)
	list := rest.NewRoute(http.MethodGet, "/players", handler)
	out = append(out, list)

	handler = generateHandler(deletePlayer, usecase)
	delete := rest.NewRoute(http.MethodDelete, "/players/:id", handler)
	out = append(out, delete)

	return out
}

// createPlayer godoc
//
//	@Summary		Create player
//	@Description	Creates a player and its homeworld.
//	@Tags			players
//	@Produce		json
//	@Param			request	body		dtos.PlayerDtoRequest	true	"Player payload"
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.PlayerDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/players [post]
func createPlayer(c *echo.Context, usecase drivingports.ForManagingPlayer) error {
	var inputDto dtos.PlayerDtoRequest
	err := c.Bind(&inputDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid player syntax")
	}

	request := mappers.ToPlayerCreationRequest(inputDto)
	player, err := usecase.Create(c.Request().Context(), request)
	if err != nil {
		if dbErr, ok := db.AsDatabaseError(err); ok && dbErr.Code == db.ErrUniqueConstraintViolation {
			return c.JSON(http.StatusConflict, "name already used")
		}

		c.Logger().Error("Failed to create player", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to create player")
	}

	out := mappers.ToPlayerResponse(player)
	return c.JSON(http.StatusCreated, out)
}

// getPlayer godoc
//
//	@Summary		Get player
//	@Description	Returns a player by id.
//	@Tags			players
//	@Produce		json
//	@Param			id	path		string	true	"Player id (UUID)"	Format(uuid)
//	@Success		200	{object}	rest.ResponseEnvelope[dtos.PlayerDtoResponse]
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/players/{id} [get]
func getPlayer(c *echo.Context, usecase drivingports.ForManagingPlayer) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	player, err := usecase.Get(c.Request().Context(), id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return c.JSON(http.StatusNotFound, "no such player")
		}

		c.Logger().Error("Failed to get player", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to get player")
	}

	out := mappers.ToPlayerResponse(player)
	return c.JSON(http.StatusOK, out)
}

// listPlayers godoc
//
//	@Summary		List players
//	@Description	Returns players, optionally filtered by api_user.
//	@Tags			players
//	@Produce		json
//	@Param			api_user	query		string	false	"API user id (UUID)"
//	@Success		200			{object}	rest.ResponseEnvelope[[]dtos.PlayerDtoResponse]
//	@Failure		400			{object}	rest.ResponseEnvelope[string]
//	@Failure		500			{object}	rest.ResponseEnvelope[string]
//	@Router			/players [get]
func listPlayers(c *echo.Context, usecase drivingports.ForManagingPlayer) error {
	exists, apiUser, err := fetchIdFromQueryParam("api_user", c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	var players []models.Player
	if exists {
		players, err = usecase.ListForApiUser(c.Request().Context(), apiUser)
	} else {
		players, err = usecase.List(c.Request().Context())
	}

	if err != nil {
		c.Logger().Error("Failed to list players", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to list players")
	}

	out := mappers.ToPlayersResponse(players)

	return c.JSON(http.StatusOK, out)
}

// deletePlayer godoc
//
//	@Summary		Delete player
//	@Description	Deletes a player by id.
//	@Tags			players
//	@Produce		json
//	@Param			id	path		string	true	"Player id (UUID)"	Format(uuid)
//	@Success		204	{string}	string
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/players/{id} [delete]
func deletePlayer(c *echo.Context, usecase drivingports.ForManagingPlayer) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	err = usecase.Delete(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error("Failed to delete player", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to delete player")
	}

	return c.NoContent(http.StatusNoContent)
}
