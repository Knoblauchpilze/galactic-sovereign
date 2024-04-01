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
}

type mockConnectionPool struct {
	db.ConnectionPool
}

var errDefault = fmt.Errorf("some error")
var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")

// TODO: Restore or delete this.
// var defaultApiKey = uuid.MustParse("cc1742fa-77b4-4f5f-ac92-058c2e47a5d6")

var defaultUserDtoRequest = communication.UserDtoRequest{
	Email:    "some-user@provider.com",
	Password: "password",
}
var defaultUser = persistence.User{
	Id:        defaultUuid,
	Email:     "e.mail@domain.com",
	Password:  "password",
	CreatedAt: time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC),
	UpdatedAt: time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC),
}

func TestUserService_Create_CallsRepositoryCreate(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	s.Create(context.Background(), defaultUserDtoRequest)

	assert.Equal(1, mur.createCalled)
}

func TestUserService_Create_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	_, err := s.Create(context.Background(), defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

// TODO: Fix this test to do what it says.
func TestUserService_Create_ReturnsCreatedUserIncludingApiKey(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	actual, err := s.Create(context.Background(), defaultUserDtoRequest)

	assert.Nil(err)

	expected := communication.UserDtoResponse{
		Id:       defaultUser.Id,
		Email:    defaultUser.Email,
		Password: defaultUser.Password,

		//ApiKeys: defaultUser.ApiKeys,
		ApiKeys: []uuid.UUID{},

		CreatedAt: defaultUser.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUserService_Get_CallsRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	s.Get(context.Background(), defaultUuid)

	assert.Equal(1, mur.getCalled)
}

func TestUserService_Get_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	_, err := s.Get(context.Background(), defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserService_Get_ReturnsUserOmitingApiKey(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	actual, err := s.Get(context.Background(), defaultUuid)

	assert.Nil(err)
	assert.Equal(defaultUuid, mur.getId)

	expected := communication.UserDtoResponse{
		Id:       defaultUser.Id,
		Email:    defaultUser.Email,
		Password: defaultUser.Password,

		ApiKeys: []uuid.UUID{},

		CreatedAt: defaultUser.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUserService_List_CallsRepositoryList(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

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
	s := NewUserService(mc, mur, mkr)

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
	s := NewUserService(mc, mur, mkr)

	actual, err := s.List(context.Background())

	assert.Nil(err)
	assert.Equal(mur.ids, actual)
}

func TestUserService_Update_CallsRepositoryGetAndUpdate(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	assert.Equal(1, mur.getCalled)
	assert.Equal(defaultUuid, mur.getId)
	assert.Equal(1, mur.updateCalled)
}

func TestUserService_Update_WhenGetFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	_, err := s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Update_CallsUpdateWithUpdatedValues(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

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
	s := NewUserService(mc, mur, mkr)

	_, err := s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Update_ReturnsUpdatedUserOmitingApiKey(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		user: defaultUser,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	actual, err := s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	assert.Nil(err)
	assert.Equal(defaultUuid, mur.getId)

	expected := communication.UserDtoResponse{
		Id:       defaultUser.Id,
		Email:    defaultUserDtoRequest.Email,
		Password: defaultUserDtoRequest.Password,

		ApiKeys: []uuid.UUID{},

		CreatedAt: defaultUser.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUserService_Delete_CallsRepositoryDelete(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	s.Delete(context.Background(), defaultUuid)

	assert.Equal(1, mur.deleteCalled)
}

func TestUserService_Delete_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{
		err: errDefault,
	}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	err := s.Delete(context.Background(), defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenRepositorySucceeds_ExpectSuccess(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUserRepository{}
	mkr := &mockApiKeyRepository{}
	mc := &mockConnectionPool{}
	s := NewUserService(mc, mur, mkr)

	err := s.Delete(context.Background(), defaultUuid)

	assert.Nil(err)
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

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	return m.err
}
