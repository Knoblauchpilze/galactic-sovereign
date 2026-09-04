package drivingadapters

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ShipActionEndpoints(usecase drivingports.ForCreatingShipAction) Routes {
	var out Routes

	handler := generateHandler(createShipAction, usecase)
	post := rest.NewRoute(http.MethodPost, "/planets/:id/ships", handler)
	out = append(out, post)

	return out
}

// createShipAction godoc
//
//	@Summary		Create ship action
//	@Description	Creates a ship action on a planet.
//	@Tags			planets
//	@Produce		json
//	@Param			request	body		dtos.ShipDtoRequest	true	"Ship payload"
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.ShipActionDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/planets/:id/ships [post]
func createShipAction(c *gin.Context, usecase drivingports.ForCreatingShipAction) {
	maybeId := c.Param("id")
	planetId, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	var inputDto dtos.ShipActionDtoRequest
	err = c.Bind(&inputDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid ship syntax")
		return
	}

	request := mappers.ToShipActionCreationRequest(planetId, inputDto)
	action, err := usecase.Create(c.Request.Context(), request)
	if err != nil {
		// TODO: Handle other errors

		logError(c.Request, "Failed to create ship", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to create ship")
		return
	}

	out := mappers.ToShipActionResponse(action)
	c.JSON(http.StatusCreated, out)
}
