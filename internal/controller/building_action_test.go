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
	"github.com/KnoblauchPilze/user-service/pkg/game"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mockBuildingActionService struct {
	action communication.BuildingActionDtoResponse
	err    error

	createCalled int

	inAction communication.BuildingActionDtoRequest
}

var defaultBuildingActionId = uuid.MustParse("694a47ab-cd58-431e-9298-e0e788bfc01e")
var defaultBuildingActionDtoRequest = communication.BuildingActionDtoRequest{
	Building: defaultBuildingId,
}
var defaultBuildingActionDtoResponse = communication.BuildingActionDtoResponse{
	Id:           defaultBuildingActionId,
	Planet:       defaultPlanetId,
	Building:     defaultBuildingId,
	CurrentLevel: 14,
	DesiredLevel: 78,
	CreatedAt:    time.Date(2024, 8, 11, 14, 12, 31, 651387243, time.UTC),
	CompletedAt:  time.Date(2024, 8, 11, 14, 12, 36, 651387243, time.UTC),
}

func TestBuildingActionEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range BuildingActionEndpoints(&mockBuildingActionService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(1, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodPost])
}

func Test_BuildingActionController(t *testing.T) {
	s := ControllerTestSuite[service.BuildingActionService]{
		generateServiceMock:      generateBuildingActionServiceMock,
		generateErrorServiceMock: generateErrorBuildingActionServiceMock,

		badInputTestCases: map[string]badInputTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:                httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not-a-dto-request")),
				idAsRouteParam:     true,
				handler:            createBuildingAction,
				expectedBodyString: "\"Invalid action syntax\"\n",
			},
		},

		noIdTestCases: map[string]noIdTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:     httptest.NewRequest(http.MethodPost, "/", nil),
				handler: createBuildingAction,
			},
		},

		badIdTestCases: map[string]badIdTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:     httptest.NewRequest(http.MethodPost, "/", nil),
				handler: createBuildingAction,
			},
		},

		errorTestCases: map[string]errorTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:                generateTestRequestWithBuildingActionBody(http.MethodPost, defaultBuildingActionDtoRequest),
				idAsRouteParam:     true,
				handler:            createBuildingAction,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"createBuildingAction_notEnoughResources": {
				req:                generateTestRequestWithBuildingActionBody(http.MethodPost, defaultBuildingActionDtoRequest),
				idAsRouteParam:     true,
				handler:            createBuildingAction,
				err:                errors.NewCode(game.NotEnoughResources),
				expectedHttpStatus: http.StatusBadRequest,
			},
			"createBuildingAction_buildingActionAlreadyInProgress": {
				req:                generateTestRequestWithBuildingActionBody(http.MethodPost, defaultBuildingActionDtoRequest),
				idAsRouteParam:     true,
				handler:            createBuildingAction,
				err:                errors.NewCode(db.DuplicatedKeySqlKey),
				expectedHttpStatus: http.StatusConflict,
			},
		},

		successTestCases: map[string]successTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:                generateTestRequestWithBuildingActionBody(http.MethodPost, defaultBuildingActionDtoRequest),
				idAsRouteParam:     true,
				handler:            createBuildingAction,
				expectedHttpStatus: http.StatusCreated,
			},
		},

		returnTestCases: map[string]returnTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:             generateTestRequestWithBuildingActionBody(http.MethodPost, defaultBuildingActionDtoRequest),
				idAsRouteParam:  true,
				handler:         createBuildingAction,
				expectedContent: defaultBuildingActionDtoResponse,
			},
		},

		serviceInteractionTestCases: map[string]serviceInteractionTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:            generateTestRequestWithBuildingActionBody(http.MethodPost, defaultBuildingActionDtoRequest),
				idAsRouteParam: true,
				handler:        createBuildingAction,

				verifyInteractions: func(s service.BuildingActionService, assert *require.Assertions) {
					m := assertBuildingActionServiceIsAMock(s, assert)

					assert.Equal(1, m.createCalled)
					actual := m.inAction
					assert.Equal(defaultUuid, actual.Planet)
					assert.Equal(defaultBuildingId, actual.Building)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateBuildingActionServiceMock() service.BuildingActionService {
	return &mockBuildingActionService{
		action: defaultBuildingActionDtoResponse,
	}
}

func generateErrorBuildingActionServiceMock(err error) service.BuildingActionService {
	return &mockBuildingActionService{
		err: err,
	}
}

func assertBuildingActionServiceIsAMock(s service.BuildingActionService, assert *require.Assertions) *mockBuildingActionService {
	m, ok := s.(*mockBuildingActionService)
	if !ok {
		assert.Fail("Provided building action service is not a mock")
	}
	return m
}

func generateTestRequestWithBuildingActionBody(method string, actionDto communication.BuildingActionDtoRequest) *http.Request {
	// Voluntarily ignoring errors
	raw, _ := json.Marshal(actionDto)
	req := httptest.NewRequest(method, "/", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (m *mockBuildingActionService) Create(ctx context.Context, actionDto communication.BuildingActionDtoRequest) (communication.BuildingActionDtoResponse, error) {
	m.createCalled++
	m.inAction = actionDto

	return m.action, m.err
}
