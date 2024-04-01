package controller

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

type mockEchoContext struct {
	echo.Context
}

func TestGenerateEchoHandler_CallsHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	in := func(_ echo.Context, _ repositories.UserRepository, _ service.UserService) error {
		handlerCalled = true
		return nil
	}

	h := generateEchoHandler(in, &mockUserRepository{}, &mockUserService{})

	err := h(mockEchoContext{})
	assert.Nil(err)
	assert.True(handlerCalled)
}

func TestGenerateEchoHandler_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	in := func(_ echo.Context, _ repositories.UserRepository, _ service.UserService) error {
		return errDefault
	}

	h := generateEchoHandler(in, &mockUserRepository{}, &mockUserService{})

	err := h(mockEchoContext{})
	assert.Equal(errDefault, err)
}
