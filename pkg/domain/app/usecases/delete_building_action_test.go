package usecases

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type deleteBuildingActionTestSuite struct {
	ctrl         *gomock.Controller
	mockExecutor *drivenportstest.MockForExecutingSimulationStep
	mockClock    *drivenportstest.MockForFetchingTime
	usecase      *DeleteBuildingActionUseCase
}

func TestUnit_DeleteBuildingAction_DeleteForPlanet(t *testing.T) {
	t.Run("calls executor with planet and timestamp", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanetWithBuildingAction(t3)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		result := models.SimulationStepResult{
			Resolution: models.PlanetUpdated,
			Planet:     &planet,
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(result, nil)

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("retries executor when fleet event is resolved", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanetWithBuildingAction(t3)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		fleetResult := models.SimulationStepResult{
			Resolution: models.FleetEventResolved,
			Fleet: &models.Fleet{
				Id: uuid.New(),
			},
		}
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			Return(fleetResult, nil)
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, planet.BuildingAction)
	})

	t.Run("persists deleted building action", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanetWithBuildingAction(t3)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, planet.BuildingAction)
	})

	t.Run("updates planet to current time", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanetWithBuildingAction(t3)
		planet.UpdatedAt = t1

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t2, planet.UpdatedAt)
	})

	t.Run("applies completed action", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanetWithBuilding()
		require.NotEqual(t, 5, planet.Buildings[0].Level)
		planet.UpdatedAt = t1
		planet.BuildingAction = &models.BuildingAction{
			Id:           uuid.New(),
			Building:     planet.Buildings[0].Building,
			DesiredLevel: 5,
			CreatedAt:    t1,
			CompletedAt:  t2,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t3)
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t3, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, planet.BuildingAction)
		assert.Equal(t, 5, planet.Buildings[0].Level)
		assert.Equal(t, t3, planet.UpdatedAt)
	})

	t.Run("returns no error when planet has no building action", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanetWithBuilding()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockExecutor.EXPECT().
			Execute(gomock.Any(), planet.Id, t2, gomock.Any()).
			Times(1).
			DoAndReturn(generatePlanetSimulatorRecorderMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, planet.BuildingAction)
	})
}

func setupDeleteBuildingActionTestSuite(t *testing.T) *deleteBuildingActionTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockExecutor := drivenportstest.NewMockForExecutingSimulationStep(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &deleteBuildingActionTestSuite{
		ctrl:         ctrl,
		mockExecutor: mockExecutor,
		mockClock:    mockClock,
		usecase:      NewDeleteBuildingActionUseCase(mockExecutor, mockClock),
	}
}
