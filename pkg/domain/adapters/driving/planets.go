package drivingadapters

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PlanetEndpoints(usecase drivingports.ForManagingPlanet) Routes {
	var out Routes

	handler := generateHandler(getPlanet, usecase)
	get := rest.NewRoute(http.MethodGet, "/planets/:id", handler)
	out = append(out, get)

	handler = generateHandler(listPlanetsForPlayer, usecase)
	list := rest.NewRoute(http.MethodGet, "/players/:id/planets", handler)
	out = append(out, list)

	handler = generateHandler(deletePlanet, usecase)
	delete := rest.NewRoute(http.MethodDelete, "/planets/:id", handler)
	out = append(out, delete)

	return out
}

// getPlanet godoc
//
//	@Summary		Get planet
//	@Description	Returns a planet and all related game data.
//	@Tags			planets
//	@Produce		json
//	@Param			id	path		string	true	"Planet id (UUID)"	Format(uuid)
//	@Success		200	{object}	rest.ResponseEnvelope[dtos.PlanetDtoResponse]
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/planets/{id} [get]
func getPlanet(c *gin.Context, usecase drivingports.ForManagingPlanet) {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	planet, err := usecase.Get(c.Request.Context(), id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, "no such planet")
			return
		}

		logError(c.Request, "Failed to get planet", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get planet")
		return
	}

	out := mappers.ToPlanetResponse(planet)
	c.JSON(http.StatusOK, out)
}

// listPlanetsForPlayer godoc
//
//	@Summary		List planets
//	@Description	Returns planets belonging to a player.
//	@Tags			players
//	@Produce		json
//	@Param			id	path		string	true	"Player id (UUID)"	Format(uuid)
//	@Success		200		{object}	rest.ResponseEnvelope[[]dtos.PlanetDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/players/{id}/planets [get]
func listPlanetsForPlayer(c *gin.Context, usecase drivingports.ForManagingPlanet) {
	maybeId := c.Param("id")
	playerId, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	planets, err := usecase.ListForPlayer(c.Request.Context(), playerId)

	if err != nil {
		logError(c.Request, "Failed to list planets", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to list planets")
		return
	}

	out := mappers.ToPlanetsResponse(planets)

	c.JSON(http.StatusOK, out)
}

// deletePlanet godoc
//
//	@Summary		Delete planet
//	@Description	Deletes a planet by id.
//	@Tags			planets
//	@Produce		json
//	@Param			id	path		string	true	"Planet id (UUID)"	Format(uuid)
//	@Success		204	{string}	string
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		409	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/planets/{id} [delete]
func deletePlanet(c *gin.Context, usecase drivingports.ForManagingPlanet) {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	err = usecase.Delete(c.Request.Context(), id)
	if err != nil {
		if err == domainerrors.ErrActionNotCompleted {
			c.AbortWithStatusJSON(http.StatusConflict, "action not completed")
			return
		}

		if err == domainerrors.ErrHomeworldCannotBeDeleted {
			c.AbortWithStatusJSON(http.StatusConflict, "homeworld cannot be deleted")
			return
		}

		logError(c.Request, "Failed to delete planet", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to delete planet")
		return
	}

	c.Status(http.StatusNoContent)
}
