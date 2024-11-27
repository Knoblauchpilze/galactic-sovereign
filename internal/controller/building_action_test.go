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

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/rest"
	"github.com/KnoblauchPilze/galactic-sovereign/internal/service"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/communication"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/game"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mockBuildingActionService struct {
	action communication.BuildingActionDtoResponse
	err    error

	createCalled int
	inAction     communication.BuildingActionDtoRequest

	deleteCalled int
	deleteId     uuid.UUID
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

func TestUnit_BuildingActionEndpoints(t *testing.T) {
	s := RouteTestSuite{
		generateRoutes: func() rest.Routes {
			return BuildingActionEndpoints(&mockBuildingActionService{}, &mockActionService{}, &mockPlanetResourceService{})
		},
		expectedRoutes: map[string]int{
			http.MethodPost:   1,
			http.MethodDelete: 1,
		},
		expectedPaths: map[string]int{
			"/planets/:id/actions": 1,
			"/actions/:id":         1,
		},

		errorTestCases: map[string]routeErrorTestCase{
			"whenActionServiceFails": {
				generateRoutes: func() rest.Routes {
					m := &mockActionService{
						err: errDefault,
					}

					return generateBuildingActionRoutesUsingGameUpdateWatcher(m, &mockPlanetResourceService{})
				},
				expectedStatusCode: http.StatusInternalServerError,
				expectedError:      "\"Failed to process actions\"\n",
			},
			"whenPlanetResourceServiceFails": {
				generateRoutes: func() rest.Routes {
					m := &mockPlanetResourceService{
						err: errDefault,
					}

					return generateBuildingActionRoutesUsingGameUpdateWatcher(&mockActionService{}, m)
				},
				expectedStatusCode: http.StatusInternalServerError,
				expectedError:      "\"Failed to update resources\"\n",
			},
		},

		interactionTestCases: []routeInteractionTestCase{
			{
				generateRoutes: func(actionService game.ActionService, planetResourceService game.PlanetResourceService) rest.Routes {
					return generateBuildingActionRoutesUsingGameUpdateWatcher(actionService, planetResourceService)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func TestUnit_BuildingActionController(t *testing.T) {
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
			"deleteBuildingAction": {
				req:     httptest.NewRequest(http.MethodDelete, "/", nil),
				handler: deleteBuildingAction,
			},
		},

		badIdTestCases: map[string]badIdTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:     httptest.NewRequest(http.MethodPost, "/", nil),
				handler: createBuildingAction,
			},
			"deleteBuildingAction": {
				req:     httptest.NewRequest(http.MethodDelete, "/", nil),
				handler: deleteBuildingAction,
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
			"deleteBuildingAction": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deleteBuildingAction,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
			"deleteBuildingAction_notFound": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deleteBuildingAction,
				err:                errors.NewCode(db.NoMatchingSqlRows),
				expectedHttpStatus: http.StatusNotFound,
			},
		},

		successTestCases: map[string]successTestCase[service.BuildingActionService]{
			"createBuildingAction": {
				req:                generateTestRequestWithBuildingActionBody(http.MethodPost, defaultBuildingActionDtoRequest),
				idAsRouteParam:     true,
				handler:            createBuildingAction,
				expectedHttpStatus: http.StatusCreated,
			},
			"deleteBuildingAction": {
				req:                httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam:     true,
				handler:            deleteBuildingAction,
				expectedHttpStatus: http.StatusNoContent,
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
			"deleteBuildingAction": {
				req:            httptest.NewRequest(http.MethodDelete, "/", nil),
				idAsRouteParam: true,
				handler:        deleteBuildingAction,

				verifyInteractions: func(s service.BuildingActionService, assert *require.Assertions) {
					m := assertBuildingActionServiceIsAMock(s, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUuid, m.deleteId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateBuildingActionRoutesUsingGameUpdateWatcher(actionService game.ActionService, planetResourceService game.PlanetResourceService) rest.Routes {
	allRoutes := BuildingActionEndpoints(&mockBuildingActionService{}, actionService, planetResourceService)

	var routes rest.Routes
	for _, route := range allRoutes {
		isDelete := route.Method() == http.MethodDelete

		if !isDelete {
			routes = append(routes, route)
		}
	}

	return routes
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

func (m *mockBuildingActionService) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id

	return m.err
}
