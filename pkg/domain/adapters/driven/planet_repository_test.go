package drivenadapters

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	crystalResourceId = uuid.MustParse("cd2ac9aa-9968-4ff5-b746-88f1f810fbb3")
	crystalMineId     = uuid.MustParse("3904d34d-9a7e-47d4-a332-091700e2c5c3")
	metalStorageId    = uuid.MustParse("22b4c0c3-c8e5-4493-89fc-522fdbb0beee")
)

func TestIT_PlanetRepository_ListForPlayer(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	p1, _, _ := insertTestPlanetForPlayer(t, conn)
	p2, player1, _ := insertTestPlanetForPlayer(t, conn)
	p3 := insertTestPlanet(t, conn, player1.Id, addPlanetResource)
	p4 := insertTestPlanet(t, conn, player1.Id, addPlanetStorage)
	p5 := insertTestPlanet(t, conn, player1.Id, addPlanetProduction)
	p6 := insertTestPlanet(t, conn, player1.Id, addPlanetProductionForBuilding)
	p7 := insertTestPlanet(t, conn, player1.Id, addPlanetBuilding)
	p8 := insertTestPlanet(t, conn, player1.Id, addPlanetBuildingAction)

	actual, err := repo.ListForPlayer(t.Context(), player1.Id)
	require.NoError(t, err, "Actual err: %v", err)

	// As all planets are registered with the same creation data their order is
	// not deterministic.
	assert.Contains(t, actual, p2.Id)
	assert.Contains(t, actual, p3.Id)
	assert.Contains(t, actual, p4.Id)
	assert.Contains(t, actual, p5.Id)
	assert.Contains(t, actual, p6.Id)
	assert.Contains(t, actual, p7.Id)
	assert.Contains(t, actual, p8.Id)
	assert.NotContains(t, actual, p1)
}

func TestIT_PlanetRepository_Delete(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)

	t.Run("deletes planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
	})

	t.Run("deletes planet with resources", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with resources", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetResource)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with storages", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetStorage)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with productions", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetProduction)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with productions for building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with productions for building", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with buildings", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetBuildingDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with buildings", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetBuilding)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetBuildingDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with building actions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuildingAction)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertBuildingActionDoesNotExist(t, conn, planet.BuildingAction.Id)
	})

	t.Run("deletes homeworld with building actions", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetBuildingAction)

		err := repo.Delete(t.Context(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertBuildingActionDoesNotExist(t, conn, planet.BuildingAction.Id)
	})

	t.Run("succeeds when the planet does not exist", func(t *testing.T) {
		nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

		err := repo.Delete(t.Context(), nonExistingId)
		require.NoError(t, err, "Actual err: %v", err)
	})
}

func newTestPlanetRepository(t *testing.T) (*PlanetRepository, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewPlanetRepository(conn), conn
}

func insertTestPlanet(
	t *testing.T,
	conn db.Connection,
	player uuid.UUID,
	modifiers ...func(*testing.T, db.Connection, *models.Planet),
) models.Planet {
	t.Helper()

	planet := models.Planet{
		Id:        uuid.New(),
		Player:    player,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		Coordinate: models.Coordinate{
			Galaxy:      1 + rand.Intn(15),
			SolarSystem: 1 + rand.Intn(421),
			Position:    1 + rand.Intn(36),
		},
		Fields:    rand.Intn(211),
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
		Version:   7,
		// This is intentional: the details (e.g. resources, etc.) are returned as empty
		// slices by the adapter
		Resources:   []models.PlanetResource{},
		Storages:    []models.PlanetResourceStorage{},
		Productions: []models.PlanetResourceProduction{},
		Buildings:   []models.PlanetBuilding{},
	}

	sqlQuery := `INSERT INTO planet (id, player, name, fields, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		planet.Id,
		planet.Player,
		planet.Name,
		planet.Fields,
		planet.CreatedAt,
		planet.UpdatedAt,
		planet.Version,
	)
	require.NoError(t, err, "Actual err: %v", err)

	sqlQuery = `INSERT INTO planet_coordinate (planet, universe, galaxy, solar_system, position)
		SELECT $1, universe, $2, $3, $4 FROM player WHERE id = $5`
	_, err = conn.Exec(
		t.Context(),
		sqlQuery,
		planet.Id,
		planet.Coordinate.Galaxy,
		planet.Coordinate.SolarSystem,
		planet.Coordinate.Position,
		player,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &planet)
	}

	return planet
}

func addPlanetHomeworld(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	sqlQuery := `INSERT INTO homeworld (player, planet) VALUES ($1, $2)`
	_, err := conn.Exec(t.Context(), sqlQuery, p.Player, p.Id)
	require.NoError(t, err, "Actual err: %v", err)

	p.Homeworld = true
}

func addPlanetResource(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	resource := models.PlanetResource{
		Resource: crystalResourceId,
		// Amount is stored with 5 decimals in the DB
		Amount: randFloat(t, 4000, 10000, 5),
	}

	sqlQuery := `INSERT INTO planet_resource (planet, resource, amount)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		p.Id,
		resource.Resource,
		resource.Amount,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Resources = append(p.Resources, resource)
}

func addPlanetStorage(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	storage := models.PlanetResourceStorage{
		Resource: crystalResourceId,
		Storage:  6233,
	}

	sqlQuery := `INSERT INTO planet_resource_storage (planet, resource, storage)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		p.Id,
		storage.Resource,
		storage.Storage,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Storages = append(p.Storages, storage)
}

func addPlanetProductionForBuilding(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	production := models.PlanetResourceProduction{
		Building:   &crystalMineId,
		Resource:   metalResourceId,
		Production: rand.Intn(784152),
	}

	sqlQuery := `INSERT INTO planet_resource_production
		(planet, building, resource, production)
		VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		p.Id,
		production.Building,
		production.Resource,
		production.Production,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Productions = append(p.Productions, production)
}

func addPlanetProduction(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	production := models.PlanetResourceProduction{
		Building:   nil,
		Resource:   metalResourceId,
		Production: rand.Intn(6589),
	}

	sqlQuery := `INSERT INTO planet_resource_production
		(planet, building, resource, production)
		VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		p.Id,
		production.Building,
		production.Resource,
		production.Production,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Productions = append(p.Productions, production)
}

func addPlanetBuilding(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	building := models.PlanetBuilding{
		Building: metalStorageId,
		Level:    0,
	}

	sqlQuery := `INSERT INTO planet_building (planet, building, level)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		p.Id,
		building.Building,
		building.Level,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Buildings = append(p.Buildings, building)
}

func addPlanetBuildingAction(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	action := models.BuildingAction{
		Id:           uuid.New(),
		Building:     metalStorageId,
		DesiredLevel: 1,
		CreatedAt:    someTime,
		CompletedAt:  someOtherTime,
		// This is intentional: the details (e.g. costs, etc.) are returned as empty
		// slices by the adapter
		Costs:       []models.BuildingActionCost{},
		Storages:    []models.BuildingActionResourceStorage{},
		Productions: []models.BuildingActionResourceProduction{},
	}

	sqlQuery := `INSERT INTO building_action
		(id, planet, building, desired_level, created_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		action.Id,
		p.Id,
		action.Building,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.BuildingAction = &action
}

// insertTestPlanetForPlayer creates a test planet. The returned planet belongs
// to a fresh player registered in a fresh universe.
// The player also has a homeworld: this is a necessary precondition to make the
// tests realistic.
func insertTestPlanetForPlayer(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Planet),
) (models.Planet, models.Player, models.Universe) {
	t.Helper()

	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, player.Id, modifiers...)

	player.Planets = append(player.Planets, planet.Id)

	return planet, player, universe
}

// insertTestHomeworldPlanetForPlayer creates a test homeworld. The returned planet
// belongs to a fresh player registered in a fresh universe. It is the homeworld of
// the player and the player has no other planet.
// The modifiers are applied to the homeworld.
func insertTestHomeworldPlanetForPlayer(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Planet),
) (models.Planet, models.Player, models.Universe) {
	t.Helper()

	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := loadPlanetFromDb(t, conn, player.Homeworld)

	for _, modifier := range modifiers {
		modifier(t, conn, &planet)
	}

	return planet, player, universe
}

func loadPlanetFromDb(t *testing.T, conn db.Connection, id uuid.UUID) models.Planet {
	t.Helper()

	tx, err := conn.BeginTx(t.Context())
	require.NoError(t, err, "Actual err: %v", err)
	defer tx.Close(t.Context())

	planet, err := loadPlanetAndDetails(t.Context(), tx, id)
	require.NoError(t, err, "Actual err: %v", err)

	return planet
}

func assertPlanetExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT id FROM planet WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, id, value)
}

func assertPlanetDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM planet WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetIsHomeworld(t *testing.T, conn db.Connection, planet uuid.UUID, player uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1 AND player = $2`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet, player)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, 1, value)
}

func assertPlanetIsNotHomeworld(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetResourceDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(resource) FROM planet_resource WHERE planet = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetStorageDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(resource) FROM planet_resource_storage WHERE planet = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetProductionDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(resource) FROM planet_resource_production WHERE planet = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetBuildingDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(building) FROM planet_building WHERE planet = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
