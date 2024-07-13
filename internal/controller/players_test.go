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

type mockPlayerService struct {
	players []communication.PlayerDtoResponse
	err     error

	createCalled int
	getCalled    int
	listCalled   int
	deleteCalled int

	inPlayer communication.PlayerDtoRequest
	inId     uuid.UUID
}

var defaultPlayerId = uuid.MustParse("bd7cb2c0-2124-4c1b-8ff8-2d3eb928ffa9")
var defaultUniverseId = uuid.MustParse("6a2b0061-360d-4cb0-92a7-ed0486499b92")
var defaultPlayerDtoRequest = communication.PlayerDtoRequest{
	ApiUser:  defaultPlayerId,
	Universe: defaultUniverseId,
	Name:     "my-player",
}
var defaultPlayerDtoResponse = communication.PlayerDtoResponse{
	Id:       defaultUuid,
	ApiUser:  defaultPlayerId,
	Universe: defaultUniverseId,
	Name:     "my-player",

	CreatedAt: time.Date(2024, 07, 13, 14, 42, 50, 651387235, time.UTC),
}

func TestPlayerEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range PlayerEndpoints(&mockPlayerService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(3, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodPost])
	assert.Equal(2, actualRoutes[http.MethodGet])
	assert.Equal(1, actualRoutes[http.MethodDelete])
}

func Test_PlayerController(t *testing.T) {
	s := ControllerTestSuite[service.PlayerService]{
		generateServiceMock:      generatePlayerServiceMock,
		generateValidServiceMock: generateValidPlayerServiceMock,

		badInputTestCases: map[string]badInputTestCase[service.PlayerService]{
			"createPlayer": {
				req:                httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-dto-request")),
				handler:            createPlayer,
				expectedBodyString: "\"Invalid player syntax\"\n",
			},
		},

		noIdTestCases: map[string]noIdTestCase[service.PlayerService]{
			"getPlayer": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: getPlayer,
			},
			"deletePlayer": {
				req:     httptest.NewRequest(http.MethodDelete, "/", nil),
				handler: deletePlayer,
			},
		},

		badIdTestCases: map[string]badIdTestCase[service.PlayerService]{
			"getPlayer": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: getPlayer,
			},
			"deletePlayer": {
				req:     httptest.NewRequest(http.MethodDelete, "/", nil),
				handler: deletePlayer,
			},
		},

		errorTestCases: map[string]errorTestCase[service.PlayerService]{
			"createPlayer": {
				req:                generateTestRequestWithDefaultPlayerBody(http.MethodPost),
				handler:            createPlayer,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"getPlayer": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:     true,
				handler:            getPlayer,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"getPlayer_notFound": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:     true,
				handler:            getPlayer,
				err:                errors.NewCode(db.NoMatchingSqlRows),
				expectedHttpStatus: http.StatusNotFound,
			},
			"listPlayers": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				handler:            listPlayers,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"deletePlayer": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deletePlayer,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"deletePlayer_notFound": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deletePlayer,
				err:                errors.NewCode(db.NoMatchingSqlRows),
				expectedHttpStatus: http.StatusNotFound,
			},
		},

		successTestCases: map[string]successTestCase[service.PlayerService]{
			"createPlayer": {
				req:                generateTestRequestWithDefaultPlayerBody(http.MethodPost),
				handler:            createPlayer,
				expectedHttpStatus: http.StatusCreated,
			},
			"getPlayer": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:     true,
				handler:            getPlayer,
				expectedHttpStatus: http.StatusOK,
			},
			"listPlayers": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				handler:            listPlayers,
				expectedHttpStatus: http.StatusOK,
			},
			"deletePlayer": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deletePlayer,
				expectedHttpStatus: http.StatusNoContent,
			},
		},

		returnTestCases: map[string]returnTestCase[service.PlayerService]{
			"createPlayer": {
				req:             generateTestRequestWithDefaultPlayerBody(http.MethodPost),
				handler:         createPlayer,
				expectedContent: defaultPlayerDtoResponse,
			},
			"getPlayer": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam:  true,
				handler:         getPlayer,
				expectedContent: defaultPlayerDtoResponse,
			},
			"listPlayers": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				handler:         listPlayers,
				expectedContent: []communication.PlayerDtoResponse{defaultPlayerDtoResponse},
			},
			"listPlayers_noData": {
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				generateValidServiceMock: func() service.PlayerService {
					return &mockPlayerService{
						players: nil,
					}
				},

				handler:         listPlayers,
				expectedContent: []communication.PlayerDtoResponse{},
			},
		},

		serviceInteractionTestCases: map[string]serviceInteractionTestCase[service.PlayerService]{
			"createPlayer": {
				req:     generateTestRequestWithDefaultPlayerBody(http.MethodPost),
				handler: createPlayer,

				verifyInteractions: func(s service.PlayerService, assert *require.Assertions) {
					m := assertPlayerServiceIsAMock(s, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultPlayerDtoRequest, m.inPlayer)
				},
			},
			"getPlayer": {
				req:            httptest.NewRequest(http.MethodGet, "/", nil),
				idAsRouteParam: true,
				handler:        getPlayer,

				verifyInteractions: func(s service.PlayerService, assert *require.Assertions) {
					m := assertPlayerServiceIsAMock(s, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
			"listPlayers": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: listPlayers,

				verifyInteractions: func(s service.PlayerService, assert *require.Assertions) {
					m := assertPlayerServiceIsAMock(s, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"deletePlayer": {
				req:            httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam: true,
				handler:        deletePlayer,

				verifyInteractions: func(s service.PlayerService, assert *require.Assertions) {
					m := assertPlayerServiceIsAMock(s, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUuid, m.inId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generatePlayerServiceMock(err error) service.PlayerService {
	return &mockPlayerService{
		err: err,
	}
}

func generateValidPlayerServiceMock() service.PlayerService {
	return &mockPlayerService{
		players: []communication.PlayerDtoResponse{defaultPlayerDtoResponse},
	}
}

func assertPlayerServiceIsAMock(s service.PlayerService, assert *require.Assertions) *mockPlayerService {
	m, ok := s.(*mockPlayerService)
	if !ok {
		assert.Fail("Provided player service is not a mock")
	}
	return m
}

func generateTestRequestWithDefaultPlayerBody(method string) *http.Request {
	return generateTestRequestWithPlayerBody(method, defaultPlayerDtoRequest)
}

func generateTestRequestWithPlayerBody(method string, playerDto communication.PlayerDtoRequest) *http.Request {
	// Voluntarily ignoring errors
	raw, _ := json.Marshal(playerDto)
	req := httptest.NewRequest(method, "/", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (m *mockPlayerService) Create(ctx context.Context, player communication.PlayerDtoRequest) (communication.PlayerDtoResponse, error) {
	m.createCalled++
	m.inPlayer = player

	var out communication.PlayerDtoResponse
	if m.players != nil {
		out = m.players[0]
	}
	return out, m.err
}

func (m *mockPlayerService) Get(ctx context.Context, id uuid.UUID) (communication.PlayerDtoResponse, error) {
	m.getCalled++
	m.inId = id

	var out communication.PlayerDtoResponse
	if m.players != nil {
		out = m.players[0]
	}
	return out, m.err
}

func (m *mockPlayerService) List(ctx context.Context) ([]communication.PlayerDtoResponse, error) {
	m.listCalled++
	return m.players, m.err
}

func (m *mockPlayerService) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	m.inId = id
	return m.err
}
