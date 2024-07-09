package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
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
			Id:          uuid.MustParse("d74072e4-c6d1-4486-9428-61f025bbf372"),
			User:        uuid.MustParse("9676fdad-6d3d-4df2-85a6-382ab4aad9dc"),
			Resource:    "resource-1",
			Permissions: []string{"POST", "GET"},
			CreatedAt:   time.Date(2024, 06, 30, 17, 44, 24, 651387237, time.UTC),
		},
	},
	Limits: []communication.LimitDtoResponse{
		{
			Name:  "limit-1",
			Value: "10",
		},
		{
			Name:  "limit-2",
			Value: "some-string",
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
		"authUser_notLoggedIn": {
			req:                generateTestRequestWithApiKey(http.MethodGet),
			handler:            authUser,
			err:                errors.NewCode(service.UserNotAuthenticated),
			expectedHttpStatus: http.StatusForbidden,
		},
		"authUser_keyExpired": {
			req:                generateTestRequestWithApiKey(http.MethodGet),
			handler:            authUser,
			err:                errors.NewCode(service.AuthenticationExpired),
			expectedHttpStatus: http.StatusForbidden,
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
			expectedHttpStatus: http.StatusNoContent,
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

func Test_WhenAuthServiceSucceeds_SetsResponseHeaderCorrectly(t *testing.T) {
	assert := assert.New(t)

	type testCaseReturn struct {
		req            *http.Request
		idAsRouteParam bool
		handler        authServiceAwareHttpHandler

		authData communication.AuthorizationDtoResponse

		expectedHeaders map[string]string
	}

	testCases := map[string]testCaseReturn{
		"authUser": {
			req:      generateTestRequestWithApiKey(http.MethodPost),
			handler:  authUser,
			authData: defaultAuthorizationResponseDto,
			expectedHeaders: map[string]string{
				"X-Acl":        `[{"id":"d74072e4-c6d1-4486-9428-61f025bbf372","user":"9676fdad-6d3d-4df2-85a6-382ab4aad9dc","resource":"resource-1","permissions":["POST","GET"],"createdAt":"2024-06-30T17:44:24.651387237Z"}]`,
				"X-User-Limit": `[{"name":"limit-1","value":"10"},{"name":"limit-2","value":"some-string"}]`,
			},
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

			assert.Equal(0, len(rw.Body.Bytes()))
			assert.Equal(len(testCase.expectedHeaders), len(rw.Header()))

			for headerKey, expectedHeader := range testCase.expectedHeaders {

				actualHeader := rw.Header().Get(headerKey)
				assert.Equal(expectedHeader, actualHeader)
			}
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
