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

func BuildingActionEndpoints(
	createUsecase drivingports.ForCreatingBuildingAction,
	deleteUsecase drivingports.ForDeletingBuildingAction,
) Routes {
	var out Routes

	handler := generateHandler(createBuildingAction, createUsecase)
	post := rest.NewRoute(http.MethodPost, "/planets/:id/actions", handler)
	out = append(out, post)

	handler = generateHandler(deleteBuildingAction, deleteUsecase)
	delete := rest.NewRoute(http.MethodDelete, "/planets/:id/actions", handler)
	out = append(out, delete)

	return out
}

// createBuildingAction godoc
//
//	@Summary		Create building action
//	@Description	Creates a building action for the planet provided in path parameter. The planet field in the body is ignored and replaced with this path value.
//	@Tags			planets
//	@Produce		json
//	@Param			id		path		string					true	"Planet id (UUID)"	Format(uuid)
//	@Param			request	body		dtos.BuildingActionDtoRequest	true	"Building action payload"
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.BuildingActionDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		404		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/planets/{id}/actions [post]
func createBuildingAction(c *gin.Context, usecase drivingports.ForCreatingBuildingAction) {
	maybeId := c.Param("id")
	planetId, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	var inputDto dtos.BuildingActionDtoRequest
	err = c.Bind(&inputDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid building action syntax")
		return
	}

	request := mappers.ToBuildingActionCreationRequest(planetId, inputDto)
	action, err := usecase.Create(c.Request.Context(), request)
	if err != nil {
		if err == domainerrors.ErrActionAlreadyInProgress {
			c.AbortWithStatusJSON(http.StatusConflict, "action already in progress")
			return
		}

		if err == domainerrors.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, "no such planet")
			return
		}

		if err == domainerrors.ErrBuildingNotFound {
			c.AbortWithStatusJSON(http.StatusBadRequest, "no such building")
			return
		}

		if err == domainerrors.ErrNotEnoughResources {
			c.AbortWithStatusJSON(http.StatusBadRequest, "not enough resources")
			return
		}

		if err == domainerrors.ErrAllFieldsUsed {
			c.AbortWithStatusJSON(http.StatusConflict, "all fields are used")
			return
		}

		logError(c.Request, "Failed to create building action", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to create building action")
		return
	}

	out := mappers.ToBuildingActionResponse(action)
	c.JSON(http.StatusCreated, out)
}

// deleteBuildingAction godoc
//
//	@Summary		Delete building action for a planet
//	@Description	Deletes an existing building action for a planet.
//	@Tags			planets
//	@Produce		json
//	@Param			id	path		string	true	"Planet id (UUID)"	Format(uuid)
//	@Success		204	{string}	string
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/planets/{id}/actions [delete]
func deleteBuildingAction(c *gin.Context, usecase drivingports.ForDeletingBuildingAction) {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid id syntax")
		return
	}

	err = usecase.DeleteForPlanet(c.Request.Context(), id)
	if err != nil {
		if err == domainerrors.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, "no such planet")
			return
		}

		logError(c.Request, "Failed to delete building action", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to delete building action")
		return
	}

	c.Status(http.StatusNoContent)
}
