package controller

import (
	"context"
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/jackc/pgx/v5/pgconn"
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

	mc := dummyEchoContext()
	mp := &mockConnectionPool{}

	healthcheck(mc, mp)

	assert.Equal(1, mp.pingCalled)
}

func TestHealthcheck_WhenPingSucceeds_SetsSatusToOk(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateTestEchoContextAndResponseRecorder()
	mp := &mockConnectionPool{}

	err := healthcheck(ctx, mp)

	assert.Nil(err)
	assert.Equal(http.StatusOK, rw.Code)
}

func TestHealthcheck_WhenPingFails_PropagatesError(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateTestEchoContextAndResponseRecorder()
	mp := &mockConnectionPool{
		err: errDefault,
	}

	err := healthcheck(ctx, mp)

	assert.Nil(err)
	assert.Equal("{\"Code\":1,\"Message\":\"Healtcheck failed\",\"Cause\":\"some error\"}\n", rw.Body.String())
}

func TestHealthcheck_WhenPingFails_SetsStatusToServiceUnavailable(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateTestEchoContextAndResponseRecorder()
	mp := &mockConnectionPool{
		err: errDefault,
	}

	err := healthcheck(ctx, mp)

	assert.Nil(err)
	assert.Equal(http.StatusServiceUnavailable, rw.Code)
}

func TestHealthcheck_WhenPingFailsWithConnectionError_SetsStatusToServiceUnavailable(t *testing.T) {
	assert := assert.New(t)

	connErr := pgconn.ConnectError{
		Config: &pgconn.Config{},
	}

	ctx, rw := generateTestEchoContextAndResponseRecorder()
	mp := &mockConnectionPool{
		err: &connErr,
	}

	err := healthcheck(ctx, mp)

	assert.Nil(err)
	assert.Equal(http.StatusServiceUnavailable, rw.Code)
	assert.Equal("{\"Code\":1,\"Message\":\"Healtcheck failed\",\"Cause\":\"failed to connect to `host= user= database=`: \"}\n", rw.Body.String())
}

func (m *mockConnectionPool) Ping(ctx context.Context) error {
	m.pingCalled++
	return m.err
}
