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

func TestBuildingActionEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range BuildingActionEndpoints(&mockBuildingActionService{}, &mockActionService{}, &mockPlanetResourceService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(2, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodPost])
	assert.Equal(1, actualRoutes[http.MethodDelete])
}

func TestBuildingActionEndpoints_WhenActionServiceFails_SetsReturnStatusInternalError(t *testing.T) {
	assert := assert.New(t)

	m := &mockActionService{
		err: errDefault,
	}

	for _, route := range BuildingActionEndpoints(&mockBuildingActionService{}, m, &mockPlanetResourceService{}) {
		ctx, rw := generateTestEchoContextWithMethodAndId(http.MethodGet)

		handler := route.Handler()
		err := handler(ctx)

		assert.Nil(err)
		assert.Equal(http.StatusInternalServerError, rw.Code)
		assert.Equal("\"Failed to process actions\"\n", rw.Body.String())
	}
}

func TestBuildingActionEndpoints_WhenNoPlanetId_ExpectPlanetResourceNotUpdated(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{}

	for _, route := range BuildingActionEndpoints(&mockBuildingActionService{}, &mockActionService{}, m) {
		ctx, _ := generateTestEchoContextWithMethod(http.MethodGet)

		handler := route.Handler()
		err := handler(ctx)

		assert.Nil(err)
		assert.Equal(0, m.updatePlanetUntilCalled)
	}
}

func TestBuildingActionEndpoints_WhenPlanetIdIsInvalid_ExpectPlanetResourceNotUpdated(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{}

	for _, route := range BuildingActionEndpoints(&mockBuildingActionService{}, &mockActionService{}, m) {
		ctx, _ := generateTestEchoContextWithMethod(http.MethodGet)
		ctx.SetParamNames("id")
		ctx.SetParamValues("not-a-uuid")

		handler := route.Handler()
		err := handler(ctx)

		assert.Nil(err)
		assert.Equal(0, m.updatePlanetUntilCalled)
	}
}
func TestBuildingActionEndpoints_WhenPlanetIdValid_ExpectPlanetResourceAreUpdated(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{}

	for _, route := range BuildingActionEndpoints(&mockBuildingActionService{}, &mockActionService{}, m) {
		ctx, _ := generateTestEchoContextWithMethodAndId(http.MethodGet)

		m.updatePlanetUntilCalled = 0

		handler := route.Handler()
		err := handler(ctx)

		assert.Nil(err)
		assert.Equal(1, m.updatePlanetUntilCalled)
		assert.Equal(defaultUuid, m.planet)
	}
}

func TestBuildingActionEndpoints_WhenPlanetResourceServiceFails_SetsReturnStatusInternalError(t *testing.T) {
	assert := assert.New(t)

	m := &mockPlanetResourceService{
		err: errDefault,
	}

	for _, route := range BuildingActionEndpoints(&mockBuildingActionService{}, &mockActionService{}, m) {
		ctx, rw := generateTestEchoContextWithMethodAndId(http.MethodGet)

		m.updatePlanetUntilCalled = 0

		handler := route.Handler()
		err := handler(ctx)

		assert.Nil(err)
		assert.Equal(1, m.updatePlanetUntilCalled)
		assert.Equal(http.StatusInternalServerError, rw.Code)
		assert.Equal("\"Failed to update resources\"\n", rw.Body.String())
	}
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
