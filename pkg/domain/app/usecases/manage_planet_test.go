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

func TestUnit_ManagePlanet_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlanets(ctrl)

	request := request.PlanetCreationRequest{
		Player: uuid.New(),
		Name:   "the-best-planet",
	}

	t.Run("persists created planet", func(t *testing.T) {
		var captured models.Planet
		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, planet models.Planet) error {
				captured = planet
				return nil
			})

		beforeInsertion := time.Now()

		usecase := NewPlanetUseCase(mockRepo)
		actual, err := usecase.Create(context.Background(), request)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, request.Player, captured.Player)
		assert.Equal(t, request.Name, captured.Name)
		assert.False(t, captured.Homeworld)
		assert.True(t, beforeInsertion.Before(captured.CreatedAt))
		assert.True(t, beforeInsertion.Before(captured.UpdatedAt))
		assert.Equal(t, []models.PlanetResource{}, captured.Resources)
		assert.Equal(t, []models.PlanetResourceStorage{}, captured.Storages)
		assert.Equal(t, []models.PlanetResourceProduction{}, captured.Productions)
		assert.Equal(t, []models.PlanetBuilding{}, captured.Buildings)
		assert.Equal(t, 0, captured.Version)
		assert.Equal(t, captured, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewPlanetUseCase(mockRepo)
		_, err := usecase.Create(context.Background(), request)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlanets(ctrl)

	t.Run("gets existing planet", func(t *testing.T) {
		expected := models.Planet{
			Id:   uuid.New(),
			Name: "my-planet",
		}

		mockRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(expected.Id)).
			Times(1).
			Return(expected, nil)

		usecase := NewPlanetUseCase(mockRepo)
		actual, err := usecase.Get(context.Background(), expected.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		mockRepo.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Planet{}, expectedErr)

		usecase := NewPlanetUseCase(mockRepo)
		_, err := usecase.Get(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlanets(ctrl)

	t.Run("lists existing planets", func(t *testing.T) {
		expected := []models.Planet{
			{
				Id:   uuid.New(),
				Name: "planet-1",
			},
			{
				Id:   uuid.New(),
				Name: "planet-2",
			},
		}

		mockRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(expected, nil)

		usecase := NewPlanetUseCase(mockRepo)
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

		usecase := NewPlanetUseCase(mockRepo)
		_, err := usecase.List(context.Background())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_ListForPlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlanets(ctrl)

	t.Run("lists existing planets", func(t *testing.T) {
		player := uuid.New()
		expected := []models.Planet{
			{
				Id:   uuid.New(),
				Name: "planet-1",
			},
			{
				Id:   uuid.New(),
				Name: "planet-2",
			},
		}

		mockRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(player)).
			Times(1).
			Return(expected, nil)

		usecase := NewPlanetUseCase(mockRepo)
		actual, err := usecase.ListForPlayer(context.Background(), player)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")

		mockRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		usecase := NewPlanetUseCase(mockRepo)
		_, err := usecase.ListForPlayer(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := drivenportstest.NewMockForManagingPlanets(ctrl)

	t.Run("deletes existing planet", func(t *testing.T) {
		id := uuid.New()

		mockRepo.EXPECT().
			Delete(gomock.Any(), gomock.Eq(id)).
			Times(1).
			Return(nil)

		usecase := NewPlanetUseCase(mockRepo)
		err := usecase.Delete(context.Background(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		mockRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		usecase := NewPlanetUseCase(mockRepo)
		err := usecase.Delete(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}
