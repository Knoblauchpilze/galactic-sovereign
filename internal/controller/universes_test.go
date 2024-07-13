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

func Test_UniverseController(t *testing.T) {
	s := ControllerTestSuite[service.UniverseService]{
		generateServiceMock:      generateUniverseServiceMock,
		generateValidServiceMock: generateValidUniverseServiceMock,

		badInputTestCases: map[string]badInputTestCase[service.UniverseService]{
			"createUniverse": {
				req:                httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-dto-request")),
				handler:            createUniverse,
				expectedBodyString: "\"Invalid universe syntax\"\n",
			},
		},

		noIdTestCases: map[string]noIdTestCase[service.UniverseService]{
			"getUniverse": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: getUniverse,
			},
			"deleteUniverse": {
				req:     httptest.NewRequest(http.MethodDelete, "/", nil),
				handler: deleteUniverse,
			},
		},

		badIdTestCases: map[string]badIdTestCase[service.UniverseService]{
			"getUniverse": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: getUniverse,
			},
			"deleteUniverse": {
				req:     httptest.NewRequest(http.MethodDelete, "/", nil),
				handler: deleteUniverse,
			},
		},

		errorTestCases: map[string]errorTestCase[service.UniverseService]{
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
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
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
		},

		successTestCases: map[string]successTestCase[service.UniverseService]{
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
		},

		returnTestCases: map[string]returnTestCase[service.UniverseService]{
			"createUniverse": {
				req:             generateTestRequestWithDefaultUniverseBody(http.MethodPost),
				handler:         createUniverse,
				expectedContent: defaultUniverseDtoResponse,
			},
			"getUniverse": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:  true,
				handler:         getUniverse,
				expectedContent: defaultUniverseDtoResponse,
			},
			"listUniverses": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				handler:         listUniverses,
				expectedContent: []communication.UniverseDtoResponse{defaultUniverseDtoResponse},
			},
			"listUniverses_noData": {
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				generateValidServiceMock: func() service.UniverseService {
					return &mockUniverseService{
						universes: nil,
					}
				},

				handler:         listUniverses,
				expectedContent: []communication.UniverseDtoResponse{},
			},
		},

		serviceInteractionTestCases: map[string]serviceInteractionTestCase[service.UniverseService]{
			"createUniverse": {
				req:     generateTestRequestWithDefaultUniverseBody(http.MethodPost),
				handler: createUniverse,

				verifyInteractions: func(us service.UniverseService, assert *require.Assertions) {
					m := assertUniverseServiceIsAMock(us, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultUniverseDtoRequest, m.inUniverse)
				},
			},
			"getUniverse": {
				req:            httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam: true,
				handler:        getUniverse,

				verifyInteractions: func(us service.UniverseService, assert *require.Assertions) {
					m := assertUniverseServiceIsAMock(us, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
			"listUniverses": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: listUniverses,

				verifyInteractions: func(us service.UniverseService, assert *require.Assertions) {
					m := assertUniverseServiceIsAMock(us, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"deleteUniverse": {
				req:            httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam: true,
				handler:        deleteUniverse,

				verifyInteractions: func(us service.UniverseService, assert *require.Assertions) {
					m := assertUniverseServiceIsAMock(us, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateUniverseServiceMock(err error) service.UniverseService {
	return &mockUniverseService{
		err: err,
	}
}

func generateValidUniverseServiceMock() service.UniverseService {
	return &mockUniverseService{
		universes: []communication.UniverseDtoResponse{defaultUniverseDtoResponse},
	}
}

func assertUniverseServiceIsAMock(s service.UniverseService, assert *require.Assertions) *mockUniverseService {
	m, ok := s.(*mockUniverseService)
	if !ok {
		assert.Fail("Provided universe service is not a mock")
	}
	return m
}

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
