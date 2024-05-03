package controller

import (
	"context"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/stretchr/testify/assert"
)

type mockConnectionPool struct {
	db.ConnectionPool

	pingCalled int
	err        error
}

func TestHealthCheckEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range HealthCheckEndpoints(&mockConnectionPool{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(1, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodGet])
}

func TestHealthcheck_CallsPoolPing(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mp := &mockConnectionPool{}

	healthcheck(mc, mp)

	assert.Equal(1, mp.pingCalled)
}

func TestHealthcheck_WhenPingSucceeds_SetsSatusToOk(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mp := &mockConnectionPool{}

	err := healthcheck(mc, mp)

	assert.Nil(err)
	assert.Equal(http.StatusOK, mc.status)
}

func TestHealthcheck_WhenPingFails_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mp := &mockConnectionPool{
		err: errDefault,
	}

	err := healthcheck(mc, mp)

	assert.Nil(err)
	assert.Equal(errDefault, mc.data)
}

func TestHealthcheck_WhenPingFails_SetsStatusToServiceUnavailable(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mp := &mockConnectionPool{
		err: errDefault,
	}

	err := healthcheck(mc, mp)

	assert.Nil(err)
	assert.Equal(http.StatusServiceUnavailable, mc.status)
}

func (m *mockConnectionPool) Ping(ctx context.Context) error {
	m.pingCalled++
	return m.err
}
