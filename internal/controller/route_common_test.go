package controller

import (
	"net/http"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/stretchr/testify/suite"
)

type routesGenerator func() rest.Routes

type routeErrorTestCase struct {
	generateRoutes     routesGenerator
	expectedStatusCode int
	expectedError      string
}

type RouteTestSuite struct {
	suite.Suite

	generateRoutes routesGenerator

	expectedRoutes map[string]int

	errorTestCases map[string]routeErrorTestCase
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
		s.T().Run(name, func(t *testing.T) {
			var routes rest.Routes
			if testCase.generateRoutes != nil {
				routes = testCase.generateRoutes()
			} else {
				routes = s.generateRoutes()
			}

			for _, route := range routes {
				ctx, rw := generateTestEchoContextWithMethodAndId(http.MethodGet)

				handler := route.Handler()
				err := handler(ctx)

				s.Require().Nil(err)
				s.Require().Equal(testCase.expectedStatusCode, rw.Code)
				s.Require().Equal(testCase.expectedError, rw.Body.String())
			}
		})
	}
}
