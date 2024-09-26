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
	ProcessActionsUntil(ctx context.Context, until time.Time) error
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
		var planet *uuid.UUID
		if id, err := uuid.Parse(maybeId); err == nil {
			planet = &id
		}

		err := data.updateGameToCurrentTime(c.Request().Context(), planet, timeStamp)
		if err != nil {
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

		return next(c)
	}
}

func (data *actionProcessingData) updateGameToCurrentTime(ctx context.Context, planet *uuid.UUID, timeStamp time.Time) error {
	data.lock.Lock()
	defer data.lock.Unlock()

	err := data.actionService.ProcessActionsUntil(ctx, timeStamp)
	if err != nil {
		return errors.WrapCode(err, actionSchedulingFailed)
	}

	if planet == nil {
		return nil
	}

	err = data.planetResourceService.UpdatePlanetUntil(ctx, *planet, timeStamp)
	if err != nil {
		return errors.WrapCode(err, planetResourceUpdateFailed)
	}

	return nil
}
