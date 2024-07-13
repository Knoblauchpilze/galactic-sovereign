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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

func TestAuthEndpoints_GeneratesExpectedRoutes(t *testing.T) {
	assert := assert.New(t)

	actualRoutes := make(map[string]int)
	for _, r := range AuthEndpoints(&mockAuthService{}) {
		actualRoutes[r.Method()]++
	}

	assert.Equal(1, len(actualRoutes))
	assert.Equal(1, actualRoutes[http.MethodGet])
}

func Test_AuthController(t *testing.T) {
	s := ControllerTestSuite[service.AuthService]{
		generateServiceMock:      generateAuthServiceMock,
		generateValidServiceMock: generateValidAuthServiceMock,

		badInputTestCases: map[string]badInputTestCase[service.AuthService]{
			"authUser": {
				req:                httptest.NewRequest(http.MethodGet, "/", nil),
				handler:            authUser,
				expectedBodyString: "\"Api key not found\"\n",
			},
			"authUser_multipleKeys": {
				req: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/", nil)

					req.Header.Add("X-Api-Key", "e6349328-543b-4b4e-8a3c-4caf7b413589")
					req.Header.Add("X-Api-Key", "de2108c2-f87b-4033-825c-4ccbbb8b778e")

					return req
				}(),
				handler:            authUser,
				expectedBodyString: "\"Api key not found\"\n",
			},
		},

		errorTestCases: map[string]errorTestCase[service.AuthService]{
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
		},

		successTestCases: map[string]successTestCase[service.AuthService]{
			"authUser": {
				req:                generateTestRequestWithApiKey(http.MethodGet),
				handler:            authUser,
				expectedHttpStatus: http.StatusNoContent,
			},
		},

		responseTestCases: map[string]responseTestCase[service.AuthService]{
			"authUser": {
				req:     generateTestRequestWithApiKey(http.MethodPost),
				handler: authUser,

				verifyResponse: func(rw *httptest.ResponseRecorder, assert *require.Assertions) {
					assert.Equal(0, len(rw.Body.Bytes()))

					expectedHeaders := map[string]string{
						"X-Acl":        `[{"id":"d74072e4-c6d1-4486-9428-61f025bbf372","user":"9676fdad-6d3d-4df2-85a6-382ab4aad9dc","resource":"resource-1","permissions":["POST","GET"],"createdAt":"2024-06-30T17:44:24.651387237Z"}]`,
						"X-User-Limit": `[{"name":"limit-1","value":"10"},{"name":"limit-2","value":"some-string"}]`,
					}

					assert.Equal(len(expectedHeaders), len(rw.Header()))

					for headerKey, expectedHeader := range expectedHeaders {
						actualHeader := rw.Header().Get(headerKey)
						assert.Equal(expectedHeader, actualHeader)
					}
				},
			},
		},

		serviceInteractionTestCases: map[string]serviceInteractionTestCase[service.AuthService]{
			"authUser": {
				req:     generateTestRequestWithApiKey(http.MethodGet),
				handler: authUser,

				verifyInteractions: func(s service.AuthService, assert *require.Assertions) {
					m := assertAuthServiceIsAMock(s, assert)

					assert.Equal(1, m.authCalled)
					assert.Equal(defaultApiKeyId, m.inApiKey)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateAuthServiceMock(err error) service.AuthService {
	return &mockAuthService{
		err: err,
	}
}

func generateValidAuthServiceMock() service.AuthService {
	return &mockAuthService{
		authData: defaultAuthorizationResponseDto,
	}
}

func assertAuthServiceIsAMock(s service.AuthService, assert *require.Assertions) *mockAuthService {
	m, ok := s.(*mockAuthService)
	if !ok {
		assert.Fail("Provided auth service is not a mock")
	}
	return m
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
