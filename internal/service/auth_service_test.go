package service

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

func Test_AuthService(t *testing.T) {
	s := ServicePoolTestSuite{
		generateRepositoriesMocks: generateValidAuthServiceMocks,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"authenticate_apiKey": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForKeyCalled)
					assert.Equal(defaultApiKey.Key, m.apiKeyId)
				},
			},
			"authenticate_getApiKeyFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						ApiKey: &mockApiKeyRepository{
							getErr: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},
				expectedError: errDefault,
			},
			"authenticate_apiKeyDoesNotExist": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						ApiKey: &mockApiKeyRepository{
							getErr: errors.NewCode(db.NoMatchingSqlRows),
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, UserNotAuthenticated))
				},
			},
			"authenticate_apiKeyExpired": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						ApiKey: &mockApiKeyRepository{
							apiKey: persistence.ApiKey{
								ValidUntil: time.Now().Add(-2 * time.Minute),
							},
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, AuthenticationExpired))
				},
			},
			"authenticate_aclForUser": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertAclRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForUserCalled)
					assert.Equal(defaultApiKey.ApiUser, m.inUserId)
				},
			},
			"authenticate_aclForUserFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{
							getForUserErr: errDefault,
						},
						ApiKey: generateApiKeyRepositoryWithValidApiKey(),
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertAclRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForUserCalled)
				},
			},
			"authenticate_acl": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertAclRepoIsAMock(repos, assert)

					assert.Equal(2, m.getCalled)
					assert.Equal(defaultAclIds[0], m.inAclIds[0])
					assert.Equal(defaultAclIds[1], m.inAclIds[1])
				},
			},
			"authenticate_aclFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{
							aclIds: defaultAclIds,
							getErr: errDefault,
						},
						ApiKey: generateApiKeyRepositoryWithValidApiKey(),
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertAclRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForUserCalled)
					assert.Equal(1, m.getCalled)
				},
			},
			"authenticate_userLimitForUser": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserLimitRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForUserCalled)
					assert.Equal(defaultApiKey.ApiUser, m.inUserId)
				},
			},
			"authenticate_userLimitForUserFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: generateApiKeyRepositoryWithValidApiKey(),
						UserLimit: &mockUserLimitRepository{
							getForUserErr: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserLimitRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForUserCalled)
				},
			},
			"authenticate_userLimit": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserLimitRepoIsAMock(repos, assert)

					assert.Equal(2, m.getCalled)
					assert.Equal(defaultUserLimitIds[0], m.inUserLimitIds[0])
					assert.Equal(defaultUserLimitIds[1], m.inUserLimitIds[1])
				},
			},
			"authenticate_userLimitFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: generateApiKeyRepositoryWithValidApiKey(),
						UserLimit: &mockUserLimitRepository{
							userLimitIds: defaultUserLimitIds,
							getErr:       errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserLimitRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForUserCalled)
					assert.Equal(1, m.getCalled)
				},
			},
		},

		returnTestCases: map[string]returnTestCase{
			"authenticate_acls": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{
							aclIds: defaultAclIds[:1],
							acl:    defaultAcl,
						},
						ApiKey:    generateApiKeyRepositoryWithValidApiKey(),
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewAuthService(pool, repos)
					out, _ := s.Authenticate(ctx, defaultApiKey.Key)
					return out
				},
				expectedContent: communication.AuthorizationDtoResponse{
					Acls: []communication.AclDtoResponse{
						{
							Id:          defaultAcl.Id,
							User:        defaultAcl.User,
							Resource:    defaultAcl.Resource,
							Permissions: defaultAcl.Permissions,
							CreatedAt:   defaultAcl.CreatedAt,
						},
					},
					Limits: []communication.LimitDtoResponse{},
				},
			},
			"authenticate_noAclsReturnsNotNilSlice": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{
							aclIds: nil,
						},
						ApiKey:    generateApiKeyRepositoryWithValidApiKey(),
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewAuthService(pool, repos)
					out, _ := s.Authenticate(ctx, defaultApiKey.Key)
					return out
				},
				expectedContent: communication.AuthorizationDtoResponse{
					Acls:   []communication.AclDtoResponse{},
					Limits: []communication.LimitDtoResponse{},
				},
			},
			"authenticate_userLimits": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: generateApiKeyRepositoryWithValidApiKey(),
						UserLimit: &mockUserLimitRepository{
							userLimitIds: defaultUserLimitIds[:1],
							userLimit:    defaultUserLimit,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewAuthService(pool, repos)
					out, _ := s.Authenticate(ctx, defaultApiKey.Key)
					return out
				},

				expectedContent: communication.AuthorizationDtoResponse{
					Acls: []communication.AclDtoResponse{},
					Limits: []communication.LimitDtoResponse{
						{
							Name:  defaultUserLimit.Limits[0].Name,
							Value: defaultUserLimit.Limits[0].Value,
						},
					},
				},
			},
			"authenticate_noUserLimitsReturnsNotNilSlice": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: generateApiKeyRepositoryWithValidApiKey(),
						UserLimit: &mockUserLimitRepository{
							userLimitIds: nil,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewAuthService(pool, repos)
					out, _ := s.Authenticate(ctx, defaultApiKey.Key)
					return out
				},

				expectedContent: communication.AuthorizationDtoResponse{
					Acls:   []communication.AclDtoResponse{},
					Limits: []communication.LimitDtoResponse{},
				},
			},
		},

		transactionTestCases: map[string]transactionTestCase{
			"authenticate": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewAuthService(pool, repos)
					_, err := s.Authenticate(ctx, defaultApiKey.Key)
					return err
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidAuthServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		Acl: &mockAclRepository{
			aclIds: defaultAclIds,
		},
		ApiKey: generateApiKeyRepositoryWithValidApiKey(),
		UserLimit: &mockUserLimitRepository{
			userLimitIds: defaultUserLimitIds,
		},
	}
}

func generateApiKeyRepositoryWithValidApiKey() *mockApiKeyRepository {
	return &mockApiKeyRepository{
		apiKey: persistence.ApiKey{
			Id:         defaultApiKey.Id,
			Key:        defaultApiKey.Key,
			ApiUser:    defaultApiKey.ApiUser,
			ValidUntil: time.Now().Add(1 * time.Hour),
		},
	}
}

func assertAclRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockAclRepository {
	m, ok := repos.Acl.(*mockAclRepository)
	if !ok {
		assert.Fail("Provided acl repository is not a mock")
	}
	return m
}

func assertApiKeyRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockApiKeyRepository {
	m, ok := repos.ApiKey.(*mockApiKeyRepository)
	if !ok {
		assert.Fail("Provided api key repository is not a mock")
	}
	return m
}

func assertUserLimitRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockUserLimitRepository {
	m, ok := repos.UserLimit.(*mockUserLimitRepository)
	if !ok {
		assert.Fail("Provided user limit repository is not a mock")
	}
	return m
}
