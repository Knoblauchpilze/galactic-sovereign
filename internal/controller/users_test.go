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

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

func Test_UserController(t *testing.T) {
	updatedUserDtoRequest := communication.UserDtoRequest{
		Email:    "some-other@e.mail",
		Password: "some-password",
	}

	s := ControllerTestSuite[service.UserService]{
		generateServiceMock:      generateUserServiceMock,
		generateValidServiceMock: generateValidUserServiceMock,

		badInputTestCases: map[string]badInputTestCase[service.UserService]{
			"createUser": {
				// https://github.com/labstack/echo/issues/2138
				req:                httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-user-dto-request")),
				handler:            createUser,
				expectedBodyString: "\"Invalid user syntax\"\n",
			},
			"updateUser": {
				req:                httptest.NewRequest(http.MethodPatch, "/", strings.NewReader("not-a-user-dto-request")),
				idAsRouteParam:     true,
				handler:            updateUser,
				expectedBodyString: "\"Invalid user syntax\"\n",
			},
			"loginUserByEmail": {
				req:                httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-user-dto-request")),
				handler:            loginUserByEmail,
				expectedBodyString: "\"Invalid user syntax\"\n",
			},
		},

		noIdTestCases: map[string]noIdTestCase[service.UserService]{
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
		},

		badIdTestCases: map[string]badIdTestCase[service.UserService]{
			"getUser": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: getUser,
			},
			"updateUser": {
				req:     httptest.NewRequest(http.MethodPatch, "/", nil),
				handler: updateUser,
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
		},

		errorTestCases: map[string]errorTestCase[service.UserService]{
			"createUser": {
				req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
				handler:            createUser,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"createUser_duplicatedKey": {
				req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
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
				req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
				handler:            listUsers,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"updateUser_notFound": {
				req:                generateTestRequestWithDefaultUserBody(http.MethodPatch),
				idAsRouteParam:     true,
				handler:            updateUser,
				err:                errors.NewCode(db.NoMatchingSqlRows),
				expectedHttpStatus: http.StatusNotFound,
			},
			"updateUser_optimisticLock": {
				req:                generateTestRequestWithDefaultUserBody(http.MethodPatch),
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
				req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
				handler:            loginUserByEmail,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"loginUserByEmail_notFound": {
				req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
				handler:            loginUserByEmail,
				err:                errors.NewCode(db.NoMatchingSqlRows),
				expectedHttpStatus: http.StatusNotFound,
			},
			"loginUserByEmail_invalidCredentials": {
				req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
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
		},

		successTestCases: map[string]successTestCase[service.UserService]{
			"createUser": {
				req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
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
				req:                generateTestRequestWithDefaultUserBody(http.MethodPatch),
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
				req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
				handler:            loginUserByEmail,
				expectedHttpStatus: http.StatusCreated,
			},
			"logoutUser": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            logoutUser,
				expectedHttpStatus: http.StatusNoContent,
			},
		},

		returnTestCases: map[string]returnTestCase[service.UserService]{
			"createUser": {
				req:             generateTestRequestWithDefaultUserBody(http.MethodPost),
				handler:         createUser,
				expectedContent: defaultUserDtoResponse,
			},
			"getUser": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:  true,
				handler:         getUser,
				expectedContent: defaultUserDtoResponse,
			},
			"listUsers": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				handler:         listUsers,
				expectedContent: []uuid.UUID{defaultUuid},
			},
			"updateUser": {
				req:            generateTestRequestWithUserBody(http.MethodPatch, updatedUserDtoRequest),
				idAsRouteParam: true,
				handler:        updateUser,
				// TODO: Was different before.
				expectedContent: defaultUserDtoResponse,
			},
			"loginUserById": {
				req:             httptest.NewRequest(http.MethodPost, "/", nil),
				idAsRouteParam:  true,
				handler:         loginUserById,
				expectedContent: defaultApiKeyDtoResponse,
			},
			"loginUserByEmail": {
				req:             generateTestRequestWithDefaultUserBody(http.MethodPost),
				idAsRouteParam:  true,
				handler:         loginUserByEmail,
				expectedContent: defaultApiKeyDtoResponse,
			},
		},

		serviceInteractionTestCases: map[string]serviceInteractionTestCase[service.UserService]{
			"createUser": {
				req:     generateTestRequestWithDefaultUserBody(http.MethodPost),
				handler: createUser,

				verifyInteractions: func(s service.UserService, assert *require.Assertions) {
					m := assertUserServiceIsAMock(s, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultUserDtoRequest, m.inUser)
				},
			},
			"getUser": {
				req:            httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam: true,
				handler:        getUser,

				verifyInteractions: func(s service.UserService, assert *require.Assertions) {
					m := assertUserServiceIsAMock(s, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
			"listUsers": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: listUsers,

				verifyInteractions: func(s service.UserService, assert *require.Assertions) {
					m := assertUserServiceIsAMock(s, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"updateUser": {
				req:            generateTestRequestWithUserBody(http.MethodPatch, updatedUserDtoRequest),
				idAsRouteParam: true,
				generateValidServiceMock: func() service.UserService {
					return &mockUserService{
						user: communication.UserDtoResponse{
							Id:       defaultUserDtoResponse.Id,
							Email:    updatedUserDtoRequest.Email,
							Password: updatedUserDtoRequest.Password,

							CreatedAt: defaultUserDtoResponse.CreatedAt,
						},
					}
				},

				handler: updateUser,

				verifyInteractions: func(s service.UserService, assert *require.Assertions) {
					m := assertUserServiceIsAMock(s, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(updatedUserDtoRequest, m.inUser)
				},
			},
			"deleteUser": {
				req:            httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam: true,
				handler:        deleteUser,

				verifyInteractions: func(s service.UserService, assert *require.Assertions) {
					m := assertUserServiceIsAMock(s, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
			"loginUserById": {
				req:            httptest.NewRequest(http.MethodPost, "/", nil),
				idAsRouteParam: true,
				handler:        loginUserById,

				verifyInteractions: func(s service.UserService, assert *require.Assertions) {
					m := assertUserServiceIsAMock(s, assert)

					assert.Equal(1, m.loginByIdCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
			"loginUserByEmail": {
				req:     generateTestRequestWithDefaultUserBody(http.MethodPost),
				handler: loginUserByEmail,

				verifyInteractions: func(s service.UserService, assert *require.Assertions) {
					m := assertUserServiceIsAMock(s, assert)

					assert.Equal(1, m.loginCalled)
					assert.Equal(defaultUserDtoRequest, m.inUser)
				},
			},
			"logoutUser": {
				req:            httptest.NewRequest(http.MethodPost, "/", nil),
				idAsRouteParam: true,
				handler:        logoutUser,

				verifyInteractions: func(s service.UserService, assert *require.Assertions) {
					m := assertUserServiceIsAMock(s, assert)

					assert.Equal(1, m.logoutCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateUserServiceMock(err error) service.UserService {
	return &mockUserService{
		err: err,
	}
}

func generateValidUserServiceMock() service.UserService {
	return &mockUserService{
		ids:    []uuid.UUID{defaultUuid},
		user:   defaultUserDtoResponse,
		apiKey: defaultApiKeyDtoResponse,
	}
}

func assertUserServiceIsAMock(s service.UserService, assert *require.Assertions) *mockUserService {
	m, ok := s.(*mockUserService)
	if !ok {
		assert.Fail("Provided user service is not a mock")
	}
	return m
}

func generateTestRequestWithDefaultUserBody(method string) *http.Request {
	return generateTestRequestWithUserBody(method, defaultUserDtoRequest)
}

func generateTestRequestWithUserBody(method string, userDto communication.UserDtoRequest) *http.Request {
	// Voluntarily ignoring errors
	raw, _ := json.Marshal(userDto)
	req := httptest.NewRequest(method, "/", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	return req
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
