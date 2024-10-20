package controller

import (
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/game"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/rest"
	"github.com/stretchr/testify/suite"
)

type routesGenerator func() rest.Routes

type routeErrorTestCase struct {
	generateRoutes     routesGenerator
	expectedStatusCode int
	expectedError      string
}

type routesWithServiceGenerator func(game.ActionService, game.PlanetResourceService) rest.Routes

type routeInteractionTestCase struct {
	generateRoutes routesWithServiceGenerator
}

type RouteTestSuite struct {
	suite.Suite

	generateRoutes routesGenerator

	expectedRoutes map[string]int

	errorTestCases       map[string]routeErrorTestCase
	interactionTestCases []routeInteractionTestCase
}

func (s *RouteTestSuite) Test_GeneratesExpectedRoutes() {
	routes := s.generateRoutes()
	actualRoutes := make(map[string]int)

	for _, r := range routes {
		actualRoutes[r.Method()]++
	}

	s.Require().Equal(len(s.expectedRoutes), len(actualRoutes))
	for method, count := range actualRoutes {
		expectedCount, ok := s.expectedRoutes[method]

		s.Require().True(ok)
		s.Require().Equal(expectedCount, count)
	}
}

func (s *RouteTestSuite) Test_WhenRouteFails_ExpectCorrectStatus() {
	for name, testCase := range s.errorTestCases {
		var routes rest.Routes
		if testCase.generateRoutes != nil {
			routes = testCase.generateRoutes()
		} else {
			routes = s.generateRoutes()
		}

		for _, route := range routes {
			s.T().Run(name+"_"+route.Method(), func(t *testing.T) {
				ctx, rw := generateTestEchoContextWithMethodAndId(http.MethodGet)

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(testCase.expectedStatusCode, rw.Code)
				s.Require().Equal(testCase.expectedError, rw.Body.String())
			})
		}
	}
}

func (s *RouteTestSuite) Test_WhenActionServiceFails_ExpectCorrectInteraction() {
	for _, testCase := range s.interactionTestCases {
		m := &mockActionService{
			err: errDefault,
		}

		s.Require().NotNil(testCase.generateRoutes)
		routes := testCase.generateRoutes(m, &mockPlanetResourceService{})

		for _, route := range routes {
			s.T().Run("whenActionServiceFails_"+route.Method(), func(t *testing.T) {
				m.processActionsCalled = 0
				ctx, _ := generateTestEchoContextWithMethodAndId(route.Method())

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(1, m.processActionsCalled)
			})
		}
	}
}

func (s *RouteTestSuite) Test_WhenNoPlanetId_ExpectNoSchedulingOfActions() {
	for _, testCase := range s.interactionTestCases {
		m := &mockActionService{}

		s.Require().NotNil(testCase.generateRoutes)
		routes := testCase.generateRoutes(m, &mockPlanetResourceService{})

		for _, route := range routes {
			s.T().Run("whenNoPlanetId_"+route.Method(), func(t *testing.T) {
				m.processActionsCalled = 0
				ctx, _ := generateTestEchoContextWithMethod(route.Method())

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(0, m.processActionsCalled)
			})
		}
	}
}

func (s *RouteTestSuite) Test_WhenPlanetIdIsInvalid_ExpectNoSchedulingOfActions() {
	for _, testCase := range s.interactionTestCases {
		m := &mockActionService{}

		s.Require().NotNil(testCase.generateRoutes)
		routes := testCase.generateRoutes(m, &mockPlanetResourceService{})

		for _, route := range routes {
			s.T().Run("whenPlanetIdIsInvalid_"+route.Method(), func(t *testing.T) {
				m.processActionsCalled = 0
				ctx, _ := generateTestEchoContextWithMethod(route.Method())
				ctx.SetParamNames("id")
				ctx.SetParamValues("not-a-uuid")

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(0, m.processActionsCalled)
			})
		}
	}
}

func (s *RouteTestSuite) Test_WhenPlanetIdIsValid_ExpectSchedulingOfActions() {
	for _, testCase := range s.interactionTestCases {
		m := &mockActionService{}

		s.Require().NotNil(testCase.generateRoutes)
		routes := testCase.generateRoutes(m, &mockPlanetResourceService{})

		for _, route := range routes {
			s.T().Run("whenPlanetIdIsValid_"+route.Method(), func(t *testing.T) {
				m.processActionsCalled = 0
				ctx, _ := generateTestEchoContextWithMethodAndId(route.Method())

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(1, m.processActionsCalled)
				s.Require().Equal(defaultUuid, m.planet)
			})
		}
	}
}

func (s *RouteTestSuite) Test_WhenPlanetResourceServiceFails_ExpectCorrectInteraction() {
	for _, testCase := range s.interactionTestCases {
		m := &mockPlanetResourceService{
			err: errDefault,
		}

		s.Require().NotNil(testCase.generateRoutes)
		routes := testCase.generateRoutes(&mockActionService{}, m)

		for _, route := range routes {
			s.T().Run("whenPlanetResourceServiceFails_"+route.Method(), func(t *testing.T) {
				m.updatePlanetUntilCalled = 0
				ctx, _ := generateTestEchoContextWithMethodAndId(route.Method())

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(1, m.updatePlanetUntilCalled)
			})
		}
	}
}

func (s *RouteTestSuite) Test_WhenNoPlanetId_ExpectNoUpdateOfResources() {
	for _, testCase := range s.interactionTestCases {
		m := &mockPlanetResourceService{}

		s.Require().NotNil(testCase.generateRoutes)
		routes := testCase.generateRoutes(&mockActionService{}, m)

		for _, route := range routes {
			s.T().Run("whenNoPlanetId_"+route.Method(), func(t *testing.T) {
				m.updatePlanetUntilCalled = 0
				ctx, _ := generateTestEchoContextWithMethod(route.Method())

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(0, m.updatePlanetUntilCalled)
			})
		}
	}
}

func (s *RouteTestSuite) Test_WhenPlanetIdIsInvalid_ExpectNoUpdateOfResources() {
	for _, testCase := range s.interactionTestCases {
		m := &mockPlanetResourceService{}

		s.Require().NotNil(testCase.generateRoutes)
		routes := testCase.generateRoutes(&mockActionService{}, m)

		for _, route := range routes {
			s.T().Run("whenPlanetIdIsInvalid_"+route.Method(), func(t *testing.T) {
				m.updatePlanetUntilCalled = 0
				ctx, _ := generateTestEchoContextWithMethod(route.Method())
				ctx.SetParamNames("id")
				ctx.SetParamValues("not-a-uuid")

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(0, m.updatePlanetUntilCalled)
			})
		}
	}
}

func (s *RouteTestSuite) Test_WhenPlanetIdIsValid_ExpectUpdateOfResources() {
	for _, testCase := range s.interactionTestCases {
		m := &mockPlanetResourceService{}

		s.Require().NotNil(testCase.generateRoutes)
		routes := testCase.generateRoutes(&mockActionService{}, m)

		for _, route := range routes {
			s.T().Run("whenPlanetIdIsValid_"+route.Method(), func(t *testing.T) {
				m.updatePlanetUntilCalled = 0
				ctx, _ := generateTestEchoContextWithMethodAndId(route.Method())

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(1, m.updatePlanetUntilCalled)
				s.Require().Equal(defaultUuid, m.planet)
			})
		}
	}
}
