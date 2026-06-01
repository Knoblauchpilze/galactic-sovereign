package driven

import (
	"context"
	"fmt"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlanetRespository_Create(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := models.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
		Version:   3,
	}

	err := repo.Create(context.Background(), planet)
	require.NoError(t, err, "Actual err: %v", err)
	assertPlanetExists(t, conn, planet.Id)

	actual, err := repo.Get(context.Background(), planet.Id)
	require.NoError(t, err, "Actual err: %v", err)

	assert.Equal(t, planet, actual)
}

func TestIT_PlanetRespository_Create_ExpectErrorWhenPlanetWithSameIdAlreadyExists(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)

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
}

func TestIT_PlanetRespository_Create_WhenHomeworld_ExpectCorrectlyMarkedAsSuch(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := models.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: true,
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
	}

	err := repo.Create(context.Background(), planet)
	require.NoError(t, err, "Actual err: %v", err)
	assertPlanetExists(t, conn, planet.Id)

	actual, err := repo.Get(context.Background(), planet.Id)
	require.NoError(t, err, "Actual err: %v", err)

	assert.Equal(t, planet, actual)
}

func TestIT_PlanetRespository_Create_WhenHomeworldAlreadyExists_ExpectFailureWhenAddingANewOne(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

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
}

func TestIT_PlanetRepository_Get(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)

	actual, err := repo.Get(context.Background(), planet.Id)
	require.NoError(t, err, "Actual err: %v", err)

	assert.Equal(t, actual, planet)
}

func TestIT_PlanetRepository_Get_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.Get(context.Background(), id)

	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_PlanetRepository_List(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	p1, player1, _ := insertTestPlanetForPlayer(t, conn)
	p2 := insertTestPlanet(t, conn, player1.Id)
	p3, player2, _ := insertTestPlanetForPlayer(t, conn)
	p4 := insertTestPlanet(t, conn, player2.Id)

	actual, err := repo.List(context.Background())
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 4)
	assert.Contains(t, actual, p1)
	assert.Contains(t, actual, p2)
	assert.Contains(t, actual, p3)
	assert.Contains(t, actual, p4)
}

func TestIT_PlanetRepository_ListForPlayer(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	p1, player1, _ := insertTestPlanetForPlayer(t, conn)
	p2 := insertTestPlanet(t, conn, player1.Id)
	p3, _, _ := insertTestPlanetForPlayer(t, conn)

	actual, err := repo.ListForPlayer(context.Background(), player1.Id)
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, p1)
	assert.Contains(t, actual, p2)
	for _, planet := range actual {
		assert.NotEqual(t, planet.Id, p3.Id)
	}
}

func TestIT_PlanetRepository_Delete(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)

	err := repo.Delete(context.Background(), planet.Id)
	require.NoError(t, err, "Actual err: %v", err)

	assertPlanetDoesNotExist(t, conn, planet.Id)
}

func TestIT_PlanetRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	err := repo.Delete(context.Background(), nonExistingId)
	require.NoError(t, err, "Actual err: %v", err)
}

func TestIT_PlanetRepository_Delete_Homeworld(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn)

	err := repo.Delete(context.Background(), planet.Id)
	require.NoError(t, err, "Actual err: %v", err)

	assertPlanetDoesNotExist(t, conn, planet.Id)
	assertPlanetIsNotHomeworld(t, conn, planet.Id)
}

func TestIT_PlanetRepository_CreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	player, _ := insertTestPlayerInUniverse(t, conn)

	planet := models.Planet{
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
		Version:   4,
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
		Id:        uuid.New(),
		Player:    player.Id,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: true,
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
		Version:   6,
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

func insertTestPlanetForPlayer(t *testing.T, conn db.Connection) (models.Planet, models.Player, models.Universe) {
	t.Helper()

	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, player.Id)
	return planet, player, universe
}

func insertTestHomeworldPlanetForPlayer(t *testing.T, conn db.Connection) (models.Planet, models.Player, models.Universe) {
	t.Helper()

	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, player.Id, addPlanetHomeworld)
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

func assertPlanetIsHomeworld(t *testing.T, conn db.Connection, planet persistence.Planet) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1 AND player = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet.Id, planet.Player)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, 1, value)
}

func assertPlanetIsNotHomeworld(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, 0, value)
}
