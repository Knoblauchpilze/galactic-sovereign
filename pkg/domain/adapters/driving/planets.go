package drivingadapters

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func PlanetEndpoints(usecase drivingports.ForManagingPlanet) rest.Routes {
	var out rest.Routes

	handler := generateHandler(createPlanet, usecase)
	post := rest.NewRoute(http.MethodPost, "/planets", handler)
	out = append(out, post)

	handler = generateHandler(getPlanet, usecase)
	get := rest.NewRoute(http.MethodGet, "/planets/:id", handler)
	out = append(out, get)

	handler = generateHandler(listPlanets, usecase)
	list := rest.NewRoute(http.MethodGet, "/planets", handler)
	out = append(out, list)

	handler = generateHandler(deletePlanet, usecase)
	delete := rest.NewRoute(http.MethodDelete, "/planets/:id", handler)
	out = append(out, delete)

	return out
}

// createPlanet godoc
//
//	@Summary		Create planet
//	@Description	Creates a planet from the provided payload.
//	@Tags			planets
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.PlanetDtoRequest	true	"Planet payload"
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.PlanetDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/planets [post]
func createPlanet(c *echo.Context, usecase drivingports.ForManagingPlanet) error {
	var inputDto dtos.PlanetDtoRequest
	err := c.Bind(&inputDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid planet syntax")
	}

	request := mappers.ToPlanetCreationRequest(inputDto)
	planet, err := usecase.Create(c.Request().Context(), request)
	if err != nil {
		if errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation) {
			return c.JSON(http.StatusConflict, "name already used")
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
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "no such planet")
		}

		c.Logger().Error("Failed to get planet", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to get planet")
	}

	out := mappers.ToPlanetResponse(planet)
	return c.JSON(http.StatusOK, out)
}

// listPlanets godoc
//
//	@Summary		List planets
//	@Description	Returns planets, optionally filtered by player id.
//	@Tags			planets
//	@Produce		json
//	@Param			player	query		string	false	"Player id (UUID)"
//	@Success		200		{object}	rest.ResponseEnvelope[[]dtos.PlanetDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/planets [get]
func listPlanets(c *echo.Context, usecase drivingports.ForManagingPlanet) error {
	exists, playerId, err := fetchIdFromQueryParam("player", c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	var planets []models.Planet
	if exists {
		planets, err = usecase.ListForPlayer(c.Request().Context(), playerId)
	} else {
		planets, err = usecase.List(c.Request().Context())
	}

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
		c.Logger().Error("Failed to delete planet", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to delete planet")
	}

	return c.NoContent(http.StatusNoContent)
}
