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

type mockUserService struct {
	ids    []uuid.UUID
	user   communication.UserDtoResponse
	apiKey communication.ApiKeyDtoResponse
	err    error

	createCalled int
	getCalled    int
	listCalled   int
	updateCalled int
	deleteCalled int
	loginCalled  int

	inUser communication.UserDtoRequest
	inId   uuid.UUID
}

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultUserRequest = communication.UserDtoRequest{
	Email:    "e.mail@domain.com",
	Password: "password",
}
var defaultUserDtoResponse = communication.UserDtoResponse{
	Id:       defaultUuid,
	Email:    "e.mail@domain.com",
	Password: "password",

	CreatedAt: time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC),
}
var defaultApiKeyDtoResponse = communication.ApiKeyDtoResponse{
	Key:        uuid.MustParse("9380e881-39c3-42f1-b594-b5d2010e67c0"),
	ValidUntil: time.Date(2024, 05, 05, 21, 32, 55, 651387237, time.UTC),
}

func TestUserEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range UserEndpoints(&mockUserService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(4, len(actualRoutes))
	assert.Equal(2, actualRoutes[http.MethodPost])
	assert.Equal(2, actualRoutes[http.MethodGet])
	assert.Equal(1, actualRoutes[http.MethodPatch])
	assert.Equal(2, actualRoutes[http.MethodDelete])
}

func TestCreateUser_WhenBindFails_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		bindErr: errDefault,
	}
	ms := &mockUserService{}

	err := createUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid user syntax", mc.data)
}

func TestCreateUser_CallsServiceCreate(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{}

	err := createUser(mc, ms)

	assert.Nil(err)
	assert.Equal(1, ms.createCalled)
}

func TestCreateUser_WhenServiceFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{
		err: errDefault,
	}

	err := createUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestCreateUser_WhenServiceFailsWithDuplicatedSqlKey_SetsStatusToConflict(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{
		err: errors.NewCode(db.DuplicatedKeySqlKey),
	}

	err := createUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusConflict, mc.status)
}

func TestCreateUser_SetsStatusToCreated(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{}

	err := createUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusCreated, mc.status)
}

func TestCreateUser_SavesExpectedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		body: defaultUserRequest,
	}
	ms := &mockUserService{
		user: defaultUserDtoResponse,
	}

	err := createUser(mc, ms)

	assert.Nil(err)
	assert.Equal(defaultUserRequest, ms.inUser)
}

func TestCreateUser_ReturnsExpectedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		body: defaultUserRequest,
	}
	ms := &mockUserService{
		user: defaultUserDtoResponse,
	}

	err := createUser(mc, ms)

	assert.Nil(err)
	actual, ok := mc.data.(communication.UserDtoResponse)
	assert.True(ok)
	assert.Equal(defaultUserDtoResponse, actual)
}

func TestGetUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{}

	err := getUser(mc, ms)

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
	ms := &mockUserService{}

	err := getUser(mc, ms)

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
	ms := &mockUserService{}

	err := getUser(mc, ms)

	assert.Nil(err)
	assert.Equal(1, ms.getCalled)
	assert.Equal(defaultUuid, ms.inId)
}

func TestGetUser_WhenServiceFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{
		err: errDefault,
	}

	err := getUser(mc, ms)

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
	ms := &mockUserService{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}

	err := getUser(mc, ms)

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
	ms := &mockUserService{}

	err := getUser(mc, ms)

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
	ms := &mockUserService{
		user: defaultUserDtoResponse,
	}

	err := getUser(mc, ms)

	assert.Nil(err)
	assert.Equal(defaultUserDtoResponse, mc.data)
}

func TestListUser_CallsServiceList(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{}

	err := listUsers(mc, ms)

	assert.Nil(err)
	assert.Equal(1, ms.listCalled)
}

func TestListUser_WhenServiceFails_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{
		err: errDefault,
	}

	err := listUsers(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestListUser_SetsStatusToOk(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{}

	err := listUsers(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusOK, mc.status)
}

func TestListUser_ReturnsExpectedIds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{
		ids: []uuid.UUID{defaultUuid},
	}

	err := listUsers(mc, ms)

	assert.Nil(err)
	assert.Equal(ms.ids, mc.data)
}

func TestUpdateUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{}

	err := updateUser(mc, ms)

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
	ms := &mockUserService{}

	err := updateUser(mc, ms)

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
	ms := &mockUserService{}

	err := updateUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid user syntax", mc.data)
}

func TestUpdateUser_CallsServiceUpdate(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
		body: defaultUserRequest,
	}
	ms := &mockUserService{}

	err := updateUser(mc, ms)

	assert.Nil(err)
	assert.Equal(1, ms.updateCalled)
	assert.Equal(defaultUuid, ms.inId)
	assert.Equal(defaultUserRequest, ms.inUser)
}

func TestUpdateUser_WhenServiceFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
		body: defaultUserRequest,
	}
	ms := &mockUserService{
		err: errDefault,
	}

	err := updateUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestUpdateUser_WhenServiceFailsWithNoSuchRows_SetsStatusToNotFound(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
		body: defaultUserRequest,
	}
	ms := &mockUserService{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}

	err := updateUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusNotFound, mc.status)
}

func TestUpdateUser_WhenServiceFailsWithOptimisticLockException_SetsStatusToConflict(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
		body: defaultUserRequest,
	}
	ms := &mockUserService{
		err: errors.NewCode(db.OptimisticLockException),
	}

	err := updateUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusConflict, mc.status)
}

func TestUpdateUser_SetsStatusToOk(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
		body: defaultUserRequest,
	}
	ms := &mockUserService{}

	err := getUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusOK, mc.status)
}

func TestUpdateUser_ReturnsExpectedUser(t *testing.T) {
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
	ms := &mockUserService{
		user: defaultUserDtoResponse,
	}

	err := updateUser(mc, ms)

	assert.Nil(err)
	assert.Equal(defaultUserDtoResponse, mc.data)
}

func TestDeleteUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{}

	err := deleteUser(mc, ms)

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
	ms := &mockUserService{}

	err := deleteUser(mc, ms)

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
	ms := &mockUserService{}

	err := deleteUser(mc, ms)

	assert.Nil(err)
	assert.Equal(1, ms.deleteCalled)
	assert.Equal(defaultUuid, ms.inId)
}

func TestDeleteUser_WhenServiceFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{
		err: errDefault,
	}

	err := deleteUser(mc, ms)

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
	ms := &mockUserService{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}

	err := deleteUser(mc, ms)

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
	ms := &mockUserService{}

	err := deleteUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusNoContent, mc.status)
}

func TestLoginUser_WhenNoId_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{}
	ms := &mockUserService{}

	err := loginUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestLoginUser_WhenIdSyntaxIsWrong_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": "not-a-valid-id",
		},
	}
	ms := &mockUserService{}

	err := loginUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, mc.status)
	assert.Equal("Invalid id syntax", mc.data)
}

func TestLoginUser_CallsServiceLogin(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{}

	err := loginUser(mc, ms)

	assert.Nil(err)
	assert.Equal(1, ms.loginCalled)
}

func TestLoginUser_WhenServiceFailsWithUnknownError_SetsStatusToInternalServerError(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{
		err: errDefault,
	}

	err := loginUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, mc.status)
}

func TestLoginUser_WhenServiceFailsWithNoMatchingRows_SetsStatusToNotFound(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{
		err: errors.NewCode(db.NoMatchingSqlRows),
	}

	err := loginUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusNotFound, mc.status)
}

func TestLoginUser_SetsStatusToCreated(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{}

	err := loginUser(mc, ms)

	assert.Nil(err)
	assert.Equal(http.StatusCreated, mc.status)
}

func TestLoginUser_LogsInExpectedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{}

	err := loginUser(mc, ms)

	assert.Nil(err)
	assert.Equal(defaultUuid, ms.inId)
}

func TestLoginUser_ReturnsUserToken(t *testing.T) {
	assert := assert.New(t)

	mc := &mockContext{
		params: map[string]string{
			"id": defaultUuid.String(),
		},
	}
	ms := &mockUserService{
		apiKey: defaultApiKeyDtoResponse,
	}

	err := loginUser(mc, ms)

	assert.Nil(err)
	actual, ok := mc.data.(communication.ApiKeyDtoResponse)
	assert.True(ok)
	assert.Equal(defaultApiKeyDtoResponse, actual)
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

func (m *mockUserService) Create(ctx context.Context, user communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	m.createCalled++
	m.inUser = user
	return m.user, m.err
}

func (m *mockUserService) Get(ctx context.Context, id uuid.UUID) (communication.UserDtoResponse, error) {
	m.getCalled++
	m.inId = id
	return m.user, m.err
}

func (m *mockUserService) List(ctx context.Context) ([]uuid.UUID, error) {
	m.listCalled++
	return m.ids, m.err
}

func (m *mockUserService) Update(ctx context.Context, id uuid.UUID, user communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	m.updateCalled++
	m.inId = id
	m.inUser = user
	return m.user, m.err
}

func (m *mockUserService) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	m.inId = id
	return m.err
}

func (m *mockUserService) Login(ctx context.Context, id uuid.UUID) (communication.ApiKeyDtoResponse, error) {
	m.loginCalled++
	m.inId = id
	return m.apiKey, m.err
}

func (m *mockUserService) Logout(ctx context.Context, id uuid.UUID) error {
	return m.err
}
