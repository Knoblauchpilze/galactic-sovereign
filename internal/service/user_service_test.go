package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
var defaultUser = persistence.User{
	Id:        defaultUserId,
	Email:     defaultUserEmail,
	Password:  defaultUserPassword,
	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func TestUserService_Create_CallsRepositoryCreate(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Create(context.Background(), defaultUserDtoRequest)

	assert.Equal(1, mur.createCalled)
	assert.Equal(defaultUserDtoRequest.Email, mur.createdUser.Email)
	assert.Equal(defaultUserDtoRequest.Password, mur.createdUser.Password)
}

func TestUserService_Create_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.Create(context.Background(), defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Create_ReturnsCreatedUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	actual, err := s.Create(context.Background(), defaultUserDtoRequest)

	assert.Nil(err)

	expected := communication.UserDtoResponse{
		Id:       defaultUser.Id,
		Email:    defaultUser.Email,
		Password: defaultUser.Password,

		CreatedAt: defaultUser.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUserService_Get_CallsRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Get(context.Background(), defaultUserId)

	assert.Equal(1, mur.getCalled)
}

func TestUserService_Get_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.Get(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Get_ReturnsUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	actual, err := s.Get(context.Background(), defaultUserId)

	assert.Nil(err)
	assert.Equal(defaultUserId, mur.getId)

	expected := communication.UserDtoResponse{
		Id:       defaultUser.Id,
		Email:    defaultUser.Email,
		Password: defaultUser.Password,

		CreatedAt: defaultUser.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUserService_List_CallsRepositoryList(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.List(context.Background())

	assert.Equal(1, mur.listCalled)
}

func TestUserService_List_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.List(context.Background())

	assert.Equal(errDefault, err)
}

func TestUserService_List_ReturnsAllUsers(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		ids: []uuid.UUID{
			uuid.MustParse("07e0eb01-c388-4148-bb45-286b09fdbc50"),
			uuid.MustParse("c759bc0d-ec75-4a55-b582-7b56b2e0710e"),
		},
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	actual, err := s.List(context.Background())

	assert.Nil(err)
	assert.Equal(mur.ids, actual)
}

func TestUserService_Update_CallsRepositoryGetAndUpdate(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Update(context.Background(), defaultUserId, defaultUserDtoRequest)

	assert.Equal(1, mur.getCalled)
	assert.Equal(defaultUserId, mur.getId)
	assert.Equal(1, mur.updateCalled)
}

func TestUserService_Update_WhenGetFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.Update(context.Background(), defaultUserId, defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Update_CallsUpdateWithUpdatedValues(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Update(context.Background(), defaultUserId, defaultUserDtoRequest)

	expected := persistence.User{
		Id:        defaultUser.Id,
		Email:     defaultUserDtoRequest.Email,
		Password:  defaultUserDtoRequest.Password,
		CreatedAt: defaultUser.CreatedAt,
		UpdatedAt: defaultUser.UpdatedAt,
	}
	assert.Equal(expected, mur.updatedUser)
}

func TestUserService_Update_WhenUpdateFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		updateErr: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.Update(context.Background(), defaultUserId, defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Update_ReturnsUpdatedUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	actual, err := s.Update(context.Background(), defaultUserId, defaultUserDtoRequest)

	assert.Nil(err)
	assert.Equal(defaultUserId, mur.getId)

	expected := communication.UserDtoResponse{
		Id:       defaultUser.Id,
		Email:    defaultUserDtoRequest.Email,
		Password: defaultUserDtoRequest.Password,

		CreatedAt: defaultUser.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUserService_Delete_CallsRepositoryDelete(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Delete(context.Background(), defaultUserId)

	assert.Equal(1, mur.deleteCalled)
	assert.Equal(1, mkr.deleteForUserCalled)
	assert.Equal(1, mar.deleteCalled)
	assert.Equal(1, mlr.deleteCalled)
}

func TestUserService_Delete_CallsTransactionClose(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Delete(context.Background(), defaultUserId)

	assert.Equal(1, mc.tx.closeCalled)
}

func TestUserService_Delete_WhenCreatingTransactionFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{
		err: errDefault,
	}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_DeletesTheRightKeys(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Delete(context.Background(), defaultUserId)

	assert.Equal(defaultUserId, mkr.deleteUserId)
}

func TestUserService_Delete_DeletesTheRightAcls(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Delete(context.Background(), defaultUserId)

	assert.Equal(defaultUserId, mar.inUserId)
}

func TestUserService_Delete_DeletesTheRightUserLimits(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Delete(context.Background(), defaultUserId)

	assert.Equal(defaultUserId, mlr.inUserId)
}

func TestUserService_Delete_WhenUserRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenApiKeyRepositoryDeleteFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{
		deleteErr: errDefault,
	}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenAclRepositoryDeleteFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mkr := &mockApiKeyRepository{}
	mar := &mockAclRepository{
		deleteErr: errDefault,
	}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenUserLimitRepositoryDeleteFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{}
	mlr := &mockUserLimitRepository{
		deleteErr: errDefault,
	}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenRepositoriesSucceeds_ExpectSuccess(t *testing.T) {
	assert := assert.New(t)

	mar := &mockAclRepository{}
	mkr := &mockApiKeyRepository{}
	mlr := &mockUserLimitRepository{}
	mur := &mockUserRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Acl:       mar,
		ApiKey:    mkr,
		UserLimit: mlr,
		User:      mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Nil(err)
}

func TestUserService_Login_FetchesUserByEmail(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Login(context.Background(), defaultUserDtoRequest)

	assert.Equal(1, mur.getEmailCalled)
	assert.Equal(defaultUserEmail, mur.getEmail)
}

func TestUserService_Login_WhenGetUserFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.Login(context.Background(), defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Login_WhenCredentialsAreWrong_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	userRequest := communication.UserDtoRequest{
		Email:    defaultUserEmail,
		Password: "not-the-right-password",
	}

	_, err := s.Login(context.Background(), userRequest)

	assert.True(errors.IsErrorWithCode(err, InvalidCredentials))
}

func TestUserService_Login_CreatesApiKeyForUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	c := Config{
		ApiKeyValidity: 1 * time.Hour,
	}
	s := NewUserService(c, mc, repos)

	s.Login(context.Background(), defaultUserDtoRequest)

	assert.Equal(1, mkr.createCalled)
	assert.Equal(defaultUserId, mkr.createdApiKey.ApiUser)
	expectedTime := time.Now().Add(59 * time.Minute)
	assert.True(expectedTime.Before(mkr.createdApiKey.ValidUntil))
}

func TestUserService_Login_WhenApiKeyCreationFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{
		createErr: errDefault,
	}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.Login(context.Background(), defaultUserDtoRequest)

	assert.Equal(1, mkr.createCalled)
	assert.Equal(errDefault, err)
}

func TestUserService_Login_ReturnsCreatedKey(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	c := Config{
		ApiKeyValidity: 1 * time.Hour,
	}
	s := NewUserService(c, mc, repos)

	actual, err := s.Login(context.Background(), defaultUserDtoRequest)

	assert.Nil(err)
	assert.Equal(mkr.createdApiKey.Key, actual.Key)
	assert.Equal(mkr.createdApiKey.ValidUntil, actual.ValidUntil)
}

func TestUserService_LoginById_FetchesUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.LoginById(context.Background(), defaultUserId)

	assert.Equal(1, mur.getCalled)
	assert.Equal(defaultUserId, mur.getId)
}

func TestUserService_LoginById_WhenGetUserFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.LoginById(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_LoginById_CreatesApiKeyForUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	c := Config{
		ApiKeyValidity: 1 * time.Hour,
	}
	s := NewUserService(c, mc, repos)

	s.LoginById(context.Background(), defaultUserId)

	assert.Equal(1, mkr.createCalled)
	assert.Equal(defaultUserId, mkr.createdApiKey.ApiUser)
	expectedTime := time.Now().Add(59 * time.Minute)
	assert.True(expectedTime.Before(mkr.createdApiKey.ValidUntil))
}

func TestUserService_LoginById_WhenApiKeyCreationFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		createErr: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	_, err := s.LoginById(context.Background(), defaultUserId)

	assert.Equal(1, mkr.createCalled)
	assert.Equal(errDefault, err)
}

func TestUserService_LoginById_ReturnsCreatedKey(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	c := Config{
		ApiKeyValidity: 1 * time.Hour,
	}
	s := NewUserService(c, mc, repos)

	actual, err := s.LoginById(context.Background(), defaultUserId)

	assert.Nil(err)
	assert.Equal(mkr.createdApiKey.Key, actual.Key)
	assert.Equal(mkr.createdApiKey.ValidUntil, actual.ValidUntil)
}

func TestUserService_Logout_FetchesGetUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	s.Logout(context.Background(), defaultUserId)

	assert.Equal(1, mur.getCalled)
	assert.Equal(defaultUserId, mur.getId)
}

func TestUserService_Logout_WhenGetUserFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Logout_WhenCreatingTransactionFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{
		err: errDefault,
	}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Logout_DeletesUserKeys(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Nil(err)
	assert.Equal(1, mkr.deleteForUserCalled)
	assert.Equal(defaultUserId, mkr.deleteUserId)
}

func TestUserService_Logout_WhenDeleteFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		deleteErr: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Equal(1, mkr.deleteForUserCalled)
	assert.Equal(errDefault, err)
}

func TestUserService_Logout_WhenAlreadyLoggedOut_StillLogsOut(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		ApiKey: mkr,
		User:   mur,
	}
	s := NewUserService(Config{}, mc, repos)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Equal(1, mkr.deleteForUserCalled)
	assert.Nil(err)
}
