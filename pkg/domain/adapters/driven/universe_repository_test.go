package drivenadapters

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_UniverseRepository_Create(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)

	t.Run("creates a universe", func(t *testing.T) {
		universe := models.Universe{
			Id:   uuid.New(),
			Name: fmt.Sprintf("universe-%s", uuid.NewString()),
			Topology: models.UniverseTopology{
				Galaxies:     12,
				SolarSystems: 487,
				Orbits:       14,
			},
			CreatedAt: someTime,
		}

		err := repo.Create(t.Context(), universe)
		require.NoError(t, err, "Actual err: %v", err)
		assertUniverseExists(t, conn, universe.Id)

		actual, err := repo.Get(t.Context(), universe.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := universe
		expected.OccupancyMap = models.OccupancyMap{
			Topology:  universe.Topology,
			UsedSlots: make(map[models.Coordinate]struct{}),
		}
		assertEqualIgnoringFields(t, actual, expected, "Buildings", "Resources")
	})

	t.Run("returns error when universe with same name already exists", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		newUniverse := models.Universe{
			Id:        uuid.New(),
			Name:      universe.Name,
			CreatedAt: someTime,
		}

		err := repo.Create(t.Context(), newUniverse)

		assert.Equal(t, domainerrors.ErrNameAlreadyTaken, err, "Actual err: %+v", err)
		assertUniverseDoesNotExist(t, conn, newUniverse.Id)
	})

}

func TestIT_UniverseRepository_Get(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)

	t.Run("gets a universe", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		actual, err := repo.Get(t.Context(), universe.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertEqualIgnoringFields(t, actual, universe, "Buildings", "Resources")
	})

	t.Run("gets a universe with resources", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)
		resource := insertTestResource(t, conn)

		actual, err := repo.Get(t.Context(), universe.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Contains(t, actual.Resources, resource)
	})

	t.Run("gets a universe with buildings", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)
		building := insertTestBuilding(t, conn)

		actual, err := repo.Get(t.Context(), universe.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Contains(t, actual.Buildings, building)
	})

	t.Run("gets a universe with occupied slots", func(t *testing.T) {
		u1 := insertTestUniverse(t, conn)
		p1 := insertTestPlayer(t, conn, u1.Id)
		planet1 := insertTestPlanet(t, conn, p1.Id)

		u2 := insertTestUniverse(t, conn)
		p2 := insertTestPlayer(t, conn, u2.Id)
		planet2 := insertTestPlanet(t, conn, p2.Id)
		require.NotEqual(t, planet1.Coordinate, planet2.Coordinate)

		actual, err := repo.Get(t.Context(), u1.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := u1
		expected.OccupancyMap = models.OccupancyMap{
			Topology: u1.Topology,
			UsedSlots: map[models.Coordinate]struct{}{
				planet1.Coordinate: {},
			},
		}

		assertEqualIgnoringFields(t, actual, expected, "Buildings", "Resources")
	})

	t.Run("returns error when universe does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func TestIT_UniverseRepository_List(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)
	u1 := insertTestUniverse(t, conn)
	u2 := insertTestUniverse(t, conn)
	resource := insertTestResource(t, conn)
	building := insertTestBuilding(t, conn)

	actual, err := repo.List(t.Context())
	require.NoError(t, err, "Actual err: %v", err)

	// The additional resources are the universes from the seed data
	assertContainsIgnoringFields(t, actual, u1, "Buildings", "Resources")
	assertContainsIgnoringFields(t, actual, u2, "Buildings", "Resources")

	for _, u := range actual {
		assert.Contains(t, u.Resources, resource)
		assert.Contains(t, u.Buildings, building)
	}
}

func TestIT_UniverseRepository_Delete(t *testing.T) {
	repo, conn := newTestUniverseRepository(t)

	t.Run("deletes universe", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		err := repo.Delete(t.Context(), universe.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertUniverseDoesNotExist(t, conn, universe.Id)
	})

	t.Run("succeeds when the universe does not exist", func(t *testing.T) {
		nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

		err := repo.Delete(t.Context(), nonExistingId)
		require.NoError(t, err, "Actual err: %v", err)
	})

	t.Run("returns error when player still exists in universe", func(t *testing.T) {
		_, universe := insertTestPlayerInUniverse(t, conn)

		err := repo.Delete(t.Context(), universe.Id)

		assert.ErrorIs(t, err, domainerrors.ErrUniverseIsNotEmpty)
	})
}

func newTestUniverseRepository(t *testing.T) (*UniverseRepository, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewUniverseRepository(conn), conn
}

func insertTestUniverse(t *testing.T, conn db.Connection) models.Universe {
	t.Helper()

	topology := models.UniverseTopology{
		Galaxies:     rand.Intn(15),
		SolarSystems: rand.Intn(800),
		Orbits:       rand.Intn(15),
	}

	universe := models.Universe{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-universe-%s", uuid.NewString()),
		Topology:  topology,
		CreatedAt: someTime,
		OccupancyMap: models.OccupancyMap{
			Topology:  topology,
			UsedSlots: make(map[models.Coordinate]struct{}),
		},
	}

	sqlQuery := `INSERT INTO universe (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		universe.Id,
		universe.Name,
		universe.CreatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	sqlQuery = `INSERT INTO universe_topology (universe, galaxies, solar_systems, orbits)
		VALUES ($1, $2, $3, $4)`
	_, err = conn.Exec(
		t.Context(),
		sqlQuery,
		universe.Id,
		universe.Topology.Galaxies,
		universe.Topology.SolarSystems,
		universe.Topology.Orbits,
	)
	require.NoError(t, err, "Actual err: %v", err)

	return universe
}

func insertTestResource(t *testing.T, conn db.Connection) models.Resource {
	t.Helper()

	resource := models.Resource{
		Id:                    uuid.New(),
		Name:                  fmt.Sprintf("my-resource-%s", uuid.NewString()),
		StartAmount:           456,
		StartProduction:       321,
		StartStorage:          778899,
		BuildTimeHoursPerUnit: 1.2564,
		CreatedAt:             someTime,
	}

	sqlQuery := `INSERT INTO resource (id, name, start_amount, start_production, start_storage, build_time_hours_per_unit, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		resource.Id,
		resource.Name,
		resource.StartAmount,
		resource.StartProduction,
		resource.StartStorage,
		resource.BuildTimeHoursPerUnit,
		resource.CreatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	return resource
}

func assertUniverseExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT id FROM universe WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, id, value)
}

func assertUniverseDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM universe WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
