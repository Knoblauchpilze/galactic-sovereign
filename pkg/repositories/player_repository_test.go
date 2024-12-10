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

func TestIT_PlayerRepository_Create(t *testing.T) {
	repo, conn, tx := newTestPlayerRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	universe := insertTestUniverse(t, conn)

	player := persistence.Player{
		Id:        uuid.New(),
		ApiUser:   uuid.New(),
		Universe:  universe.Id,
		Name:      fmt.Sprintf("player-%s", uuid.NewString()),
		CreatedAt: time.Now(),
	}

	actual, err := repo.Create(context.Background(), tx, player)
	assert.Nil(t, err)
	tx.Close(context.Background())

	assert.True(t, eassert.EqualsIgnoringFields(actual, player, "UpdatedAt"))
	assert.True(t, eassert.AreTimeCloserThan(actual.UpdatedAt, actual.CreatedAt, 1*time.Second))
	assertPlayerExists(t, conn, player.Id)
}

func TestIT_PlayerRepository_Create_WhenDuplicateName_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestPlayerRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	player, universe := insertTestPlayerInUniverse(t, conn)

	newPlayer := persistence.Player{
		Id:        uuid.New(),
		ApiUser:   uuid.New(),
		Universe:  universe.Id,
		Name:      player.Name,
		CreatedAt: time.Now(),
	}

	_, err := repo.Create(context.Background(), tx, newPlayer)

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertPlayerDoesNotExist(t, conn, newPlayer.Id)
}

func TestIT_PlayerRepository_Get(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())
	player, _ := insertTestPlayerInUniverse(t, conn)

	actual, err := repo.Get(context.Background(), player.Id)
	assert.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, player))
}

func TestIT_PlayerRepository_Get_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.Get(context.Background(), id)
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_PlayerRepository_List(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())
	p1, universe := insertTestPlayerInUniverse(t, conn)
	p2 := insertTestPlayer(t, conn, universe.Id)

	actual, err := repo.List(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 2)
	assert.True(t, eassert.ContainsIgnoringFields(actual, p1))
	assert.True(t, eassert.ContainsIgnoringFields(actual, p2))
}

func TestIT_PlayerRepository_ListForApiUser(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	defer conn.Close(context.Background())
	p1, universe := insertTestPlayerInUniverse(t, conn)
	insertTestPlayer(t, conn, universe.Id)

	actual, err := repo.ListForApiUser(context.Background(), p1.ApiUser)

	assert.Nil(t, err)
	assert.Equal(t, len(actual), 1)
	assert.True(t, eassert.ContainsIgnoringFields(actual, p1))
}

func TestIT_PlayerRepository_Delete(t *testing.T) {
	repo, conn, tx := newTestPlayerRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	player, _ := insertTestPlayerInUniverse(t, conn)

	err := repo.Delete(context.Background(), tx, player.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertPlayerDoesNotExist(t, conn, player.Id)
}

func TestIT_PlayerRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn, tx := newTestPlayerRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	err := repo.Delete(context.Background(), tx, nonExistingId)
	tx.Close(context.Background())

	assert.Nil(t, err)
}

func newTestPlayerRepository(t *testing.T) (PlayerRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewPlayerRepository(conn), conn
}

func newTestPlayerRepositoryAndTransaction(t *testing.T) (PlayerRepository, db.Connection, db.Transaction) {
	repo, conn := newTestPlayerRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestPlayer(t *testing.T, conn db.Connection, universe uuid.UUID) persistence.Player {
	someTime := time.Date(2024, 11, 29, 17, 56, 02, 0, time.UTC)

	player := persistence.Player{
		Id:        uuid.New(),
		ApiUser:   uuid.New(),
		Universe:  universe,
		Name:      fmt.Sprintf("my-player-%s", uuid.NewString()),
		CreatedAt: someTime,
	}

	sqlQuery := `INSERT INTO player (id, api_user, universe, name, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING updated_at`
	updatedAt, err := db.QueryOne[time.Time](
		context.Background(),
		conn,
		sqlQuery,
		player.Id,
		player.ApiUser,
		player.Universe,
		player.Name,
		player.CreatedAt,
	)
	require.Nil(t, err)

	player.UpdatedAt = updatedAt

	return player
}

func insertTestPlayerInUniverse(t *testing.T, conn db.Connection) (persistence.Player, persistence.Universe) {
	universe := insertTestUniverse(t, conn)
	player := insertTestPlayer(t, conn, universe.Id)
	return player, universe
}

func assertPlayerExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	sqlQuery := `SELECT id FROM player WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertPlayerDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	sqlQuery := `SELECT COUNT(id) FROM player WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.Nil(t, err)
	require.Zero(t, value)
}
