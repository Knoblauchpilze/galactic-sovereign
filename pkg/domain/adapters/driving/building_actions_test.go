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

func TestUnit_BuildingActions_CreateBuildingAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingBuildingAction(ctrl)

	t.Run("returns 400 when planet id is invalid", func(t *testing.T) {
		dto := dtos.BuildingActionDtoRequest{Building: uuid.New()}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("returns 400 when body is invalid", func(t *testing.T) {
		req := generateTestRequestWithJsonBody(t, http.MethodPost, "not-a-dto-request")
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid building action syntax", actual)
	})

	t.Run("forwards creation to use case", func(t *testing.T) {
		dto := dtos.BuildingActionDtoRequest{Building: uuid.New()}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		expectedRequest := request.BuildingActionCreationRequest{Planet: sampleUuid, Building: dto.Building}
		action := models.BuildingAction{
			Id:           uuid.New(),
			Building:     dto.Building,
			DesiredLevel: 6,
			CreatedAt:    someTime,
			CompletedAt:  someOtherTime,
			Costs: []models.BuildingActionCost{
				{
					Resource: uuid.New(),
					Amount:   1478,
				},
			},
			Storages: []models.BuildingActionResourceStorage{
				{
					Resource: uuid.New(),
					Storage:  48790,
				},
			},
			Productions: []models.BuildingActionResourceProduction{
				{
					Resource:   uuid.New(),
					Production: 12,
				},
				{
					Resource:   uuid.New(),
					Production: 8917,
				},
			},
		}

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Eq(expectedRequest)).
			Times(1).
			Return(action, nil)

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.BuildingActionDtoResponse](t, rw)
		expected := dtos.BuildingActionDtoResponse{
			Id:           action.Id,
			Building:     action.Building,
			DesiredLevel: action.DesiredLevel,
			CreatedAt:    action.CreatedAt,
			CompletedAt:  action.CompletedAt,
			Costs: []dtos.BuildingActionCostDtoResponse{
				{Resource: action.Costs[0].Resource, Amount: 1478},
			},
			Storages: []dtos.BuildingActionStorageDtoResponse{
				{Resource: action.Storages[0].Resource, Storage: 48790},
			},
			Productions: []dtos.BuildingActionProductionDtoResponse{
				{
					Resource:   action.Productions[0].Resource,
					Production: 12,
				},
				{
					Resource:   action.Productions[1].Resource,
					Production: 8917,
				},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 409 when action already exists", func(t *testing.T) {
		dto := dtos.BuildingActionDtoRequest{Building: uuid.New()}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.BuildingAction{}, domainerrors.ErrActionAlreadyInProgress)

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusConflict, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "action already in progress", actual)
	})

	t.Run("returns 404 when planet is not found", func(t *testing.T) {
		dto := dtos.BuildingActionDtoRequest{Building: uuid.New()}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.BuildingAction{}, domainerrors.ErrNotFound)

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such planet", actual)
	})

	t.Run("returns 400 when building is not found", func(t *testing.T) {
		dto := dtos.BuildingActionDtoRequest{Building: uuid.New()}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.BuildingAction{}, domainerrors.ErrBuildingNotFound)

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such building", actual)
	})

	t.Run("returns 400 when not enough resources are on the planet", func(t *testing.T) {
		dto := dtos.BuildingActionDtoRequest{Building: uuid.New()}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.BuildingAction{}, domainerrors.ErrNotEnoughResources)

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "not enough resources", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		dto := dtos.BuildingActionDtoRequest{Building: uuid.New()}
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.BuildingAction{}, errors.New("stubbed error"))

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to create building action", actual)
	})
}

func TestUnit_BuildingActions_DeleteBuildingAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingBuildingAction(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := deleteBuildingAction(ctx, mockUsecase)
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

		err := deleteBuildingAction(ctx, mockUsecase)
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

		err := deleteBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to delete building action", actual)
	})
}
