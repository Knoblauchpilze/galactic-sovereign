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
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUnit_Fleets_CreateFleet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForCreatingFleet(ctrl)

	t.Run("returns 400 when planet id is invalid", func(t *testing.T) {
		dto := dtos.FleetDtoRequest{Mission: dtos.MissionColonize}
		handler := generateHandler[drivingports.ForCreatingFleet](createFleet, mockUsecase)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", "not-a-uuid")
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("returns 400 when body is invalid", func(t *testing.T) {
		handler := generateHandler[drivingports.ForCreatingFleet](createFleet, mockUsecase)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, "not-a-dto-request")
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid fleet syntax", actual)
	})

	t.Run("returns 400 when ship count is zero", func(t *testing.T) {
		handler := generateHandler[drivingports.ForCreatingFleet](createFleet, mockUsecase)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		dto := dtos.FleetDtoRequest{
			Mission: dtos.MissionColonize,
			Ships: []dtos.FleetShipDtoRequest{
				{Ship: uuid.New(), Count: 0},
			},
		}

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid fleet syntax", actual)
	})

	t.Run("returns 400 when ship count is negative", func(t *testing.T) {
		handler := generateHandler[drivingports.ForCreatingFleet](createFleet, mockUsecase)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		dto := dtos.FleetDtoRequest{
			Mission: dtos.MissionColonize,
			Ships: []dtos.FleetShipDtoRequest{
				{Ship: uuid.New(), Count: -1},
			},
		}

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid fleet syntax", actual)
	})

	t.Run("forwards creation to use case", func(t *testing.T) {
		dto := generateSampleFleetDtoRequest()

		expectedRequest := models.FleetOrder{
			Mission: models.MissionColonize,
			Destination: models.Coordinate{
				Galaxy:      *dto.Destination.Galaxy,
				SolarSystem: *dto.Destination.SolarSystem,
				Position:    *dto.Destination.Position,
			},
			Ships: []models.FleetShip{
				{
					Ship:  dto.Ships[0].Ship,
					Count: dto.Ships[0].Count,
				},
				{
					Ship:  dto.Ships[1].Ship,
					Count: dto.Ships[1].Count,
				},
			},
		}

		fleet := models.Fleet{
			Id:        uuid.New(),
			CreatedAt: someTime,
			Ships: []models.FleetShip{
				{
					Ship:  dto.Ships[0].Ship,
					Count: dto.Ships[0].Count,
				},
				{
					Ship:  dto.Ships[1].Ship,
					Count: dto.Ships[1].Count,
				},
			},
		}

		mockUsecase.EXPECT().
			Create(gomock.Any(), sampleUuid, expectedRequest).
			Times(1).
			Return(fleet, nil)

		handler := generateHandler[drivingports.ForCreatingFleet](createFleet, mockUsecase)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.FleetDtoResponse](t, rw)
		expected := dtos.FleetDtoResponse{
			Id:        fleet.Id,
			CreatedAt: fleet.CreatedAt,
			Ships: []dtos.FleetShipDtoResponse{
				{Ship: fleet.Ships[0].Ship, Count: fleet.Ships[0].Count},
				{Ship: fleet.Ships[1].Ship, Count: fleet.Ships[1].Count},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("allows creation for 0 coordinates", func(t *testing.T) {
		dto := dtos.FleetDtoRequest{
			Mission: dtos.MissionColonize,
			Destination: dtos.FleetDestinationDtoRequest{
				Galaxy:      new(0),
				SolarSystem: new(0),
				Position:    new(0),
			},
			Ships: []dtos.FleetShipDtoRequest{
				{Ship: uuid.New(), Count: 2},
			},
		}

		expectedRequest := models.FleetOrder{
			Mission: models.MissionColonize,
			Destination: models.Coordinate{
				Galaxy:      0,
				SolarSystem: 0,
				Position:    0,
			},
			Ships: []models.FleetShip{
				{Ship: dto.Ships[0].Ship, Count: dto.Ships[0].Count},
			},
		}

		fleet := models.Fleet{
			Id:        uuid.New(),
			CreatedAt: someTime,
			Ships: []models.FleetShip{
				{Ship: dto.Ships[0].Ship, Count: dto.Ships[0].Count},
			},
		}

		mockUsecase.EXPECT().
			Create(gomock.Any(), sampleUuid, expectedRequest).
			Times(1).
			Return(fleet, nil)

		handler := generateHandler[drivingports.ForCreatingFleet](createFleet, mockUsecase)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.FleetDtoResponse](t, rw)
		expected := dtos.FleetDtoResponse{
			Id:        fleet.Id,
			CreatedAt: fleet.CreatedAt,
			Ships: []dtos.FleetShipDtoResponse{
				{Ship: fleet.Ships[0].Ship, Count: fleet.Ships[0].Count},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 409 when not enough ships are on the planet", func(t *testing.T) {
		dto := generateSampleFleetDtoRequest()

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Fleet{}, domainerrors.ErrNotEnoughShips)

		handler := generateHandler[drivingports.ForCreatingFleet](
			createFleet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "not enough ships", actual)
	})

	t.Run("returns 409 when not enough ships are on the planet", func(t *testing.T) {
		dto := generateSampleFleetDtoRequest()

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Fleet{}, domainerrors.ErrFleetDestinationInvalid)

		handler := generateHandler[drivingports.ForCreatingFleet](
			createFleet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid fleet destination", actual)
	})

	t.Run("returns 400 when fleet has invalid configuration", func(t *testing.T) {
		dto := generateSampleFleetDtoRequest()

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Fleet{}, domainerrors.ErrInvalidFleetConfiguration)

		handler := generateHandler[drivingports.ForCreatingFleet](
			createFleet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid fleet configuration", actual)

	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		dto := generateSampleFleetDtoRequest()

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Fleet{}, errors.New("stubbed error"))

		handler := generateHandler[drivingports.ForCreatingFleet](createFleet, mockUsecase)
		r := createTestGinRouter(t, http.MethodPost, "/planets/:id/fleets", handler)

		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		addRequestPath(t, req, "/planets/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to create fleet", actual)
	})
}

func generateSampleFleetDtoRequest() dtos.FleetDtoRequest {
	return dtos.FleetDtoRequest{
		Mission: dtos.MissionColonize,
		Destination: dtos.FleetDestinationDtoRequest{
			Galaxy:      new(8),
			SolarSystem: new(4),
			Position:    new(10),
		},
		Ships: []dtos.FleetShipDtoRequest{
			{
				Ship:  uuid.New(),
				Count: 2,
			},
			{
				Ship:  uuid.New(),
				Count: 3,
			},
		},
	}
}
