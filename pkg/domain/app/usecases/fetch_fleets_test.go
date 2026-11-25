package usecases

import (
	"errors"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fetchFleetsTestSuite struct {
	ctrl           *gomock.Controller
	mockPlayerRepo *drivenportstest.MockForManagingPlayers
	mockExecutor   *drivenportstest.MockForExecutingSimulationStep
	mockFleetRepo  *drivenportstest.MockForListingFleets
	mockClock      *drivenportstest.MockForFetchingTime
	usecase        *FetchFleetsUseCase
}

// TODO: Extend those tests when the fleets are fetched
func TestUnit_FetchFleetsUseCase_ListForPlayer(t *testing.T) {
	t.Run("calls executor with universe and timestamp", func(t *testing.T) {
		suite := setupFetchFleetsTestSuite(t)

		player := models.Player{
			Id:       uuid.New(),
			Universe: uuid.New(),
			Name:     "the-player",
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), player.Id).
			Times(1).
			Return(player, nil)
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), player.Universe, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, domainerrors.ErrNoFleetEvent)
		fleets := []models.Fleet{
			{Id: uuid.New()},
			{Id: uuid.New()},
		}
		suite.mockFleetRepo.EXPECT().
			ListForPlayer(gomock.Any(), player.Id).
			Times(1).
			Return(fleets, nil)

		actual, err := suite.usecase.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, fleets, actual)
	})

	t.Run("retries executor when fleet event is resolved", func(t *testing.T) {
		suite := setupFetchFleetsTestSuite(t)

		player := models.Player{
			Id:       uuid.New(),
			Universe: uuid.New(),
			Name:     "the-player",
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), player.Id).
			Times(1).
			Return(player, nil)
		result := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Fleet: &models.Fleet{
				Id: uuid.New(),
			},
		}
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), player.Universe, t2, gomock.Any()).
			Times(1).
			Return(result, nil)
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), player.Universe, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, domainerrors.ErrNoFleetEvent)
		fleets := []models.Fleet{
			{Id: uuid.New()},
			{Id: uuid.New()},
		}
		suite.mockFleetRepo.EXPECT().
			ListForPlayer(gomock.Any(), player.Id).
			Times(1).
			Return(fleets, nil)

		actual, err := suite.usecase.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, fleets, actual)
	})

	t.Run("returns error when executor indicates that a planet has been updated", func(t *testing.T) {
		suite := setupFetchFleetsTestSuite(t)

		player := models.Player{
			Id:       uuid.New(),
			Universe: uuid.New(),
			Name:     "the-player",
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), player.Id).
			Times(1).
			Return(player, nil)
		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet: &models.Planet{
				Id:     uuid.New(),
				Player: player.Id,
			},
		}
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), player.Universe, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		_, err := suite.usecase.ListForPlayer(t.Context(), player.Id)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})

	t.Run("returns error when executor indicates that a fleet has been created", func(t *testing.T) {
		suite := setupFetchFleetsTestSuite(t)

		player := models.Player{
			Id:       uuid.New(),
			Universe: uuid.New(),
			Name:     "the-player",
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), player.Id).
			Times(1).
			Return(player, nil)
		result := models.SimulationStepResult{
			Resolution: models.FleetCreated,
			Fleet: &models.Fleet{
				Id:     uuid.New(),
				Player: player.Id,
			},
		}
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), player.Universe, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		_, err := suite.usecase.ListForPlayer(t.Context(), player.Id)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})

	t.Run("returns not found when player does not exist", func(t *testing.T) {
		suite := setupFetchFleetsTestSuite(t)

		player := models.Player{
			Id:       uuid.New(),
			Universe: uuid.New(),
			Name:     "the-player",
		}

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), player.Id).
			Times(1).
			Return(models.Player{}, domainerrors.ErrNotFound)

		_, err := suite.usecase.ListForPlayer(t.Context(), player.Id)

		assert.ErrorIs(t, err, domainerrors.ErrPlayerNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when fetching player fails", func(t *testing.T) {
		suite := setupFetchFleetsTestSuite(t)

		player := models.Player{
			Id:       uuid.New(),
			Universe: uuid.New(),
			Name:     "the-player",
		}

		errSample := errors.New("sample err")
		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), player.Id).
			Times(1).
			Return(models.Player{}, errSample)

		_, err := suite.usecase.ListForPlayer(t.Context(), player.Id)

		assert.ErrorIs(t, err, errSample, "Actual err: %v", err)
	})

	t.Run("returns error when fetching fleets fails", func(t *testing.T) {
		suite := setupFetchFleetsTestSuite(t)

		player := models.Player{
			Id:       uuid.New(),
			Universe: uuid.New(),
			Name:     "the-player",
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)

		suite.mockPlayerRepo.EXPECT().
			Get(gomock.Any(), player.Id).
			Times(1).
			Return(player, nil)
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), player.Universe, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, domainerrors.ErrNoFleetEvent)
		errSample := errors.New("sample err")
		suite.mockFleetRepo.EXPECT().
			ListForPlayer(gomock.Any(), player.Id).
			Times(1).
			Return([]models.Fleet{}, errSample)

		_, err := suite.usecase.ListForPlayer(t.Context(), player.Id)

		assert.ErrorIs(t, err, errSample, "Actual err: %v", err)
	})
}

func setupFetchFleetsTestSuite(t *testing.T) *fetchFleetsTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockPlayerRepo := drivenportstest.NewMockForManagingPlayers(ctrl)
	mockExecutor := drivenportstest.NewMockForExecutingSimulationStep(ctrl)
	mockFleetRepo := drivenportstest.NewMockForListingFleets(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &fetchFleetsTestSuite{
		ctrl:           ctrl,
		mockPlayerRepo: mockPlayerRepo,
		mockExecutor:   mockExecutor,
		mockFleetRepo:  mockFleetRepo,
		mockClock:      mockClock,
		usecase:        NewFetchFleetsUseCase(mockPlayerRepo, mockExecutor, mockFleetRepo, mockClock),
	}
}
