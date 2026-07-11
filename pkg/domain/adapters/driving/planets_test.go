package drivingadapters

import (
	"net/http"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/drivingportstest"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUnit_Planets_CreatePlanet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForCreatingPlanet(ctrl)

	planetId := uuid.New()

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodPost)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := createPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards creation to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodPost)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		expectedRequest := request.PlanetCreationRequest{
			Player: sampleUuid,
		}
		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Eq(expectedRequest)).
			Times(1).
			Return(models.Planet{
				Id:        planetId,
				Player:    sampleUuid,
				Name:      "my-planet",
				CreatedAt: someTime,
				Version:   0,
			}, nil)

		err := createPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.PlanetDtoResponse](t, rw)
		expected := dtos.PlanetDtoResponse{
			Id:          planetId,
			Player:      sampleUuid,
			Name:        "my-planet",
			CreatedAt:   someTime,
			Resources:   []dtos.PlanetResourceDtoResponse{},
			Storages:    []dtos.PlanetResourceStorageDtoResponse{},
			Productions: []dtos.PlanetResourceProductionDtoResponse{},
			Buildings:   []dtos.PlanetBuildingDtoResponse{},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 400 when player is not found", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodPost)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Planet{}, domainerrors.ErrPlayerNotFound)

		err := createPlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such player", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodPost)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

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
			Return(models.Planet{}, domainerrors.ErrNotFound)

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

func TestUnit_Planets_ListPlanetsForPlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlanet(ctrl)

	t.Run("returns 400 when player id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := listPlanetsForPlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards listing to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		planets := []models.Planet{
			{Id: uuid.New(), Name: "planet-1", CreatedAt: someTime},
			{Id: uuid.New(), Name: "planet-2", CreatedAt: someOtherTime},
		}
		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(planets, nil)

		err := listPlanetsForPlayer(ctx, mockUsecase)
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
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return([]models.Planet{}, nil)

		err := listPlanetsForPlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlanetDtoResponse{}, actual)
	})

	t.Run("return empty slice when use case returns nil response", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(nil, nil)

		err := listPlanetsForPlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlanetDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlanetDtoResponse{}, actual)
	})

	t.Run("returns 500 when use cas fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return([]models.Planet{}, errors.New("stubbed error"))

		err := listPlanetsForPlayer(ctx, mockUsecase)
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

	t.Run("returns 409 when use case returns action is not completed", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(domainerrors.ErrActionNotCompleted)

		err := deletePlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "action not completed", actual)
	})

	t.Run("returns 409 when use case returns homeworld cannot be deleted", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(domainerrors.ErrHomeworldCannotBeDeleted)

		err := deletePlanet(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "homeworld cannot be deleted", actual)
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
