package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/usecases/drivenportstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type DeleterMock func(context.Context, uuid.UUID, drivenports.PlanetDeleter) error

type deletePlanetTestSuite struct {
	ctrl        *gomock.Controller
	mockDeleter *drivenportstest.MockForDeletingPlanet
	mockClock   *drivenportstest.MockForFetchingTime
	usecase     *DeletePlanetUseCase
}

func TestUnit_DeletePlanet_Delete(t *testing.T) {
	t.Run("deletes existing planet through deleter", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		suite.mockDeleter.EXPECT().
			Delete(gomock.Any(), id, gomock.Any()).
			Times(1).
			Return(nil)

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
		suite.mockDeleter.EXPECT().
			Delete(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingDeleterMock(&planet))

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
		suite.mockDeleter.EXPECT().
			Delete(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingDeleterMock(&planet))

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
		suite.mockDeleter.EXPECT().
			Delete(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingDeleterMock(&planet))

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
		suite.mockDeleter.EXPECT().
			Delete(gomock.Any(), planet.Id, gomock.Any()).
			Times(1).
			DoAndReturn(generateApplyingDeleterMock(&planet))

		err := suite.usecase.Delete(t.Context(), planet.Id)

		assert.ErrorIs(t, err, domainerrors.ErrShipActionNotCompleted, "Actual err: %v", err)
	})

	t.Run("returns error when deleter fails", func(t *testing.T) {
		suite := setupDeletePlanetTestSuite(t)
		id := uuid.New()

		suite.mockClock.EXPECT().Now(gomock.Any()).Times(1).Return(t2)
		expected := errors.New("stubbed error")
		suite.mockDeleter.EXPECT().
			Delete(gomock.Any(), id, gomock.Any()).
			Times(1).
			Return(expected)

		err := suite.usecase.Delete(t.Context(), id)

		assert.ErrorIs(t, err, expected, "Actual err: %v", err)
	})
}

func setupDeletePlanetTestSuite(t *testing.T) *deletePlanetTestSuite {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockDeleter := drivenportstest.NewMockForDeletingPlanet(ctrl)
	mockClock := drivenportstest.NewMockForFetchingTime(ctrl)

	return &deletePlanetTestSuite{
		ctrl:        ctrl,
		mockDeleter: mockDeleter,
		mockClock:   mockClock,
		usecase:     NewDeletePlanetUseCase(mockDeleter, mockClock),
	}
}

// generateApplyingDeleterMock generates a function mock for the planet deleter
// which applies the provided deleter to a known planet.
func generateApplyingDeleterMock(p *models.Planet) DeleterMock {
	return func(
		ctx context.Context, id uuid.UUID, m drivenports.PlanetDeleter,
	) error {
		err := m(p)
		return err
	}
}
