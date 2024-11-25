package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUnit_Throttle_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next, called := createHandlerFuncWithCalledBoolean()

	em, close := Throttle(1, 1, 1)
	callable := em(next)
	callable(ctx)

	assert.True(*called)
	close <- true
}

func TestUnit_Throttle_WhenNoTokensLeft_ExpectTooManyRequests(t *testing.T) {
	assert := assert.New(t)
	ctx, rw := generateTestEchoContext()
	next := createHandlerFuncReturning(nil)

	em, close := Throttle(0, 0, 0)
	callable := em(next)
	callable(ctx)

	assert.Equal(http.StatusTooManyRequests, rw.Code)
	close <- true
}

func TestUnit_Throttle_WhenWaitingForRefill_ExpectOk(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()
	next := createHandlerFuncReturningCode(http.StatusOK)

	em, close := Throttle(0, 2, 2)
	callable := em(next)

	ctx, rw := generateNewEchoContext(e)
	callable(ctx)
	assert.Equal(http.StatusTooManyRequests, rw.Code)

	time.Sleep(2 * time.Second)

	ctx, rw = generateNewEchoContext(e)
	callable(ctx)
	assert.Equal(http.StatusOK, rw.Code)

	close <- true
}

func TestUnit_Throttle_PropagatesError(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContext()
	next := createHandlerFuncReturning(errDefault)

	em, close := Throttle(1, 1, 1)

	callable := em(next)
	actual := callable(ctx)

	assert.Equal(errDefault, actual)
	close <- true
}

func TestUnit_Throttle_ConcurrentUse_ExpectFirstServed(t *testing.T) {
	assert := assert.New(t)

	em, close := Throttle(1, 1, 1)
	next := createHandlerFuncReturningCode(http.StatusOK)
	handler := em(next)

	e := echo.New()

	var wg sync.WaitGroup
	wg.Add(2)

	c1 := func() {
		defer wg.Done()

		ctx, rw := generateNewEchoContext(e)
		handler(ctx)
		assert.Equal(http.StatusOK, rw.Code)

		time.Sleep(1500 * time.Millisecond)

		ctx, rw = generateNewEchoContext(e)
		handler(ctx)
		assert.Equal(http.StatusTooManyRequests, rw.Code)
	}

	c2 := func() {
		defer wg.Done()

		time.Sleep(1100 * time.Millisecond)

		ctx, rw := generateNewEchoContext(e)
		handler(ctx)
		assert.Equal(http.StatusOK, rw.Code)
	}

	go c1()
	go c2()

	wg.Wait()

	close <- true
}

func generateNewEchoContext(e *echo.Echo) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()
	return e.NewContext(req, rw), rw
}
