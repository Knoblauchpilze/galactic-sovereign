package controllers

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var errDefault = fmt.Errorf("some error")

type mockUserRepository struct {
	repositories.UserRepository
}

type mockEchoContext struct {
	echo.Context
}

func TestGenerateEchoHandler_CallsHandler(t *testing.T) {
	assert := assert.New(t)

	handlerCalled := false
	in := func(c echo.Context, repo repositories.UserRepository) error {
		handlerCalled = true
		return nil
	}

	h := generateEchoHandler(in, mockUserRepository{})

	err := h(mockEchoContext{})
	assert.Nil(err)
	assert.True(handlerCalled)
}

func TestGenerateEchoHandler_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	in := func(c echo.Context, repo repositories.UserRepository) error {
		return errDefault
	}

	h := generateEchoHandler(in, mockUserRepository{})

	err := h(mockEchoContext{})
	assert.Equal(errDefault, err)
}
