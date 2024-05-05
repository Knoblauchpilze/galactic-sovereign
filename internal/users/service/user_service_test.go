package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	repositories.UserRepository

	user      persistence.User
	ids       []uuid.UUID
	err       error
	updateErr error

	createCalled int
	createdUser  persistence.User
	getCalled    int
	getId        uuid.UUID
	listCalled   int
	updateCalled int
	updatedUser  persistence.User
	deleteCalled int
}

type mockApiKeyRepository struct {
	repositories.ApiKeyRepository

	apiKeyIds []uuid.UUID
	createErr error
	getErr    error
	deleteErr error

	createCalled       int
	createdApiKey      persistence.ApiKey
	getForUserCalled   int
	getForUserTxCalled int
	userId             uuid.UUID
	deleteCalled       int
	deleteTxCalled     int
	deleteIds          []uuid.UUID
}

type mockConnectionPool struct {
	db.ConnectionPool

	tx  mockTransaction
	err error
}

type mockTransaction struct {
	db.Transaction

	timeStamp time.Time

	closeCalled int
}

var errDefault = fmt.Errorf("some error")
var defaultUserId = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultApiKeyId = uuid.MustParse("cc1742fa-77b4-4f5f-ac92-058c2e47a5d6")
var testDate = time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC)

var defaultUserDtoRequest = communication.UserDtoRequest{
	Email:    "some-user@provider.com",
	Password: "password",
}
var defaultUser = persistence.User{
	Id:        defaultUserId,
	Email:     "e.mail@domain.com",
	Password:  "password",
	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func TestUserService_Create_CallsRepositoryCreate(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	s.Create(context.Background(), defaultUserDtoRequest)

	assert.Equal(1, mur.createCalled)
}

func TestUserService_Create_WhenUserRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

	actual, err := s.List(context.Background())

	assert.Nil(err)
	assert.Equal(mur.ids, actual)
}

func TestUserService_Update_CallsRepositoryGetAndUpdate(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

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

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	s.Delete(context.Background(), defaultUserId)

	assert.Equal(1, mur.deleteCalled)
	assert.Equal(1, mkr.deleteTxCalled)
}

func TestUserService_Delete_CallsTransactionClose(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_FetchesUsersKeys(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	s.Delete(context.Background(), defaultUserId)

	assert.Equal(1, mkr.getForUserTxCalled)
	assert.Equal(defaultUserId, mkr.userId)
}

func TestUserService_Delete_DeletesTheRightKeys(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		apiKeyIds: []uuid.UUID{defaultApiKeyId},
	}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	s.Delete(context.Background(), defaultUserId)

	assert.Equal(mkr.apiKeyIds, mkr.deleteIds)
}

func TestUserService_Delete_WhenUserRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenApiKeyRepositoryGetFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		getErr: errDefault,
	}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenApiKeyRepositoryDeleteFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		deleteErr: errDefault,
	}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenRepositoriesSucceeds_ExpectSuccess(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Delete(context.Background(), defaultUserId)

	assert.Nil(err)
}

func TestUserService_Login_FetchesGetUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	s.Login(context.Background(), defaultUserId)

	assert.Equal(1, mur.getCalled)
	assert.Equal(defaultUserId, mur.getId)
}

func TestUserService_Login_WhenGetUserFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	_, err := s.Login(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Login_CreatesApiKeyForUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	c := Config{
		ApiKeyValidity: 1 * time.Hour,
	}
	s := NewUserService(c, mc, mur, mkr)

	s.Login(context.Background(), defaultUserId)

	assert.Equal(1, mkr.createCalled)
	assert.Equal(defaultUserId, mkr.createdApiKey.ApiUser)
	expectedTime := time.Now().Add(59 * time.Minute)
	assert.True(expectedTime.Before(mkr.createdApiKey.ValidUntil))
}

func TestUserService_Login_WhenApiKeyCreationFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		createErr: errDefault,
	}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	_, err := s.Login(context.Background(), defaultUserId)

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
	c := Config{
		ApiKeyValidity: 1 * time.Hour,
	}
	s := NewUserService(c, mc, mur, mkr)

	actual, err := s.Login(context.Background(), defaultUserId)

	assert.Nil(err)
	assert.Equal(mkr.createdApiKey.Key, actual.Key)
	assert.Equal(mkr.createdApiKey.ValidUntil, actual.ValidUntil)
}

func TestUserService_Logout_FetchesGetUser(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

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
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Logout_FetchesUserKeys(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	s.Logout(context.Background(), defaultUserId)

	assert.Equal(1, mkr.getForUserCalled)
	assert.Equal(defaultUserId, mkr.userId)
}

func TestUserService_Logout_WhenGetForUserFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		getErr: errDefault,
	}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserService_Logout_DeletesUserKeys(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		apiKeyIds: []uuid.UUID{defaultApiKeyId},
	}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Nil(err)
	assert.Equal(1, mkr.deleteCalled)
	assert.Equal([]uuid.UUID{defaultApiKeyId}, mkr.deleteIds)
}

func TestUserService_Logout_WhenDeleteFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{
		deleteErr: errDefault,
	}
	mc := &mockConnectionPool{}
	s := NewUserService(Config{}, mc, mur, mkr)

	err := s.Logout(context.Background(), defaultUserId)

	assert.Equal(1, mkr.deleteCalled)
	assert.Equal(errDefault, err)
}

func (m *mockUserRepository) Create(ctx context.Context, user persistence.User) (persistence.User, error) {
	m.createCalled++
	m.createdUser = user
	return m.user, m.err
}

func (m *mockUserRepository) Get(ctx context.Context, id uuid.UUID) (persistence.User, error) {
	m.getCalled++
	m.getId = id
	return m.user, m.err
}

func (m *mockUserRepository) List(ctx context.Context) ([]uuid.UUID, error) {
	m.listCalled++
	return m.ids, m.err
}

func (m *mockUserRepository) Update(ctx context.Context, user persistence.User) (persistence.User, error) {
	m.updateCalled++
	m.updatedUser = user
	return m.updatedUser, m.updateErr
}

func (m *mockUserRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	return m.err
}

func (m *mockApiKeyRepository) Create(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error) {
	m.createCalled++
	m.createdApiKey = apiKey
	return apiKey, m.createErr
}

func (m *mockApiKeyRepository) GetForUser(ctx context.Context, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserCalled++
	m.userId = user
	return m.apiKeyIds, m.getErr
}

func (m *mockApiKeyRepository) GetForUserTx(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserTxCalled++
	m.userId = user
	return m.apiKeyIds, m.getErr
}

func (m *mockApiKeyRepository) Delete(ctx context.Context, ids []uuid.UUID) error {
	m.deleteCalled++
	m.deleteIds = ids
	return m.deleteErr
}

func (m *mockApiKeyRepository) DeleteTx(ctx context.Context, tx db.Transaction, ids []uuid.UUID) error {
	m.deleteTxCalled++
	m.deleteIds = ids
	return m.deleteErr
}

func (m *mockConnectionPool) StartTransaction(ctx context.Context) (db.Transaction, error) {
	return &m.tx, m.err
}

func (m *mockTransaction) Close(ctx context.Context) {
	m.closeCalled++
}

func (m *mockTransaction) TimeStamp() time.Time {
	return m.timeStamp
}
