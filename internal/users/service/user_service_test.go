package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
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

var errDefault = fmt.Errorf("some error")
var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultApiKey = uuid.MustParse("cc1742fa-77b4-4f5f-ac92-058c2e47a5d6")

var defaultUserDtoRequest = communication.UserDtoRequest{
	Email:    "some-user@provider.com",
	Password: "password",
}
var defaultUser = persistence.User{
	Id:        defaultUuid,
	Email:     "e.mail@domain.com",
	Password:  "password",
	ApiKeys:   []uuid.UUID{defaultApiKey},
	CreatedAt: time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC),
	UpdatedAt: time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC),
}

func TestUserService_Create_CallsRepositoryCreate(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{}
	s := NewUserService(mr)

	s.Create(context.Background(), defaultUserDtoRequest)

	assert.Equal(1, mr.createCalled)
}

func TestUserService_Create_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		err: errDefault,
	}
	s := NewUserService(mr)

	_, err := s.Create(context.Background(), defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Create_ReturnsCreatedUserIncludingApiKey(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		user: defaultUser,
	}
	s := NewUserService(mr)

	actual, err := s.Create(context.Background(), defaultUserDtoRequest)

	assert.Nil(err)

	expected := communication.UserDtoResponse{
		Id:       defaultUser.Id,
		Email:    defaultUser.Email,
		Password: defaultUser.Password,

		ApiKeys: defaultUser.ApiKeys,

		CreatedAt: defaultUser.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUserService_Get_CallsRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{}
	s := NewUserService(mr)

	s.Get(context.Background(), defaultUuid)

	assert.Equal(1, mr.getCalled)
}

func TestUserService_Get_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		err: errDefault,
	}
	s := NewUserService(mr)

	_, err := s.Get(context.Background(), defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserService_Get_ReturnsUserOmitingApiKey(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		user: defaultUser,
	}
	s := NewUserService(mr)

	actual, err := s.Get(context.Background(), defaultUuid)

	assert.Nil(err)
	assert.Equal(defaultUuid, mr.getId)

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

	mr := &mockUserRepository{}
	s := NewUserService(mr)

	s.List(context.Background())

	assert.Equal(1, mr.listCalled)
}

func TestUserService_List_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		err: errDefault,
	}
	s := NewUserService(mr)

	_, err := s.List(context.Background())

	assert.Equal(errDefault, err)
}

func TestUserService_List_ReturnsAllUsers(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		ids: []uuid.UUID{
			uuid.MustParse("07e0eb01-c388-4148-bb45-286b09fdbc50"),
			uuid.MustParse("c759bc0d-ec75-4a55-b582-7b56b2e0710e"),
		},
	}
	s := NewUserService(mr)

	actual, err := s.List(context.Background())

	assert.Nil(err)
	assert.Equal(mr.ids, actual)
}

func TestUserService_Update_CallsRepositoryGetAndUpdate(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{}
	s := NewUserService(mr)

	s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	assert.Equal(1, mr.getCalled)
	assert.Equal(defaultUuid, mr.getId)
	assert.Equal(1, mr.updateCalled)
}

func TestUserService_Update_WhenGetFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		err: errDefault,
	}
	s := NewUserService(mr)

	_, err := s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Update_CallsUpdateWithUpdatedValues(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		user: defaultUser,
	}
	s := NewUserService(mr)

	s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	expected := persistence.User{
		Id:        defaultUser.Id,
		Email:     defaultUserDtoRequest.Email,
		Password:  defaultUserDtoRequest.Password,
		ApiKeys:   []uuid.UUID{defaultApiKey},
		CreatedAt: defaultUser.CreatedAt,
		UpdatedAt: defaultUser.UpdatedAt,
	}
	assert.Equal(expected, mr.updatedUser)
}

func TestUserService_Update_WhenUpdateFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		updateErr: errDefault,
	}
	s := NewUserService(mr)

	_, err := s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUserService_Update_ReturnsUpdatedUserOmitingApiKey(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		user: defaultUser,
	}
	s := NewUserService(mr)

	actual, err := s.Update(context.Background(), defaultUuid, defaultUserDtoRequest)

	assert.Nil(err)
	assert.Equal(defaultUuid, mr.getId)

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

	mr := &mockUserRepository{}
	s := NewUserService(mr)

	s.Delete(context.Background(), defaultUuid)

	assert.Equal(1, mr.deleteCalled)
}

func TestUserService_Delete_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{
		err: errDefault,
	}
	s := NewUserService(mr)

	err := s.Delete(context.Background(), defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserService_Delete_WhenRepositorySucceeds_ExpectSuccess(t *testing.T) {
	assert := assert.New(t)

	mr := &mockUserRepository{}
	s := NewUserService(mr)

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
