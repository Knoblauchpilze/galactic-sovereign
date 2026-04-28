package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/game"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func BuildingActionEndpoints(buildingActionService service.BuildingActionService,
	actionService game.ActionService,
	planetResourceService game.PlanetResourceService) rest.Routes {
	var out rest.Routes

	postHandler := fromBuildingActionServiceAwareHttpHandler(createBuildingAction, buildingActionService)
	post := game.NewResourceRoute(http.MethodPost, "/planets/:id/actions", postHandler, actionService, planetResourceService)
	out = append(out, post)

	deleteHandler := fromBuildingActionServiceAwareHttpHandler(deleteBuildingAction, buildingActionService)
	delete := game.NewResourceRoute(http.MethodDelete, "/actions", deleteHandler, actionService, planetResourceService)
	out = append(out, delete)

	return out
}

// createBuildingAction godoc
//
//	@Summary		Create building action
//	@Description	Creates a building action for the planet provided in path parameter. The planet field in the body is ignored and replaced with this path value.
//	@Tags			actions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Planet id (UUID)"	Format(uuid)
//	@Param			request	body		communication.BuildingActionDtoRequest	true	"Building action payload"
//	@Success		201		{object}	rest.ResponseEnvelope[communication.BuildingActionDtoResponse]
//	@Failure		400		{object}	rest.ResponseEnvelope[string]
//	@Failure		409		{object}	rest.ResponseEnvelope[string]
//	@Failure		500		{object}	rest.ResponseEnvelope[string]
//	@Router			/planets/{id}/actions [post]
func createBuildingAction(c *echo.Context, s service.BuildingActionService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	var actionDtoRequest communication.BuildingActionDtoRequest
	err = c.Bind(&actionDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid action syntax")
	}

	actionDtoRequest.Planet = id

	out, err := s.Create(c.Request().Context(), actionDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, game.NotEnoughResources) {
			return c.JSON(http.StatusBadRequest, "Not enough resources")
		}
		if errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation) {
			return c.JSON(http.StatusConflict, "Building action already exists")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

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
//	@Failure		404	{object}	rest.ResponseEnvelope[string]
//	@Failure		500	{object}	rest.ResponseEnvelope[string]
//	@Router			/actions/{id} [delete]
func deleteBuildingAction(c *echo.Context, s service.BuildingActionService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such action")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
