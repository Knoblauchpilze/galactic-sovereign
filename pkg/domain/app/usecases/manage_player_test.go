package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/request"
	drivingports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driving"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type playerTestSuite struct {
	ctrl             *gomock.Controller
	mockPlayerRepo   *drivenportstest.MockForManagingPlayers
	mockResourceRepo *drivenportstest.MockForListingResources
	mockBuildingRepo *drivenportstest.MockForListingBuildings
	mockPlanetRepo   *drivenportstest.MockForManagingPlanets
	usecase          drivingports.ForManagingPlayer
}

func TestUnit_ManagePlayer_Create(t *testing.T) {
	suite := setupPlayerTestSuite(t)

	resources := []models.Resource{
		{
			Id:              metalResourceId,
			StartAmount:     145,
			StartStorage:    226,
			StartProduction: 897,
		},
	}

	buildings := []models.Building{{Id: uuid.New()}}

	request := request.PlayerCreationRequest{
		ApiUser:  uuid.New(),
		Universe: uuid.New(),
		Name:     "the-best-player",
	}

	t.Run("persists created player", func(t *testing.T) {
		var captured models.Player
		var capturedHomeworld models.Planet
		suite.mockResourceRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(resources, nil)
		suite.mockBuildingRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(buildings, nil)
		suite.mockPlayerRepo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(captured), gomock.AssignableToTypeOf(capturedHomeworld)).
			Times(1).
			DoAndReturn(func(ctx context.Context, player models.Player, planet models.Planet) error {
				captured = player
				capturedHomeworld = planet
				return nil
			})

		beforeInsertion := time.Now()

		actual, err := suite.usecase.Create(context.Background(), request)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, request.Name, captured.Name)
		assert.True(t, beforeInsertion.Before(captured.CreatedAt))
		assert.Equal(t, 0, captured.Version)
		assert.Equal(t, captured, actual)
		assert.Equal(t, []uuid.UUID{capturedHomeworld.Id}, captured.Planets)

		assert.Equal(t, captured.Id, capturedHomeworld.Player)
		assert.True(t, capturedHomeworld.Homeworld)
		expectedResources := []models.PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   145,
			},
		}
		assert.Equal(t, expectedResources, capturedHomeworld.Resources)
		expectedStorages := []models.PlanetResourceStorage{
			{
				Resource: metalResourceId,
				Storage:  226,
			},
		}
		assert.Equal(t, expectedStorages, capturedHomeworld.Storages)
		expectedProductions := []models.PlanetResourceProduction{
			{
				Resource:   metalResourceId,
				Building:   nil,
				Production: 897,
			},
		}
		assert.Equal(t, expectedProductions, capturedHomeworld.Productions)
		expectedBuildings := []models.PlanetBuilding{
			{
				Building: buildings[0].Id,
				Level:    0,
			},
		}
		assert.Equal(t, expectedBuildings, capturedHomeworld.Buildings)
		assert.Nil(t, capturedHomeworld.BuildingAction)
	})

	t.Run("returns error when creation fails", func(t *testing.T) {
		suite.mockResourceRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(resources, nil)
		suite.mockBuildingRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(buildings, nil)
		expectedErr := errors.New("stubbed error")
		suite.mockPlayerRepo.EXPECT().
			Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		_, err := suite.usecase.Create(context.Background(), request)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlayer_Get(t *testing.T) {
	suite := setupPlayerTestSuite(t)

	t.Run("gets existing player", func(t *testing.T) {
		expected := models.Player{
			Id:       uuid.New(),
			ApiUser:  uuid.New(),
			Universe: uuid.New(),
			Name:     "my-player",
		}

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(expected.Id)).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.Get(context.Background(), expected.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")
		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.Player{}, expectedErr)

		_, err := suite.usecase.Get(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlayer_List(t *testing.T) {
	suite := setupPlayerTestSuite(t)

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

		suite.mockPlayerRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.List(context.Background())
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")

		suite.mockPlayerRepo.EXPECT().
			List(gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		_, err := suite.usecase.List(context.Background())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlayer_ListForApiUser(t *testing.T) {
	suite := setupPlayerTestSuite(t)

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

		suite.mockPlayerRepo.EXPECT().
			ListForApiUser(gomock.Any(), gomock.Eq(apiUser)).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.ListForApiUser(context.Background(), apiUser)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := errors.New("stubbed error")

		suite.mockPlayerRepo.EXPECT().
			ListForApiUser(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, expectedErr)

		_, err := suite.usecase.ListForApiUser(context.Background(), uuid.New())

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func TestUnit_ManagePlayer_Delete(t *testing.T) {
	suite := setupPlayerTestSuite(t)

	t.Run("deletes existing player", func(t *testing.T) {
		player := models.Player{Id: uuid.New()}

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(player.Id)).
			Times(1).
			Return(player, nil)
		suite.mockPlayerRepo.EXPECT().
			Delete(gomock.Any(), gomock.Eq(player)).
			Times(1).
			Return(nil)

		err := suite.usecase.Delete(context.Background(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("succeeds when building action is not found", func(t *testing.T) {
		playerId := uuid.New()

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(playerId)).
			Times(1).
			Return(models.Player{}, domainerrors.ErrNotFound)

		err := suite.usecase.Delete(context.Background(), playerId)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		player := models.Player{Id: uuid.New()}

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), gomock.Eq(player.Id)).
			Times(1).
			Return(player, nil)
		expectedErr := errors.New("stubbed error")
		suite.mockPlayerRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Times(1).
			Return(expectedErr)

		err := suite.usecase.Delete(context.Background(), player.Id)

		assert.ErrorIs(t, expectedErr, err, "Actual err: %v", err)
	})
}

func setupPlayerTestSuite(t *testing.T) *playerTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockPlayerRepo := drivenportstest.NewMockForManagingPlayers(ctrl)
	mockResourceRepo := drivenportstest.NewMockForListingResources(ctrl)
	mockBuildingRepo := drivenportstest.NewMockForListingBuildings(ctrl)
	mockPlanetRepo := drivenportstest.NewMockForManagingPlanets(ctrl)

	return &playerTestSuite{
		ctrl:             ctrl,
		mockPlayerRepo:   mockPlayerRepo,
		mockResourceRepo: mockResourceRepo,
		mockBuildingRepo: mockBuildingRepo,
		mockPlanetRepo:   mockPlanetRepo,
		usecase: NewPlayerUseCase(
			mockPlayerRepo,
			mockResourceRepo,
			mockBuildingRepo,
			mockPlanetRepo,
		),
	}
}
