package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type createPlanetTestSuite struct {
	ctrl             *gomock.Controller
	mockPlayerRepo   *drivenportstest.MockForManagingPlayers
	mockUniverseRepo *drivenportstest.MockForManagingUniverses
	mockPlanetRepo   *drivenportstest.MockForCreatingPlanets
	usecase          *CreatePlanetUseCase
}

func TestUnit_CreatePlanet_Create(t *testing.T) {
	request := request.PlanetCreationRequest{
		Player: uuid.New(),
	}

	player := models.Player{
		Id:       request.Player,
		Universe: uuid.New(),
	}

	universe := models.Universe{
		Id: player.Universe,
	}

	t.Run("persists created planet", func(t *testing.T) {
		suite := setupCreatePlanetTestSuite(t)
		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Player)).
			Times(1).
			Return(player, nil)
		suite.mockUniverseRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(player.Universe)).
			Times(1).
			Return(universe, nil)
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
		assert.Equal(t, "colony", captured.Name)
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

	t.Run("returns error when player is not found", func(t *testing.T) {
		suite := setupCreatePlanetTestSuite(t)
		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Player{}, domainerrors.ErrNotFound)

		_, err := suite.usecase.Create(context.Background(), request)

		assert.ErrorIs(t, domainerrors.ErrPlayerNotFound, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		suite := setupCreatePlanetTestSuite(t)
		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(request.Player)).
			Times(1).
			Return(player, nil)
		suite.mockUniverseRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(player.Universe)).
			Times(1).
			Return(universe, nil)
		expectedErr := errors.New("stubbed error")
		suite.mockPlanetRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		_, err := suite.usecase.Create(context.Background(), request)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func setupCreatePlanetTestSuite(t *testing.T) *createPlanetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockPlayerRepo := drivenportstest.NewMockForManagingPlayers(ctrl)
	mockUniverseRepo := drivenportstest.NewMockForManagingUniverses(ctrl)
	mockPlanetRepo := drivenportstest.NewMockForCreatingPlanets(ctrl)

	return &createPlanetTestSuite{
		ctrl:             ctrl,
		mockPlayerRepo:   mockPlayerRepo,
		mockUniverseRepo: mockUniverseRepo,
		mockPlanetRepo:   mockPlanetRepo,
		usecase: NewCreatePlanetUseCase(
			mockPlayerRepo,
			mockUniverseRepo,
			mockPlanetRepo,
		),
	}
}
