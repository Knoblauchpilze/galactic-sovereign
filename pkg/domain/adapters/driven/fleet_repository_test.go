package drivenadapters

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_FleetRepository_ListForPlayer(t *testing.T) {
	t.Run("returns no fleet when player has none", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		_, player, _ := insertTestPlanetForPlayer(t, conn)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Empty(t, actual)
	})

	t.Run("returns fleet departing from a planet of the player", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		planet, player, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Fleet{fleet}, actual)
	})

	t.Run("returns fleet targeting a planet of the player", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		target, player, universe := insertTestPlanetForPlayer(t, conn)
		source, _ := insertTestPlanetInUniverse(t, conn, universe.Id)
		fleet := insertTestFleet(t, conn, source, addFleetShip)
		fleet.Destination.Target = &target.Id
		upsertFleetDestination(t, conn, fleet)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Fleet{fleet}, actual)
	})

	t.Run("returns fleet targeting another planet", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		source, player, universe := insertTestPlanetForPlayer(t, conn)
		destination, _ := insertTestPlanetInUniverse(t, conn, universe.Id)

		fleet := insertTestFleet(t, conn, source, addFleetShip)
		fleet.Destination.Target = &destination.Id
		upsertFleetDestination(t, conn, fleet)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Fleet{fleet}, actual)
	})

	t.Run("returns fleet returning to a planet of the player", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		planet, player, _ := insertTestPlanetForPlayer(t, conn)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		fleet.UpdatedAt = fleet.Schedule.ArrivalAt
		fleet.Schedule.Status = models.StatusReturning
		upsertFleetSchedule(t, conn, fleet)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Fleet{fleet}, actual)
	})

	t.Run("returns fleet returning to a planet of the player with a destination", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		source, player, universe := insertTestPlanetForPlayer(t, conn)
		destination, _ := insertTestPlanetInUniverse(t, conn, universe.Id)

		fleet := insertTestFleet(t, conn, source, addFleetShip)
		fleet.Destination.Target = &destination.Id
		upsertFleetDestination(t, conn, fleet)
		fleet.UpdatedAt = fleet.Schedule.ArrivalAt
		fleet.Schedule.Status = models.StatusReturning
		upsertFleetSchedule(t, conn, fleet)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Fleet{fleet}, actual)
	})

	t.Run("does not return fleet of another player", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		_, player, universe := insertTestPlanetForPlayer(t, conn)
		other, _ := insertTestPlanetInUniverse(t, conn, universe.Id)
		insertTestFleet(t, conn, other, addFleetShip)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Empty(t, actual)
	})

	t.Run("returns fleet only once when source and target belong to the player", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		planet, player, _ := insertTestPlanetForPlayer(t, conn)
		target := insertTestPlanet(t, conn, player.Id)
		fleet := insertTestFleet(t, conn, planet, addFleetShip)
		fleet.Destination.Target = &target.Id
		upsertFleetDestination(t, conn, fleet)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Fleet{fleet}, actual)
	})

	t.Run("returns fleets ordered by next event", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		planet, player, _ := insertTestPlanetForPlayer(t, conn)
		later := insertTestFleet(t, conn, planet, addFleetShip)
		earlier := insertTestFleet(t, conn, planet)
		earlier.Schedule.ArrivalAt = later.Schedule.ArrivalAt.Add(-30 * time.Minute)
		upsertFleetSchedule(t, conn, earlier)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Fleet{earlier, later}, actual)
	})

	t.Run("returns fleets ordered by creation time when next event are the same", func(t *testing.T) {
		repo, conn := newTestFleetRepository(t)

		planet, player, _ := insertTestPlanetForPlayer(t, conn)
		later := insertTestFleet(t, conn, planet, addFleetShip)
		earlier := insertTestFleet(t, conn, planet)
		earlier.Schedule.ArrivalAt = later.Schedule.ArrivalAt
		earlier.CreatedAt = later.CreatedAt.Add(-1 * time.Second)
		upsertFleetSchedule(t, conn, earlier)

		actual, err := repo.ListForPlayer(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Fleet{earlier, later}, actual)
	})
}

func newTestFleetRepository(t *testing.T) (*FleetRepository, database.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewFleetRepository(conn), conn
}
