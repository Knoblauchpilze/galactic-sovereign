package domainservices

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_DetermineFleetSpeed(t *testing.T) {
	t.Run("returns speed of single ship when fleet only has one", func(t *testing.T) {
		order := models.FleetOrder{
			Ships: []models.FleetShip{
				{Ship: lightFighterId, Count: 1},
			},
		}
		catalog := []models.Ship{
			{Id: lightFighterId, BaseSpeed: 1000},
		}

		actual, err := DetermineFleetSpeed(order, catalog)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, 1000, actual)
	})

	t.Run("returns ship with lowest speed", func(t *testing.T) {
		order := models.FleetOrder{
			Ships: []models.FleetShip{
				{Ship: lightFighterId, Count: 1},
				{Ship: smallCargoId, Count: 1},
			},
		}
		catalog := []models.Ship{
			{Id: smallCargoId, BaseSpeed: 1000},
			{Id: lightFighterId, BaseSpeed: 700},
		}

		actual, err := DetermineFleetSpeed(order, catalog)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, 700, actual)
	})

	t.Run("returns error when fleet has no ships", func(t *testing.T) {
		order := models.FleetOrder{
			Ships: []models.FleetShip{},
		}
		catalog := []models.Ship{
			{Id: smallCargoId, BaseSpeed: 1000},
		}

		_, err := DetermineFleetSpeed(order, catalog)

		assert.ErrorIs(t, err, domainerrors.ErrNoShipInFleet, "Actual err: %v", err)
	})

	t.Run("returns error when catalog does not contain all ships", func(t *testing.T) {
		order := models.FleetOrder{
			Ships: []models.FleetShip{
				{Ship: lightFighterId, Count: 1},
			},
		}
		catalog := []models.Ship{
			{Id: smallCargoId, BaseSpeed: 1000},
		}

		_, err := DetermineFleetSpeed(order, catalog)

		assert.ErrorIs(t, err, domainerrors.ErrShipNotFound, "Actual err: %v", err)
	})
}
