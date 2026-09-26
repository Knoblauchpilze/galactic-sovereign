package drivenadapters

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_SolarSystemRepository_Get(t *testing.T) {
	repo, conn := newTestSolarSystemRepository(t)

	t.Run("gets a solar system", func(t *testing.T) {
		player, universe := insertTestPlayerInUniverse(t, conn)
		planet := insertTestPlanet(t, conn, player.Id)
		planet.Coordinate = models.Coordinate{
			Galaxy:      0,
			SolarSystem: 1,
			Position:    2,
		}
		upsertPlanetCoordinate(t, conn, planet)

		actual, err := repo.GetSolarSystem(
			t.Context(),
			universe.Id,
			planet.Coordinate.Galaxy,
			planet.Coordinate.SolarSystem,
		)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SolarSystem{
			Universe: universe.Id,
			Galaxy:   planet.Coordinate.Galaxy,
			Number:   planet.Coordinate.SolarSystem,
			Orbits:   universe.Topology.Orbits,
			Planets: []models.SolarSystemPlanet{
				{
					Id:         planet.Id,
					PlayerName: player.Name,
					Name:       planet.Name,
					Homeworld:  planet.Homeworld,
					Position:   planet.Coordinate.Position,
				},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("orders planets by position in a solar system", func(t *testing.T) {
		player1, universe := insertTestPlayerInUniverse(t, conn)
		require.Greater(t, universe.Topology.SolarSystems, 1)
		require.GreaterOrEqual(t, universe.Topology.Orbits, 3)

		planet1 := insertTestPlanet(t, conn, player1.Id)
		planet1.Coordinate = models.Coordinate{
			Galaxy:      0,
			SolarSystem: 1,
			Position:    2,
		}
		upsertPlanetCoordinate(t, conn, planet1)

		player2 := insertTestPlayer(t, conn, universe.Id)
		planet2 := insertTestPlanet(t, conn, player2.Id)
		planet2.Coordinate = models.Coordinate{
			Galaxy:      0,
			SolarSystem: 1,
			Position:    1,
		}
		upsertPlanetCoordinate(t, conn, planet2)

		actual, err := repo.GetSolarSystem(
			t.Context(),
			universe.Id,
			planet1.Coordinate.Galaxy,
			planet1.Coordinate.SolarSystem,
		)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SolarSystem{
			Universe: universe.Id,
			Galaxy:   planet1.Coordinate.Galaxy,
			Number:   planet1.Coordinate.SolarSystem,
			Orbits:   universe.Topology.Orbits,
			Planets: []models.SolarSystemPlanet{
				{
					Id:         planet2.Id,
					PlayerName: player2.Name,
					Name:       planet2.Name,
					Homeworld:  planet2.Homeworld,
					Position:   planet2.Coordinate.Position,
				},
				{
					Id:         planet1.Id,
					PlayerName: player1.Name,
					Name:       planet1.Name,
					Homeworld:  planet1.Homeworld,
					Position:   planet1.Coordinate.Position,
				},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("gets an empty solar system", func(t *testing.T) {
		_, universe := insertTestPlayerInUniverse(t, conn)

		actual, err := repo.GetSolarSystem(
			t.Context(),
			universe.Id,
			0,
			0,
		)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.SolarSystem{
			Universe: universe.Id,
			Galaxy:   0,
			Number:   0,
			Orbits:   universe.Topology.Orbits,
			Planets:  []models.SolarSystemPlanet{},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when universe does not exist", func(t *testing.T) {
		_, err := repo.GetSolarSystem(
			t.Context(),
			uuid.New(),
			0,
			0,
		)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when galaxy is out of bounds for universe", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		_, err := repo.GetSolarSystem(
			t.Context(),
			universe.Id,
			// The last existing galaxy is one less than the topology
			// value as it's 0 based indexing
			universe.Topology.Galaxies,
			0,
		)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})

	t.Run("returns error when solar system is out of bounds for universe", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		_, err := repo.GetSolarSystem(
			t.Context(),
			universe.Id,
			0,
			// The last existing solar system is one less than the topology
			// value as it's 0 based indexing
			universe.Topology.SolarSystems,
		)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func newTestSolarSystemRepository(t *testing.T) (*SolarSystemRepository, database.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewSolarSystemRepository(conn), conn
}
