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

func FleetEndpoints(
	createUsecase drivingports.ForCreatingFleet,
) Routes {
	var out Routes

	handler := generateHandler(createFleet, createUsecase)
	post := rest.NewRoute(http.MethodPost, "/planets/:id/fleets", handler)
	out = append(out, post)

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

	request := mappers.ToFleetCreationRequest(planetId, inputDto)
	fleet, err := usecase.Create(c.Request.Context(), request)
	if err != nil {
		logError(c.Request, "Failed to create fleet", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to create fleet")
		return
	}

	out := mappers.ToFleetResponse(fleet)
	c.JSON(http.StatusCreated, out)
}
