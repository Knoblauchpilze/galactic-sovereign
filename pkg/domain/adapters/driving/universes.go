package driving

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// createUniverse godoc
//
//	@Summary		Create universe
//	@Description	Creates a universe.
//	@Tags			universes
//	@Accept			json
//	@Produce		json
//	@Param			request	body		communication.UniverseDtoRequest	true	"Universe payload"
//	@Success		201		{object}	rest.ResponseEnvelope[communication.UniverseDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
func CreateUniverse(c *echo.Context, usecase driving.ForManagingUniverse) error {
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

// ListUniverses godoc
//
//	@Summary		List universes
//	@Description	Returns all universes.
//	@Tags			universes
//	@Produce		json
//	@Success		200	{object}	rest.ResponseEnvelope[[]communication.UniverseDtoResponse]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes [get]
func ListUniverses(c *echo.Context, usecase driving.ForManagingUniverse) error {
	universes, err := usecase.List(c.Request().Context())
	if err != nil {
		c.Logger().Error("Failed to list universes", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to list universes")
	}

	out := mappers.ToUniversesResponse(universes)

	return c.JSON(http.StatusOK, out)
}

// DeleteUniverse godoc
//
//	@Summary		Delete universe
//	@Description	Deletes a universe by id.
//	@Tags			universes
//	@Produce		json
//	@Param			id	path		string	true	"Universe id (UUID)"	Format(uuid)
//	@Success		204	{string}	string
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/universes/{id} [delete]
func DeleteUniverse(c *echo.Context, usecase driving.ForManagingUniverse) error {
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
