package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type handlerFunc[Service any] func(echo.Context, Service) error
type generateServiceMock[Service any] func(err error) Service
type generateValidServiceMock[Service any] func() Service

type badInputTestCase[Service any] struct {
	req                *http.Request
	idAsRouteParam     bool
	handler            handlerFunc[Service]
	expectedBodyString string
}

type noIdTestCase[Service any] struct {
	req     *http.Request
	handler handlerFunc[Service]
}

type badIdTestCase[Service any] noIdTestCase[Service]

type errorTestCase[Service any] struct {
	req                *http.Request
	idAsRouteParam     bool
	handler            handlerFunc[Service]
	err                error
	expectedHttpStatus int
}

type successTestCase[Service any] struct {
	req                *http.Request
	idAsRouteParam     bool
	handler            handlerFunc[Service]
	expectedHttpStatus int
}

type returnTestCase[Service any] struct {
	req            *http.Request
	idAsRouteParam bool
	handler        handlerFunc[Service]

	expectedContent interface{}
}

type verifyResponse func(*httptest.ResponseRecorder, *require.Assertions)

type responseTestCase[Service any] struct {
	req            *http.Request
	idAsRouteParam bool
	handler        handlerFunc[Service]

	verifyResponse verifyResponse
}

type verifyMockInteractions[Service any] func(Service, *require.Assertions)

type serviceInteractionTestCase[Service any] struct {
	req            *http.Request
	idAsRouteParam bool
	handler        handlerFunc[Service]

	verifyInteractions verifyMockInteractions[Service]
}

type ControllerTestSuite[Service any] struct {
	suite.Suite

	generateServiceMock      generateServiceMock[Service]
	generateValidServiceMock generateValidServiceMock[Service]

	badInputTestCases map[string]badInputTestCase[Service]
	noIdTestCases     map[string]noIdTestCase[Service]
	badIdTestCases    map[string]badIdTestCase[Service]
	errorTestCases    map[string]errorTestCase[Service]
	successTestCases  map[string]successTestCase[Service]

	returnTestCases             map[string]returnTestCase[Service]
	responseTestCases           map[string]responseTestCase[Service]
	serviceInteractionTestCases map[string]serviceInteractionTestCase[Service]
}

func (s *ControllerTestSuite[Service]) TestWhenBadInputProvided_Expect400Status() {
	for name, testCase := range s.badInputTestCases {
		s.T().Run(name, func(t *testing.T) {
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			m := s.generateServiceMock(nil)
			err := testCase.handler(ctx, m)

			s.Require().Nil(err)
			s.Require().Equal(http.StatusBadRequest, rw.Code)
			s.Require().Equal(testCase.expectedBodyString, rw.Body.String())
		})
	}
}

func (s *ControllerTestSuite[Service]) TestWhenNoIdProvided_Expect400Status() {
	for name, testCase := range s.noIdTestCases {
		s.T().Run(name, func(t *testing.T) {
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)

			m := s.generateServiceMock(nil)
			err := testCase.handler(ctx, m)

			s.Require().Nil(err)
			s.Require().Equal(http.StatusBadRequest, rw.Code)
			s.Require().Equal("\"Invalid id syntax\"\n", rw.Body.String())
		})
	}
}

func (s *ControllerTestSuite[Service]) TestIdSyntaxIsWrong_Expect400Status() {
	for name, testCase := range s.badIdTestCases {
		s.T().Run(name, func(t *testing.T) {
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			ctx.SetParamNames("id")
			ctx.SetParamValues("not-a-uuid")

			m := s.generateServiceMock(nil)
			err := testCase.handler(ctx, m)

			s.Require().Nil(err)
			s.Require().Equal(http.StatusBadRequest, rw.Code)
			s.Require().Equal("\"Invalid id syntax\"\n", rw.Body.String())
		})
	}
}

func (s *ControllerTestSuite[Service]) TestWhenServiceFails_ExpectCorrectStatus() {
	for name, testCase := range s.errorTestCases {
		s.T().Run(name, func(t *testing.T) {
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			m := s.generateServiceMock(testCase.err)
			err := testCase.handler(ctx, m)

			s.Require().Nil(err)
			s.Require().Equal(testCase.expectedHttpStatus, rw.Code)
		})
	}
}

func (s *ControllerTestSuite[Service]) TestWhenServiceSucceeds_ExpectCorrectStatus() {
	for name, testCase := range s.successTestCases {
		s.T().Run(name, func(t *testing.T) {
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			m := s.generateServiceMock(nil)
			err := testCase.handler(ctx, m)

			s.Require().Nil(err)
			s.Require().Equal(testCase.expectedHttpStatus, rw.Code)
		})
	}
}

func (s *ControllerTestSuite[Service]) TestWhenServiceSucceeds_ReturnsExpectedValue() {
	for name, testCase := range s.returnTestCases {
		s.T().Run(name, func(t *testing.T) {
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			m := s.generateValidServiceMock()
			err := testCase.handler(ctx, m)

			s.Require().Nil(err)

			actual := strings.Trim(rw.Body.String(), "\n")
			expected, err := json.Marshal(testCase.expectedContent)
			s.Require().Nil(err)
			s.Require().Equal(string(expected), actual)
		})
	}
}

func (s *ControllerTestSuite[Service]) TestWhenServiceSucceeds_ReturnsExpectedResponse() {
	for name, testCase := range s.responseTestCases {
		s.T().Run(name, func(t *testing.T) {
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			m := s.generateValidServiceMock()
			err := testCase.handler(ctx, m)

			s.Require().Nil(err)
			testCase.verifyResponse(rw, s.Require())
		})
	}
}

func (s *ControllerTestSuite[Service]) TestWhenServiceSucceeds_ExpectCorrectInteraction() {
	for name, testCase := range s.serviceInteractionTestCases {
		s.T().Run(name, func(t *testing.T) {
			ctx, _ := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			m := s.generateValidServiceMock()
			err := testCase.handler(ctx, m)

			s.Require().Nil(err)
			testCase.verifyInteractions(m, s.Require())
		})
	}
}
