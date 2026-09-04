package drivingadapters

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/drivingportstest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
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

		dto := dtos.ShipActionDtoRequest{Ship: uuid.New(), Count: 26}
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

		dto := dtos.ShipActionDtoRequest{Ship: uuid.New(), Count: 0}
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

		dto := dtos.ShipActionDtoRequest{Ship: uuid.New(), Count: -1}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid ship syntax", actual)
	})

	t.Run("forwards creation to use case", func(t *testing.T) {
		dto := dtos.ShipActionDtoRequest{Ship: uuid.New(), Count: 4}

		expectedRequest := request.ShipActionCreationRequest{
			Planet: sampleUuid,
			Ship:   dto.Ship,
			Count:  dto.Count,
		}
		action := models.ShipAction{
			Id:                 uuid.New(),
			Ship:               dto.Ship,
			Count:              dto.Count,
			CreatedAt:          someTime,
			NextCompletionAt:   someOtherTime,
			UnitCompletionTime: someOtherTime.Sub(someTime),
			Costs: []models.ShipActionCost{
				{
					Resource: uuid.New(),
					Amount:   1478,
				},
			},
		}

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Eq(expectedRequest)).
			Times(1).
			Return(action, nil)

		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.ShipActionDtoResponse](t, rw)
		expected := dtos.ShipActionDtoResponse{
			Id:               action.Id,
			Ship:             action.Ship,
			Count:            action.Count,
			CreatedAt:        action.CreatedAt,
			NextCompletionAt: action.NextCompletionAt,
			// Corresponds to someOtherTime - someTime
			UnitCompletionTime: "PT20H20M26S",
			Costs: []dtos.ShipActionCostDtoResponse{
				{Resource: action.Costs[0].Resource, Amount: 1478},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 404 when planet is not found", func(t *testing.T) {
		dto := dtos.ShipActionDtoRequest{Ship: uuid.New(), Count: 1}

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.ShipAction{}, domainerrors.ErrNotFound)

		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such planet", actual)
	})

	t.Run("returns 400 when ship is not found", func(t *testing.T) {
		dto := dtos.ShipActionDtoRequest{Ship: uuid.New(), Count: 1}

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.ShipAction{}, domainerrors.ErrShipNotFound)

		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such ship", actual)
	})

	t.Run("returns 400 when not enough resources are on the planet", func(t *testing.T) {
		dto := dtos.ShipActionDtoRequest{Ship: uuid.New(), Count: 1}

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.ShipAction{}, domainerrors.ErrNotEnoughResources)

		handler := generateHandler[drivingports.ForCreatingShipAction](
			createShipAction,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/ships", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "not enough resources", actual)
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

		dto := dtos.ShipActionDtoRequest{Ship: uuid.New(), Count: 2}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/ships", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to create ship", actual)
	})
}
