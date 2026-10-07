package domainservices

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_PlanetDeletionGuard(t *testing.T) {
	t.Run("advances planet to current time", func(t *testing.T) {
		p := models.Planet{
			Id:        uuid.New(),
			Homeworld: false,
			CreatedAt: t1,
			UpdatedAt: t1,
		}

		err := PlanetDeletionGuard(&p, t2)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t2, p.UpdatedAt)
	})

	t.Run("returns error when planet is homeworld", func(t *testing.T) {
		p := models.Planet{
			Id:        uuid.New(),
			Homeworld: true,
			CreatedAt: t1,
			UpdatedAt: t1,
		}

		err := PlanetDeletionGuard(&p, t2)

		assert.ErrorIs(t, err, domainerrors.ErrHomeworldCannotBeDeleted, "Actual err: %v", err)
	})

	t.Run("returns error when planet has building action", func(t *testing.T) {
		p := models.Planet{
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

		err := PlanetDeletionGuard(&p, t2)

		assert.ErrorIs(t, err, domainerrors.ErrBuildingActionNotCompleted, "Actual err: %v", err)
	})

	t.Run("returns error when planet has ship actions", func(t *testing.T) {
		p := models.Planet{
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

		err := PlanetDeletionGuard(&p, t2)

		assert.ErrorIs(t, err, domainerrors.ErrShipActionNotCompleted, "Actual err: %v", err)
	})
}
