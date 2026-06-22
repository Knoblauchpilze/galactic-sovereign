package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type planetTestSuite struct {
	ctrl           *gomock.Controller
	mockPlayerRepo *drivenportstest.MockForManagingPlayers
	mockPlanetRepo *drivenportstest.MockForManagingPlanets
	usecase        drivingports.ForManagingPlanet
}

func TestUnit_ManagePlanet_Create(t *testing.T) {
	request := request.PlanetCreationRequest{
		Player: uuid.New(),
		Name:   "the-best-planet",
	}

	t.Run("persists created planet", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		var captured models.Planet
		suite.mockPlanetRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured)).
			Times(1).
			DoAndReturn(func(ctx context.Context, planet models.Planet) error {
				captured = planet
				return nil
			})

		beforeInsertion := time.Now()

		actual, err := suite.usecase.Create(context.Background(), request)
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
		suite := setupPlanetTestSuite(t)
		expectedErr := errors.New("stubbed error")
		suite.mockPlanetRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		_, err := suite.usecase.Create(context.Background(), request)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_Get(t *testing.T) {
	t.Run("gets existing planet", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expected := models.Planet{
			Id:   uuid.New(),
			Name: "my-planet",
		}

		suite.mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(expected.Id)).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.Get(context.Background(), expected.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expectedErr := errors.New("stubbed error")
		suite.mockPlanetRepo.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Planet{}, expectedErr)

		_, err := suite.usecase.Get(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_List(t *testing.T) {
	t.Run("lists existing planets", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
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

		suite.mockPlanetRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.List(context.Background())
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expectedErr := errors.New("stubbed error")

		suite.mockPlanetRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		_, err := suite.usecase.List(context.Background())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_ListForPlayer(t *testing.T) {
	t.Run("lists existing planets", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
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

		suite.mockPlanetRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Eq(player)).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.ListForPlayer(context.Background(), player)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expectedErr := errors.New("stubbed error")

		suite.mockPlanetRepo.EXPECT().
			ListForPlayer(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		_, err := suite.usecase.ListForPlayer(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlanet_Delete(t *testing.T) {
	t.Run("deletes existing planet", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		id := uuid.New()

		suite.mockPlanetRepo.EXPECT().
			Delete(gomock.Any(), gomock.Eq(id)).
			Times(1).
			Return(nil)

		err := suite.usecase.Delete(context.Background(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupPlanetTestSuite(t)
		expectedErr := errors.New("stubbed error")
		suite.mockPlanetRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		err := suite.usecase.Delete(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func setupPlanetTestSuite(t *testing.T) *planetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockPlayerRepo := drivenportstest.NewMockForManagingPlayers(ctrl)
	mockPlanetRepo := drivenportstest.NewMockForManagingPlanets(ctrl)

	return &planetTestSuite{
		ctrl:           ctrl,
		mockPlayerRepo: mockPlayerRepo,
		mockPlanetRepo: mockPlanetRepo,
		usecase:        NewPlanetUseCase(mockPlayerRepo, mockPlanetRepo),
	}
}
