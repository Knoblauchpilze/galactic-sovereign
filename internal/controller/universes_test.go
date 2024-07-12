package controller

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/google/uuid"
)

type mockUniverseService struct {
	universes []communication.UniverseDtoResponse
	err       error

	createCalled int
	getCalled    int
	listCalled   int
	deleteCalled int

	inUniverse communication.UniverseDtoRequest
	inId       uuid.UUID
}

var defaultUniverseDtoRequest = communication.UniverseDtoRequest{
	Name: "my-universe",
}
var defaultUniverseDtoResponse = communication.UniverseDtoResponse{
	Id:   defaultUuid,
	Name: "my-universe",

	CreatedAt: time.Date(2024, 07, 12, 16, 40, 05, 651387232, time.UTC),
}

// type testCase struct {
// 	req            *http.Request
// 	idAsRouteParam bool
// 	handler        userServiceAwareHttpHandler
// }

// func TestUserEndpoints_GeneratesExpectedRoutes(t *testing.T) {
// 	assert := assert.New(t)

// 	actualRoutes := make(map[string]int)
// 	for _, r := range UserEndpoints(&mockUserService{}) {
// 		actualRoutes[r.Method()]++
// 	}

// 	assert.Equal(4, len(actualRoutes))
// 	assert.Equal(3, actualRoutes[http.MethodPost])
// 	assert.Equal(2, actualRoutes[http.MethodGet])
// 	assert.Equal(1, actualRoutes[http.MethodPatch])
// 	assert.Equal(2, actualRoutes[http.MethodDelete])
// }

// func Test_WhenBodyIsNotAUserDto_SetsStatusTo400(t *testing.T) {
// 	assert := assert.New(t)

// 	// https://github.com/labstack/echo/issues/2138
// 	postReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-user-dto-request"))

// 	testCases := map[string]testCase{
// 		"createUser": {
// 			req:     postReq,
// 			handler: createUser,
// 		},
// 		"updateUser": {
// 			req:            httptest.NewRequest(http.MethodPatch, "/", strings.NewReader("not-a-user-dto-request")),
// 			idAsRouteParam: true,
// 			handler:        updateUser,
// 		},
// 		"loginUserByEmail": {
// 			req:     postReq,
// 			handler: loginUserByEmail,
// 		},
// 	}

// 	for name, testCase := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			mock := &mockUserService{}
// 			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
// 			if testCase.idAsRouteParam {
// 				// https://echo.labstack.com/docs/testing#getuser
// 				ctx.SetParamNames("id")
// 				ctx.SetParamValues(defaultUuid.String())
// 			}

// 			err := testCase.handler(ctx, mock)

// 			assert.Nil(err)
// 			assert.Equal(http.StatusBadRequest, rw.Code)
// 			assert.Equal("\"Invalid user syntax\"\n", rw.Body.String())
// 		})
// 	}
// }

// func Test_WhenNoId_SetsStatusTo400(t *testing.T) {
// 	assert := assert.New(t)

// 	testCases := map[string]testCase{
// 		"getUser": {
// 			req:     httptest.NewRequest(http.MethodGet, "/", nil),
// 			handler: getUser,
// 		},
// 		"updateUser": {
// 			req:     httptest.NewRequest(http.MethodPatch, "/", nil),
// 			handler: getUser,
// 		},
// 		"deleteUser": {
// 			req:     httptest.NewRequest(http.MethodDelete, "/", nil),
// 			handler: deleteUser,
// 		},
// 		"loginUserById": {
// 			req:     httptest.NewRequest(http.MethodPost, "/", nil),
// 			handler: loginUserById,
// 		},
// 		"logoutUser": {
// 			req:     httptest.NewRequest(http.MethodPost, "/", nil),
// 			handler: logoutUser,
// 		},
// 	}

// 	for name, testCase := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			mock := &mockUserService{}
// 			ctx, rw := generateTestEchoContextFromRequest(testCase.req)

// 			err := testCase.handler(ctx, mock)

// 			assert.Nil(err)
// 			assert.Equal(http.StatusBadRequest, rw.Code)
// 			assert.Equal("\"Invalid id syntax\"\n", rw.Body.String())
// 		})
// 	}
// }

// func Test_WhenIdSyntaxIsWrong_SetsStatusTo400(t *testing.T) {
// 	assert := assert.New(t)

// 	testCases := map[string]testCase{
// 		"getUser": {
// 			req:     httptest.NewRequest(http.MethodGet, "/", nil),
// 			handler: getUser,
// 		},
// 		"updateUser": {
// 			req:     httptest.NewRequest(http.MethodPatch, "/", nil),
// 			handler: updateUser,
// 		},
// 		"deleteUser": {
// 			req:     httptest.NewRequest(http.MethodDelete, "/", nil),
// 			handler: deleteUser,
// 		},
// 		"loginUserById": {
// 			req:     httptest.NewRequest(http.MethodPost, "/", nil),
// 			handler: loginUserById,
// 		},
// 		"logoutUser": {
// 			req:     httptest.NewRequest(http.MethodPost, "/", nil),
// 			handler: logoutUser,
// 		},
// 	}

// 	for name, testCase := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			mock := &mockUserService{}
// 			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
// 			ctx.SetParamNames("id")
// 			ctx.SetParamValues("not-a-uuid")

// 			err := testCase.handler(ctx, mock)

// 			assert.Nil(err)
// 			assert.Equal(http.StatusBadRequest, rw.Code)
// 			assert.Equal("\"Invalid id syntax\"\n", rw.Body.String())
// 		})
// 	}
// }

// func Test_WhenServiceFails_SetsExpectedStatus(t *testing.T) {
// 	assert := assert.New(t)

// 	type testCaseError struct {
// 		req                *http.Request
// 		idAsRouteParam     bool
// 		handler            userServiceAwareHttpHandler
// 		err                error
// 		expectedHttpStatus int
// 	}

// 	testCases := map[string]testCaseError{
// 		"createUser": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			handler:            createUser,
// 			err:                errDefault,
// 			expectedHttpStatus: http.StatusInternalServerError,
// 		},
// 		"createUser_duplicatedKey": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			handler:            createUser,
// 			err:                errors.NewCode(db.DuplicatedKeySqlKey),
// 			expectedHttpStatus: http.StatusConflict,
// 		},
// 		"getUser": {
// 			req:                httptest.NewRequest(http.MethodGet, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            getUser,
// 			err:                errDefault,
// 			expectedHttpStatus: http.StatusInternalServerError,
// 		},
// 		"getUser_notFound": {
// 			req:                httptest.NewRequest(http.MethodGet, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            getUser,
// 			err:                errors.NewCode(db.NoMatchingSqlRows),
// 			expectedHttpStatus: http.StatusNotFound,
// 		},
// 		"listUsers": {
// 			req:                generateTestPostRequest(),
// 			handler:            listUsers,
// 			err:                errDefault,
// 			expectedHttpStatus: http.StatusInternalServerError,
// 		},
// 		"updateUser_notFound": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPatch),
// 			idAsRouteParam:     true,
// 			handler:            updateUser,
// 			err:                errors.NewCode(db.NoMatchingSqlRows),
// 			expectedHttpStatus: http.StatusNotFound,
// 		},
// 		"updateUser_optimisticLock": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPatch),
// 			idAsRouteParam:     true,
// 			handler:            updateUser,
// 			err:                errors.NewCode(db.OptimisticLockException),
// 			expectedHttpStatus: http.StatusConflict,
// 		},
// 		"deleteUser": {
// 			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            deleteUser,
// 			err:                errDefault,
// 			expectedHttpStatus: http.StatusInternalServerError,
// 		},
// 		"deleteUser_notFound": {
// 			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            deleteUser,
// 			err:                errors.NewCode(db.NoMatchingSqlRows),
// 			expectedHttpStatus: http.StatusNotFound,
// 		},
// 		"loginUserById": {
// 			req:                httptest.NewRequest(http.MethodPost, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            loginUserById,
// 			err:                errDefault,
// 			expectedHttpStatus: http.StatusInternalServerError,
// 		},
// 		"loginUserById_notFound": {
// 			req:                httptest.NewRequest(http.MethodPost, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            loginUserById,
// 			err:                errors.NewCode(db.NoMatchingSqlRows),
// 			expectedHttpStatus: http.StatusNotFound,
// 		},
// 		"loginUserByEmail": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			handler:            loginUserByEmail,
// 			err:                errDefault,
// 			expectedHttpStatus: http.StatusInternalServerError,
// 		},
// 		"loginUserByEmail_notFound": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			handler:            loginUserByEmail,
// 			err:                errors.NewCode(db.NoMatchingSqlRows),
// 			expectedHttpStatus: http.StatusNotFound,
// 		},
// 		"loginUserByEmail_invalidCredentials": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			handler:            loginUserByEmail,
// 			err:                errors.NewCode(service.InvalidCredentials),
// 			expectedHttpStatus: http.StatusUnauthorized,
// 		},
// 		"logoutUser": {
// 			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            logoutUser,
// 			err:                errDefault,
// 			expectedHttpStatus: http.StatusInternalServerError,
// 		},
// 		"logoutUser_notFound": {
// 			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            logoutUser,
// 			err:                errors.NewCode(db.NoMatchingSqlRows),
// 			expectedHttpStatus: http.StatusNotFound,
// 		},
// 	}

// 	for name, testCase := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			mock := &mockUserService{
// 				err: testCase.err,
// 			}

// 			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
// 			if testCase.idAsRouteParam {
// 				ctx.SetParamNames("id")
// 				ctx.SetParamValues(defaultUuid.String())
// 			}

// 			err := testCase.handler(ctx, mock)

// 			assert.Nil(err)
// 			assert.Equal(testCase.expectedHttpStatus, rw.Code)
// 		})
// 	}
// }

// func Test_WhenServiceSucceeds_SetsExpectedStatus(t *testing.T) {
// 	assert := assert.New(t)

// 	type testCaseSuccess struct {
// 		req                *http.Request
// 		idAsRouteParam     bool
// 		handler            userServiceAwareHttpHandler
// 		expectedHttpStatus int
// 	}

// 	testCases := map[string]testCaseSuccess{
// 		"createUser": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			handler:            createUser,
// 			expectedHttpStatus: http.StatusCreated,
// 		},
// 		"getUser": {
// 			req:                httptest.NewRequest(http.MethodGet, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            getUser,
// 			expectedHttpStatus: http.StatusOK,
// 		},
// 		"listUser": {
// 			req:                httptest.NewRequest(http.MethodGet, "/", nil),
// 			handler:            listUsers,
// 			expectedHttpStatus: http.StatusOK,
// 		},
// 		"updateUser": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPatch),
// 			idAsRouteParam:     true,
// 			handler:            updateUser,
// 			expectedHttpStatus: http.StatusOK,
// 		},
// 		"deleteUser": {
// 			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            deleteUser,
// 			expectedHttpStatus: http.StatusNoContent,
// 		},
// 		"loginUserById": {
// 			req:                httptest.NewRequest(http.MethodPost, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            loginUserById,
// 			expectedHttpStatus: http.StatusCreated,
// 		},
// 		"loginUserByEmail": {
// 			req:                generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			handler:            loginUserByEmail,
// 			expectedHttpStatus: http.StatusCreated,
// 		},
// 		"logoutUser": {
// 			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
// 			idAsRouteParam:     true,
// 			handler:            logoutUser,
// 			expectedHttpStatus: http.StatusNoContent,
// 		},
// 	}

// 	for name, testCase := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			mock := &mockUserService{}

// 			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
// 			if testCase.idAsRouteParam {
// 				ctx.SetParamNames("id")
// 				ctx.SetParamValues(defaultUuid.String())
// 			}

// 			err := testCase.handler(ctx, mock)

// 			assert.Nil(err)
// 			assert.Equal(testCase.expectedHttpStatus, rw.Code)
// 		})
// 	}
// }

// func Test_WhenServiceSucceeds_ReturnsExpectedValue(t *testing.T) {
// 	assert := assert.New(t)

// 	type testCaseReturn struct {
// 		req            *http.Request
// 		idAsRouteParam bool
// 		handler        userServiceAwareHttpHandler

// 		ids       []uuid.UUID
// 		userDto   communication.UserDtoResponse
// 		apiKeyDto communication.ApiKeyDtoResponse

// 		expectedContent interface{}
// 	}

// 	updatedUser := communication.UserDtoRequest{
// 		Email:    "some-other@e.mail",
// 		Password: "some-password",
// 	}
// 	updatedUserResponse := communication.UserDtoResponse{
// 		Id:       defaultUserDtoResponse.Id,
// 		Email:    updatedUser.Email,
// 		Password: updatedUser.Password,

// 		CreatedAt: defaultUserDtoResponse.CreatedAt,
// 	}

// 	testCases := map[string]testCaseReturn{
// 		"createUser": {
// 			req:             generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			handler:         createUser,
// 			userDto:         defaultUserDtoResponse,
// 			expectedContent: defaultUserDtoResponse,
// 		},
// 		"getUser": {
// 			req:             httptest.NewRequest(http.MethodGet, "/", nil),
// 			idAsRouteParam:  true,
// 			handler:         getUser,
// 			userDto:         defaultUserDtoResponse,
// 			expectedContent: defaultUserDtoResponse,
// 		},
// 		"listUsers": {
// 			req:             httptest.NewRequest(http.MethodGet, "/", nil),
// 			handler:         listUsers,
// 			ids:             []uuid.UUID{defaultUuid},
// 			expectedContent: []uuid.UUID{defaultUuid},
// 		},
// 		"updateUser": {
// 			req:             generateTestRequestWithUserBody(http.MethodPatch, updatedUser),
// 			idAsRouteParam:  true,
// 			handler:         updateUser,
// 			userDto:         updatedUserResponse,
// 			expectedContent: updatedUserResponse,
// 		},
// 		"loginUserById": {
// 			req:             httptest.NewRequest(http.MethodPost, "/", nil),
// 			idAsRouteParam:  true,
// 			handler:         loginUserById,
// 			apiKeyDto:       defaultApiKeyDtoResponse,
// 			expectedContent: defaultApiKeyDtoResponse,
// 		},
// 		"loginUserByEmail": {
// 			req:             generateTestRequestWithDefaultUserBody(http.MethodPost),
// 			idAsRouteParam:  true,
// 			handler:         loginUserByEmail,
// 			apiKeyDto:       defaultApiKeyDtoResponse,
// 			expectedContent: defaultApiKeyDtoResponse,
// 		},
// 	}

// 	for name, testCase := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			mock := &mockUserService{
// 				ids:    testCase.ids,
// 				user:   testCase.userDto,
// 				apiKey: testCase.apiKeyDto,
// 			}

// 			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
// 			if testCase.idAsRouteParam {
// 				ctx.SetParamNames("id")
// 				ctx.SetParamValues(defaultUuid.String())
// 			}

// 			err := testCase.handler(ctx, mock)

// 			assert.Nil(err)

// 			actual := strings.Trim(rw.Body.String(), "\n")
// 			expected, err := json.Marshal(testCase.expectedContent)
// 			assert.Nil(err)
// 			assert.Equal(string(expected), actual)
// 		})
// 	}
// }

// func TestCreateUser_CallsServiceCreate(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateTestEchoContextFromRequest(generateTestPostRequest())
// 	ms := &mockUserService{}

// 	err := createUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(1, ms.createCalled)
// }

// func TestCreateUser_SavesExpectedUser(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateTestEchoContextFromRequest(generateTestPostRequest())
// 	ms := &mockUserService{
// 		user: defaultUserDtoResponse,
// 	}

// 	err := createUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(defaultUserDtoRequest, ms.inUser)
// }

// func TestGetUser_CallsServiceGet(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithValidUuid(http.MethodGet)
// 	ms := &mockUserService{}

// 	err := getUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(1, ms.getCalled)
// }

// func TestGetUser_GetsExpectedUser(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithValidUuid(http.MethodGet)
// 	ms := &mockUserService{}

// 	err := getUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(defaultUuid, ms.inId)
// }

// func TestListUser_CallsServiceList(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateTestEchoContextWithMethod(http.MethodGet)
// 	ms := &mockUserService{}

// 	err := listUsers(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(1, ms.listCalled)
// }

// func TestUpdateUser_CallsServiceUpdate(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithUuidAndBody(http.MethodPatch)
// 	ms := &mockUserService{}

// 	err := updateUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(1, ms.updateCalled)
// 	assert.Equal(defaultUuid, ms.inId)
// 	assert.Equal(defaultUserDtoRequest, ms.inUser)
// }

// func TestDeleteUser_CallsServiceDelete(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithValidUuid(http.MethodDelete)
// 	ms := &mockUserService{}

// 	err := deleteUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(1, ms.deleteCalled)
// }

// func TestDeleteUser_LogsOutExpectedUser(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithValidUuid(http.MethodDelete)
// 	ms := &mockUserService{}

// 	err := deleteUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(defaultUuid, ms.inId)
// }

// func TestLoginUserById_CallsServiceLoginById(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithValidUuid(http.MethodPost)
// 	ms := &mockUserService{}

// 	err := loginUserById(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(1, ms.loginByIdCalled)
// }

// func TestLoginUserById_LogsInExpectedUser(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithValidUuid(http.MethodPost)
// 	ms := &mockUserService{}

// 	err := loginUserById(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(defaultUuid, ms.inId)
// }

// func TestLoginUserByEmail_CallsServiceLogin(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithBody(http.MethodPost)
// 	ms := &mockUserService{}

// 	err := loginUserByEmail(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(1, ms.loginCalled)
// }

// func TestLoginUserByEmail_LogsInExpectedUser(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithBody(http.MethodPost)
// 	ms := &mockUserService{}

// 	err := loginUserByEmail(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(defaultUserDtoRequest, ms.inUser)
// }

// func TestLogoutUser_CallsServiceLogout(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithValidUuid(http.MethodPost)
// 	ms := &mockUserService{}

// 	err := logoutUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(1, ms.logoutCalled)
// }

// func TestLogoutUser_LogsOutExpectedUser(t *testing.T) {
// 	assert := assert.New(t)

// 	ctx, _ := generateEchoContextWithValidUuid(http.MethodPost)
// 	ms := &mockUserService{}

// 	err := logoutUser(ctx, ms)

// 	assert.Nil(err)
// 	assert.Equal(defaultUuid, ms.inId)
// }

// func generateTestPostRequest() *http.Request {
// 	return generateTestRequestWithDefaultUserBody(http.MethodPost)
// }

// func generateTestRequestWithDefaultUserBody(method string) *http.Request {
// 	return generateTestRequestWithUserBody(method, defaultUserDtoRequest)
// }

// func generateTestRequestWithUserBody(method string, userDto communication.UserDtoRequest) *http.Request {
// 	// Voluntarily ignoring errors
// 	raw, _ := json.Marshal(userDto)
// 	req := httptest.NewRequest(method, "/", bytes.NewReader(raw))
// 	req.Header.Set("Content-Type", "application/json")

// 	return req
// }

// func generateEchoContextWithValidUuid(method string) (echo.Context, *httptest.ResponseRecorder) {
// 	ctx, rw := generateTestEchoContextWithMethod(method)
// 	ctx.SetParamNames("id")
// 	ctx.SetParamValues(defaultUuid.String())
// 	return ctx, rw
// }

// func generateEchoContextWithBody(method string) (echo.Context, *httptest.ResponseRecorder) {
// 	req := generateTestRequestWithDefaultUserBody(method)
// 	return generateTestEchoContextFromRequest(req)
// }

// func generateEchoContextWithUuidAndBody(method string) (echo.Context, *httptest.ResponseRecorder) {
// 	req := generateTestRequestWithDefaultUserBody(method)

// 	ctx, rw := generateTestEchoContextFromRequest(req)
// 	ctx.SetParamNames("id")
// 	ctx.SetParamValues(defaultUuid.String())
// 	return ctx, rw
// }

func (m *mockUniverseService) Create(ctx context.Context, universe communication.UniverseDtoRequest) (communication.UniverseDtoResponse, error) {
	m.createCalled++
	m.inUniverse = universe
	return m.universes[0], m.err
}

func (m *mockUniverseService) Get(ctx context.Context, id uuid.UUID) (communication.UniverseDtoResponse, error) {
	m.getCalled++
	m.inId = id
	return m.universes[0], m.err
}

func (m *mockUniverseService) List(ctx context.Context) ([]communication.UniverseDtoResponse, error) {
	m.listCalled++
	return m.universes, m.err
}

func (m *mockUniverseService) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	m.inId = id
	return m.err
}
