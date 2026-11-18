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

func FleetEndpoints(
	createUsecase drivingports.ForCreatingFleet,
	fetchUsecase drivingports.ForFetchingFleets,
) Routes {
	var out Routes

	handler := generateHandler(createFleet, createUsecase)
	post := rest.NewRoute(http.MethodPost, "/planets/:id/fleets", handler)
	out = append(out, post)

	handler = generateHandler(listFleetsForPlayer, fetchUsecase)
	list := rest.NewRoute(http.MethodGet, "/players/:id/fleets", handler)
	out = append(out, list)

	return out
}

// createFleet godoc
//
//	@Summary		Create fleet
//	@Description	Creates a fleet for the planet provided in path parameter.
//	@Tags			fleets
//	@Produce		json
//	@Param			id		path		string					true	"Planet id (UUID)"	Format(uuid)
//	@Param			request	body		dtos.FleetDtoRequest	true	"Fleet payload"
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.FleetDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/planets/{id}/fleets [post]
func createFleet(c *gin.Context, usecase drivingports.ForCreatingFleet) {
	maybeId := c.Param("id")
	planetId, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	var inputDto dtos.FleetDtoRequest
	err = c.Bind(&inputDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid fleet syntax")
		return
	}

	request, err := mappers.ToFleetOrder(inputDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid fleet syntax")
		return
	}

	fleet, err := usecase.Create(c.Request.Context(), planetId, request)
	if err != nil {
		if err == domainerrors.ErrNotEnoughShips {
			c.AbortWithStatusJSON(http.StatusConflict, "not enough ships")
			return
		}

		if err == domainerrors.ErrFleetDestinationInvalid {
			c.AbortWithStatusJSON(http.StatusConflict, "invalid fleet destination")
			return
		}

		if err == domainerrors.ErrInvalidFleetConfiguration {
			c.AbortWithStatusJSON(http.StatusBadRequest, "invalid fleet configuration")
			return
		}

		if err == domainerrors.ErrCoordinatesOutOfBound {
			c.AbortWithStatusJSON(http.StatusBadRequest, "invalid fleet coordinates")
			return
		}

		logError(c.Request, "Failed to create fleet", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to create fleet")
		return
	}

	out := mappers.ToFleetResponse(fleet)
	c.JSON(http.StatusCreated, out)
}

// listFleetsForPlayer godoc
//
//	@Summary		List fleets
//	@Description	Returns fleets belonging to a player or targeting a player's planet.
//	@Tags			players
//	@Produce		json
//	@Param			id	path		string	true	"Player id (UUID)"	Format(uuid)
//	@Success		200		{object}	rest.ResponseEnvelope[[]dtos.FleetDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/players/{id}/fleets [get]
func listFleetsForPlayer(c *gin.Context, usecase drivingports.ForFetchingFleets) {
	maybeId := c.Param("id")
	playerId, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	fleets, err := usecase.ListForPlayer(c.Request.Context(), playerId)

	if err != nil {
		logError(c.Request, "Failed to list fleets", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to list fleets")
		return
	}

	out := mappers.ToFleetsResponse(fleets)

	c.JSON(http.StatusOK, out)
}
