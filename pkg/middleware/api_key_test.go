package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockApiKeyUserRepository struct {
	repositories.ApiKeyRepository

	apiKey persistence.ApiKey
	err    error

	getApiKeyCalled int
	inApiKey        uuid.UUID
}

func TestApiKey_WhenApiKeyNotDefined_Fails(t *testing.T) {
	assert := assert.New(t)
	ctx, rw := generateTestEchoContextWithMethod(http.MethodGet)
	mr := &mockApiKeyUserRepository{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(ctx)

	assert.False(*called)
	assert.Equal(http.StatusBadRequest, rw.Code)

	var actual string
	err := json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal("API key not found", actual)
}

func TestApiKey_WhenMoreThanOneApiKeyNotDefined_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(apiKeyHeaderKey, defaultApiKey1.String())
	req.Header.Add(apiKeyHeaderKey, defaultApiKey2.String())
	ctx, rw := generateTestEchoContextFromRequest(req)
	mr := &mockApiKeyUserRepository{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(ctx)

	assert.False(*called)
	assert.Equal(http.StatusBadRequest, rw.Code)
}

func TestApiKey_WhenApiKeyIsNotAUuid_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(apiKeyHeaderKey, "not-a-uuid")
	ctx, rw := generateTestEchoContextFromRequest(req)
	mr := &mockApiKeyUserRepository{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(ctx)

	assert.False(*called)
	assert.Equal(http.StatusBadRequest, rw.Code)

	var actual string
	err := json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal("API key has wrong format", actual)
}

func TestApiKey_AttemptsToFetchApiKeyFromRepository(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContextWithApiKeyAndMethod(http.MethodGet)
	mr := &mockApiKeyUserRepository{}
	next := createHandlerFuncReturning(nil)

	em := ApiKey(mr)
	callable := em(next)
	callable(ctx)

	assert.Equal(defaultApiKey1, mr.inApiKey)
}

func TestApiKey_WhenFetchingApiKeyFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)
	ctx, rw := generateTestEchoContextWithApiKeyAndMethod(http.MethodGet)
	mr := &mockApiKeyUserRepository{
		err: errDefault,
	}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(ctx)

	assert.False(*called)
	assert.Equal(http.StatusInternalServerError, rw.Code)
}

func TestApiKey_WhenApiKeyIsNotFound_SetsStatusToUnauthorized(t *testing.T) {
	assert := assert.New(t)
	ctx, rw := generateTestEchoContextWithApiKeyAndMethod(http.MethodGet)
	mr := &mockApiKeyUserRepository{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(ctx)

	assert.False(*called)
	assert.Equal(http.StatusUnauthorized, rw.Code)

	var actual string
	err := json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal("Invalid API key", actual)
}

var defaultApiKeyId = uuid.MustParse("5bda15f9-85f1-4700-867c-0a7cbda0f82c")

func TestApiKey_WhenApiKeyIsExpired_SetsStatusToUnauthorized(t *testing.T) {
	assert := assert.New(t)
	ctx, rw := generateTestEchoContextWithApiKeyAndMethod(http.MethodGet)
	mr := &mockApiKeyUserRepository{
		apiKey: persistence.ApiKey{
			Id:         defaultApiKeyId,
			Key:        defaultApiKey1,
			ApiUser:    defaultUuid,
			ValidUntil: time.Now().Add(-1 * time.Minute),
		},
	}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(ctx)

	assert.False(*called)
	assert.Equal(http.StatusUnauthorized, rw.Code)

	var actual string
	err := json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal("API key expired", actual)
}

func TestApiKey_WhenApiKeyIsValid_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := generateTestEchoContextWithApiKeyAndMethod(http.MethodGet)
	mr := &mockApiKeyUserRepository{
		apiKey: persistence.ApiKey{
			Id:         defaultApiKeyId,
			Key:        defaultApiKey1,
			ApiUser:    defaultUuid,
			ValidUntil: time.Now().Add(1 * time.Minute),
		},
	}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(ctx)

	assert.True(*called)
}

func (m *mockApiKeyUserRepository) GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error) {
	m.getApiKeyCalled++
	m.inApiKey = apiKey

	return m.apiKey, m.err
}
