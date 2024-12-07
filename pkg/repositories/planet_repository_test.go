package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/db/pgx"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	eassert "github.com/KnoblauchPilze/easy-assert/assert"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlanetRespository_Create(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 14, 48, 47, 0, time.UTC),
	}

	actual, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, planet, "UpdatedAt"))
	assert.Equal(t, actual.CreatedAt, actual.UpdatedAt)
	assertPlanetExists(t, conn, planet.Id)
}

func TestIT_PlanetRespository_Create_WhenHomeworld_ExpectCorrectlyMarkedAsSuch(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: true,
		CreatedAt: time.Date(2024, 11, 30, 14, 51, 57, 0, time.UTC),
	}

	actual, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, planet, "UpdatedAt"))
	assert.Equal(t, actual.CreatedAt, actual.UpdatedAt)
	assertPlanetIsHomeworld(t, conn, planet)
}

func TestIT_PlanetRespository_Create_WhenHomeworldAlreadyExists_ExpectFailureWhenAddingANewOne(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	_, player, _ := insertTestHomeworldPlanetForPlayer(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: true,
		CreatedAt: time.Date(2024, 11, 30, 14, 57, 49, 0, time.UTC),
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)

	assertPlanetDoesNotExist(t, conn, planet.Id)
}

func TestIT_PlanetRepository_Create_RegistersBuildingsForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)
	b1 := insertTestBuilding(t, conn)
	b2 := insertTestBuilding(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 15, 12, 53, 0, time.UTC),
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assertPlanetBuildingExists(t, conn, planet.Id, b1.Id)
	assertPlanetBuildingExists(t, conn, planet.Id, b2.Id)
}

func TestIT_PlanetRepository_Create_RegistersBuildingWithLevel0(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)
	building := insertTestBuilding(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 15, 17, 53, 0, time.UTC),
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assertPlanetBuildingLevel(t, conn, planet.Id, building.Id, 0)
}

func TestIT_PlanetRepository_Create_RegistersResourcesForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)
	r1 := insertTestResource(t, conn)
	r2 := insertTestResource(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 15, 21, 38, 0, time.UTC),
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assertPlanetResourceExists(t, conn, planet.Id, r1.Id)
	assertPlanetResourceExists(t, conn, planet.Id, r2.Id)
}

func TestIT_PlanetRepository_Create_RegistersResourceWithStartAmount(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)
	resource := insertTestResource(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 15, 22, 36, 0, time.UTC),
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assertPlanetResourceAmount(t, conn, planet.Id, resource.Id, float64(resource.StartAmount))
}

func TestIT_PlanetRepository_Create_RegistersResourceProductionsForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)
	r1 := insertTestResource(t, conn)
	r2 := insertTestResource(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 15, 24, 47, 0, time.UTC),
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assertPlanetResourceProductionExists(t, conn, planet.Id, r1.Id)
	assertPlanetResourceProductionExists(t, conn, planet.Id, r2.Id)
}

func TestIT_PlanetRepository_Create_RegistersResourceWithStartProductionAndNoBuilding(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)
	resource := insertTestResource(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 15, 26, 17, 0, time.UTC),
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assertPlanetResourceProduction(t, conn, planet.Id, resource.Id, resource.StartProduction)
}

func TestIT_PlanetRepository_Create_RegistersResourceStoragesForPlanet(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)
	r1 := insertTestResource(t, conn)
	r2 := insertTestResource(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		CreatedAt: time.Date(2024, 11, 30, 15, 30, 18, 0, time.UTC),
		Homeworld: false,
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assertPlanetResourceStorageExists(t, conn, planet.Id, r1.Id)
	assertPlanetResourceStorageExists(t, conn, planet.Id, r2.Id)
}

func TestIT_PlanetRepository_Create_RegistersResourceWithStartStorageAndNoBuilding(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)
	resource := insertTestResource(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 15, 33, 03, 0, time.UTC),
	}

	_, err := repo.Create(context.Background(), tx, planet)
	tx.Close(context.Background())
	require.Nil(t, err)

	assertPlanetResourceStorage(t, conn, planet.Id, resource.Id, resource.StartStorage)
}

func TestIT_PlanetRepository_Get(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)

	actual, err := repo.Get(context.Background(), tx, planet.Id)
	assert.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, planet))
}

func TestIT_PlanetRepository_Get_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.Get(context.Background(), tx, id)
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_PlanetRepository_List(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	p1, player, _ := insertTestPlanetForPlayer(t, conn)
	p2 := insertTestPlanet(t, conn, player.Id, false)

	actual, err := repo.List(context.Background(), tx)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, p1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, p2))
}

func TestIT_PlanetRepository_ListForPlayer(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	p1, player1, _ := insertTestPlanetForPlayer(t, conn)
	p2 := insertTestPlanet(t, conn, player1.Id, false)
	p3, _, _ := insertTestPlanetForPlayer(t, conn)

	actual, err := repo.ListForPlayer(context.Background(), tx, player1.Id)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, p1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, p2))
	for _, planet := range actual {
		assert.NotEqual(t, planet.Id, p3.Id)
	}
}

func TestIT_PlanetRepository_Delete(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)

	err := repo.Delete(context.Background(), tx, planet.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertPlanetDoesNotExist(t, conn, planet.Id)
}

func TestIT_PlanetRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	err := repo.Delete(context.Background(), tx, nonExistingId)
	tx.Close(context.Background())

	assert.Nil(t, err)
}

func TestIT_PlanetRepository_Delete_Homeworld(t *testing.T) {
	repo, conn, tx := newTestPlanetRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn)

	err := repo.Delete(context.Background(), tx, planet.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertPlanetDoesNotExist(t, conn, planet.Id)
	assertPlanetIsNotHomeworld(t, conn, planet.Id)
}

func TestIT_PlanetRepository_CreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: time.Date(2024, 11, 30, 15, 50, 30, 0, time.UTC),
	}

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		_, err = repo.Create(context.Background(), tx, planet)
		require.Nil(t, err)
	}()

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		planetFromDb, err := repo.Get(context.Background(), tx, planet.Id)
		require.Nil(t, err)

		assert.True(t, eassert.EqualsIgnoringFields(planet, planetFromDb, "UpdatedAt"))
		assert.True(t, eassert.AreTimeCloserThan(planet.CreatedAt, planetFromDb.UpdatedAt, 1*time.Second))
	}()

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		err = repo.Delete(context.Background(), tx, planet.Id)
		require.Nil(t, err)
	}()

	assertPlanetDoesNotExist(t, conn, planet.Id)
}

func TestIT_PlanetRepository_HomeWorldCreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: true,
		CreatedAt: time.Date(2024, 11, 30, 15, 50, 55, 0, time.UTC),
	}

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		_, err = repo.Create(context.Background(), tx, planet)
		require.Nil(t, err)
	}()

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		planetFromDb, err := repo.Get(context.Background(), tx, planet.Id)
		require.Nil(t, err)

		assert.True(t, eassert.EqualsIgnoringFields(planet, planetFromDb, "UpdatedAt"))
		assert.True(t, eassert.AreTimeCloserThan(planet.CreatedAt, planetFromDb.UpdatedAt, 1*time.Second))
	}()

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		err = repo.Delete(context.Background(), tx, planet.Id)
		require.Nil(t, err)
	}()

	assertPlanetDoesNotExist(t, conn, planet.Id)
}

func newTestPlanetRepository(t *testing.T) (PlanetRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewPlanetRepository(conn), conn
}

func newTestPlanetRepositoryAndTransaction(t *testing.T) (PlanetRepository, db.Connection, db.Transaction) {
	repo, conn := newTestPlanetRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestPlanet(t *testing.T, conn db.Connection, player uuid.UUID, homeworld bool) persistence.Planet {
	someTime := time.Date(2024, 11, 30, 11, 31, 58, 0, time.UTC)

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    player,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: homeworld,
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	sqlQuery := `INSERT INTO planet (id, player, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planet.Id,
		planet.Player,
		planet.Name,
		planet.CreatedAt,
		planet.UpdatedAt,
	)
	require.Nil(t, err)

	if homeworld {
		sqlQuery := `INSERT INTO homeworld (player, planet) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), sqlQuery, planet.Player, planet.Id)
		require.Nil(t, err)
	}

	return planet
}

func insertTestPlanetForPlayer(t *testing.T, conn db.Connection) (persistence.Planet, persistence.Player, persistence.Universe) {
	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, player.Id, false)
	return planet, player, universe
}

func insertTestHomeworldPlanetForPlayer(t *testing.T, conn db.Connection) (persistence.Planet, persistence.Player, persistence.Universe) {
	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, player.Id, true)
	return planet, player, universe
}

func assertPlanetExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	sqlQuery := `SELECT id FROM planet WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertPlanetDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	sqlQuery := `SELECT COUNT(id) FROM planet WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertPlanetIsHomeworld(t *testing.T, conn db.Connection, planet persistence.Planet) {
	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1 AND player = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet.Id, planet.Player)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertPlanetIsNotHomeworld(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.Nil(t, err)
	require.Equal(t, 0, value)
}
