package controller

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

type mockEchoContext struct {
	echo.Context
}

func TestFromRepositoriesAwareHttpHandler_CallsHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	in := func(_ echo.Context, _ service.UserService) error {
		handlerCalled = true
		return nil
	}

	h := fromRepositoriesAwareHttpHandler(in, &mockUserService{})

	err := h(mockEchoContext{})
	assert.Nil(err)
	assert.True(handlerCalled)
}

func TestFromRepositoriesAwareHttpHandler_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	in := func(_ echo.Context, _ service.UserService) error {
		return errDefault
	}

	h := fromRepositoriesAwareHttpHandler(in, &mockUserService{})

	err := h(mockEchoContext{})
	assert.Equal(errDefault, err)
}
