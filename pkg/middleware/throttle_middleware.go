package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type throttle struct {
	refillPerSecond int
	capacity        int

	lock   sync.RWMutex
	tokens int
}

func ThrottleMiddleware(init int, refillPerSecond int, capacity int) (handler echo.MiddlewareFunc, close chan bool) {
	ticker := time.NewTicker(time.Second)

	data := throttle{
		refillPerSecond: refillPerSecond,
		capacity:        capacity,
		tokens:          init,
	}

	close = make(chan bool)

	go func() {
		defer ticker.Stop()

		running := true
		for running {
			select {
			case <-ticker.C:
				data.refill()
			case <-close:
				running = false
			}
		}
	}()

	handler = func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			allowed := false
			func() {
				data.lock.RLock()
				defer data.lock.RUnlock()

				allowed = data.tokens > 0
				if allowed {
					data.tokens--
				}
			}()

			if !allowed {
				c.Logger().Warnf("Throttling %s", c.Request().Host)
				return c.JSON(http.StatusTooManyRequests, "Throttled")
			}

			return next(c)
		}
	}

	return
}

func (t *throttle) refill() {
	t.lock.Lock()
	defer t.lock.Unlock()

	updated := t.tokens + t.refillPerSecond
	if updated > t.capacity {
		updated = t.capacity
	}

	t.tokens = updated
}
