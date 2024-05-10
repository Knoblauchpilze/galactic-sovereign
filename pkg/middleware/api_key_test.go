package middleware

import (
	"context"
	"net/http"
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
	mc := newMockEchoContext(http.StatusOK)
	mr := &mockApiKeyUserRepository{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(mc)

	assert.False(*called)
	assert.Equal(http.StatusBadRequest, mc.reportedCode)
	actual, ok := mc.jsonContent.(string)
	assert.True(ok)
	assert.Equal("API key not found", actual)
}

var defaultApiKey1 = uuid.MustParse("f847c203-1c56-43ad-9ac1-46f27d650917")
var defaultApiKey2 = uuid.MustParse("297d3309-d88b-4b83-8d82-9c6aae8a9d7a")

func TestApiKey_WhenMoreThanOneApiKeyNotDefined_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {
			defaultApiKey1.String(),
			defaultApiKey2.String(),
		},
	}
	mr := &mockApiKeyUserRepository{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(mc)

	assert.False(*called)
	assert.Equal(http.StatusBadRequest, mc.reportedCode)
}

func TestApiKey_WhenApiKeyIsNotAUuid_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {"not-a-uuid"},
	}
	mr := &mockApiKeyUserRepository{}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(mc)

	assert.False(*called)
	assert.Equal(http.StatusBadRequest, mc.reportedCode)
	actual, ok := mc.jsonContent.(string)
	assert.True(ok)
	assert.Equal("API key has wrong format", actual)
}

func TestApiKey_AttemptsToFetchApiKeyFromRepository(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {defaultApiKey1.String()},
	}
	mr := &mockApiKeyUserRepository{}
	next := createHandlerFuncReturning(nil)

	em := ApiKey(mr)
	callable := em(next)
	callable(mc)

	assert.Equal(defaultApiKey1, mr.inApiKey)
}

func TestApiKey_WhenFetchingApiKeyFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {defaultApiKey1.String()},
	}
	mr := &mockApiKeyUserRepository{
		err: errDefault,
	}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(mc)

	assert.False(*called)
	assert.Equal(http.StatusInternalServerError, mc.reportedCode)
}

func TestApiKey_WhenApiKeyIsNotFound_SetsStatusToUnauthorized(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {defaultApiKey1.String()},
	}
	mr := &mockApiKeyUserRepository{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}
	next, called := createHandlerFuncWithCalledBoolean()

	em := ApiKey(mr)
	callable := em(next)
	callable(mc)

	assert.False(*called)
	assert.Equal(http.StatusUnauthorized, mc.reportedCode)
	actual, ok := mc.jsonContent.(string)
	assert.True(ok)
	assert.Equal("Invalid API key", actual)
}

var defaultApiKeyId = uuid.MustParse("5bda15f9-85f1-4700-867c-0a7cbda0f82c")

func TestApiKey_WhenApiKeyIsExpired_SetsStatusToUnauthorized(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {defaultApiKey1.String()},
	}
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
	callable(mc)

	assert.False(*called)
	assert.Equal(http.StatusUnauthorized, mc.reportedCode)
	actual, ok := mc.jsonContent.(string)
	assert.True(ok)
	assert.Equal("API key expired", actual)
}

func TestApiKey_WhenApiKeyIsValid_CallsNextMiddleware(t *testing.T) {
	assert := assert.New(t)
	mc := newMockEchoContext(http.StatusOK)
	mc.request.Header = map[string][]string{
		apiKeyHeaderKey: {defaultApiKey1.String()},
	}
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
	callable(mc)

	assert.True(*called)
}

func (m *mockApiKeyUserRepository) GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error) {
	m.getApiKeyCalled++
	m.inApiKey = apiKey

	return m.apiKey, m.err
}