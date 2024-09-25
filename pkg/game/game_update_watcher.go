package game

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type ActionService interface {
	ProcessActionsUntil(ctx context.Context, until time.Time) error
}

type actionProcessingData struct {
	lock    sync.Mutex
	service ActionService
}

func GameUpdateWatcher(service ActionService, next echo.HandlerFunc) echo.HandlerFunc {
	data := actionProcessingData{
		service: service,
	}

	return func(c echo.Context) error {
		err := data.schedulePendingActions(c.Request().Context())
		if err != nil {
			c.Logger().Errorf("Failed to scheduled pending actions %v", err)
			return c.JSON(http.StatusInternalServerError, "Failed to process actions")
		}

		return next(c)
	}
}

func (data *actionProcessingData) schedulePendingActions(ctx context.Context) error {
	data.lock.Lock()
	defer data.lock.Unlock()

	return data.service.ProcessActionsUntil(ctx, time.Now())
}
