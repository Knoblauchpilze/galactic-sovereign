package drivingadapters

import (
	"log/slog"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/mappers"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func BuildingActionEndpoints(usecase drivingports.ForManagingBuildingAction) rest.Routes {
	var out rest.Routes

	handler := generateHandler(createBuildingAction, usecase)
	post := rest.NewRoute(http.MethodPost, "/planets/:id/actions", handler)
	out = append(out, post)

	handler = generateHandler(deleteBuildingAction, usecase)
	delete := rest.NewRoute(http.MethodDelete, "/actions/:id", handler)
	out = append(out, delete)

	return out
}

// createBuildingAction godoc
//
//	@Summary		Create building action
//	@Description	Creates a building action for the planet provided in path parameter. The planet field in the body is ignored and replaced with this path value.
//	@Tags			actions
//	@Produce		json
//	@Param			id		path		string					true	"Planet id (UUID)"	Format(uuid)
//	@Param			request	body		dtos.BuildingActionDtoRequest	true	"Building action payload"
//	@Success		201		{object}	rest.ResponseEnvelope[dtos.BuildingActionDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		404		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/planets/{id}/actions [post]
func createBuildingAction(c *echo.Context, usecase drivingports.ForManagingBuildingAction) error {
	maybeId := c.Param("id")
	planetId, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	var inputDto dtos.BuildingActionDtoRequest
	err = c.Bind(&inputDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid building action syntax")
	}

	request := mappers.ToBuildingActionCreationRequest(planetId, inputDto)
	action, err := usecase.Create(c.Request().Context(), request)
	if err != nil {
		if err == domainerrors.ErrActionAlreadyInProgress {
			return c.JSON(http.StatusConflict, "action already in progress")
		}

		if err == domainerrors.ErrNotFound {
			return c.JSON(http.StatusNotFound, "no such planet")
		}

		if err == domainerrors.ErrBuildingNotFound {
			return c.JSON(http.StatusBadRequest, "no such building")
		}

		if err == domainerrors.ErrNotEnoughResources {
			return c.JSON(http.StatusBadRequest, "not enough resources")
		}

		c.Logger().Error("Failed to create building action", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to create building action")
	}

	out := mappers.ToBuildingActionResponse(action)
	return c.JSON(http.StatusCreated, out)
}

// deleteBuildingAction godoc
//
//	@Summary		Delete building action
//	@Description	Deletes an existing building action.
//	@Tags			actions
//	@Produce		json
//	@Param			id	path		string	true	"Action id (UUID)"	Format(uuid)
//	@Success		204	{string}	string
//	@Failure		400	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/actions/{id} [delete]
func deleteBuildingAction(c *echo.Context, usecase drivingports.ForManagingBuildingAction) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id syntax")
	}

	err = usecase.Delete(c.Request().Context(), id)
	if err != nil {
		c.Logger().Error("Failed to delete building action", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, "failed to delete building action")
	}

	return c.NoContent(http.StatusNoContent)
}
