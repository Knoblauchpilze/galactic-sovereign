package drivingadapters

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PlayerEndpoints(usecase drivingports.ForManagingPlayer) Routes {
	var out Routes

	handler := generateHandler(createPlayer, usecase)
	post := rest.NewRoute(http.MethodPost, "/players", handler)
	out = append(out, post)

	handler = generateHandler(getPlayer, usecase)
	get := rest.NewRoute(http.MethodGet, "/players/:id", handler)
	out = append(out, get)

	handler = generateHandler(listPlayersForApiUser, usecase)
	list := rest.NewRoute(http.MethodGet, "/users/:id/players", handler)
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
func createPlayer(c *gin.Context, usecase drivingports.ForManagingPlayer) {
	var inputDto dtos.PlayerDtoRequest
	err := c.Bind(&inputDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid player syntax")
		return
	}

	request := mappers.ToPlayerCreationRequest(inputDto)
	player, err := usecase.Create(c.Request.Context(), request)
	if err != nil {
		if err == domainerrors.ErrNameAlreadyTaken {
			c.AbortWithStatusJSON(http.StatusConflict, "name already used")
			return
		}

		if err == domainerrors.ErrUniverseNotFound {
			c.AbortWithStatusJSON(http.StatusBadRequest, "no such universe")
			return
		}

		logError(c.Request, "Failed to create player", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to create player")
		return
	}

	out := mappers.ToPlayerResponse(player)
	c.JSON(http.StatusCreated, out)
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
func getPlayer(c *gin.Context, usecase drivingports.ForManagingPlayer) {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	player, err := usecase.Get(c.Request.Context(), id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, "no such player")
			return
		}

		logError(c.Request, "Failed to get player", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get player")
		return
	}

	out := mappers.ToPlayerResponse(player)
	c.JSON(http.StatusOK, out)
}

// listPlayersForApiUser godoc
//
//	@Summary		List players belonging to a user
//	@Description	Returns players associated to an API user.
//	@Tags			users
//	@Produce		json
//	@Success		200			{object}	rest.ResponseEnvelope[[]dtos.PlayerDtoResponse]
//	@Failure		400			{object}	rest.ResponseEnvelope[string]
//	@Failure		500			{object}	rest.ResponseEnvelope[string]
//	@Router			/users/{id}/players [get]
func listPlayersForApiUser(c *gin.Context, usecase drivingports.ForManagingPlayer) {
	maybeId := c.Param("id")
	apiUserId, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	players, err := usecase.ListForApiUser(c.Request.Context(), apiUserId)
	if err != nil {
		logError(c.Request, "Failed to list players", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to list players")
		return
	}

	out := mappers.ToPlayersResponse(players)

	c.JSON(http.StatusOK, out)
}

// deletePlayer godoc
//
//	@Summary		Delete player
//	@Description	Deletes a player by id.
//	@Tags			players
//	@Produce		json
//	@Param			id	path		string	true	"Player id (UUID)"	Format(uuid)
//	@Success		204
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/players/{id} [delete]
func deletePlayer(c *gin.Context, usecase drivingports.ForManagingPlayer) {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	err = usecase.Delete(c.Request.Context(), id)
	if err != nil {
		logError(c.Request, "Failed to delete player", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to delete player")
		return
	}

	c.Status(http.StatusNoContent)
}
