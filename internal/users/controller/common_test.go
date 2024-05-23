package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

func TestFromRepositoriesAwareHttpHandler_CallsHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	in := func(_ echo.Context, _ service.UserService) error {
		handlerCalled = true
		return nil
	}

	h := fromRepositoriesAwareHttpHandler(in, &mockUserService{})

	err := h(dummyEchoContext())
	assert.Nil(err)
	assert.True(handlerCalled)
}

func TestFromRepositoriesAwareHttpHandler_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	in := func(_ echo.Context, _ service.UserService) error {
		return errDefault
	}

	h := fromRepositoriesAwareHttpHandler(in, &mockUserService{})

	err := h(dummyEchoContext())
	assert.Equal(errDefault, err)
}

func TestFromDbAwareHttpHandler_CallsHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	in := func(_ echo.Context, _ db.ConnectionPool) error {
		handlerCalled = true
		return nil
	}

	h := fromDbAwareHttpHandler(in, &mockConnectionPool{})

	err := h(dummyEchoContext())
	assert.Nil(err)
	assert.True(handlerCalled)
}

func TestFromDbAwareHttpHandler_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	in := func(_ echo.Context, _ db.ConnectionPool) error {
		return errDefault
	}

	h := fromDbAwareHttpHandler(in, &mockConnectionPool{})

	err := h(dummyEchoContext())
	assert.Equal(errDefault, err)
}

func dummyEchoContext() echo.Context {
	ctx, _ := generateTestEchoContext()
	return ctx
}

func generateTestEchoContext() (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	return generateTestEchoContextFromRequest(req)
}

func generateTestEchoContextFromRequest(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rw := httptest.NewRecorder()

	ctx := e.NewContext(req, rw)
	return ctx, rw
}
