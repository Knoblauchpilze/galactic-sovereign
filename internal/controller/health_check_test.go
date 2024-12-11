package controller

import (
	"context"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mockConnection struct {
	db.Connection

	pingCalled int
	err        error
}

func TestUnit_HealthCheckEndpoints(t *testing.T) {
	s := RouteTestSuite{
		generateRoutes: func() rest.Routes {
			return HealthCheckEndpoints(&mockConnection{})
		},
		expectedRoutes: map[string]int{
			http.MethodGet: 1,
		},
		expectedPaths: map[string]int{
			"/healthcheck": 1,
		},
	}

	suite.Run(t, &s)
}

func TestUnit_Healthcheck_CallsPoolPing(t *testing.T) {
	assert := assert.New(t)

	mc := dummyEchoContext()
	mp := &mockConnection{}

	healthcheck(mc, mp)

	assert.Equal(1, mp.pingCalled)
}

func TestUnit_Healthcheck_WhenPingSucceeds_SetsSatusToOk(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateTestEchoContextWithMethod(http.MethodGet)
	mp := &mockConnection{}

	err := healthcheck(ctx, mp)

	assert.Nil(err)
	assert.Equal(http.StatusOK, rw.Code)
}

func TestUnit_Healthcheck_WhenPingFails_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateTestEchoContextWithMethod(http.MethodGet)
	mp := &mockConnection{
		err: errDefault,
	}

	err := healthcheck(ctx, mp)

	assert.Nil(err)
	expectedResponse := `
	{
		"Cause": "some error",
		"Code": 260,
		"Message": "An unexpected error occurred"
	}`
	assert.JSONEq(expectedResponse, rw.Body.String())
}

func TestUnit_Healthcheck_WhenPingFails_SetsStatusToServiceUnavailable(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateTestEchoContextWithMethod(http.MethodGet)
	mp := &mockConnection{
		err: errDefault,
	}

	err := healthcheck(ctx, mp)

	assert.Nil(err)
	assert.Equal(http.StatusServiceUnavailable, rw.Code)
}

func (m *mockConnection) Ping(ctx context.Context) error {
	m.pingCalled++
	return m.err
}
