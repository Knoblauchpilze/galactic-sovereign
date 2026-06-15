package drivenadapters

import (
	"context"
	"fmt"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlayerRepository_Create(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())

	t.Run("creates a player", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		player := models.Player{
			Id:        uuid.New(),
			ApiUser:   uuid.New(),
			Universe:  universe.Id,
			Name:      fmt.Sprintf("player-%s", uuid.NewString()),
			CreatedAt: someTime,
			Planets:   []uuid.UUID{},
		}

		err := repo.Create(context.Background(), player)
		require.NoError(t, err, "Actual err: %v", err)
		assertPlayerExists(t, conn, player.Id)

		actual, err := repo.Get(context.Background(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, player, actual)
	})

	t.Run("does not create planet for a player", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		player := models.Player{
			Id:        uuid.New(),
			ApiUser:   uuid.New(),
			Universe:  universe.Id,
			Name:      fmt.Sprintf("player-%s", uuid.NewString()),
			CreatedAt: someTime,
			Planets:   []uuid.UUID{uuid.New()},
		}

		err := repo.Create(context.Background(), player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlayerExists(t, conn, player.Id)
		assertPlanetDoesNotExist(t, conn, player.Planets[0])
	})

	t.Run("returns error when player with same name already exists", func(t *testing.T) {
		player, universe := insertTestPlayerInUniverse(t, conn)

		newPlayer := models.Player{
			Id:        uuid.New(),
			ApiUser:   uuid.New(),
			Universe:  universe.Id,
			Name:      player.Name,
			CreatedAt: someTime,
		}

		err := repo.Create(context.Background(), newPlayer)

		actual, ok := db.AsDatabaseError(err)
		require.True(t, ok)
		assert.Equal(t, db.ErrUniqueConstraintViolation, actual.Code, "Actual err: %v", err)
		assertPlayerDoesNotExist(t, conn, newPlayer.Id)
	})
}

func TestIT_PlayerRepository_Get(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())

	t.Run("gets a player", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)

		actual, err := repo.Get(context.Background(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, player)
	})

	t.Run("gets a player with planets", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn, addPlayerPlanet)

		actual, err := repo.Get(context.Background(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, player)
	})

	t.Run("returns error when player does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(context.Background(), id)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func TestIT_PlayerRepository_List(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())

	p1, universe := insertTestPlayerInUniverse(t, conn)
	p2 := insertTestPlayer(t, conn, universe.Id)
	p3 := insertTestPlayer(t, conn, universe.Id, addPlayerPlanet)

	actual, err := repo.List(context.Background())
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 2)
	assert.Contains(t, actual, p1)
	assert.Contains(t, actual, p2)
	assert.Contains(t, actual, p3)
}

func TestIT_PlayerRepository_ListForApiUser(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())

	t.Run("list player for an API user", func(t *testing.T) {
		p1, universe := insertTestPlayerInUniverse(t, conn)
		insertTestPlayer(t, conn, universe.Id)

		actual, err := repo.ListForApiUser(context.Background(), p1.ApiUser)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Player{p1}, actual)
	})

	t.Run("list player with planets for an API user", func(t *testing.T) {
		p1, universe := insertTestPlayerInUniverse(t, conn, addPlayerPlanet)
		insertTestPlayer(t, conn, universe.Id)

		actual, err := repo.ListForApiUser(context.Background(), p1.ApiUser)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Player{p1}, actual)
	})
}

func TestIT_PlayerRepository_Delete(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())

	t.Run("deletes a player", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)

		err := repo.Delete(context.Background(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlayerDoesNotExist(t, conn, player.Id)
	})

	t.Run("succeeds when the player does not exist", func(t *testing.T) {
		nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

		err := repo.Delete(context.Background(), nonExistingId)
		require.NoError(t, err, "Actual err: %v", err)
	})
}

func newTestPlayerRepository(t *testing.T) (drivenports.ForManagingPlayers, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewPlayerRepository(conn), conn
}

func insertTestPlayer(
	t *testing.T,
	conn db.Connection,
	universe uuid.UUID,
	modifiers ...func(*testing.T, db.Connection, *models.Player),
) models.Player {
	t.Helper()

	player := models.Player{
		Id:        uuid.New(),
		ApiUser:   uuid.New(),
		Universe:  universe,
		Name:      fmt.Sprintf("my-player-%s", uuid.NewString()),
		CreatedAt: someTime,
		// This is intentional: the details (e.g. planets, etc.) are returned as empty
		// slices by the adapter
		Planets: []uuid.UUID{},
	}

	sqlQuery := `INSERT INTO player (id, api_user, universe, name, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		player.Id,
		player.ApiUser,
		player.Universe,
		player.Name,
		player.CreatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &player)
	}

	return player
}

func insertTestPlayerInUniverse(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Player),
) (models.Player, models.Universe) {
	universe := insertTestUniverse(t, conn)
	player := insertTestPlayer(t, conn, universe.Id, modifiers...)
	return player, universe
}

func addPlayerPlanet(t *testing.T, conn db.Connection, p *models.Player) {
	t.Helper()

	planetId := uuid.New()

	sqlQuery := `INSERT INTO planet (id, player, name, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planetId,
		p.Id,
		fmt.Sprintf("my-planet-%s", planetId.String()),
		someTime,
		someOtherTime,
		8,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Planets = append(p.Planets, planetId)
}

func assertPlayerExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT id FROM player WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, id, value)
}

func assertPlayerDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM player WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
