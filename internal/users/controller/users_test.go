package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockContext struct {
	echo.Context

	params  map[string]string
	body    communication.UserDtoRequest
	bindErr error

	status int
	data   interface{}
}

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

type mockDbConnection struct {
	db.Connection
}

type mockUserService struct {
	ids  []uuid.UUID
	user communication.UserDtoResponse
	err  error

	getCalled    int
	listCalled   int
	deleteCalled int
}

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultApiKey = uuid.MustParse("cc1742fa-77b4-4f5f-ac92-058c2e47a5d6")
var defaultUserRequest = communication.UserDtoRequest{
	Email:    "e.mail@domain.com",
	Password: "password",
}
var defaultUser = persistence.User{
	Id:        defaultUuid,
	Email:     "e.mail@domain.com",
	Password:  "password",
	ApiKeys:   []uuid.UUID{defaultApiKey},
	CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	UpdatedAt: time.Date(2009, 11, 17, 20, 34, 59, 651387237, time.UTC),
}
var defaultUserDtoResponse = communication.UserDtoResponse{
	Id:       defaultUuid,
	Email:    "e.mail@domain.com",
	Password: "password",

	ApiKeys: []uuid.UUID{},

	CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
}

func TestUserEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range UserEndpoints(&mockDbConnection{}, &mockUserService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(4, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodPost])
	assert.Equal(2, actualRoutes[http.MethodGet])
	assert.Equal(1, actualRoutes[http.MethodPatch])
	assert.Equal(1, actualRoutes[http.MethodDelete])
}

func TestCreateUser_WhenBindFails_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		bindErr: errDefault,
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := createUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid user syntax", mc.data)
}

func TestCreateUser_CallsRepositoryCreate(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	createUser(mc, mr, ms)

	assert.Equal(1, mr.createCalled)
}

func TestCreateUser_WhenRepositorySucceeds_SetsStatusToCreated(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := createUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusCreated, mc.status)
}

func TestCreateUser_WhenRepositorySucceeds_ReturnsExpectedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		body: defaultUserRequest,
	}
	mr := &mockUserRepository{
		user: defaultUser,
	}
	ms := &mockUserService{}

	createUser(mc, mr, ms)

	actual, ok := mc.data.(communication.UserDtoResponse)
	assert.True(ok)

	_, err := uuid.Parse(actual.Id.String())
	assert.Nil(err)

	expected := communication.UserDtoResponse{
		Id:       defaultUuid,
		Email:    "e.mail@domain.com",
		Password: "password",

		ApiKeys: []uuid.UUID{defaultApiKey},

		CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	}
	assert.Equal(expected, actual)
}

func TestCreateUser_WhenRepositorySucceeds_SavesExpectedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		body: defaultUserRequest,
	}
	mr := &mockUserRepository{
		user: defaultUser,
	}
	ms := &mockUserService{}

	createUser(mc, mr, ms)

	actual := mr.createdUser
	_, err := uuid.Parse(actual.Id.String())
	assert.Nil(err)
	assert.Equal(defaultUserRequest.Email, actual.Email)
	assert.Equal(defaultUserRequest.Password, actual.Password)
	n := time.Now()
	assert.True(actual.CreatedAt.Before(n))
	assert.True(actual.UpdatedAt.Before(n))
}

func TestCreateUser_WhenRepositoryFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{
		err: errDefault,
	}
	ms := &mockUserService{}

	err := createUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestGetUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := getUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestGetUser_WhenIdSyntaxIsWrong_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": "not-a-valid-id",
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := getUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestGetUser_CallsServiceGet(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	getUser(mc, mr, ms)

	assert.Equal(1, ms.getCalled)
}

func TestGetUser_WhenServiceFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{
		err: errDefault,
	}

	err := getUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestGetUser_WhenServiceFailsWithNoMatchingRows_SetsStatusToNotFound(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}

	err := getUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusNotFound, mc.status)
}

func TestGetUser_SetsStatusToOk(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := getUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusOK, mc.status)
}

func TestGetUser_ReturnsExpectedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{
		user: defaultUserDtoResponse,
	}

	getUser(mc, mr, ms)

	assert.Equal(defaultUserDtoResponse, mc.data)
}

func TestListUser_CallsServiceList(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	listUsers(mc, mr, ms)

	assert.Equal(1, ms.listCalled)
}

func TestListUser_WhenServiceFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}
	ms := &mockUserService{
		err: errDefault,
	}

	err := listUsers(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestListUser_WhenServiceSucceeds_SetsStatusToOk(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := listUsers(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusOK, mc.status)
}

func TestListUser_ReturnsExpectedIds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}
	ms := &mockUserService{
		ids: []uuid.UUID{defaultUuid},
	}

	listUsers(mc, mr, ms)

	assert.Equal(ms.ids, mc.data)
}

func TestUpdateUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := updateUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestUpdateUser_WhenIdSyntaxIsWrong_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": "not-a-valid-id",
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := updateUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestUpdateUser_WhenIdIsCorrectButBindFails_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
		bindErr: errDefault,
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := updateUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid user syntax", mc.data)
}

func TestUpdateUser_AttemptsToFetchUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	updateUser(mc, mr, ms)

	assert.Equal(1, mr.getCalled)
	assert.Equal(defaultUuid, mr.getId)
}

func TestUpdateUser_WhenGetUserFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		err: errDefault,
	}
	ms := &mockUserService{}

	err := updateUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestUpdateUser_WhenGetUserFailsWithNoWithNoMatchingRows_SetsStatusToNotFound(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}
	ms := &mockUserService{}

	err := updateUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusNotFound, mc.status)
}

func TestUpdateUser_CallsRepositoryUpdate(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	updateUser(mc, mr, ms)

	assert.Equal(1, mr.updateCalled)
}

func TestUpdateUser_WhenRepositorySucceeds_SetsStatusToOk(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := getUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusOK, mc.status)
}

func TestUpdateUser_WhenRepositorySucceeds_ReturnsExpectedUser(t *testing.T) {
	assert := assert.New(t)

	updatedUsed := communication.UserDtoRequest{
		Email:    "some-other-email",
		Password: "some-password",
	}

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
		body: updatedUsed,
	}
	mr := &mockUserRepository{
		user: defaultUser,
	}
	ms := &mockUserService{}

	updateUser(mc, mr, ms)

	assert.Equal(defaultUserDtoResponse, mc.data)
}

func TestUpdateUser_UpdatesUserWithBodyInfo(t *testing.T) {
	assert := assert.New(t)

	updatedUser := communication.UserDtoRequest{
		Email:    "some-other-email",
		Password: "some-password",
	}

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
		body: updatedUser,
	}
	mr := &mockUserRepository{
		user: defaultUser,
	}
	ms := &mockUserService{}

	updateUser(mc, mr, ms)

	assert.Equal(defaultUuid, mr.updatedUser.Id)
	assert.Equal(updatedUser.Email, mr.updatedUser.Email)
	assert.Equal(updatedUser.Password, mr.updatedUser.Password)
}

func TestUpdateUser_WhenRepositoryFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		updateErr: errDefault,
	}
	ms := &mockUserService{}

	err := updateUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestUpdateUser_WhenRepositoryFailsWithOptimisticLocking_SetsStatusToConflict(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{
		updateErr: errors.NewCode(db.OptimisticLockException),
	}
	ms := &mockUserService{}

	err := updateUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusConflict, mc.status)
}

func TestDeleteUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := deleteUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestDeleteUser_WhenIdSyntaxIsWrong_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": "not-a-valid-id",
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := deleteUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestDeleteUser_CallsServiceDelete(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	deleteUser(mc, mr, ms)

	assert.Equal(1, ms.deleteCalled)
}

func TestDeleteUser_WhenServiceFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{
		err: errDefault,
	}

	err := deleteUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestDeleteUser_WhenServiceFailsWithNoMatchingRows_SetsStatusToNotFound(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}

	err := deleteUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusNotFound, mc.status)
}

func TestDeleteUser_SetsStatusToNoContent(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	mr := &mockUserRepository{}
	ms := &mockUserService{}

	err := deleteUser(mc, mr, ms)

	assert.Nil(err)
	assert.Equal(http.StatusNoContent, mc.status)
}

func (m *mockContext) Request() *http.Request {
	return httptest.NewRequest(http.MethodGet, "http://localhost:3000", nil)
}

func (m *mockContext) Param(key string) string {
	if m.params == nil {
		return ""
	}

	if value, ok := m.params[key]; ok {
		return value
	}

	return ""
}

func (m *mockContext) Bind(i interface{}) error {
	dto := i.(*communication.UserDtoRequest)
	*dto = m.body
	return m.bindErr
}

func (m *mockContext) JSON(status int, message interface{}) error {
	m.status = status
	m.data = message
	return nil
}

func (m *mockContext) NoContent(status int) error {
	m.status = status
	return nil
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
	return m.user, m.updateErr
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	return m.err
}

func (m *mockUserService) Create(ctx context.Context, user communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	return communication.UserDtoResponse{}, errors.NewCode(errors.NotImplementedCode)
}

func (m *mockUserService) Get(ctx context.Context, id uuid.UUID) (communication.UserDtoResponse, error) {
	m.getCalled++
	return m.user, m.err
}

func (m *mockUserService) List(ctx context.Context) ([]uuid.UUID, error) {
	m.listCalled++
	return m.ids, m.err
}

func (m *mockUserService) Update(ctx context.Context, id uuid.UUID, user communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	return communication.UserDtoResponse{}, errors.NewCode(errors.NotImplementedCode)
}

func (m *mockUserService) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	return m.err
}
