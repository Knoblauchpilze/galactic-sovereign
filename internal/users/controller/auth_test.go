package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockAuthService struct {
	service.AuthService

	authData communication.AuthorizationDtoResponse
	err      error

	authCalled int

	inApiKey uuid.UUID
}

var defaultApiKeyId = uuid.MustParse("4c2a950c-ce65-4fb4-87b3-ce588dcfc1ea")
var defaultAuthorizationResponseDto = communication.AuthorizationDtoResponse{
	Acls: []communication.AclDtoResponse{
		{
			Resource:    "resource-1",
			Permissions: []string{"POST", "GET"},
		},
	},
	Limits: []communication.LimitDtoResponse{
		{
			Name:  "limit-1",
			Value: "10",
		},
	},
}

type authTestCase struct {
	req     *http.Request
	handler authServiceAwareHttpHandler
}

func TestAuthEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range AuthEndpoints(&mockAuthService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(1, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodGet])
}

func Test_WhenNoApiKey_SetsStatusTo400(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]authTestCase{
		"authUser": {
			req:     httptest.NewRequest(http.MethodGet, "/", nil),
			handler: authUser,
		},
		"authUser_multipleKeys": {
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)

				req.Header.Add("X-Api-Key", "e6349328-543b-4b4e-8a3c-4caf7b413589")
				req.Header.Add("X-Api-Key", "de2108c2-f87b-4033-825c-4ccbbb8b778e")

				return req
			}(),
			handler: authUser,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockAuthService{}
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(http.StatusBadRequest, rw.Code)
			assert.Equal("\"Api key not found\"\n", rw.Body.String())
		})
	}
}

func Test_WhenApiKeySyntaxIsWrong_SetsStatusTo400(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]authTestCase{
		"authUser": {
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)

				req.Header.Add("X-Api-Key", "not-a-uuid")

				return req
			}(),
			handler: authUser,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockAuthService{}
			ctx, rw := generateTestEchoContextFromRequest(testCase.req)

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(http.StatusBadRequest, rw.Code)
			assert.Equal("\"Invalid api key syntax\"\n", rw.Body.String())
		})
	}
}

func Test_WhenAuthServiceFails_SetsExpectedStatus(t *testing.T) {
	assert := assert.New(t)

	type testCaseError struct {
		req                *http.Request
		handler            authServiceAwareHttpHandler
		err                error
		expectedHttpStatus int
	}

	testCases := map[string]testCaseError{
		"authUser": {
			req:                generateTestRequestWithApiKey(http.MethodGet),
			handler:            authUser,
			err:                errDefault,
			expectedHttpStatus: http.StatusInternalServerError,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockAuthService{
				err: testCase.err,
			}

			ctx, rw := generateTestEchoContextFromRequest(testCase.req)

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(testCase.expectedHttpStatus, rw.Code)
		})
	}
}

func Test_WhenAuthServiceSucceeds_SetsExpectedStatus(t *testing.T) {
	assert := assert.New(t)

	type testCaseSuccess struct {
		req                *http.Request
		idAsRouteParam     bool
		handler            authServiceAwareHttpHandler
		expectedHttpStatus int
	}

	testCases := map[string]testCaseSuccess{
		"authUser": {
			req:                generateTestRequestWithApiKey(http.MethodGet),
			handler:            authUser,
			expectedHttpStatus: http.StatusOK,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockAuthService{}

			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			err := testCase.handler(ctx, mock)

			assert.Nil(err)
			assert.Equal(testCase.expectedHttpStatus, rw.Code)
		})
	}
}

func Test_WhenAuthServiceSucceeds_ReturnsExpectedValue(t *testing.T) {
	assert := assert.New(t)

	type testCaseReturn struct {
		req            *http.Request
		idAsRouteParam bool
		handler        authServiceAwareHttpHandler

		authData communication.AuthorizationDtoResponse

		expectedContent interface{}
	}

	testCases := map[string]testCaseReturn{
		"authUser": {
			req:             generateTestRequestWithApiKey(http.MethodPost),
			handler:         authUser,
			authData:        defaultAuthorizationResponseDto,
			expectedContent: defaultAuthorizationResponseDto,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			mock := &mockAuthService{
				authData: testCase.authData,
			}

			ctx, rw := generateTestEchoContextFromRequest(testCase.req)
			if testCase.idAsRouteParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues(defaultUuid.String())
			}

			err := testCase.handler(ctx, mock)

			assert.Nil(err)

			actual := strings.Trim(rw.Body.String(), "\n")
			expected, err := json.Marshal(testCase.expectedContent)
			assert.Nil(err)
			assert.Equal(string(expected), actual)
		})
	}
}

func TestAuthUser_CallsServiceCreate(t *testing.T) {
	assert := assert.New(t)

	req := generateTestRequestWithApiKey(http.MethodGet)
	ctx, _ := generateTestEchoContextFromRequest(req)
	ms := &mockAuthService{}

	err := authUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(1, ms.authCalled)
}

func TestAuthUser_GetsExpectedApiKey(t *testing.T) {
	assert := assert.New(t)

	req := generateTestRequestWithApiKey(http.MethodGet)
	ctx, _ := generateTestEchoContextFromRequest(req)
	ms := &mockAuthService{}

	err := authUser(ctx, ms)

	assert.Nil(err)
	assert.Equal(defaultApiKeyId, ms.inApiKey)
}

func generateTestRequestWithApiKey(method string) *http.Request {
	req := httptest.NewRequest(method, "/", nil)
	req.Header.Set("X-Api-Key", defaultApiKeyId.String())

	return req
}

func (m *mockAuthService) Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationDtoResponse, error) {
	m.authCalled++
	m.inApiKey = apiKey
	return m.authData, m.err
}
