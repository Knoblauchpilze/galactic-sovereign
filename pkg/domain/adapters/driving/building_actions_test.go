package drivingadapters

import (
	"net/http"
	"testing"

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
		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Eq(expectedRequest)).
			Times(1).
			Return(models.BuildingAction{
				Id:        sampleUuid,
				Planet:    sampleUuid,
				CreatedAt: someTime,
				Version:   0,
			}, nil)

		err := createBuildingAction(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.BuildingActionDtoResponse](t, rw)
		expected := dtos.BuildingActionDtoResponse{
			Id:        sampleUuid,
			Planet:    sampleUuid,
			CreatedAt: someTime,
		}
		assert.Equal(t, expected, actual)
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
