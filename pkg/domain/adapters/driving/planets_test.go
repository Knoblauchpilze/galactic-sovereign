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

func TestUnit_Planets_GetPlanet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		handler := generateHandler[drivingports.ForManagingPlanet](
			getPlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/:id", handler)

		req := generateTestRequest(t, http.MethodGet, addInvalidUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards fetching to use case", func(t *testing.T) {
		planet := models.Planet{
			Id:   uuid.New(),
			Name: "planet-1",
			Coordinate: models.Coordinate{
				Galaxy:      36,
				SolarSystem: 151,
				Position:    12,
			},
			Fields:    147,
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
			Resources: []models.PlanetResource{
				{
					Resource: uuid.New(),
					Amount:   1478,
				},
			},
			Storages: []models.PlanetResourceStorage{
				{
					Resource: uuid.New(),
					Storage:  48790,
				},
			},
			Productions: []models.PlanetResourceProduction{
				{
					Resource:   uuid.New(),
					Production: 12,
				},
				{
					Resource:   uuid.New(),
					Building:   ptrFor(uuid.New()),
					Production: 8917,
				},
			},
			Buildings: []models.PlanetBuilding{
				{
					Building: uuid.New(),
					Level:    14,
				},
			},
			BuildingAction: &models.BuildingAction{
				Id: sampleUuid,
			},
		}
		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(planet, nil)

		handler := generateHandler[drivingports.ForManagingPlanet](
			getPlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/:id", handler)

		req := generateTestRequest(t, http.MethodGet, addSampleUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[dtos.PlanetDtoResponse](t, rw)
		expected := dtos.PlanetDtoResponse{
			Id:   planet.Id,
			Name: planet.Name,
			Coordinate: dtos.CoordinateDtoResponse{
				Galaxy:      36,
				SolarSystem: 151,
				Position:    12,
			},
			Fields:    147,
			CreatedAt: planet.CreatedAt,
			UpdatedAt: planet.UpdatedAt,
			Resources: []dtos.PlanetResourceDtoResponse{
				{Resource: planet.Resources[0].Resource, Amount: 1478},
			},
			Storages: []dtos.PlanetResourceStorageDtoResponse{
				{Resource: planet.Storages[0].Resource, Storage: 48790},
			},
			Productions: []dtos.PlanetResourceProductionDtoResponse{
				{
					Resource:   planet.Productions[0].Resource,
					Production: 12,
				},
				{
					Building:   planet.Productions[1].Building,
					Resource:   planet.Productions[1].Resource,
					Production: 8917,
				},
			},
			Buildings: []dtos.PlanetBuildingDtoResponse{
				{
					Building: planet.Buildings[0].Building,
					Level:    planet.Buildings[0].Level,
				},
			},
			BuildingAction: &dtos.BuildingActionDtoResponse{
				Id:          sampleUuid,
				Costs:       []dtos.BuildingActionCostDtoResponse{},
				Storages:    []dtos.BuildingActionStorageDtoResponse{},
				Productions: []dtos.BuildingActionProductionDtoResponse{},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("correctly ignores building action when not provided", func(t *testing.T) {
		planet := models.Planet{
			Id:             uuid.New(),
			Name:           "planet-1",
			CreatedAt:      someTime,
			UpdatedAt:      someOtherTime,
			Resources:      []models.PlanetResource{},
			Storages:       []models.PlanetResourceStorage{},
			Productions:    []models.PlanetResourceProduction{},
			Buildings:      []models.PlanetBuilding{},
			BuildingAction: nil,
		}
		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(planet, nil)

		handler := generateHandler[drivingports.ForManagingPlanet](
			getPlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/:id", handler)

		req := generateTestRequest(t, http.MethodGet, addSampleUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[dtos.PlanetDtoResponse](t, rw)
		expected := dtos.PlanetDtoResponse{
			Id:             planet.Id,
			Name:           planet.Name,
			CreatedAt:      planet.CreatedAt,
			UpdatedAt:      planet.UpdatedAt,
			Resources:      []dtos.PlanetResourceDtoResponse{},
			Storages:       []dtos.PlanetResourceStorageDtoResponse{},
			Productions:    []dtos.PlanetResourceProductionDtoResponse{},
			Buildings:      []dtos.PlanetBuildingDtoResponse{},
			BuildingAction: nil,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 404 when planet does not exist", func(t *testing.T) {
		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(models.Planet{}, domainerrors.ErrNotFound)

		handler := generateHandler[drivingports.ForManagingPlanet](
			getPlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/:id", handler)

		req := generateTestRequest(t, http.MethodGet, addSampleUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such planet", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Planet{}, errors.New("stubbed error"))

		handler := generateHandler[drivingports.ForManagingPlanet](
			getPlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/:id", handler)

		req := generateTestRequest(t, http.MethodGet, addSampleUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to get planet", actual)
	})
}

func TestUnit_Planets_ListPlanetsForPlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	t.Run("returns 400 when player id is invalid", func(t *testing.T) {
		handler := generateHandler[drivingports.ForManagingPlanet](
			listPlanetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/planets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/not-a-uuid/planets")
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards listing to use case", func(t *testing.T) {
		planets := []models.Planet{
			{Id: uuid.New(), Name: "planet-1", CreatedAt: someTime},
			{Id: uuid.New(), Name: "planet-2", CreatedAt: someOtherTime},
		}
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(planets, nil)

		handler := generateHandler[drivingports.ForManagingPlanet](
			listPlanetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/planets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/%s/planets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		expected := []dtos.PlanetDtoResponse{
			{
				Id:          planets[0].Id,
				Name:        planets[0].Name,
				CreatedAt:   planets[0].CreatedAt,
				Resources:   []dtos.PlanetResourceDtoResponse{},
				Storages:    []dtos.PlanetResourceStorageDtoResponse{},
				Productions: []dtos.PlanetResourceProductionDtoResponse{},
				Buildings:   []dtos.PlanetBuildingDtoResponse{},
			},
			{
				Id:          planets[1].Id,
				Name:        planets[1].Name,
				CreatedAt:   planets[1].CreatedAt,
				Resources:   []dtos.PlanetResourceDtoResponse{},
				Storages:    []dtos.PlanetResourceStorageDtoResponse{},
				Productions: []dtos.PlanetResourceProductionDtoResponse{},
				Buildings:   []dtos.PlanetBuildingDtoResponse{},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("return empty slice when use case returns no planet", func(t *testing.T) {
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return([]models.Planet{}, nil)

		handler := generateHandler[drivingports.ForManagingPlanet](
			listPlanetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/planets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/%s/planets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlanetDtoResponse{}, actual)
	})

	t.Run("return empty slice when use case returns nil response", func(t *testing.T) {
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(nil, nil)

		handler := generateHandler[drivingports.ForManagingPlanet](
			listPlanetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/planets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/%s/planets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlanetDtoResponse{}, actual)
	})

	t.Run("returns 500 when use cas fails", func(t *testing.T) {
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return([]models.Planet{}, errors.New("stubbed error"))

		handler := generateHandler[drivingports.ForManagingPlanet](
			listPlanetsForPlayer,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodGet, "/players/:id/planets", handler)

		req := generateTestRequest(t, http.MethodGet)
		addRequestPath(t, req, "/players/%s/planets", sampleUuid)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to list planets", actual)
	})
}

func TestUnit_Planets_DeletePlanet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		handler := generateHandler[drivingports.ForManagingPlanet](
			deletePlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

		req := generateTestRequest(t, http.MethodDelete, addInvalidUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards deletion to use case", func(t *testing.T) {
		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(nil)

		handler := generateHandler[drivingports.ForManagingPlanet](
			deletePlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

		req := generateTestRequest(t, http.MethodDelete, addSampleUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	t.Run("returns 409 when use case returns action is not completed", func(t *testing.T) {
		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(domainerrors.ErrActionNotCompleted)

		handler := generateHandler[drivingports.ForManagingPlanet](
			deletePlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

		req := generateTestRequest(t, http.MethodDelete, addSampleUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "action not completed", actual)
	})

	t.Run("returns 409 when use case returns homeworld cannot be deleted", func(t *testing.T) {
		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(domainerrors.ErrHomeworldCannotBeDeleted)

		handler := generateHandler[drivingports.ForManagingPlanet](
			deletePlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

		req := generateTestRequest(t, http.MethodDelete, addSampleUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "homeworld cannot be deleted", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(errors.New("stubbed error"))

		handler := generateHandler[drivingports.ForManagingPlanet](
			deletePlanet,
			mockUsecase,
		)
		r := createTestGinRouter(t, http.MethodDelete, "/:id", handler)

		req := generateTestRequest(t, http.MethodDelete, addSampleUuidPathParam)
		rw := httptest.NewRecorder()
		r.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to delete planet", actual)
	})
}
