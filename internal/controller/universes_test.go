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

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/galactic-sovereign/internal/service"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mockUniverseService struct {
	universes []communication.FullUniverseDtoResponse
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
var defaultResourceDtoResponse = communication.ResourceDtoResponse{
	Id:   defaultUuid,
	Name: "my-resource",

	CreatedAt: time.Date(2024, 8, 3, 14, 29, 31, 651387240, time.UTC),
}
var defaultBuildingDtoResponse = communication.BuildingDtoResponse{
	Id:   defaultUuid,
	Name: "my-building",

	CreatedAt: time.Date(2024, 8, 8, 21, 42, 03, 651387242, time.UTC),
}
var defaultBuildingCostDtoResponse = communication.BuildingCostDtoResponse{
	Building: defaultUuid,
	Resource: defaultUuid,
	Cost:     58,
}
var defaultFullUniverseDtoResponse = communication.FullUniverseDtoResponse{
	UniverseDtoResponse: defaultUniverseDtoResponse,
	Resources: []communication.ResourceDtoResponse{
		defaultResourceDtoResponse,
	},
	Buildings: []communication.FullBuildingDtoResponse{
		{
			BuildingDtoResponse: defaultBuildingDtoResponse,
			Costs: []communication.BuildingCostDtoResponse{
				defaultBuildingCostDtoResponse,
			},
		},
	},
}

func TestUnit_UniverseEndpoints(t *testing.T) {
	s := RouteTestSuite{
		generateRoutes: func() rest.Routes {
			return UniverseEndpoints(&mockUniverseService{})
		},
		expectedRoutes: map[string]int{
			http.MethodPost:   1,
			http.MethodGet:    2,
			http.MethodDelete: 1,
		},
		expectedPaths: map[string]int{
			"/universes":     2,
			"/universes/:id": 2,
		},
	}

	suite.Run(t, &s)
}

func TestUnit_UniverseController(t *testing.T) {
	s := ControllerTestSuite[service.UniverseService]{
		generateServiceMock:      generateUniverseServiceMock,
		generateErrorServiceMock: generateErrorUniverseServiceMock,

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
				err:                errors.NewCode(pgx.UniqueConstraintViolation),
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
				err:                errors.NewCode(db.NoMatchingRows),
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
				err:                errors.NewCode(db.NoMatchingRows),
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
				expectedContent: defaultFullUniverseDtoResponse,
			},
			"listUniverses": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				handler:         listUniverses,
				expectedContent: []communication.UniverseDtoResponse{defaultUniverseDtoResponse},
			},
			"listUniverses_noData": {
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				generateServiceMock: func() service.UniverseService {
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

				verifyInteractions: func(s service.UniverseService, assert *require.Assertions) {
					m := assertUniverseServiceIsAMock(s, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultUniverseDtoRequest, m.inUniverse)
				},
			},
			"getUniverse": {
				req:            httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam: true,
				handler:        getUniverse,

				verifyInteractions: func(s service.UniverseService, assert *require.Assertions) {
					m := assertUniverseServiceIsAMock(s, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
			"listUniverses": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: listUniverses,

				verifyInteractions: func(s service.UniverseService, assert *require.Assertions) {
					m := assertUniverseServiceIsAMock(s, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"deleteUniverse": {
				req:            httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam: true,
				handler:        deleteUniverse,

				verifyInteractions: func(s service.UniverseService, assert *require.Assertions) {
					m := assertUniverseServiceIsAMock(s, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateUniverseServiceMock() service.UniverseService {
	return &mockUniverseService{
		universes: []communication.FullUniverseDtoResponse{defaultFullUniverseDtoResponse},
	}
}

func generateErrorUniverseServiceMock(err error) service.UniverseService {
	return &mockUniverseService{
		err: err,
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

	var out communication.FullUniverseDtoResponse
	if m.universes != nil {
		out = m.universes[0]
	}
	return out.UniverseDtoResponse, m.err
}

func (m *mockUniverseService) Get(ctx context.Context, id uuid.UUID) (communication.FullUniverseDtoResponse, error) {
	m.getCalled++
	m.inId = id

	var out communication.FullUniverseDtoResponse
	if m.universes != nil {
		out = m.universes[0]
	}
	return out, m.err
}

func (m *mockUniverseService) List(ctx context.Context) ([]communication.UniverseDtoResponse, error) {
	m.listCalled++

	var out []communication.UniverseDtoResponse
	for _, fullDto := range m.universes {
		out = append(out, fullDto.UniverseDtoResponse)
	}

	return out, m.err
}

func (m *mockUniverseService) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	m.inId = id
	return m.err
}
