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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUnit_Fleets_CreateFleet(t *testing.T) {
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
			Id:      uuid.New(),
			Player:  uuid.New(),
			Source:  sampleUuid,
			Mission: expectedRequest.Mission,
			Destination: models.FleetDestination{
				Coordinate: expectedRequest.Destination,
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
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: someTime,
				ReturnAt:  someOtherTime,
				Status:    models.StatusOutbound,
			},
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
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
			Id:      fleet.Id,
			Player:  fleet.Player,
			Source:  fleet.Source,
			Mission: fleet.Mission,
			Destination: dtos.FleetDestinationDtoResponse{
				Galaxy:      fleet.Destination.Coordinate.Galaxy,
				SolarSystem: fleet.Destination.Coordinate.SolarSystem,
				Position:    fleet.Destination.Coordinate.Position,
				Target:      nil,
			},
			Ships: []dtos.FleetShipDtoResponse{
				{Ship: fleet.Ships[0].Ship, Count: fleet.Ships[0].Count},
				{Ship: fleet.Ships[1].Ship, Count: fleet.Ships[1].Count},
			},
			Schedule: dtos.FleetScheduleDtoResponse{
				ArrivalAt: fleet.Schedule.ArrivalAt,
				ReturnAt:  fleet.Schedule.ReturnAt,
				Status:    fleet.Schedule.Status,
			},
			CreatedAt: fleet.CreatedAt,
			UpdatedAt: fleet.UpdatedAt,
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
			Id:      uuid.New(),
			Player:  uuid.New(),
			Source:  sampleUuid,
			Mission: expectedRequest.Mission,
			Destination: models.FleetDestination{
				Coordinate: expectedRequest.Destination,
			},
			Ships: []models.FleetShip{
				{
					Ship:  dto.Ships[0].Ship,
					Count: dto.Ships[0].Count,
				},
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: someTime,
				ReturnAt:  someOtherTime,
				Status:    models.StatusOutbound,
			},
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
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
			Id:      fleet.Id,
			Player:  fleet.Player,
			Source:  fleet.Source,
			Mission: fleet.Mission,
			Destination: dtos.FleetDestinationDtoResponse{
				Galaxy:      fleet.Destination.Coordinate.Galaxy,
				SolarSystem: fleet.Destination.Coordinate.SolarSystem,
				Position:    fleet.Destination.Coordinate.Position,
				Target:      nil,
			},
			Ships: []dtos.FleetShipDtoResponse{
				{Ship: fleet.Ships[0].Ship, Count: fleet.Ships[0].Count},
			},
			Schedule: dtos.FleetScheduleDtoResponse{
				ArrivalAt: fleet.Schedule.ArrivalAt,
				ReturnAt:  fleet.Schedule.ReturnAt,
				Status:    fleet.Schedule.Status,
			},
			CreatedAt: fleet.CreatedAt,
			UpdatedAt: fleet.UpdatedAt,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns target when available", func(t *testing.T) {
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

		destination := uuid.New()
		fleet := models.Fleet{
			Id:      uuid.New(),
			Player:  uuid.New(),
			Source:  sampleUuid,
			Mission: expectedRequest.Mission,
			Destination: models.FleetDestination{
				Coordinate: expectedRequest.Destination,
				Target:     &destination,
			},
			Ships: []models.FleetShip{
				{
					Ship:  dto.Ships[0].Ship,
					Count: dto.Ships[0].Count,
				},
			},
			Schedule: models.FleetFlightSchedule{
				ArrivalAt: someTime,
				ReturnAt:  someOtherTime,
				Status:    models.StatusOutbound,
			},
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
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
			Id:      fleet.Id,
			Player:  fleet.Player,
			Source:  fleet.Source,
			Mission: fleet.Mission,
			Destination: dtos.FleetDestinationDtoResponse{
				Galaxy:      fleet.Destination.Coordinate.Galaxy,
				SolarSystem: fleet.Destination.Coordinate.SolarSystem,
				Position:    fleet.Destination.Coordinate.Position,
				Target:      &destination,
			},
			Ships: []dtos.FleetShipDtoResponse{
				{Ship: fleet.Ships[0].Ship, Count: fleet.Ships[0].Count},
			},
			Schedule: dtos.FleetScheduleDtoResponse{
				ArrivalAt: fleet.Schedule.ArrivalAt,
				ReturnAt:  fleet.Schedule.ReturnAt,
				Status:    fleet.Schedule.Status,
			},
			CreatedAt: fleet.CreatedAt,
			UpdatedAt: fleet.UpdatedAt,
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

	t.Run("returns 400 when fleet has invalid coordinates", func(t *testing.T) {
		dto := generateSampleFleetDtoRequest()

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Fleet{}, domainerrors.ErrCoordinatesOutOfBound)

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
		assert.Equal(t, "invalid fleet coordinates", actual)

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

func TestUnit_Fleets_ListFleetsForPlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForFetchingFleets(ctrl)

	t.Run("returns 400 when player id is invalid", func(t *testing.T) {
		handler := generateHandler[drivingports.ForFetchingFleets](
			listFleetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/fleets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/not-a-uuid/fleets")
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards listing to use case", func(t *testing.T) {
		fleets := []models.Fleet{
			{Id: uuid.New(), Source: uuid.New(), CreatedAt: someTime},
			{Id: uuid.New(), Source: uuid.New(), CreatedAt: someOtherTime},
		}
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), sampleUuid).
			Times(1).
			Return(fleets, nil)

		handler := generateHandler[drivingports.ForFetchingFleets](
			listFleetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/fleets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.FleetDtoResponse](t, rw)
		expected := []dtos.FleetDtoResponse{
			{
				Id:        fleets[0].Id,
				Source:    fleets[0].Source,
				CreatedAt: fleets[0].CreatedAt,
				Ships:     []dtos.FleetShipDtoResponse{},
			},
			{
				Id:        fleets[1].Id,
				Source:    fleets[1].Source,
				CreatedAt: fleets[1].CreatedAt,
				Ships:     []dtos.FleetShipDtoResponse{},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("return empty slice when use case returns no planet", func(t *testing.T) {
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), sampleUuid).
			Times(1).
			Return([]models.Fleet{}, nil)

		handler := generateHandler[drivingports.ForFetchingFleets](
			listFleetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/fleets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.FleetDtoResponse](t, rw)
		assert.Equal(t, []dtos.FleetDtoResponse{}, actual)
	})

	t.Run("return empty slice when use case returns nil response", func(t *testing.T) {
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), sampleUuid).
			Times(1).
			Return(nil, nil)

		handler := generateHandler[drivingports.ForFetchingFleets](
			listFleetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/fleets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.FleetDtoResponse](t, rw)
		assert.Equal(t, []dtos.FleetDtoResponse{}, actual)
	})

	t.Run("returns 404 when universe does not exist", func(t *testing.T) {
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), sampleUuid).
			Times(1).
			Return([]models.Fleet{}, domainerrors.ErrPlayerNotFound)

		handler := generateHandler[drivingports.ForFetchingFleets](
			listFleetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/fleets", handler)

		req := generateTestRequest(t, http.MethodGet, addSampleUuidPathParam)
		addRequestPath(t, req, "/players/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such player", actual)
	})

	t.Run("returns 500 when use cas fails", func(t *testing.T) {
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), sampleUuid).
			Times(1).
			Return([]models.Fleet{}, errors.New("stubbed error"))

		handler := generateHandler[drivingports.ForFetchingFleets](
			listFleetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/fleets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/%s/fleets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to list fleets", actual)
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
