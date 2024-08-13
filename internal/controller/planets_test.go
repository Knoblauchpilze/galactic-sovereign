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

type mockPlanetService struct {
	planets []communication.FullPlanetDtoResponse
	err     error

	createCalled        int
	getCalled           int
	listCalled          int
	listForPlayerCalled int
	deleteCalled        int

	inPlanet   communication.PlanetDtoRequest
	inId       uuid.UUID
	inPlayerId uuid.UUID
}

var defaultPlanetId = uuid.MustParse("080f5a2b-800a-458d-9806-7660bde4db00")
var defaultResourceId = uuid.MustParse("84e01480-5433-4ebe-a078-7e1cd18c86c4")
var defaultBuildingId = uuid.MustParse("ba846861-b015-4726-b9a8-3fe1cf2621e7")
var defaultPlanetDtoRequest = communication.PlanetDtoRequest{
	Player: defaultPlayerId,
	Name:   "my-planet",
}
var defaultPlanetDtoResponse = communication.PlanetDtoResponse{
	Id:     defaultUuid,
	Player: defaultPlayerId,
	Name:   "my-planet",

	CreatedAt: time.Date(2024, 07, 13, 10, 53, 10, 651387238, time.UTC),
}
var defaultPlanetResourceDtoResponse = communication.PlanetResourceDtoResponse{
	Planet:   defaultPlanetId,
	Resource: defaultResourceId,
	Amount:   654.321,

	CreatedAt: time.Date(2024, 07, 31, 18, 10, 30, 651387238, time.UTC),
	UpdatedAt: time.Date(2024, 07, 31, 18, 10, 59, 651387238, time.UTC),
}
var defaultPlanetBuildingDtoResponse = communication.PlanetBuildingDtoResponse{
	Planet:   defaultPlanetId,
	Building: defaultBuildingId,
	Level:    36,

	CreatedAt: time.Date(2024, 8, 6, 21, 45, 30, 651387241, time.UTC),
	UpdatedAt: time.Date(2024, 8, 6, 21, 45, 33, 651387241, time.UTC),
}
var defaultFullPlanetDtoResponse = communication.FullPlanetDtoResponse{
	PlanetDtoResponse: defaultPlanetDtoResponse,
	Resources: []communication.PlanetResourceDtoResponse{
		defaultPlanetResourceDtoResponse,
	},
	Buildings: []communication.PlanetBuildingDtoResponse{
		defaultPlanetBuildingDtoResponse,
	},
	BuildingActions: []communication.BuildingActionDtoResponse{
		defaultBuildingActionDtoResponse,
	},
}

func TestPlanetEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range PlanetEndpoints(&mockPlanetService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(3, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodPost])
	assert.Equal(2, actualRoutes[http.MethodGet])
	assert.Equal(1, actualRoutes[http.MethodDelete])
}

func Test_PlanetController(t *testing.T) {
	s := ControllerTestSuite[service.PlanetService]{
		generateServiceMock:      generatePlanetServiceMock,
		generateErrorServiceMock: generateErrorPlanetServiceMock,

		badInputTestCases: map[string]badInputTestCase[service.PlanetService]{
			"createPlanet": {
				req:                httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-dto-request")),
				handler:            createPlanet,
				expectedBodyString: "\"Invalid planet syntax\"\n",
			},
		},

		noIdTestCases: map[string]noIdTestCase[service.PlanetService]{
			"getPlanet": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: getPlanet,
			},
			"deletePlanet": {
				req:     httptest.NewRequest(http.MethodDelete, "/", nil),
				handler: deletePlanet,
			},
		},

		badIdTestCases: map[string]badIdTestCase[service.PlanetService]{
			"getPlanet": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: getPlanet,
			},
			"listPlanets_badPlayerId": {
				req:     generateTestRequestWithPlayerAsQueryParam("not-a-uuid"),
				handler: listPlanets,
			},
			"deletePlanet": {
				req:     httptest.NewRequest(http.MethodDelete, "/", nil),
				handler: deletePlanet,
			},
		},

		errorTestCases: map[string]errorTestCase[service.PlanetService]{
			"createPlanet": {
				req:                generateTestRequestWithDefaultPlanetBody(http.MethodPost),
				handler:            createPlanet,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"getPlanet": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:     true,
				handler:            getPlanet,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"getPlanet_notFound": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:     true,
				handler:            getPlanet,
				err:                errors.NewCode(db.NoMatchingSqlRows),
				expectedHttpStatus: http.StatusNotFound,
			},
			"listPlanets": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				handler:            listPlanets,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"deletePlanet": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deletePlanet,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"deletePlanet_notFound": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deletePlanet,
				err:                errors.NewCode(db.NoMatchingSqlRows),
				expectedHttpStatus: http.StatusNotFound,
			},
		},

		successTestCases: map[string]successTestCase[service.PlanetService]{
			"createPlanet": {
				req:                generateTestRequestWithDefaultPlanetBody(http.MethodPost),
				handler:            createPlanet,
				expectedHttpStatus: http.StatusCreated,
			},
			"getPlanet": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:     true,
				handler:            getPlanet,
				expectedHttpStatus: http.StatusOK,
			},
			"listPlanets": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				handler:            listPlanets,
				expectedHttpStatus: http.StatusOK,
			},
			"deletePlanet": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deletePlanet,
				expectedHttpStatus: http.StatusNoContent,
			},
		},

		returnTestCases: map[string]returnTestCase[service.PlanetService]{
			"createPlanet": {
				req:             generateTestRequestWithDefaultPlanetBody(http.MethodPost),
				handler:         createPlanet,
				expectedContent: defaultPlanetDtoResponse,
			},
			"getPlanet": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:  true,
				handler:         getPlanet,
				expectedContent: defaultFullPlanetDtoResponse,
			},
			"listPlanets": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				handler:         listPlanets,
				expectedContent: []communication.PlanetDtoResponse{defaultPlanetDtoResponse},
			},
			"listPlanets_noData": {
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				generateServiceMock: func() service.PlanetService {
					return &mockPlanetService{
						planets: nil,
					}
				},
				handler:         listPlanets,
				expectedContent: []communication.PlanetDtoResponse{},
			},
		},

		serviceInteractionTestCases: map[string]serviceInteractionTestCase[service.PlanetService]{
			"createPlanet": {
				req:     generateTestRequestWithDefaultPlanetBody(http.MethodPost),
				handler: createPlanet,

				verifyInteractions: func(s service.PlanetService, assert *require.Assertions) {
					m := assertPlanetServiceIsAMock(s, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultPlanetDtoRequest, m.inPlanet)
				},
			},
			"getPlanet": {
				req:            httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam: true,
				handler:        getPlanet,

				verifyInteractions: func(s service.PlanetService, assert *require.Assertions) {
					m := assertPlanetServiceIsAMock(s, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
			"listPlanets": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: listPlanets,

				verifyInteractions: func(s service.PlanetService, assert *require.Assertions) {
					m := assertPlanetServiceIsAMock(s, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"listPlanets_withPlayerId": {
				req:     generateTestRequestWithPlayerAsQueryParam(defaultPlayerId.String()),
				handler: listPlanets,

				verifyInteractions: func(s service.PlanetService, assert *require.Assertions) {
					m := assertPlanetServiceIsAMock(s, assert)

					assert.Equal(defaultPlayerId, m.inPlayerId)
					assert.Equal(1, m.listForPlayerCalled)
				},
			},
			"deletePlanet": {
				req:            httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam: true,
				handler:        deletePlanet,

				verifyInteractions: func(s service.PlanetService, assert *require.Assertions) {
					m := assertPlanetServiceIsAMock(s, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generatePlanetServiceMock() service.PlanetService {
	return &mockPlanetService{
		planets: []communication.FullPlanetDtoResponse{defaultFullPlanetDtoResponse},
	}
}

func generateErrorPlanetServiceMock(err error) service.PlanetService {
	return &mockPlanetService{
		err: err,
	}
}

func assertPlanetServiceIsAMock(s service.PlanetService, assert *require.Assertions) *mockPlanetService {
	m, ok := s.(*mockPlanetService)
	if !ok {
		assert.Fail("Provided planet service is not a mock")
	}
	return m
}

func generateTestRequestWithDefaultPlanetBody(method string) *http.Request {
	return generateTestRequestWithPlanetBody(method, defaultPlanetDtoRequest)
}

func generateTestRequestWithPlanetBody(method string, planetDto communication.PlanetDtoRequest) *http.Request {
	// Voluntarily ignoring errors
	raw, _ := json.Marshal(planetDto)
	req := httptest.NewRequest(method, "/", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func generateTestRequestWithPlayerAsQueryParam(playerId string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	query := req.URL.Query()
	query.Add("player", playerId)
	req.URL.RawQuery = query.Encode()
	return req
}

func (m *mockPlanetService) Create(ctx context.Context, planet communication.PlanetDtoRequest) (communication.PlanetDtoResponse, error) {
	m.createCalled++
	m.inPlanet = planet

	var out communication.FullPlanetDtoResponse
	if m.planets != nil {
		out = m.planets[0]
	}
	return out.PlanetDtoResponse, m.err
}

func (m *mockPlanetService) Get(ctx context.Context, id uuid.UUID) (communication.FullPlanetDtoResponse, error) {
	m.getCalled++
	m.inId = id

	var out communication.FullPlanetDtoResponse
	if m.planets != nil {
		out = m.planets[0]
	}
	return out, m.err
}

func (m *mockPlanetService) List(ctx context.Context) ([]communication.PlanetDtoResponse, error) {
	m.listCalled++
	var out []communication.PlanetDtoResponse
	for _, planet := range m.planets {
		out = append(out, planet.PlanetDtoResponse)
	}
	return out, m.err
}

func (m *mockPlanetService) ListForPlayer(ctx context.Context, player uuid.UUID) ([]communication.PlanetDtoResponse, error) {
	m.listForPlayerCalled++
	m.inPlayerId = player
	var out []communication.PlanetDtoResponse
	for _, planet := range m.planets {
		out = append(out, planet.PlanetDtoResponse)
	}
	return out, m.err
}

func (m *mockPlanetService) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	m.inId = id
	return m.err
}
