package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mockResourceService struct {
	resources []communication.ResourceDtoResponse
	err       error

	listCalled int
}

var defaultResourceDtoResponse = communication.ResourceDtoResponse{
	Id:   defaultUuid,
	Name: "my-resource",

	CreatedAt: time.Date(2024, 8, 3, 14, 29, 31, 651387240, time.UTC),
}

func TestResourceEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range ResourceEndpoints(&mockResourceService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(1, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodGet])
}

func Test_ResourceController(t *testing.T) {
	s := ControllerTestSuite[service.ResourceService]{
		generateServiceMock:      generateResourceServiceMock,
		generateValidServiceMock: generateValidResourceServiceMock,

		errorTestCases: map[string]errorTestCase[service.ResourceService]{
			"listResources": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				handler:            listResources,
				err:                errDefault,
				expectedHttpStatus: http.StatusInternalServerError,
			},
		},

		successTestCases: map[string]successTestCase[service.ResourceService]{
			"listResources": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				handler:            listResources,
				expectedHttpStatus: http.StatusOK,
			},
		},

		returnTestCases: map[string]returnTestCase[service.ResourceService]{
			"listResources": {
				req:             httptest.NewRequest(http.MethodGet, "/", nil),
				handler:         listResources,
				expectedContent: []communication.ResourceDtoResponse{defaultResourceDtoResponse},
			},
			"listResources_noData": {
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				generateValidServiceMock: func() service.ResourceService {
					return &mockResourceService{
						resources: nil,
					}
				},

				handler:         listResources,
				expectedContent: []communication.ResourceDtoResponse{},
			},
		},

		serviceInteractionTestCases: map[string]serviceInteractionTestCase[service.ResourceService]{
			"listResources": {
				req:     httptest.NewRequest(http.MethodGet, "/", nil),
				handler: listResources,

				verifyInteractions: func(s service.ResourceService, assert *require.Assertions) {
					m := assertResourceServiceIsAMock(s, assert)

					assert.Equal(1, m.listCalled)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateResourceServiceMock(err error) service.ResourceService {
	return &mockResourceService{
		err: err,
	}
}

func generateValidResourceServiceMock() service.ResourceService {
	return &mockResourceService{
		resources: []communication.ResourceDtoResponse{defaultResourceDtoResponse},
	}
}

func assertResourceServiceIsAMock(s service.ResourceService, assert *require.Assertions) *mockResourceService {
	m, ok := s.(*mockResourceService)
	if !ok {
		assert.Fail("Provided resource service is not a mock")
	}
	return m
}

func (m *mockResourceService) List(ctx context.Context) ([]communication.ResourceDtoResponse, error) {
	m.listCalled++
	return m.resources, m.err
}
