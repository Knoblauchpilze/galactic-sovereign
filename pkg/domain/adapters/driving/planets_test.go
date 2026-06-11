package drivingadapters

import (
	"net/http"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/drivingportstest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUnit_Planets_CreatePlanet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	dto := dtos.PlanetDtoRequest{
		Player: uuid.New(),
		Name:   "my-planet",
	}

	t.Run("returns 400 when body is invalid", func(t *testing.T) {
		req := generateTestRequestWithJsonBody(t, http.MethodPost, "not-a-dto-request")
		ctx, rw := generateTestContextFromRequest(t, req)

		err := createPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid planet syntax", actual)
	})

	t.Run("forwards creation to use case", func(t *testing.T) {
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req)

		expectedRequest := request.PlanetCreationRequest{
			Player: dto.Player,
			Name:   dto.Name,
		}
		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Eq(expectedRequest)).
			Times(1).
			Return(models.Planet{
				Id:        sampleUuid,
				Name:      dto.Name,
				CreatedAt: someTime,
				Version:   0,
			}, nil)

		err := createPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.PlanetDtoResponse](t, rw)
		expected := dtos.PlanetDtoResponse{
			Id:          sampleUuid,
			Name:        dto.Name,
			CreatedAt:   someTime,
			Resources:   []dtos.PlanetResourceDtoResponse{},
			Storages:    []dtos.PlanetResourceStorageDtoResponse{},
			Productions: []dtos.PlanetResourceProductionDtoResponse{},
			Buildings:   []dtos.PlanetBuildingDtoResponse{},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Planet{}, errors.New("stubbed error"))

		err := createPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to create planet", actual)
	})
}

func TestUnit_Planets_GetPlanet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := getPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards fetching to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		planet := models.Planet{
			Id:        uuid.New(),
			Name:      "planet-1",
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
					Building:  uuid.New(),
					Level:     14,
					CreatedAt: someTime,
					UpdatedAt: someOtherTime,
				},
			},
			BuildingAction: &sampleUuid,
		}
		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(planet, nil)

		err := getPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[dtos.PlanetDtoResponse](t, rw)
		expected := dtos.PlanetDtoResponse{
			Id:        planet.Id,
			Name:      planet.Name,
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
					Building:  planet.Buildings[0].Building,
					Level:     planet.Buildings[0].Level,
					CreatedAt: someTime,
					UpdatedAt: someOtherTime,
				},
			},
			BuildingAction: &sampleUuid,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("correctly ignores building action when not provided", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

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

		err := getPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

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
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(models.Planet{}, errors.NewCode(db.NoMatchingRows))

		err := getPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such planet", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Planet{}, errors.New("stubbed error"))

		err := getPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to get planet", actual)
	})
}

func TestUnit_Planets_ListPlanets(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	t.Run("forwards listing to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		planets := []models.Planet{
			{Id: uuid.New(), Name: "planet-1", CreatedAt: someTime},
			{Id: uuid.New(), Name: "planet-2", CreatedAt: someOtherTime},
		}
		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(planets, nil)

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

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
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return([]models.Planet{}, nil)

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlanetDtoResponse{}, actual)
	})

	t.Run("return empty slice when use case returns nil response", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(nil, nil)

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlanetDtoResponse{}, actual)
	})

	t.Run("returns 500 when use cas fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return([]models.Planet{}, errors.New("stubbed error"))

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to list planets", actual)
	})
}

func TestUnit_Planets_ListPlanets_ForPlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	t.Run("returns 400 when api user id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "player", "not-a-uuid")
		ctx, rw := generateTestContextFromRequest(t, req)

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards listing to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "player", sampleUuid.String())
		ctx, rw := generateTestContextFromRequest(t, req)

		planets := []models.Planet{
			{Id: uuid.New(), Name: "planet-1", CreatedAt: someTime},
			{Id: uuid.New(), Name: "planet-2", CreatedAt: someOtherTime},
		}
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(planets, nil)

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

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
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "player", sampleUuid.String())
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return([]models.Planet{}, nil)

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlanetDtoResponse{}, actual)
	})

	t.Run("return empty slice when use case returns nil response", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "player", sampleUuid.String())
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(nil, nil)

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlanetDtoResponse{}, actual)
	})

	t.Run("returns 500 when use cas fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "player", sampleUuid.String())
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return([]models.Planet{}, errors.New("stubbed error"))

		err := listPlanets(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to list planets", actual)
	})
}

func TestUnit_Planets_DeletePlanet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := deletePlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards deletion to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(nil)

		err := deletePlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(errors.New("stubbed error"))

		err := deletePlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to delete planet", actual)
	})
}
