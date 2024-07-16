package service

import (
	"context"
	"fmt"
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

var errDefault = fmt.Errorf("some error")
var defaultUserId = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultUserEmail = "some-user@provider.com"
var defaultUserPassword = "password"
var testDate = time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC)

var defaultUserDtoRequest = communication.UserDtoRequest{
	Email:    defaultUserEmail,
	Password: defaultUserPassword,
}
var defaultUpdatedUserDtoRequest = communication.UserDtoRequest{
	Email:    "my-updated-email",
	Password: "my-updated-password",
}
var defaultUser = persistence.User{
	Id:        defaultUserId,
	Email:     defaultUserEmail,
	Password:  defaultUserPassword,
	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func Test_UserService(t *testing.T) {
	s := ServiceTestSuite{
		generateValidRepositoriesMock: generateValidUserRepositoryMock,
		generateErrorRepositoriesMock: generateErrorUserRepositoryMock,

		errorTestCases: map[string]errorTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Create(ctx, defaultUserDtoRequest)
					return err
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Get(ctx, defaultUserId)
					return err
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.List(ctx)
					return err
				},
			},
			"update_getUserFails": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Update(ctx, defaultUserId, defaultUpdatedUserDtoRequest)
					return err
				},
			},
			"update": {
				generateErrorRepositoriesMock: func(err error) repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User: &mockUserRepository{
							updateErr: err,
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Update(ctx, defaultUserId, defaultUpdatedUserDtoRequest)
					return err
				},
			},
			"delete_aclFails": {
				generateErrorRepositoriesMock: func(err error) repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{
							deleteErr: err,
						},
						ApiKey:    &mockApiKeyRepository{},
						User:      &mockUserRepository{},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
			},
			"delete_apiKeyFails": {
				generateErrorRepositoriesMock: func(err error) repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{
							deleteErr: err,
						},
						User:      &mockUserRepository{},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
			},
			"delete_userFails": {
				generateErrorRepositoriesMock: func(err error) repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User: &mockUserRepository{
							err: err,
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
			},
			"delete_userLimitFails": {
				generateErrorRepositoriesMock: func(err error) repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User:   &mockUserRepository{},
						UserLimit: &mockUserLimitRepository{
							deleteErr: err,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
			},
			"login_userFails": {
				generateErrorRepositoriesMock: func(err error) repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User: &mockUserRepository{
							err: err,
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Login(ctx, defaultUserDtoRequest)
					return err
				},
			},
			"login_wrongCredentials": {
				generateErrorRepositoriesMock: func(err error) repositories.Repositories {
					return repositories.Repositories{
						Acl:       &mockAclRepository{},
						ApiKey:    &mockApiKeyRepository{},
						User:      &mockUserRepository{},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					userRequest := communication.UserDtoRequest{
						Email:    defaultUserEmail,
						Password: "not-the-right-password",
					}
					_, err := s.Login(ctx, userRequest)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, InvalidCredentials))
				},
			},
			"loginById_userFails": {
				generateErrorRepositoriesMock: func(err error) repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User: &mockUserRepository{
							err: err,
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.LoginById(ctx, defaultUserId)
					return err
				},
			},
		},

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Create(ctx, defaultUserDtoRequest)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultUserDtoRequest.Email, m.createdUser.Email)
					assert.Equal(defaultUserDtoRequest.Password, m.createdUser.Password)
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Get(ctx, defaultUserId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUserId, m.getId)
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.List(ctx)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"update": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Update(ctx, defaultUserId, defaultUpdatedUserDtoRequest)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUserId, m.getId)
					assert.Equal(1, m.updateCalled)
					expectedUpdatedUser := persistence.User{
						Id:        defaultUser.Id,
						Email:     "my-updated-email",
						Password:  "my-updated-password",
						CreatedAt: defaultUser.CreatedAt,
						UpdatedAt: defaultUser.UpdatedAt,
					}
					assert.Equal(expectedUpdatedUser, m.updatedUser)
				},
			},
			"delete_acl": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertAclRepoIsAMock(repos, assert)
					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUserId, m.inUserId)
				},
			},
			"delete_apiKey": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)
					assert.Equal(1, m.deleteForUserCalled)
					assert.Equal(defaultUserId, m.deleteUserId)
				},
			},
			"delete_user": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)
					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUserId, m.deleteId)
				},
			},
			"delete_userLimit": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserLimitRepoIsAMock(repos, assert)
					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUserId, m.inUserId)
				},
			},
			"delete_aclFails": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{
							deleteErr: errDefault,
						},
						ApiKey:    &mockApiKeyRepository{},
						User:      &mockUserRepository{},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertAclRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
				},
			},
			"delete_apiKeyFails": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{
							deleteErr: errDefault,
						},
						User:      &mockUserRepository{},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForUserCalled)
				},
			},
			"delete_userFails": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User: &mockUserRepository{
							err: errDefault,
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
				},
			},
			"delete_userLimitFails": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User:   &mockUserRepository{},
						UserLimit: &mockUserLimitRepository{
							deleteErr: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserLimitRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
				},
			},
			"login": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Login(ctx, defaultUserDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.getEmailCalled)
					assert.Equal(defaultUser.Email, m.getEmail)
				},
			},
			"login_apiKey": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					c := Config{
						ApiKeyValidity: 1 * time.Hour,
					}
					s := NewUserService(c, pool, repos)
					_, err := s.Login(ctx, defaultUserDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultUser.Id, m.createdApiKey.ApiUser)
					expectedTime := time.Now().Add(59 * time.Minute)
					assert.True(expectedTime.Before(m.createdApiKey.ValidUntil))
				},
			},
			"login_apiKeyFails": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{
							createErr: errDefault,
						},
						User: &mockUserRepository{
							user: defaultUser,
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.Login(ctx, defaultUserDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
				},
			},
			"loginById": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.LoginById(ctx, defaultUserId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUser.Id, m.getId)
				},
			},
			"loginById_apiKey": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					c := Config{
						ApiKeyValidity: 1 * time.Hour,
					}
					s := NewUserService(c, pool, repos)
					_, err := s.LoginById(ctx, defaultUserId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultUser.Id, m.createdApiKey.ApiUser)
					expectedTime := time.Now().Add(59 * time.Minute)
					assert.True(expectedTime.Before(m.createdApiKey.ValidUntil))
				},
			},
			"loginById_apiKeyFails": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{
							createErr: errDefault,
						},
						User: &mockUserRepository{
							user: defaultUser,
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					_, err := s.LoginById(ctx, defaultUserId)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
				},
			},
			"logout": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Logout(ctx, defaultUserId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUserId, m.getId)
				},
			},
			"logout_apiKey": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Logout(ctx, defaultUserId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForUserCalled)
					assert.Equal(defaultUserId, m.deleteUserId)
				},
			},
			"logout_userFails": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User: &mockUserRepository{
							err: errDefault,
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Logout(ctx, defaultUserId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUserRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
				},
			},
			"logout_apiKeyFails": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl: &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{
							deleteErr: errDefault,
						},
						User:      &mockUserRepository{},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Logout(ctx, defaultUserId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForUserCalled)
				},
			},
		},

		returnTestCases: map[string]returnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewUserService(Config{}, pool, repos)
					out, _ := s.Create(ctx, defaultUserDtoRequest)
					return out
				},

				expectedContent: communication.UserDtoResponse{
					Id:       defaultUser.Id,
					Email:    defaultUser.Email,
					Password: defaultUser.Password,

					CreatedAt: defaultUser.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewUserService(Config{}, pool, repos)
					out, _ := s.Get(ctx, defaultUserId)
					return out
				},

				expectedContent: communication.UserDtoResponse{
					Id:       defaultUser.Id,
					Email:    defaultUser.Email,
					Password: defaultUser.Password,

					CreatedAt: defaultUser.CreatedAt,
				},
			},
			"list": {
				generateValidRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Acl:    &mockAclRepository{},
						ApiKey: &mockApiKeyRepository{},
						User: &mockUserRepository{
							ids: []uuid.UUID{
								uuid.MustParse("07e0eb01-c388-4148-bb45-286b09fdbc50"),
								uuid.MustParse("c759bc0d-ec75-4a55-b582-7b56b2e0710e"),
							},
						},
						UserLimit: &mockUserLimitRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewUserService(Config{}, pool, repos)
					out, _ := s.List(ctx)
					return out
				},

				expectedContent: []uuid.UUID{
					uuid.MustParse("07e0eb01-c388-4148-bb45-286b09fdbc50"),
					uuid.MustParse("c759bc0d-ec75-4a55-b582-7b56b2e0710e"),
				},
			},
			"update": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewUserService(Config{}, pool, repos)
					out, _ := s.Update(ctx, defaultUserId, defaultUpdatedUserDtoRequest)
					return out
				},

				expectedContent: communication.UserDtoResponse{
					Id:        defaultUser.Id,
					Email:     defaultUpdatedUserDtoRequest.Email,
					Password:  defaultUpdatedUserDtoRequest.Password,
					CreatedAt: defaultUser.CreatedAt,
				},
			},
			"login": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					c := Config{
						ApiKeyValidity: 1 * time.Hour,
					}
					s := NewUserService(c, pool, repos)
					out, _ := s.Login(ctx, defaultUserDtoRequest)
					return out
				},

				verifyContent: func(in interface{}, repos repositories.Repositories, assert *require.Assertions) {
					actual := assertInterfaceIsAnApiKeyDtoResponse(in, assert)
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(actual.Key, m.createdApiKey.Key)
					assert.Equal(actual.ValidUntil, m.createdApiKey.ValidUntil)
				},
			},
			"loginById": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					c := Config{
						ApiKeyValidity: 1 * time.Hour,
					}
					s := NewUserService(c, pool, repos)
					out, _ := s.LoginById(ctx, defaultUserId)
					return out
				},

				verifyContent: func(in interface{}, repos repositories.Repositories, assert *require.Assertions) {
					actual := assertInterfaceIsAnApiKeyDtoResponse(in, assert)
					m := assertApiKeyRepoIsAMock(repos, assert)

					assert.Equal(actual.Key, m.createdApiKey.Key)
					assert.Equal(actual.ValidUntil, m.createdApiKey.ValidUntil)
				},
			},
		},

		transactionTestCases: map[string]transactionTestCase{
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Delete(ctx, defaultUserId)
				},
			},
			"logout": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUserService(Config{}, pool, repos)
					return s.Logout(ctx, defaultUserId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidUserRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Acl:    &mockAclRepository{},
		ApiKey: &mockApiKeyRepository{},
		User: &mockUserRepository{
			user: defaultUser,
		},
		UserLimit: &mockUserLimitRepository{},
	}
}

func generateErrorUserRepositoryMock(err error) repositories.Repositories {
	return repositories.Repositories{
		Acl:    &mockAclRepository{},
		ApiKey: &mockApiKeyRepository{},
		User: &mockUserRepository{
			err: err,
		},
		UserLimit: &mockUserLimitRepository{},
	}
}

func assertUserRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockUserRepository {
	m, ok := repos.User.(*mockUserRepository)
	if !ok {
		assert.Fail("Provided user repository is not a mock")
	}
	return m
}

func assertInterfaceIsAnApiKeyDtoResponse(in interface{}, assert *require.Assertions) communication.ApiKeyDtoResponse {
	out, ok := in.(communication.ApiKeyDtoResponse)
	if !ok {
		assert.Fail("Provided input is not an ApiKeyDtoResponse")
	}
	return out
}
