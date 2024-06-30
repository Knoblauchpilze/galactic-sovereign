package service

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultApiKey = persistence.ApiKey{
	Id:      uuid.MustParse("1fe10052-0d94-4127-ab12-ef25038c689e"),
	Key:     uuid.MustParse("60873f25-54b0-45e9-b920-2bd7d82cd438"),
	ApiUser: defaultUserId,
}

var defaultAclIds = []uuid.UUID{
	uuid.MustParse("e4667ff7-1ed5-4ce0-ac06-10668eab8a70"),
	uuid.MustParse("0bfa5491-b0df-4976-ac8e-c916fb750874"),
}
var defaultUserLimitIds = []uuid.UUID{
	uuid.MustParse("cbdd762b-1c20-4992-90fd-f190494e5525"),
	uuid.MustParse("f256f190-bdab-443e-9b83-f6a5b992f632"),
}

var defaultAcl = persistence.Acl{
	Id:   defaultAclIds[0],
	User: defaultUserId,

	Resource:    "my-resource",
	Permissions: []string{"GET"},

	CreatedAt: time.Date(2024, 06, 28, 15, 11, 20, 651387237, time.UTC),
	UpdatedAt: time.Date(2024, 06, 28, 15, 11, 22, 651387237, time.UTC),
}
var defaultUserLimit = persistence.UserLimit{
	Id:   defaultUserLimitIds[0],
	Name: "my-limit",
	User: defaultUserId,

	Limits: []persistence.Limit{
		{
			Id: uuid.MustParse("2efb5dd3-9951-4afe-8a57-1a31cda39373"),

			Name:  "my-name-1",
			Value: "my-value-1",

			CreatedAt: time.Date(2024, 06, 28, 15, 19, 25, 651387237, time.UTC),
			UpdatedAt: time.Date(2024, 06, 28, 15, 19, 27, 651387237, time.UTC),
		},
	},

	CreatedAt: time.Date(2024, 06, 28, 15, 18, 10, 651387237, time.UTC),
	UpdatedAt: time.Date(2024, 06, 28, 15, 18, 12, 651387237, time.UTC),
}

func TestAuthService_Authenticate_FetchesTheApiKeyDetails(t *testing.T) {
	assert := assert.New(t)

	repos, _, mockApi, _, _ := createMockRepositories()
	s := NewAuthService(&mockConnectionPool{}, repos)

	s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(1, mockApi.getForKeyCalled)
	assert.Equal(defaultApiKey.Id, mockApi.apiKeyId)
}

func TestAuthService_Authenticate_WhenApiKeyDoesNotExist_ReturnsNotAuthenticated(t *testing.T) {
	assert := assert.New(t)

	repos, _, mockApi, _, _ := createMockRepositories()
	mockApi.getErr = errors.NewCode(db.NoMatchingSqlRows)
	s := NewAuthService(&mockConnectionPool{}, repos)

	_, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.True(errors.IsErrorWithCode(err, UserNotAuthenticated))
}

func TestAuthService_Authenticate_WhenFetchingApiKeyFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	repos, _, mockApi, _, _ := createMockRepositories()
	mockApi.getErr = errDefault
	s := NewAuthService(&mockConnectionPool{}, repos)

	_, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(errDefault, err)
}

func TestAuthService_Authenticate_WhenApiKeyExpired_ReturnsAuthenticationExpired(t *testing.T) {
	assert := assert.New(t)

	repos, _, mockApi, _, _ := createMockRepositories()
	mockApi.apiKey = persistence.ApiKey{
		ValidUntil: time.Now().Add(-2 * time.Minute),
	}
	s := NewAuthService(&mockConnectionPool{}, repos)

	_, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.True(errors.IsErrorWithCode(err, AuthenticationExpired))
}

func TestAuthService_Authenticate_WhenTransactionFailsToBeCreated_ExpectError(t *testing.T) {
	assert := assert.New(t)

	repos, _, _, _, _ := createMockRepositoriesWithValidApiKey()
	mockPool := &mockConnectionPool{
		err: errDefault,
	}
	s := NewAuthService(mockPool, repos)

	_, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(errDefault, err)
}

func TestAuthService_Authenticate_FetchesAclsForUser(t *testing.T) {
	assert := assert.New(t)

	repos, mockAcl, _, _, _ := createMockRepositoriesWithValidApiKey()
	s := NewAuthService(&mockConnectionPool{}, repos)

	s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(1, mockAcl.getForUserCalled)
	assert.Equal(defaultApiKey.ApiUser, mockAcl.inUserId)
}

func TestAuthService_Authenticate_WhenFetchingAclsForUserFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	repos, mockAcl, _, _, _ := createMockRepositoriesWithValidApiKey()
	mockAcl.getForUserErr = errDefault
	s := NewAuthService(&mockConnectionPool{}, repos)

	_, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(errDefault, err)
}

func TestAuthService_Authenticate_FetchesAclsFromRepository(t *testing.T) {
	assert := assert.New(t)

	repos, mockAcl, _, _, _ := createMockRepositoriesWithValidApiKey()
	mockAcl.aclIds = defaultAclIds
	s := NewAuthService(&mockConnectionPool{}, repos)

	s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(2, mockAcl.getCalled)
	assert.Equal(defaultAclIds[0], mockAcl.inAclIds[0])
	assert.Equal(defaultAclIds[1], mockAcl.inAclIds[1])
}

func TestAuthService_Authenticate_WhenFetchingOfAclsFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	repos, mockAcl, _, _, _ := createMockRepositoriesWithValidApiKey()
	mockAcl.aclIds = defaultAclIds
	mockAcl.getErr = errDefault
	s := NewAuthService(&mockConnectionPool{}, repos)

	_, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(1, mockAcl.getCalled)
	assert.Equal(errDefault, err)
}

func TestAuthService_Authenticate_ReturnsExpectedAcls(t *testing.T) {
	assert := assert.New(t)

	repos, mockAcl, _, _, _ := createMockRepositoriesWithValidApiKey()
	mockAcl.aclIds = defaultAclIds[:1]
	mockAcl.acl = defaultAcl
	s := NewAuthService(&mockConnectionPool{}, repos)

	actual, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Nil(err)
	assert.Equal(1, len(actual.Acls))
	assert.Equal(0, len(actual.Limits))

	assert.Equal(defaultAcl.Id, actual.Acls[0].Id)
	assert.Equal(defaultAcl.User, actual.Acls[0].User)
	assert.Equal(defaultAcl.Resource, actual.Acls[0].Resource)
	assert.Equal(defaultAcl.Permissions, actual.Acls[0].Permissions)
	assert.Equal(defaultAcl.CreatedAt, actual.Acls[0].CreatedAt)
}

func TestAuthService_Authenticate_WhenUserDoesNotHaveAcls_ExpectNonNilSliceReturned(t *testing.T) {
	assert := assert.New(t)

	repos, mockAcl, _, _, _ := createMockRepositoriesWithValidApiKey()
	mockAcl.aclIds = nil
	s := NewAuthService(&mockConnectionPool{}, repos)

	actual, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Nil(err)
	assert.NotNil(actual.Acls)
	assert.Equal(0, len(actual.Acls))
}

func TestAuthService_Authenticate_FetchesUserLimitsForUser(t *testing.T) {
	assert := assert.New(t)

	repos, _, _, _, mockUserLimit := createMockRepositoriesWithValidApiKey()
	s := NewAuthService(&mockConnectionPool{}, repos)

	s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(1, mockUserLimit.getForUserCalled)
	assert.Equal(defaultApiKey.ApiUser, mockUserLimit.inUserId)
}

func TestAuthService_Authenticate_WhenFetchingUserLimitsForUserFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	repos, _, _, _, mockUserLimit := createMockRepositoriesWithValidApiKey()
	mockUserLimit.getForUserErr = errDefault
	s := NewAuthService(&mockConnectionPool{}, repos)

	_, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(errDefault, err)
}

func TestAuthService_Authenticate_FetchesUserLimitsFromRepository(t *testing.T) {
	assert := assert.New(t)

	repos, _, _, _, mockUserLimit := createMockRepositoriesWithValidApiKey()
	mockUserLimit.userLimitIds = defaultUserLimitIds
	s := NewAuthService(&mockConnectionPool{}, repos)

	s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(2, mockUserLimit.getCalled)
	assert.Equal(defaultUserLimitIds[0], mockUserLimit.inUserLimitIds[0])
	assert.Equal(defaultUserLimitIds[1], mockUserLimit.inUserLimitIds[1])
}

func TestAuthService_Authenticate_WhenFetchingOfUserLimitsFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	repos, _, _, _, mockUserLimit := createMockRepositoriesWithValidApiKey()
	mockUserLimit.userLimitIds = defaultUserLimitIds
	mockUserLimit.getErr = errDefault
	s := NewAuthService(&mockConnectionPool{}, repos)

	_, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Equal(1, mockUserLimit.getCalled)
	assert.Equal(errDefault, err)
}

func TestAuthService_Authenticate_ReturnsExpectedUserLimits(t *testing.T) {
	assert := assert.New(t)

	repos, _, _, _, mockUserLimit := createMockRepositoriesWithValidApiKey()
	mockUserLimit.userLimitIds = defaultUserLimitIds[:1]
	mockUserLimit.userLimit = defaultUserLimit
	s := NewAuthService(&mockConnectionPool{}, repos)

	actual, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Nil(err)
	assert.Equal(0, len(actual.Acls))
	assert.Equal(1, len(actual.Limits))

	assert.Equal(defaultUserLimit.Limits[0].Name, actual.Limits[0].Name)
	assert.Equal(defaultUserLimit.Limits[0].Value, actual.Limits[0].Value)
}

func TestAuthService_Authenticate_WhenUserDoesNotHaveUserLimits_ExpectNonNilSliceReturned(t *testing.T) {
	assert := assert.New(t)

	repos, _, _, _, mockUserLimit := createMockRepositoriesWithValidApiKey()
	mockUserLimit.userLimitIds = nil
	s := NewAuthService(&mockConnectionPool{}, repos)

	actual, err := s.Authenticate(context.Background(), defaultApiKey.Id)

	assert.Nil(err)
	assert.NotNil(actual.Limits)
	assert.Equal(0, len(actual.Limits))
}

func createMockRepositories() (repositories.Repositories, *mockAclRepository, *mockApiKeyRepository, *mockUserRepository, *mockUserLimitRepository) {
	mockAcl := &mockAclRepository{}
	mockApi := &mockApiKeyRepository{}
	mockUser := &mockUserRepository{}
	mockUserLimit := &mockUserLimitRepository{}

	repos := createAllRepositories(mockAcl, mockApi, mockUser, mockUserLimit)

	return repos, mockAcl, mockApi, mockUser, mockUserLimit
}

func createMockRepositoriesWithValidApiKey() (repositories.Repositories, *mockAclRepository, *mockApiKeyRepository, *mockUserRepository, *mockUserLimitRepository) {
	mockAcl := &mockAclRepository{}
	mockApi := &mockApiKeyRepository{
		apiKey: defaultApiKey,
	}
	mockUser := &mockUserRepository{}
	mockUserLimit := &mockUserLimitRepository{}

	mockApi.apiKey.ValidUntil = time.Now().Add(1 * time.Hour)

	repos := createAllRepositories(mockAcl, mockApi, mockUser, mockUserLimit)

	return repos, mockAcl, mockApi, mockUser, mockUserLimit
}
