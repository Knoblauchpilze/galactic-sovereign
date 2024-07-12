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

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

type universeTestCase struct {
	req            *http.Request
	idAsRouteParam bool
	handler        universeServiceAwareHttpHandler
}

type universeErrorTestCase struct {
	req                *http.Request
	idAsRouteParam     bool
	handler            universeServiceAwareHttpHandler
	err                error
	expectedHttpStatus int
}

type universeSuccessTestCase struct {
	req                *http.Request
	idAsRouteParam     bool
	handler            universeServiceAwareHttpHandler
	expectedHttpStatus int
}

type universeTestCaseReturn struct {
	req            *http.Request
	idAsRouteParam bool
	handler        universeServiceAwareHttpHandler

	universeDto communication.UniverseDtoResponse

	expectedContent interface{}
}

func TestUniverseEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range UniverseEndpoints(&mockUniverseService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(3, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodPost])
	assert.Equal(2, actualRoutes[http.MethodGet])
	assert.Equal(1, actualRoutes[http.MethodDelete])
}

func Test_Universes_WhenBodyIsNotAValidUniverseDto_SetsStatusTo400(t *testing.T) {
	assert := assert.New(t)

	postReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-dto-request"))

	testCases := map[string]universeTestCase{
		"createUniverse": {
			req:     postReq,
			handler: createUniverse,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUniverseService{}
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(http.StatusBadRequest, rw.Code)
			assert.Equal("\"Invalid universe syntax\"\n", rw.Body.String())
		})
	}
}

func Test_Universes_WhenNoId_SetsStatusTo400(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]universeTestCase{
		"getUniverse": {
			req:     httptest.NewRequest(http.MethodGet, "/", nil),
			handler: getUniverse,
		},
		"deleteUniverse": {
			req:     httptest.NewRequest(http.MethodDelete, "/", nil),
			handler: deleteUniverse,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUniverseService{}
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(http.StatusBadRequest, rw.Code)
			assert.Equal("\"Invalid id syntax\"\n", rw.Body.String())
		})
	}
}

func Test_Universes_WhenIdSyntaxIsWrong_SetsStatusTo400(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]universeTestCase{
		"getUniverse": {
			req:     httptest.NewRequest(http.MethodGet, "/", nil),
			handler: getUniverse,
		},
		"deleteUniverse": {
			req:     httptest.NewRequest(http.MethodDelete, "/", nil),
			handler: deleteUniverse,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUniverseService{}
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

func Test_Universes_WhenServiceFails_SetsExpectedStatus(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]universeErrorTestCase{
		"createUniverse": {
			req:                generateTestRequestWithDefaultUniverseBody(http.MethodPost),
			handler:            createUniverse,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"createUniverse_duplicatedKey": {
			req:                generateTestRequestWithDefaultUniverseBody(http.MethodPost),
			handler:            createUniverse,
			err:                errors.NewCode(db.DuplicatedKeySqlKey),
			expectedHttpStatus: http.StatusConflict,
		},
		"getUniverse": {
			req:                httptest.NewRequest(http.MethodGet, "/", nil),
			idAsRouteParam:     true,
			handler:            getUniverse,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"getUniverse_notFound": {
			req:                httptest.NewRequest(http.MethodGet, "/", nil),
			idAsRouteParam:     true,
			handler:            getUniverse,
			err:                errors.NewCode(db.NoMatchingSqlRows),
			expectedHttpStatus: http.StatusNotFound,
		},
		"listUniverses": {
			req:                generateTestPostRequest(),
			handler:            listUniverses,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"deleteUniverse": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            deleteUniverse,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
		"deleteUniverse_notFound": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            deleteUniverse,
			err:                errors.NewCode(db.NoMatchingSqlRows),
			expectedHttpStatus: http.StatusNotFound,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUniverseService{
				err: testCase.err,
			}

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

func Test_Universes_WhenServiceSucceeds_SetsExpectedStatus(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]universeSuccessTestCase{
		"createUniverse": {
			req:                generateTestRequestWithDefaultUniverseBody(http.MethodPost),
			handler:            createUniverse,
			expectedHttpStatus: http.StatusCreated,
		},
		"getUniverse": {
			req:                httptest.NewRequest(http.MethodGet, "/", nil),
			idAsRouteParam:     true,
			handler:            getUniverse,
			expectedHttpStatus: http.StatusOK,
		},
		"listUniverses": {
			req:                httptest.NewRequest(http.MethodGet, "/", nil),
			handler:            listUniverses,
			expectedHttpStatus: http.StatusOK,
		},
		"deleteUniverse": {
			req:                httptest.NewRequest(http.MethodDelete, "/", nil),
			idAsRouteParam:     true,
			handler:            deleteUniverse,
			expectedHttpStatus: http.StatusNoContent,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUniverseService{}

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

func Test_Universes_WhenServiceSucceeds_ReturnsExpectedValue(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]universeTestCaseReturn{
		"createUniverse": {
			req:             generateTestRequestWithDefaultUniverseBody(http.MethodPost),
			handler:         createUniverse,
			universeDto:     defaultUniverseDtoResponse,
			expectedContent: defaultUniverseDtoResponse,
		},
		"getUniverse": {
			req:             httptest.NewRequest(http.MethodGet, "/", nil),
			idAsRouteParam:  true,
			handler:         getUniverse,
			universeDto:     defaultUniverseDtoResponse,
			expectedContent: defaultUniverseDtoResponse,
		},
		"listUniverses": {
			req:             httptest.NewRequest(http.MethodGet, "/", nil),
			handler:         listUniverses,
			universeDto:     defaultUniverseDtoResponse,
			expectedContent: []communication.UniverseDtoResponse{defaultUniverseDtoResponse},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockUniverseService{
				universes: []communication.UniverseDtoResponse{testCase.universeDto},
			}

			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			err := testCase.handler(ctx, mock)

			assert.Nil(err)

			actual := strings.Trim(rw.Body.String(), "\n")
			expected, err := json.Marshal(testCase.expectedContent)
			assert.Nil(err)
			assert.Equal(string(expected), actual)
		})
	}
}

func TestCreateUniverse_CallsServiceCreate(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateTestEchoContextFromRequest(generateTestPostRequest())
	ms := &mockUniverseService{}

	err := createUniverse(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.createCalled)
}

func TestCreateUniverse_SavesExpectedUniverse(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateTestEchoContextFromRequest(generateTestRequestWithDefaultUniverseBody(http.MethodPost))
	ms := &mockUniverseService{
		universes: []communication.UniverseDtoResponse{defaultUniverseDtoResponse},
	}

	err := createUniverse(ctx, ms)

	assert.Nil(err)
	assert.Equal(defaultUniverseDtoRequest, ms.inUniverse)
}

func TestGetUniverse_CallsServiceGet(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodGet)
	ms := &mockUniverseService{}

	err := getUniverse(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.getCalled)
}

func TestGetUniverse_GetsExpectedUniverse(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodGet)
	ms := &mockUniverseService{}

	err := getUniverse(ctx, ms)

	assert.Nil(err)
	assert.Equal(defaultUuid, ms.inId)
}

func TestListUniverse_CallsServiceList(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateTestEchoContextWithMethod(http.MethodGet)
	ms := &mockUniverseService{}

	err := listUniverses(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.listCalled)
}

func TestDeleteUniverse_CallsServiceDelete(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodDelete)
	ms := &mockUniverseService{}

	err := deleteUniverse(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.deleteCalled)
}

func TestDeleteUniverse_DeletesExpectedUniverse(t *testing.T) {
	assert := assert.New(t)

	ctx, _ := generateEchoContextWithValidUuid(http.MethodDelete)
	ms := &mockUniverseService{}

	err := deleteUniverse(ctx, ms)

	assert.Nil(err)
	assert.Equal(defaultUuid, ms.inId)
}

// func generateTestPostRequest() *http.Request {
// 	return generateTestRequestWithDefaultUserBody(http.MethodPost)
// }

func generateTestRequestWithDefaultUniverseBody(method string) *http.Request {
	return generateTestRequestWithUniverseBody(method, defaultUniverseDtoRequest)
}

func generateTestRequestWithUniverseBody(method string, universeDto communication.UniverseDtoRequest) *http.Request {
	// Voluntarily ignoring errors
	raw, _ := json.Marshal(universeDto)
	req := httptest.NewRequest(method, "/", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	return req
}

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

	var out communication.UniverseDtoResponse
	if m.universes != nil {
		out = m.universes[0]
	}
	return out, m.err
}

func (m *mockUniverseService) Get(ctx context.Context, id uuid.UUID) (communication.UniverseDtoResponse, error) {
	m.getCalled++
	m.inId = id

	var out communication.UniverseDtoResponse
	if m.universes != nil {
		out = m.universes[0]
	}
	return out, m.err
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
