package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	ids    []uuid.UUID
	user   communication.UserDtoResponse
	apiKey communication.ApiKeyDtoResponse
	err    error

	createCalled    int
	getCalled       int
	listCalled      int
	updateCalled    int
	deleteCalled    int
	loginCalled     int
	loginByIdCalled int
	logoutCalled    int

	inUser communication.UserDtoRequest
	inId   uuid.UUID
}

type testCase struct {
	req     *http.Request
	handler userServiceAwareHttpHandler
}

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultUserDtoRequest = communication.UserDtoRequest{
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
	assert.Equal(3, actualRoutes[http.MethodPost])
	assert.Equal(2, actualRoutes[http.MethodGet])
	assert.Equal(1, actualRoutes[http.MethodPatch])
	assert.Equal(2, actualRoutes[http.MethodDelete])
}

func Test_WhenBodyIsNotAUserDto_SetsStatusTo400(t *testing.T) {
	assert := assert.New(t)

	// https://github.com/labstack/echo/issues/2138
	postReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-user-dto-request"))

	testCases := map[string]testCase{
		"createUser": {
			req:     postReq,
			handler: createUser,
		},
		"loginUserByEmail": {
			req:     postReq,
			handler: loginUserByEmail,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUserService{}
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(http.StatusBadRequest, rw.Code)
			assert.Equal("\"Invalid user syntax\"\n", rw.Body.String())
		})
	}
}

func Test_WhenNoId_SetsStatusTo400(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]testCase{
		"getUser": {
			req:     httptest.NewRequest(http.MethodGet, "/", nil),
			handler: getUser,
		},
		"updateUser": {
			req:     httptest.NewRequest(http.MethodPatch, "/", nil),
			handler: getUser,
		},
		"deleteUser": {
			req:     httptest.NewRequest(http.MethodDelete, "/", nil),
			handler: deleteUser,
		},
		"loginUserById": {
			req:     httptest.NewRequest(http.MethodPost, "/", nil),
			handler: loginUserById,
		},
		"logoutUser": {
			req:     httptest.NewRequest(http.MethodPost, "/", nil),
			handler: logoutUser,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUserService{}
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(http.StatusBadRequest, rw.Code)
			assert.Equal("\"Invalid id syntax\"\n", rw.Body.String())
		})
	}
}

func Test_WhenIdSyntaxIsWrong_SetsStatusTo400(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]testCase{
		"getUser": {
			req:     httptest.NewRequest(http.MethodGet, "/", nil),
			handler: getUser,
		},
		"updateUser": {
			req:     httptest.NewRequest(http.MethodPatch, "/", nil),
			handler: getUser,
		},
		"deleteUser": {
			req:     httptest.NewRequest(http.MethodDelete, "/", nil),
			handler: deleteUser,
		},
		"loginUserById": {
			req:     httptest.NewRequest(http.MethodPost, "/", nil),
			handler: loginUserById,
		},
		"logoutUser": {
			req:     httptest.NewRequest(http.MethodPost, "/", nil),
			handler: logoutUser,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUserService{}
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			ctx.SetParamNames("id")
			ctx.SetParamValues("not-a-uuid")

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(http.StatusBadRequest, rw.Code)
			assert.Equal("\"Invalid id syntax\"\n", rw.Body.String())
		})
	}
}

func Test_WhenServiceFails_SetsExpectedStatus(t *testing.T) {
	assert := assert.New(t)

	type testCaseError struct {
		req                *http.Request
		idAsRouteParam     bool
		handler            userServiceAwareHttpHandler
		err                error
		expectedHttpStatus int
	}

	testCases := map[string]testCaseError{
		"createUser": {
			req:                generateTestRequestWithUserBody(http.MethodPost),
			handler:            createUser,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"createUser_duplicatedKey": {
			req:                generateTestRequestWithUserBody(http.MethodPost),
			handler:            createUser,
			err:                errors.NewCode(db.DuplicatedKeySqlKey),
			expectedHttpStatus: http.StatusConflict,
		},
		"getUser": {
			req:                httptest.NewRequest(http.MethodGet, "/", nil),
			idAsRouteParam:     true,
			handler:            getUser,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"getUser_notFound": {
			req:                httptest.NewRequest(http.MethodGet, "/", nil),
			idAsRouteParam:     true,
			handler:            getUser,
			err:                errors.NewCode(db.NoMatchingSqlRows),
			expectedHttpStatus: http.StatusNotFound,
		},
		"listUsers": {
			req:                generateTestPostRequest(),
			handler:            listUsers,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"updateUser_notFound": {
			req:                generateTestRequestWithUserBody(http.MethodPatch),
			idAsRouteParam:     true,
			handler:            updateUser,
			err:                errors.NewCode(db.NoMatchingSqlRows),
			expectedHttpStatus: http.StatusNotFound,
		},
		"updateUser_optimisticLock": {
			req:                generateTestRequestWithUserBody(http.MethodPatch),
			idAsRouteParam:     true,
			handler:            updateUser,
			err:                errors.NewCode(db.OptimisticLockException),
			expectedHttpStatus: http.StatusConflict,
		},
		"deleteUser": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            deleteUser,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"deleteUser_notFound": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            deleteUser,
			err:                errors.NewCode(db.NoMatchingSqlRows),
			expectedHttpStatus: http.StatusNotFound,
		},
		"loginUserById": {
			req:                httptest.NewRequest(http.MethodPost, "/", nil),
			idAsRouteParam:     true,
			handler:            loginUserById,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"loginUserById_notFound": {
			req:                httptest.NewRequest(http.MethodPost, "/", nil),
			idAsRouteParam:     true,
			handler:            loginUserById,
			err:                errors.NewCode(db.NoMatchingSqlRows),
			expectedHttpStatus: http.StatusNotFound,
		},
		"loginUserByEmail": {
			req:                generateTestRequestWithUserBody(http.MethodPost),
			handler:            loginUserByEmail,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"loginUserByEmail_notFound": {
			req:                generateTestRequestWithUserBody(http.MethodPost),
			handler:            loginUserByEmail,
			err:                errors.NewCode(db.NoMatchingSqlRows),
			expectedHttpStatus: http.StatusNotFound,
		},
		"loginUserByEmail_invalidCredentials": {
			req:                generateTestRequestWithUserBody(http.MethodPost),
			handler:            loginUserByEmail,
			err:                errors.NewCode(service.InvalidCredentials),
			expectedHttpStatus: http.StatusUnauthorized,
		},
		"logoutUser": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            logoutUser,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"logoutUser_notFound": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            logoutUser,
			err:                errors.NewCode(db.NoMatchingSqlRows),
			expectedHttpStatus: http.StatusNotFound,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUserService{
				err: testCase.err,
			}

			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				// https://echo.labstack.com/docs/testing#getuser
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(testCase.expectedHttpStatus, rw.Code)
		})
	}
}

func Test_WhenServiceSucceeds_SetsExpectedStatus(t *testing.T) {
	assert := assert.New(t)

	type testCaseError struct {
		req                *http.Request
		idAsRouteParam     bool
		handler            userServiceAwareHttpHandler
		expectedHttpStatus int
	}

	testCases := map[string]testCaseError{
		"createUser": {
			req:                generateTestRequestWithUserBody(http.MethodPost),
			handler:            createUser,
			expectedHttpStatus: http.StatusCreated,
		},
		"getUser": {
			req:                httptest.NewRequest(http.MethodGet, "/", nil),
			idAsRouteParam:     true,
			handler:            getUser,
			expectedHttpStatus: http.StatusOK,
		},
		"listUser": {
			req:                httptest.NewRequest(http.MethodGet, "/", nil),
			handler:            listUsers,
			expectedHttpStatus: http.StatusOK,
		},
		"updateUser": {
			req:                generateTestRequestWithUserBody(http.MethodPatch),
			idAsRouteParam:     true,
			handler:            updateUser,
			expectedHttpStatus: http.StatusOK,
		},
		"deleteUser": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            deleteUser,
			expectedHttpStatus: http.StatusNoContent,
		},
		"loginUserById": {
			req:                httptest.NewRequest(http.MethodPost, "/", nil),
			idAsRouteParam:     true,
			handler:            loginUserById,
			expectedHttpStatus: http.StatusCreated,
		},
		"loginUserByEmail": {
			req:                generateTestRequestWithUserBody(http.MethodPost),
			handler:            loginUserByEmail,
			expectedHttpStatus: http.StatusCreated,
		},
		"logoutUser": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            logoutUser,
			expectedHttpStatus: http.StatusNoContent,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUserService{}

			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(testCase.expectedHttpStatus, rw.Code)
		})
	}
}

func TestCreateUser_CallsServiceCreate(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateTestEchoContextFromRequest(generateTestPostRequest())
	ms := &mockUserService{}

	err := createUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.createCalled)
}

func TestCreateUser_SavesExpectedUser(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateTestEchoContextFromRequest(generateTestPostRequest())
	ms := &mockUserService{
		user: defaultUserDtoResponse,
	}

	err := createUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(defaultUserDtoRequest, ms.inUser)
}

func TestCreateUser_ReturnsExpectedUser(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateTestEchoContextFromRequest(generateTestPostRequest())
	ms := &mockUserService{
		user: defaultUserDtoResponse,
	}

	err := createUser(ctx, ms)

	assert.Nil(err)

	var actual communication.UserDtoResponse
	err = json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal(defaultUserDtoResponse, actual)
}

func TestGetUser_CallsServiceGet(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodGet)
	ms := &mockUserService{}

	err := getUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.getCalled)
	assert.Equal(defaultUuid, ms.inId)
}

func TestGetUser_ReturnsExpectedUser(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateEchoContextWithValidUuid(http.MethodGet)
	ms := &mockUserService{
		user: defaultUserDtoResponse,
	}

	err := getUser(ctx, ms)

	assert.Nil(err)

	var actual communication.UserDtoResponse
	err = json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal(defaultUserDtoResponse, actual)
}

func TestListUser_CallsServiceList(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateTestEchoContextWithMethod(http.MethodGet)
	ms := &mockUserService{}

	err := listUsers(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.listCalled)
}

func TestListUser_ReturnsExpectedIds(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateTestEchoContextWithMethod(http.MethodGet)
	ms := &mockUserService{
		ids: []uuid.UUID{defaultUuid},
	}

	err := listUsers(ctx, ms)

	assert.Nil(err)

	var ids []uuid.UUID
	err = json.Unmarshal(rw.Body.Bytes(), &ids)
	assert.Nil(err)
	assert.Equal(ms.ids, ids)
}

func TestUpdateUser_WhenIdIsCorrectButBodyIsNotAUserDto_SetsStatusToBadRequest(t *testing.T) {
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader("not-a-user-dto-request"))
	ctx, rw := generateEchoContextWithValidUuid(http.MethodPatch)
	ctx.SetRequest(req)

	ms := &mockUserService{}

	err := updateUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, rw.Code)
	assert.Equal("\"Invalid user syntax\"\n", rw.Body.String())
}

func TestUpdateUser_CallsServiceUpdate(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithUuidAndBody(http.MethodPatch)
	ms := &mockUserService{}

	err := updateUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.updateCalled)
	assert.Equal(defaultUuid, ms.inId)
	assert.Equal(defaultUserDtoRequest, ms.inUser)
}

func TestUpdateUser_ReturnsExpectedUser(t *testing.T) {
	assert := assert.New(t)

	updatedUser := communication.UserDtoRequest{
		Email:    "some-other-email",
		Password: "some-password",
	}
	updatedResponse := communication.UserDtoResponse{
		Id:       defaultUserDtoResponse.Id,
		Email:    updatedUser.Email,
		Password: updatedUser.Password,

		CreatedAt: defaultUserDtoResponse.CreatedAt,
	}

	raw, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	ctx, rw := generateTestEchoContextFromRequest(req)
	ctx.SetParamNames("id")
	ctx.SetParamValues(defaultUuid.String())

	ms := &mockUserService{
		user: updatedResponse,
	}

	err := updateUser(ctx, ms)

	assert.Nil(err)

	var actual communication.UserDtoResponse
	err = json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal(updatedResponse, actual)
}

func TestDeleteUser_CallsServiceDelete(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodDelete)
	ms := &mockUserService{}

	err := deleteUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.deleteCalled)
	assert.Equal(defaultUuid, ms.inId)
}

func TestLoginUserById_CallsServiceLoginById(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodPost)
	ms := &mockUserService{}

	err := loginUserById(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.loginByIdCalled)
}

func TestLoginUserById_LogsInExpectedUser(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodPost)
	ms := &mockUserService{}

	err := loginUserById(ctx, ms)

	assert.Nil(err)
	assert.Equal(defaultUuid, ms.inId)
}

func TestLoginUserById_ReturnsUserToken(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateEchoContextWithValidUuid(http.MethodPost)
	ms := &mockUserService{
		apiKey: defaultApiKeyDtoResponse,
	}

	err := loginUserById(ctx, ms)

	assert.Nil(err)

	var actual communication.ApiKeyDtoResponse
	err = json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal(defaultApiKeyDtoResponse, actual)
}

func TestLoginUserByEmail_CallsServiceLogin(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithBody(http.MethodPost)
	ms := &mockUserService{}

	err := loginUserByEmail(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.loginCalled)
}

func TestLoginUserByEmail_LogsInExpectedUser(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithBody(http.MethodPost)
	ms := &mockUserService{}

	err := loginUserByEmail(ctx, ms)

	assert.Nil(err)
	assert.Equal(defaultUserDtoRequest, ms.inUser)
}

func TestLoginUserByEmail_ReturnsUserToken(t *testing.T) {
	assert := assert.New(t)

	ctx, rw := generateEchoContextWithBody(http.MethodPost)
	ms := &mockUserService{
		apiKey: defaultApiKeyDtoResponse,
	}

	err := loginUserByEmail(ctx, ms)

	assert.Nil(err)

	var actual communication.ApiKeyDtoResponse
	err = json.Unmarshal(rw.Body.Bytes(), &actual)
	assert.Nil(err)
	assert.Equal(defaultApiKeyDtoResponse, actual)
}

func TestLogoutUser_CallsServiceLogout(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodPost)
	ms := &mockUserService{}

	err := logoutUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.logoutCalled)
}

func TestLogoutUser_LogsOutExpectedUser(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodPost)
	ms := &mockUserService{}

	err := logoutUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(defaultUuid, ms.inId)
}

func generateTestPostRequest() *http.Request {
	return generateTestRequestWithUserBody(http.MethodPost)
}

func generateTestRequestWithUserBody(method string) *http.Request {
	// Voluntarily ignoring errors
	raw, _ := json.Marshal(defaultUserDtoRequest)
	req := httptest.NewRequest(method, "/", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func generateEchoContextWithValidUuid(method string) (echo.Context, *httptest.ResponseRecorder) {
	return generateEchoContextWithUuid(method, defaultUuid.String())
}

func generateEchoContextWithUuid(method string, id string) (echo.Context, *httptest.ResponseRecorder) {
	ctx, rw := generateTestEchoContextWithMethod(method)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id)
	return ctx, rw
}

func generateEchoContextWithBody(method string) (echo.Context, *httptest.ResponseRecorder) {
	req := generateTestRequestWithUserBody(method)
	return generateTestEchoContextFromRequest(req)
}

func generateEchoContextWithUuidAndBody(method string) (echo.Context, *httptest.ResponseRecorder) {
	req := generateTestRequestWithUserBody(method)

	ctx, rw := generateTestEchoContextFromRequest(req)
	ctx.SetParamNames("id")
	ctx.SetParamValues(defaultUuid.String())
	return ctx, rw
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

func (m *mockUserService) Login(ctx context.Context, user communication.UserDtoRequest) (communication.ApiKeyDtoResponse, error) {
	m.loginCalled++
	m.inUser = user
	return m.apiKey, m.err
}

func (m *mockUserService) LoginById(ctx context.Context, id uuid.UUID) (communication.ApiKeyDtoResponse, error) {
	m.loginByIdCalled++
	m.inId = id
	return m.apiKey, m.err
}

func (m *mockUserService) Logout(ctx context.Context, id uuid.UUID) error {
	m.logoutCalled++
	m.inId = id
	return m.err
}
