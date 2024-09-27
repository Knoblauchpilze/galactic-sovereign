package game

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ActionService interface {
	ProcessActionsUntil(ctx context.Context, planet uuid.UUID, until time.Time) error
}

type PlanetResourceService interface {
	UpdatePlanetUntil(ctx context.Context, planet uuid.UUID, until time.Time) error
}

type actionProcessingData struct {
	lock                  sync.Mutex
	actionService         ActionService
	planetResourceService PlanetResourceService
}

func GameUpdateWatcher(actionService ActionService, planetResourceService PlanetResourceService, next echo.HandlerFunc) echo.HandlerFunc {
	data := actionProcessingData{
		actionService:         actionService,
		planetResourceService: planetResourceService,
	}

	return func(c echo.Context) error {
		timeStamp := time.Now()

		maybeId := c.Param("id")
		if id, err := uuid.Parse(maybeId); err == nil {
			err := data.updateGameToCurrentTime(c.Request().Context(), id, timeStamp)

			if err != nil {
				return handleError(err, c)
			}
		}

		return next(c)
	}
}

func handleError(err error, c echo.Context) error {
	if errors.IsErrorWithCode(err, actionSchedulingFailed) {
		c.Logger().Errorf("Failed to scheduled pending actions %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to process actions")
	}
	if errors.IsErrorWithCode(err, planetResourceUpdateFailed) {
		c.Logger().Errorf("Failed to update planet resources %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to update resources")
	}

	c.Logger().Errorf("Failed to update game to current time %v", err)
	return c.JSON(http.StatusInternalServerError, "Failed to update game")
}

func (data *actionProcessingData) updateGameToCurrentTime(ctx context.Context, planet uuid.UUID, timeStamp time.Time) error {
	data.lock.Lock()
	defer data.lock.Unlock()

	err := data.actionService.ProcessActionsUntil(ctx, planet, timeStamp)
	if err != nil {
		return errors.WrapCode(err, actionSchedulingFailed)
	}

	err = data.planetResourceService.UpdatePlanetUntil(ctx, planet, timeStamp)
	if err != nil {
		return errors.WrapCode(err, planetResourceUpdateFailed)
	}

	return nil
}
