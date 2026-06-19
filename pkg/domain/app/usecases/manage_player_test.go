package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUnit_ManagePlayer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlayers(ctrl)

	request := request.PlayerCreationRequest{
		ApiUser:  uuid.New(),
		Universe: uuid.New(),
		Name:     "the-best-player",
	}

	t.Run("persists created player", func(t *testing.T) {
		var captured models.Player
		var capturedHomeworld models.Planet
		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured), gomock.AssignableToTypeOf(capturedHomeworld)).
			Times(1).
			DoAndReturn(func(ctx context.Context, player models.Player, planet models.Planet) error {
				captured = player
				capturedHomeworld = planet
				return nil
			})

		beforeInsertion := time.Now()

		usecase := NewPlayerUseCase(mockRepo)
		actual, err := usecase.Create(context.Background(), request)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, request.Name, captured.Name)
		assert.True(t, beforeInsertion.Before(captured.CreatedAt))
		assert.Equal(t, 0, captured.Version)
		assert.Equal(t, captured, actual)
		assert.Equal(t, []uuid.UUID{capturedHomeworld.Id}, captured.Planets)

		assert.Equal(t, captured.Id, capturedHomeworld.Player)
		assert.True(t, capturedHomeworld.Homeworld)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewPlayerUseCase(mockRepo)
		_, err := usecase.Create(context.Background(), request)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlayer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlayers(ctrl)

	t.Run("gets existing player", func(t *testing.T) {
		expected := models.Player{
			Id:       uuid.New(),
			ApiUser:  uuid.New(),
			Universe: uuid.New(),
			Name:     "my-player",
		}

		mockRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(expected.Id)).
			Times(1).
			Return(expected, nil)

		usecase := NewPlayerUseCase(mockRepo)
		actual, err := usecase.Get(context.Background(), expected.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		mockRepo.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Player{}, expectedErr)

		usecase := NewPlayerUseCase(mockRepo)
		_, err := usecase.Get(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlayer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlayers(ctrl)

	t.Run("lists existing players", func(t *testing.T) {
		expected := []models.Player{
			{
				Id:       uuid.New(),
				ApiUser:  uuid.New(),
				Universe: uuid.New(),
				Name:     "player-1",
			},
			{
				Id:       uuid.New(),
				ApiUser:  uuid.New(),
				Universe: uuid.New(),
				Name:     "player-2",
			},
		}

		mockRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(expected, nil)

		usecase := NewPlayerUseCase(mockRepo)
		actual, err := usecase.List(context.Background())
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")

		mockRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		usecase := NewPlayerUseCase(mockRepo)
		_, err := usecase.List(context.Background())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlayer_ListForApiUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlayers(ctrl)

	t.Run("lists existing players", func(t *testing.T) {
		apiUser := uuid.New()
		expected := []models.Player{
			{
				Id:       uuid.New(),
				ApiUser:  apiUser,
				Universe: uuid.New(),
				Name:     "player-1",
			},
			{
				Id:       uuid.New(),
				ApiUser:  apiUser,
				Universe: uuid.New(),
				Name:     "player-2",
			},
		}

		mockRepo.EXPECT().
			ListForApiUser(gomock.Any(), gomock.Eq(apiUser)).
			Times(1).
			Return(expected, nil)

		usecase := NewPlayerUseCase(mockRepo)
		actual, err := usecase.ListForApiUser(context.Background(), apiUser)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")

		mockRepo.EXPECT().
			ListForApiUser(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		usecase := NewPlayerUseCase(mockRepo)
		_, err := usecase.ListForApiUser(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlayer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlayers(ctrl)

	t.Run("deletes existing player", func(t *testing.T) {
		id := uuid.New()

		mockRepo.EXPECT().
			Delete(gomock.Any(), gomock.Eq(id)).
			Times(1).
			Return(nil)

		usecase := NewPlayerUseCase(mockRepo)
		err := usecase.Delete(context.Background(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		mockRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewPlayerUseCase(mockRepo)
		err := usecase.Delete(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}
