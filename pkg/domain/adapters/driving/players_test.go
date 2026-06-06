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

func TestUnit_Players_CreatePlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlayer(ctrl)

	dto := dtos.PlayerDtoRequest{
		ApiUser:  uuid.New(),
		Universe: uuid.New(),
		Name:     "my-Player",
	}

	t.Run("returns 400 when body is invalid", func(t *testing.T) {
		req := generateTestRequestWithJsonBody(t, http.MethodPost, "not-a-dto-request")
		ctx, rw := generateTestContextFromRequest(t, req)

		err := CreatePlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid player syntax", actual)
	})

	t.Run("forwards creation to use case", func(t *testing.T) {
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req)

		expectedRequest := request.PlayerCreationRequest{
			ApiUser:  dto.ApiUser,
			Universe: dto.Universe,
			Name:     dto.Name,
		}
		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Eq(expectedRequest)).
			Times(1).
			Return(models.Player{
				Id:        sampleUuid,
				ApiUser:   dto.ApiUser,
				Universe:  dto.Universe,
				Name:      dto.Name,
				CreatedAt: someTime,
				Version:   0,
			}, nil)

		err := CreatePlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusCreated, rw.Code)
		actual := decodeResponseBody[dtos.PlayerDtoResponse](t, rw)
		expected := dtos.PlayerDtoResponse{
			Id:        sampleUuid,
			ApiUser:   dto.ApiUser,
			Universe:  dto.Universe,
			Name:      dto.Name,
			CreatedAt: someTime,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		req := generateTestRequestWithJsonBody(t, http.MethodPost, dto)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Player{}, errors.New("stubbed error"))

		err := CreatePlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to create player", actual)
	})
}

func TestUnit_Players_GetPlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlayer(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := GetPlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards fetching to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		player := models.Player{
			Id:        uuid.New(),
			ApiUser:   uuid.New(),
			Universe:  uuid.New(),
			Name:      "player-1",
			CreatedAt: someTime,
		}
		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(player, nil)

		err := GetPlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[dtos.PlayerDtoResponse](t, rw)
		expected := dtos.PlayerDtoResponse{
			Id:        player.Id,
			ApiUser:   player.ApiUser,
			Universe:  player.Universe,
			Name:      player.Name,
			CreatedAt: player.CreatedAt,
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns 404 when player does not exist", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(models.Player{}, errors.NewCode(db.NoMatchingRows))

		err := GetPlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "no such player", actual)
	})

	t.Run("returns 500 when use case fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req, addIdPathParam)

		mockUsecase.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Player{}, errors.New("stubbed error"))

		err := GetPlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to get player", actual)
	})
}

func TestUnit_Players_ListPlayers(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlayer(ctrl)

	t.Run("forwards listing to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		players := []models.Player{
			{Id: uuid.New(), Name: "player-1", CreatedAt: someTime},
			{Id: uuid.New(), Name: "player-2", CreatedAt: someOtherTime},
		}
		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(players, nil)

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlayerDtoResponse](t, rw)
		expected := []dtos.PlayerDtoResponse{
			{Id: players[0].Id, Name: players[0].Name, CreatedAt: players[0].CreatedAt},
			{Id: players[1].Id, Name: players[1].Name, CreatedAt: players[1].CreatedAt},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("return empty slice when use case returns no player", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return([]models.Player{}, nil)

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlayerDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlayerDtoResponse{}, actual)
	})

	t.Run("return empty slice when use case returns nil response", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(nil, nil)

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlayerDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlayerDtoResponse{}, actual)
	})

	t.Run("returns 500 when use cas fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return([]models.Player{}, errors.New("stubbed error"))

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to list players", actual)
	})
}

func TestUnit_Players_ListPlayers_ForApiUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlayer(ctrl)

	t.Run("returns 400 when api user id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "api_user", "not-a-uuid")
		ctx, rw := generateTestContextFromRequest(t, req)

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "invalid id syntax", actual)
	})

	t.Run("forwards listing to use case", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "api_user", sampleUuid.String())
		ctx, rw := generateTestContextFromRequest(t, req)

		players := []models.Player{
			{Id: uuid.New(), Name: "player-1", CreatedAt: someTime},
			{Id: uuid.New(), Name: "player-2", CreatedAt: someOtherTime},
		}
		mockUsecase.EXPECT().
			ListForApiUser(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(players, nil)

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlayerDtoResponse](t, rw)
		expected := []dtos.PlayerDtoResponse{
			{Id: players[0].Id, Name: players[0].Name, CreatedAt: players[0].CreatedAt},
			{Id: players[1].Id, Name: players[1].Name, CreatedAt: players[1].CreatedAt},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("return empty slice when use case returns no player", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "api_user", sampleUuid.String())
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			ListForApiUser(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return([]models.Player{}, nil)

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlayerDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlayerDtoResponse{}, actual)
	})

	t.Run("return empty slice when use case returns nil response", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "api_user", sampleUuid.String())
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			ListForApiUser(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return(nil, nil)

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusOK, rw.Code)
		actual := decodeResponseBody[[]dtos.PlayerDtoResponse](t, rw)
		assert.Equal(t, []dtos.PlayerDtoResponse{}, actual)
	})

	t.Run("returns 500 when use cas fails", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodGet)
		addQueryParam(t, req, "api_user", sampleUuid.String())
		ctx, rw := generateTestContextFromRequest(t, req)

		mockUsecase.EXPECT().
			ListForApiUser(gomock.Any(), gomock.Eq(sampleUuid)).
			Times(1).
			Return([]models.Player{}, errors.New("stubbed error"))

		err := ListPlayers(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to list players", actual)
	})
}

func TestUnit_Players_DeletePlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsecase := drivingportstest.NewMockForManagingPlayer(ctrl)

	t.Run("returns 400 when id is invalid", func(t *testing.T) {
		req := generateTestRequest(t, http.MethodDelete)
		ctx, rw := generateTestContextFromRequest(t, req)
		ctx.SetPathValues([]echo.PathValue{{Name: "id", Value: "not-a-uuid"}})

		err := DeletePlayer(ctx, mockUsecase)
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

		err := DeletePlayer(ctx, mockUsecase)
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

		err := DeletePlayer(ctx, mockUsecase)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		actual := decodeResponseBody[string](t, rw)
		assert.Equal(t, "failed to delete player", actual)
	})
}
