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
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func UniverseEndpoints(usecase drivingports.ForManagingUniverse) rest.Routes {
	var out rest.Routes

	handler := generateHandler(createUniverse, usecase)
	post := rest.NewRoute(http.MethodPost, "/universes", handler)
	out = append(out, post)

	handler = generateHandler(getUniverse, usecase)
	get := rest.NewRoute(http.MethodGet, "/universes/:id", handler)
	out = append(out, get)

	handler = generateHandler(listUniverses, usecase)
	list := rest.NewRoute(http.MethodGet, "/universes", handler)
	out = append(out, list)

	handler = generateHandler(deleteUniverse, usecase)
	delete := rest.NewRoute(http.MethodDelete, "/universes/:id", handler)
	out = append(out, delete)

	return out
}

// createUniverse godoc
//
//	@Summary		Create universe
//	@Description	Creates a universe.
//	@Tags			universes
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.UniverseDtoRequest	true	"Universe payload"
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.UniverseDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/universes [post]
func createUniverse(c *echo.Context, usecase drivingports.ForManagingUniverse) error {
	var inputDto dtos.UniverseDtoRequest
	err := c.Bind(&inputDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid universe syntax")
	}

	request := mappers.ToUniverseCreationRequest(inputDto)
	universe, err := usecase.Create(c.Request().Context(), request)
	if err != nil {
		if errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation) {
			return c.JSON(http.StatusConflict, "name already used")
		}

		c.Logger().Error("Failed to create universe", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to create universe")
	}

	out := mappers.ToUniverseResponse(universe)
	return c.JSON(http.StatusCreated, out)
}

// getUniverse godoc
//
//	@Summary		Get universe
//	@Description	Returns a universe and related resources/buildings.
//	@Tags			universes
//	@Produce		json
//	@Param			id	path		string	true	"Universe id (UUID)"	Format(uuid)
//	@Success		200	{object}	rest.ResponseEnvelope[dtos.UniverseDtoResponse]
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes/{id} [get]
func getUniverse(c *echo.Context, usecase drivingports.ForManagingUniverse) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	universe, err := usecase.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "no such universe")
		}

		c.Logger().Error("Failed to get universe", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to get universe")
	}

	out := mappers.ToUniverseResponse(universe)
	return c.JSON(http.StatusOK, out)
}

// listUniverses godoc
//
//	@Summary		List universes
//	@Description	Returns all universes.
//	@Tags			universes
//	@Produce		json
//	@Success		200	{object}	rest.ResponseEnvelope[[]dtos.UniverseDtoResponse]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes [get]
func listUniverses(c *echo.Context, usecase drivingports.ForManagingUniverse) error {
	universes, err := usecase.List(c.Request().Context())
	if err != nil {
		c.Logger().Error("Failed to list universes", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to list universes")
	}

	out := mappers.ToUniversesResponse(universes)

	return c.JSON(http.StatusOK, out)
}

// deleteUniverse godoc
//
//	@Summary		Delete universe
//	@Description	Deletes a universe by id.
//	@Tags			universes
//	@Produce		json
//	@Param			id	path		string	true	"Universe id (UUID)"	Format(uuid)
//	@Success		204	{string}	string
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes/{id} [delete]
func deleteUniverse(c *echo.Context, usecase drivingports.ForManagingUniverse) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	err = usecase.Delete(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error("Failed to delete universe", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to delete universe")
	}

	return c.NoContent(http.StatusNoContent)
}
