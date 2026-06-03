package driven

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	crystalResourceId = uuid.MustParse("cd2ac9aa-9968-4ff5-b746-88f1f810fbb3")
	crystalMineId     = uuid.MustParse("3904d34d-9a7e-47d4-a332-091700e2c5c3")
)

func TestIT_PlanetRespository_Create(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)

	t.Run("create", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)

		planet := models.Planet{
			Id:          uuid.New(),
			Player:      player.Id,
			Name:        fmt.Sprintf("my-planet-%s", uuid.NewString()),
			Homeworld:   false,
			CreatedAt:   someTime,
			UpdatedAt:   someOtherTime,
			Version:     3,
			Resources:   []models.PlanetResource{},
			Storages:    []models.PlanetResourceStorage{},
			Productions: []models.PlanetResourceProduction{},
		}

		err := repo.Create(context.Background(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertPlanetExists(t, conn, planet.Id)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, actual)
	})

	t.Run("expect error when planet with same id already exists", func(t *testing.T) {
		planet, player, _ := insertTestPlanetForPlayer(t, conn)

		duplicatedPlanet := models.Planet{
			Id:        planet.Id,
			Player:    player.Id,
			Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
			Homeworld: false,
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
			Version:   4,
		}

		err := repo.Create(context.Background(), duplicatedPlanet)
		assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	})

	t.Run("creates homeworld and marks it as such", func(t *testing.T) {

		player, _ := insertTestPlayerInUniverse(t, conn)

		planet := models.Planet{
			Id:          uuid.New(),
			Player:      player.Id,
			Name:        fmt.Sprintf("my-planet-%s", uuid.NewString()),
			Homeworld:   true,
			CreatedAt:   someTime,
			UpdatedAt:   someOtherTime,
			Resources:   []models.PlanetResource{},
			Storages:    []models.PlanetResourceStorage{},
			Productions: []models.PlanetResourceProduction{},
		}

		err := repo.Create(context.Background(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertPlanetExists(t, conn, planet.Id)
		assertPlanetIsHomeworld(t, conn, planet.Id, planet.Player)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, actual)
	})

	t.Run("expect error when player already has a homeworld", func(t *testing.T) {
		_, player, _ := insertTestHomeworldPlanetForPlayer(t, conn)

		planet := models.Planet{
			Id:        uuid.New(),
			Player:    player.Id,
			Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
			Homeworld: true,
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
		}

		err := repo.Create(context.Background(), planet)
		assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
		assertPlanetDoesNotExist(t, conn, planet.Id)
	})
}

func TestIT_PlanetRepository_Get(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	t.Run("get planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, planet)
	})

	t.Run("get planet with resources", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, planet)
	})

	t.Run("get planet with storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, planet)
	})

	t.Run("returns error when planet does not exist", func(t *testing.T) {
		// Non-existent id
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(context.Background(), id)

		assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
	})
}

func TestIT_PlanetRepository_List(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	p1, player1, _ := insertTestPlanetForPlayer(t, conn)
	p2 := insertTestPlanet(t, conn, player1.Id)
	p3, player2, _ := insertTestPlanetForPlayer(t, conn)
	p4 := insertTestPlanet(t, conn, player2.Id)
	p5 := insertTestPlanet(t, conn, player2.Id, addPlanetResource)
	p6 := insertTestPlanet(t, conn, player2.Id, addPlanetStorage)

	actual, err := repo.List(context.Background())
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 5)
	assert.Contains(t, actual, p1)
	assert.Contains(t, actual, p2)
	assert.Contains(t, actual, p3)
	assert.Contains(t, actual, p4)
	assert.Contains(t, actual, p5)
	assert.Contains(t, actual, p6)
}

func TestIT_PlanetRepository_ListForPlayer(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	p1, player1, _ := insertTestPlanetForPlayer(t, conn)
	p2 := insertTestPlanet(t, conn, player1.Id, addPlanetResource)
	p3 := insertTestPlanet(t, conn, player1.Id, addPlanetStorage)
	p4, _, _ := insertTestPlanetForPlayer(t, conn)

	actual, err := repo.ListForPlayer(context.Background(), player1.Id)
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, p1)
	assert.Contains(t, actual, p2)
	assert.Contains(t, actual, p3)
	for _, planet := range actual {
		assert.NotEqual(t, planet.Id, p4.Id)
	}
}

func TestIT_PlanetRepository_Delete(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	t.Run("delete planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
	})

	t.Run("delete homeworld", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
	})

	t.Run("delete planet with resources", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("delete homeworld with resources", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetResource)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("delete planet with storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("delete homeworld with storages", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetStorage)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})
}

func TestIT_PlanetRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	err := repo.Delete(context.Background(), nonExistingId)
	require.NoError(t, err, "Actual err: %v", err)
}

func TestIT_PlanetRepository_CreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := models.Planet{
		Id:          uuid.New(),
		Player:      player.Id,
		Name:        fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld:   false,
		CreatedAt:   someTime,
		UpdatedAt:   someOtherTime,
		Version:     4,
		Resources:   []models.PlanetResource{},
		Storages:    []models.PlanetResourceStorage{},
		Productions: []models.PlanetResourceProduction{},
	}

	func() {
		err := repo.Create(context.Background(), planet)
		require.NoError(t, err, "Actual err: %v", err)
	}()

	func() {
		planetFromDb, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, planetFromDb)
	}()

	func() {
		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
	}()

	assertPlanetDoesNotExist(t, conn, planet.Id)
}

func TestIT_PlanetRepository_HomeWorldCreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := models.Planet{
		Id:          uuid.New(),
		Player:      player.Id,
		Name:        fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld:   true,
		CreatedAt:   someTime,
		UpdatedAt:   someOtherTime,
		Version:     6,
		Resources:   []models.PlanetResource{},
		Storages:    []models.PlanetResourceStorage{},
		Productions: []models.PlanetResourceProduction{},
	}

	func() {
		err := repo.Create(context.Background(), planet)
		require.NoError(t, err, "Actual err: %v", err)
	}()

	func() {
		planetFromDb, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, planetFromDb)
	}()

	func() {
		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)
	}()

	assertPlanetDoesNotExist(t, conn, planet.Id)
}

func TestIT_PlanetRepository_DeleteForPlayer_ExpectHomeworldToBeDeleted(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	t.Run("delete homeworld", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
	})

	t.Run("delete planet with resources", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("delete planet with storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})
}

func newTestPlanetRepository(t *testing.T) (driven.ForManagingPlanets, db.Connection) {
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
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
		Version:   7,
		// This is intentional: the details (e.g. resources, etc.) are returned as empty
		// slices by the adapter
		Resources:   []models.PlanetResource{},
		Storages:    []models.PlanetResourceStorage{},
		Productions: []models.PlanetResourceProduction{},
	}

	sqlQuery := `INSERT INTO planet (id, player, name, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planet.Id,
		planet.Player,
		planet.Name,
		planet.CreatedAt,
		planet.UpdatedAt,
		planet.Version,
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
	_, err := conn.Exec(context.Background(), sqlQuery, p.Player, p.Id)
	require.NoError(t, err, "Actual err: %v", err)

	p.Homeworld = true
}

func addPlanetResource(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	resource := models.PlanetResource{
		Resource: crystalResourceId,
		// Amount is stored with 5 decimals in the DB
		Amount:    randFloat(5),
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
	}

	sqlQuery := `INSERT INTO planet_resource (planet, resource, amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		resource.Resource,
		resource.Amount,
		resource.CreatedAt,
		resource.UpdatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Resources = append(p.Resources, resource)
}

func addPlanetStorage(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	storage := models.PlanetResourceStorage{
		Resource:  crystalResourceId,
		Storage:   6233,
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	sqlQuery := `INSERT INTO planet_resource_storage (planet, resource, storage, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		storage.Resource,
		storage.Storage,
		storage.CreatedAt,
		storage.UpdatedAt,
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
		CreatedAt:  someTime,
		UpdatedAt:  someOtherTime,
	}

	sqlQuery := `INSERT INTO planet_resource_production
		(planet, building, resource, production, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		production.Building,
		production.Resource,
		production.Production,
		production.CreatedAt,
		production.UpdatedAt,
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
		CreatedAt:  someTime,
		UpdatedAt:  someOtherTime,
	}

	sqlQuery := `INSERT INTO planet_resource_production
		(planet, building, resource, production, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		production.Building,
		production.Resource,
		production.Production,
		production.CreatedAt,
		production.UpdatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Productions = append(p.Productions, production)
}

func insertTestPlanetForPlayer(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Planet),
) (models.Planet, models.Player, models.Universe) {
	t.Helper()

	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, player.Id, modifiers...)
	return planet, player, universe
}

func insertTestHomeworldPlanetForPlayer(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Planet),
) (models.Planet, models.Player, models.Universe) {
	t.Helper()

	player, universe := insertTestPlayerInUniverse(t, conn)
	modifiers = append(modifiers, addPlanetHomeworld)
	planet := insertTestPlanet(t, conn, player.Id, modifiers...)
	return planet, player, universe
}

func assertPlanetExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT id FROM planet WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, id, value)
}

func assertPlanetDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM planet WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetIsHomeworld(t *testing.T, conn db.Connection, planet uuid.UUID, player uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1 AND player = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, player)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, 1, value)
}

func assertPlanetIsNotHomeworld(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetResourceDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(resource) FROM planet_resource WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetStorageDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(resource) FROM planet_resource_storage WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
