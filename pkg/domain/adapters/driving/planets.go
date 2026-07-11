package drivingadapters

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func PlanetEndpoints(
	createUsecase drivingports.ForCreatingPlanet,
	usecase drivingports.ForManagingPlanet,
) rest.Routes {
	var out rest.Routes

	handler := generateHandler(createPlanet, createUsecase)
	post := rest.NewRoute(http.MethodPost, "/players/:id/planets", handler)
	out = append(out, post)

	handler = generateHandler(getPlanet, usecase)
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

// createPlanet godoc
//
//	@Summary		Create planet
//	@Description	Creates a planet for the player
//	@Tags			players
//	@Produce		json
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.PlanetDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/players/{id}/planets [post]
func createPlanet(c *echo.Context, usecase drivingports.ForCreatingPlanet) error {
	maybeId := c.Param("id")
	playerId, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	request := request.PlanetCreationRequest{Player: playerId}
	planet, err := usecase.Create(c.Request().Context(), request)
	if err != nil {
		if dbErr, ok := db.AsDatabaseError(err); ok && dbErr.Code == db.ErrUniqueConstraintViolation {
			return c.JSON(http.StatusConflict, "name already used")
		}

		if err == domainerrors.ErrPlayerNotFound {
			return c.JSON(http.StatusBadRequest, "no such player")
		}

		c.Logger().Error("Failed to create planet", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to create planet")
	}

	out := mappers.ToPlanetResponse(planet)
	return c.JSON(http.StatusCreated, out)
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
func getPlanet(c *echo.Context, usecase drivingports.ForManagingPlanet) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	planet, err := usecase.Get(c.Request().Context(), id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			return c.JSON(http.StatusNotFound, "no such planet")
		}

		c.Logger().Error("Failed to get planet", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to get planet")
	}

	out := mappers.ToPlanetResponse(planet)
	return c.JSON(http.StatusOK, out)
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
func listPlanetsForPlayer(c *echo.Context, usecase drivingports.ForManagingPlanet) error {
	maybeId := c.Param("id")
	playerId, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	planets, err := usecase.ListForPlayer(c.Request().Context(), playerId)

	if err != nil {
		c.Logger().Error("Failed to list planets", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to list planets")
	}

	out := mappers.ToPlanetsResponse(planets)

	return c.JSON(http.StatusOK, out)
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
func deletePlanet(c *echo.Context, usecase drivingports.ForManagingPlanet) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	err = usecase.Delete(c.Request().Context(), id)
	if err != nil {
		if err == domainerrors.ErrActionNotCompleted {
			return c.JSON(http.StatusConflict, "action not completed")
		}

		if err == domainerrors.ErrHomeworldCannotBeDeleted {
			return c.JSON(http.StatusConflict, "homeworld cannot be deleted")
		}

		c.Logger().Error("Failed to delete planet", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to delete planet")
	}

	return c.NoContent(http.StatusNoContent)
}
