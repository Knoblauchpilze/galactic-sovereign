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

type deletePlanetTestSuite struct {
	ctrl              *gomock.Controller
	mockPlanetMutator *drivenportstest.MockForMutatingPlanet
	mockClock         *drivenportstest.MockForFetchingTime
	usecase           *DeletePlanetUseCase
}

func TestUnit_DeletePlanet_Delete(t *testing.T) {
	t.Run("deletes existing planet through mutator", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(id), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: true}, nil)

		err := suite.usecase.Delete(t.Context(), id)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when planet has a building action", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(id), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{}, domainerrors.ErrBuildingActionNotCompleted)

		err := suite.usecase.Delete(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingActionNotCompleted, "Actual err: %v", err)
	})

	t.Run("returns error when planet has ship action", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(id), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{}, domainerrors.ErrShipActionNotCompleted)

		err := suite.usecase.Delete(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrShipActionNotCompleted, "Actual err: %v", err)
	})

	t.Run("returns error when mutator returns no error but does not mark the planet as deleted", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(id), gomock.Any()).
			Times(1).
			Return(models.PlanetMutationResult{Deleted: false}, nil)

		err := suite.usecase.Delete(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrPlanetDeletionFailed, "Actual err: %v", err)
	})

	t.Run("returns error when homeworld is deleted", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		planet := models.Planet{
			Id:        uuid.New(),
			Player:    uuid.New(),
			Name:      "my-planet",
			Homeworld: true,
			CreatedAt: t1,
			UpdatedAt: t1,
			Version:   2,
		}

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockPlanetMutator.EXPECT().
			Mutate(gomock.Any(), gomock.Eq(planet.Id), gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingMutatorMock(&planet))

		err := suite.usecase.Delete(t.Context(), planet.Id)

		assert.ErrorIs(t, err, domainerrors.ErrHomeworldCannotBeDeleted, "Actual err: %v", err)
	})
}

func setupDeletePlanetTestSuite(t *testing.T) *deletePlanetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockPlanetMutator := drivenportstest.NewMockForMutatingPlanet(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &deletePlanetTestSuite{
		ctrl:              ctrl,
		mockPlanetMutator: mockPlanetMutator,
		mockClock:         mockClock,
		usecase:           NewDeletePlanetUseCase(mockPlanetMutator, mockClock),
	}
}
