package drivingadapters

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/drivingportstest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUnit_Ships_CreateShipAction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForCreatingShipAction(ctrl)

	t.Run("returns 400 when planet id is invalid", func(t *testing.T) {
		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		dto := dtos.ShipDtoRequest{Ship: uuid.New(), Count: 26}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", "not-a-uuid")
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("returns 400 when body is invalid", func(t *testing.T) {
		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, "not-a-dto-request")
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid ship syntax", actual)
	})

	t.Run("returns 400 when ship count is zero", func(t *testing.T) {
		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		dto := dtos.ShipDtoRequest{Ship: uuid.New(), Count: 0}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid ship syntax", actual)
	})

	t.Run("returns 400 when ship count is negative", func(t *testing.T) {
		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		dto := dtos.ShipDtoRequest{Ship: uuid.New(), Count: -1}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid ship syntax", actual)
	})

	t.Run("forwards creation to use case", func(t *testing.T) {
		dto := dtos.ShipDtoRequest{Ship: uuid.New(), Count: 4}

		expectedRequest := request.ShipActionCreationRequest{
			Planet: sampleUuid,
			Ship:   dto.Ship,
			Count:  dto.Count,
		}

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Eq(expectedRequest)).
			Times(1).
			Return(models.ShipAction{}, nil)

		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.ShipAction{}, errors.New("stubbed error"))

		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		dto := dtos.ShipDtoRequest{Ship: uuid.New(), Count: 2}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to create ship", actual)
	})
}
