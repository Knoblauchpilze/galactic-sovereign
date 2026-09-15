package drivenadapters

import (
	"math/rand"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
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
		modifier := generatePlanetCoordinatesModifier(t, universe)
		planet := insertTestPlanet(t, conn, player.Id, modifier)

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
					Id:        planet.Id,
					Player:    player.Id,
					Name:      planet.Name,
					Homeworld: planet.Homeworld,
					Position:  planet.Coordinate.Position,
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
}

func newTestSolarSystemRepository(t *testing.T) (*SolarSystemRepository, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewSolarSystemRepository(conn), conn
}

func generatePlanetCoordinatesModifier(
	t *testing.T,
	universe models.Universe,
) func(*testing.T, db.Connection, *models.Planet) {
	t.Helper()

	return func(t *testing.T, conn db.Connection, p *models.Planet) {
		t.Helper()

		coordinates := models.Coordinate{
			Galaxy:      rand.Intn(universe.Topology.Galaxies),
			SolarSystem: rand.Intn(universe.Topology.SolarSystems),
			Position:    rand.Intn(universe.Topology.Orbits),
		}

		sqlQuery := `UPDATE planet_coordinate SET
			galaxy = $1, solar_system = $2, position = $3
			WHERE planet = $4`
		_, err := conn.Exec(
			t.Context(),
			sqlQuery,
			coordinates.Galaxy,
			coordinates.SolarSystem,
			coordinates.Position,
			p.Id,
		)
		require.NoError(t, err, "Actual err: %v", err)

		p.Coordinate = coordinates
	}
}
