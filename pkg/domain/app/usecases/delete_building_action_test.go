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

type deleteBuildingActionTestSuite struct {
	ctrl        *gomock.Controller
	mockMutator *drivenportstest.MockForMutatingPlanet
	mockClock   *drivenportstest.MockForFetchingTime
	usecase     *DeleteBuildingActionUseCase
}

func TestUnit_DeleteBuildingAction_DeleteForPlanet(t *testing.T) {
	t.Run("persists deleted building action", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanetWithAction(t2)

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, planet.BuildingAction)
	})

	t.Run("updates planet to current time", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanetWithAction(t2)
		planet.UpdatedAt = t1

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t2, planet.UpdatedAt)
	})

	t.Run("applies completed action", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanet()
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
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, planet.BuildingAction)
		assert.Equal(t, 5, planet.Buildings[0].Level)
		assert.Equal(t, t3, planet.UpdatedAt)
	})

	t.Run("returns no error when planet has no building action", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanet()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, planet.BuildingAction)
	})

	t.Run("returns error when planet is deleted", func(t *testing.T) {
		suite := setupDeleteBuildingActionTestSuite(t)

		planet := generateTestPlanet()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockMutator.EXPECT().
			Mutate(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: true}, nil)

		err := suite.usecase.DeleteForPlanet(t.Context(), planet.Id)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func setupDeleteBuildingActionTestSuite(t *testing.T) *deleteBuildingActionTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockMutator := drivenportstest.NewMockForMutatingPlanet(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &deleteBuildingActionTestSuite{
		ctrl:        ctrl,
		mockMutator: mockMutator,
		mockClock:   mockClock,
		usecase:     NewDeleteBuildingActionUseCase(mockMutator, mockClock),
	}
}
