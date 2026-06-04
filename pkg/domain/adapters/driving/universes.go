package driving

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
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
func createUniverse(c *echo.Context, usecase driving.ForCreatingUniverse) error {
	var inputDto dtos.UniverseDtoRequest
	err := c.Bind(&inputDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid universe syntax")
	}

	request := request.UniverseCreationRequest{
		Name: inputDto.Name,
	}

	out, err := usecase.Create(c.Request().Context(), request)
	if err != nil {
		if errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation) {
			return c.JSON(http.StatusConflict, "Name already used")
		}

		c.Logger().Error("Failed to create universe", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to create universe")
	}

	return c.JSON(http.StatusCreated, out)
}
