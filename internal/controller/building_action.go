package controller

import (
	"fmt"
	"net/http"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/game"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func BuildingActionEndpoints(service service.BuildingActionService) rest.Routes {
	var out rest.Routes

	path := fmt.Sprintf("/planets/%s/actions", rest.RouteIdPlaceholder)
	postHandler := fromBuildingActionServiceAwareHttpHandler(createBuildingAction, service)
	post := rest.NewResourceRoute(http.MethodPost, false, path, postHandler)
	out = append(out, post)

	deleteHandler := fromBuildingActionServiceAwareHttpHandler(createBuildingAction, service)
	delete := rest.NewResourceRoute(http.MethodDelete, false, "/actions", deleteHandler)
	out = append(out, delete)

	return out
}

func createBuildingAction(c echo.Context, s service.BuildingActionService) error {
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
		if errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey) {
			return c.JSON(http.StatusConflict, "Building action already exists")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func deleteBuildingAction(c echo.Context, s service.BuildingActionService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such action")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
