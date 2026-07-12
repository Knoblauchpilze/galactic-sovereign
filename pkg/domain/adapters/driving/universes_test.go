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

func TestUnit_Universes_CreateUniverse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingUniverse(ctrl)

	t.Run("returns 400 when body is invalid", func(t *testing.T) {
		req := generateTestRequestWithJsonBody(t, http.MethodPost, "not-a-dto-request")
		ctx, rw := generateTestContextFromRequest(t, req)

		err := createUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid universe syntax", actual)
	})

	t.Run("forwards creation to use case", func(t *testing.T) {
		dto := dtos.UniverseDtoRequest{Name: "my-universe"}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req)

		expectedRequest := request.UniverseCreationRequest{Name: dto.Name}
		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Eq(expectedRequest)).
			Times(1).
			Return(models.Universe{
				Id:        sampleUuid,
				Name:      dto.Name,
				CreatedAt: someTime,
				Version:   0,
				Resources: []models.Resource{
					{
						Id:              sampleResourceId,
						Name:            "resource",
						StartAmount:     26,
						StartProduction: 47,
						StartStorage:    1055,
						CreatedAt:       someOtherTime,
					},
				},
			}, nil)

		err := createUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.UniverseDtoResponse](t, rw)
		expected := dtos.UniverseDtoResponse{
			Id:        sampleUuid,
			Name:      dto.Name,
			CreatedAt: someTime,
			Resources: []dtos.ResourceDtoResponse{
				{
					Id:              sampleResourceId,
					Name:            "resource",
					StartAmount:     26,
					StartProduction: 47,
					StartStorage:    1055,
					CreatedAt:       someOtherTime,
				},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 409 when name is already taken", func(t *testing.T) {
		dto := dtos.UniverseDtoRequest{Name: "my-universe"}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Universe{}, domainerrors.ErrNameAlreadyTaken)

		err := createUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "name already used", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		dto := dtos.UniverseDtoRequest{Name: "my-universe"}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Universe{}, errors.New("stubbed error"))

		err := createUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to create universe", actual)
	})
}

func TestUnit_Universes_GetUniverse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingUniverse(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := getUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards fetching to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		buildingId := uuid.New()
		buildingCostResourceId := uuid.New()
		buildingProductionResourceId := uuid.New()
		buildingStorageResourceId := uuid.New()

		universe := models.Universe{
			Id:        uuid.New(),
			Name:      "universe-1",
			CreatedAt: someTime,
			Resources: []models.Resource{
				{
					Id:              sampleResourceId,
					Name:            "resource",
					StartAmount:     26,
					StartProduction: 47,
					StartStorage:    1055,
					CreatedAt:       someOtherTime,
				},
			},
			Buildings: []models.Building{
				{
					Id:        buildingId,
					Name:      "building",
					CreatedAt: someOtherTime,
					Costs: []models.BuildingCost{
						{
							Resource: buildingCostResourceId,
							Cost:     42,
							Progress: 1.5,
						},
					},
					Productions: []models.BuildingResourceProduction{
						{
							Resource: buildingProductionResourceId,
							Base:     30,
							Progress: 1.1,
						},
					},
					Storages: []models.BuildingResourceStorage{
						{
							Resource: buildingStorageResourceId,
							Base:     5000,
							Scale:    2.5,
							Progress: 1.833,
						},
					},
				},
			},
		}
		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(universe, nil)

		err := getUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[dtos.UniverseDtoResponse](t, rw)
		expected := dtos.UniverseDtoResponse{
			Id:        universe.Id,
			Name:      universe.Name,
			CreatedAt: universe.CreatedAt,
			Resources: []dtos.ResourceDtoResponse{
				{
					Id:              sampleResourceId,
					Name:            "resource",
					StartAmount:     26,
					StartProduction: 47,
					StartStorage:    1055,
					CreatedAt:       someOtherTime,
				},
			},
			Buildings: []dtos.BuildingDtoResponse{
				{
					Id:        buildingId,
					Name:      "building",
					CreatedAt: someOtherTime,
					Costs: []dtos.BuildingCostDtoResponse{
						{
							Resource: buildingCostResourceId,
							Cost:     42,
							Progress: 1.5,
						},
					},
					Productions: []dtos.BuildingResourceProductionDtoResponse{
						{
							Resource: buildingProductionResourceId,
							Base:     30,
							Progress: 1.1,
						},
					},
					Storages: []dtos.BuildingResourceStorageDtoResponse{
						{
							Resource: buildingStorageResourceId,
							Base:     5000,
							Scale:    2.5,
							Progress: 1.833,
						},
					},
				},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 404 when universe does not exist", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(models.Universe{}, domainerrors.ErrNotFound)

		err := getUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such universe", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Universe{}, errors.New("stubbed error"))

		err := getUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to get universe", actual)
	})
}

func TestUnit_Universes_ListUniverses(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingUniverse(ctrl)

	t.Run("forwards listing to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		universes := []models.Universe{
			{
				Id:        uuid.New(),
				Name:      "universe-1",
				CreatedAt: someTime,
				Resources: []models.Resource{},
			},
			{
				Id:        uuid.New(),
				Name:      "universe-2",
				CreatedAt: someOtherTime,
				Resources: []models.Resource{},
			},
		}
		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(universes, nil)

		err := listUniverses(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.UniverseDtoResponse](t, rw)
		expected := []dtos.UniverseDtoResponse{
			{
				Id:        universes[0].Id,
				Name:      universes[0].Name,
				CreatedAt: universes[0].CreatedAt,
				Resources: []dtos.ResourceDtoResponse{},
			},
			{
				Id:        universes[1].Id,
				Name:      universes[1].Name,
				CreatedAt: universes[1].CreatedAt,
				Resources: []dtos.ResourceDtoResponse{},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("return empty slice when use case returns no universe", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return([]models.Universe{}, nil)

		err := listUniverses(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.UniverseDtoResponse](t, rw)
		assert.Equal(t, []dtos.UniverseDtoResponse{}, actual)
	})

	t.Run("return empty slice when use case returns nil response", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(nil, nil)

		err := listUniverses(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.UniverseDtoResponse](t, rw)
		assert.Equal(t, []dtos.UniverseDtoResponse{}, actual)
	})

	t.Run("returns 500 when use cas fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return([]models.Universe{}, errors.New("stubbed error"))

		err := listUniverses(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to list universes", actual)
	})
}

func TestUnit_Universes_DeleteUniverse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingUniverse(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := deleteUniverse(ctx, mockUsecase)
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

		err := deleteUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	t.Run("returns 409 when universe is not empty", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(domainerrors.ErrUniverseIsNotEmpty)

		err := deleteUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "universe is not empty", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(errors.New("stubbed error"))

		err := deleteUniverse(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to delete universe", actual)
	})
}
