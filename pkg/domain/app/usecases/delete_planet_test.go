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

type deletePlanetTestSuite struct {
	ctrl         *gomock.Controller
	mockExecutor *drivenportstest.MockForExecutingSimulationStep
	mockClock    *drivenportstest.MockForFetchingTime
	usecase      *DeletePlanetUseCase
}

func TestUnit_DeletePlanet_Delete(t *testing.T) {
	t.Run("calls executor with planet and timestamp", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		result := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
		}
		suite.mockExecutor.EXPECT().
			ExecuteForDelete(gomock.Any(), id, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		err := suite.usecase.Delete(t.Context(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("retries executor when fleet event is resolved", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		fleetResult := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Fleet: &models.Fleet{
				Id: uuid.New(),
			},
		}
		suite.mockExecutor.EXPECT().
			ExecuteForDelete(gomock.Any(), id, t2, gomock.Any()).
			Times(1).
			Return(fleetResult, nil)
		result := models.SimulationStepResult{
			Resolution: models.PlanetDeleted,
		}
		suite.mockExecutor.EXPECT().
			ExecuteForDelete(gomock.Any(), id, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		err := suite.usecase.Delete(t.Context(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("updates planet to current time", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)

		planet := models.Planet{
			Id:        uuid.New(),
			CreatedAt: t1,
			UpdatedAt: t1,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			ExecuteForDelete(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t2, planet.UpdatedAt)
	})

	t.Run("returns error when planet is homeworld", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)

		planet := models.Planet{
			Id:        uuid.New(),
			Homeworld: true,
			CreatedAt: t1,
			UpdatedAt: t1,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			ExecuteForDelete(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.Delete(t.Context(), planet.Id)

		assert.ErrorIs(t, err, domainerrors.ErrHomeworldCannotBeDeleted, "Actual err: %v", err)
	})

	t.Run("returns error when planet has building action", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)

		planet := models.Planet{
			Id:        uuid.New(),
			Homeworld: false,
			CreatedAt: t1,
			UpdatedAt: t1,
			BuildingAction: &models.BuildingAction{
				Id:          uuid.New(),
				CreatedAt:   t1,
				CompletedAt: t3,
			},
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			ExecuteForDelete(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.Delete(t.Context(), planet.Id)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingActionNotCompleted, "Actual err: %v", err)
	})

	t.Run("returns error when planet has ship actions", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)

		planet := models.Planet{
			Id:        uuid.New(),
			Homeworld: false,
			CreatedAt: t1,
			UpdatedAt: t1,
			ShipActions: []models.ShipAction{
				{
					Id:                 uuid.New(),
					Ship:               uuid.New(),
					Count:              1,
					CreatedAt:          t1,
					NextCompletionAt:   t3,
					UnitCompletionTime: t3.Sub(t1),
				},
			},
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			ExecuteForDelete(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.Delete(t.Context(), planet.Id)

		assert.ErrorIs(t, err, domainerrors.ErrShipActionNotCompleted, "Actual err: %v", err)
	})

	t.Run("returns error when deleter fails", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		expected := errors.New("stubbed error")
		suite.mockExecutor.EXPECT().
			ExecuteForDelete(gomock.Any(), id, t2, gomock.Any()).
			Times(1).
			Return(models.SimulationStepResult{}, expected)

		err := suite.usecase.Delete(t.Context(), id)

		assert.ErrorIs(t, err, expected, "Actual err: %v", err)
	})
}

func setupDeletePlanetTestSuite(t *testing.T) *deletePlanetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockExecutor := drivenportstest.NewMockForExecutingSimulationStep(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &deletePlanetTestSuite{
		ctrl:         ctrl,
		mockExecutor: mockExecutor,
		mockClock:    mockClock,
		usecase:      NewDeletePlanetUseCase(mockExecutor, mockClock),
	}
}
