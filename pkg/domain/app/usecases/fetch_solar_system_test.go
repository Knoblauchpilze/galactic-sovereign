package usecases

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type solarSystemTestSuite struct {
	ctrl                *gomock.Controller
	mockExecutor        *drivenportstest.MockForExecutingSimulationStep
	mockSolarSystemRepo *drivenportstest.MockForFetchingSolarSystem
	mockClock           *drivenportstest.MockForFetchingTime
	usecase             *FetchSolarSystemUseCase
}

func TestUnit_FetchSolarSystem_GetSolarSystem(t *testing.T) {
	t.Run("calls executor with universe and timestamp", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		universeId := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), universeId, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, domainerrors.ErrNoFleetEvent)
		suite.mockSolarSystemRepo.EXPECT().
			GetSolarSystem(gomock.Any(), universeId, 1, 3).
			Times(1).
			Return(models.SolarSystem{}, nil)

		actual, err := suite.usecase.GetSolarSystem(t.Context(), universeId, 1, 3)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, models.SolarSystem{}, actual)
	})

	t.Run("retries executor when fleet event is resolved", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		universeId := uuid.New()

		expected := models.SolarSystem{
			Universe: universeId,
			Galaxy:   1,
			Number:   3,
			Orbits:   5,
			Planets:  []models.SolarSystemPlanet{},
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		result := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Fleet: &models.Fleet{
				Id: uuid.New(),
			},
		}
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), universeId, t2, gomock.Any()).
			Times(1).
			Return(result, nil)
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), universeId, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, domainerrors.ErrNoFleetEvent)
		suite.mockSolarSystemRepo.EXPECT().
			GetSolarSystem(gomock.Any(), universeId, 1, 3).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.GetSolarSystem(t.Context(), universeId, 1, 3)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when executor indicates that a planet has been updated", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		universeId := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet: &models.Planet{
				Id:     uuid.New(),
				Player: uuid.New(),
			},
		}
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), universeId, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		_, err := suite.usecase.GetSolarSystem(t.Context(), universeId, 1, 3)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})

	t.Run("returns error when executor indicates that a fleet has been created", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		universeId := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		result := models.SimulationStepResult{
			Resolution: models.FleetCreated,
			Fleet: &models.Fleet{
				Id:     uuid.New(),
				Player: uuid.New(),
			},
		}
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), universeId, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		_, err := suite.usecase.GetSolarSystem(t.Context(), universeId, 1, 3)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidSimulationStepResult, "Actual err: %v", err)
	})

	t.Run("gets solar system in existing universe", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		universeId := uuid.New()

		expected := models.SolarSystem{
			Universe: universeId,
			Galaxy:   1,
			Number:   3,
			Orbits:   5,
			Planets:  []models.SolarSystemPlanet{},
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), universeId, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, domainerrors.ErrNoFleetEvent)
		suite.mockSolarSystemRepo.EXPECT().
			GetSolarSystem(gomock.Any(), universeId, 1, 3).
			Times(1).
			Return(expected, nil)

		actual, err := suite.usecase.GetSolarSystem(t.Context(), universeId, 1, 3)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when universe does not exist", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		universeId := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		// The executor does return no fleet event when the universe does not exist
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), universeId, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, domainerrors.ErrNoFleetEvent)
		suite.mockSolarSystemRepo.EXPECT().
			GetSolarSystem(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(models.SolarSystem{}, domainerrors.ErrNotFound)

		_, err := suite.usecase.GetSolarSystem(t.Context(), uuid.New(), 2, 4)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when requested coordinates are out of bound", func(t *testing.T) {
		suite := setupSolarSystemTestSuite(t)

		universeId := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			ExecuteFleetForUniverse(gomock.Any(), universeId, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, domainerrors.ErrNoFleetEvent)
		suite.mockSolarSystemRepo.EXPECT().
			GetSolarSystem(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			// This is the same error as when the universe does not exist
			Return(models.SolarSystem{}, domainerrors.ErrNotFound)

		_, err := suite.usecase.GetSolarSystem(t.Context(), universeId, 2, 4)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func setupSolarSystemTestSuite(t *testing.T) *solarSystemTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockExecutor := drivenportstest.NewMockForExecutingSimulationStep(ctrl)
	mockSolarSystemRepo := drivenportstest.NewMockForFetchingSolarSystem(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &solarSystemTestSuite{
		ctrl:                ctrl,
		mockExecutor:        mockExecutor,
		mockSolarSystemRepo: mockSolarSystemRepo,
		mockClock:           mockClock,
		usecase:             NewFetchSolarSystemUseCase(mockExecutor, mockSolarSystemRepo, mockClock),
	}
}
