package repositories

import (
	"context"
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

func TestIT_BuildingActionRepository_Create(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building := insertTestBuilding(t, conn)

	action := persistence.BuildingAction{
		Id:           uuid.New(),
		Planet:       planet.Id,
		Building:     building.Id,
		CurrentLevel: 2,
		DesiredLevel: 3,
		CreatedAt:    time.Now(),
		CompletedAt:  time.Now().Add(1 * time.Hour),
	}

	actual, err := repo.Create(context.Background(), tx, action)
	tx.Close(context.Background())
	assert.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, action))
	assertBuildingActionExists(t, conn, action.Id)
}

func TestIT_BuildingActionRepository_Create_WhenDuplicatePlanet_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	_, planet := insertTestBuildingAction(t, conn)
	building := insertTestBuilding(t, conn)

	newAction := persistence.BuildingAction{
		Id:           uuid.New(),
		Planet:       planet.Id,
		Building:     building.Id,
		CurrentLevel: 4,
		DesiredLevel: 5,
		CreatedAt:    time.Now(),
		CompletedAt:  time.Now().Add(1 * time.Minute),
	}

	_, err := repo.Create(context.Background(), tx, newAction)

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertBuildingActionDoesNotExist(t, conn, newAction.Id)
}

func TestIT_BuildingActionRepository_Get(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	action, _ := insertTestBuildingAction(t, conn)

	actual, err := repo.Get(context.Background(), tx, action.Id)
	assert.Nil(t, err)

	assert.True(t, eassert.EqualsIgnoringFields(actual, action))
}

func TestIT_BuildingActionRepository_Get_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.Get(context.Background(), tx, id)
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_BuildingActionRepository_ListForPlanet(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	action1, planet1 := insertTestBuildingAction(t, conn)
	insertTestBuildingAction(t, conn)

	actual, err := repo.ListForPlanet(context.Background(), tx, planet1.Id)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 1)
	assert.True(t, eassert.ContainsIgnoringFields(actual, action1))
}

func TestIT_BuildingActionRepository_ListBeforeCompletionTime_WhenCompletionTimeBeforeQueriedTime_ExpectActionReturned(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	action, planet := insertTestBuildingAction(t, conn)

	until := action.CompletedAt.Add(1 * time.Second)

	actual, err := repo.ListBeforeCompletionTime(context.Background(), tx, planet.Id, until)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(actual), 1)
	assert.True(t, eassert.ContainsIgnoringFields(actual, action))
}

func TestIT_BuildingActionRepository_ListBeforeCompletionTime_WhenCompletionTimeInTheFuture_ExpectActionNotReturned(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	defer tx.Close(context.Background())
	action, planet := insertTestBuildingAction(t, conn)

	until := action.CompletedAt.Add(-1 * time.Second)

	actual, err := repo.ListBeforeCompletionTime(context.Background(), tx, planet.Id, until)

	assert.Nil(t, err)
	assert.Equal(t, len(actual), 0)
}

func TestIT_BuildingActionRepository_Delete(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	action, _ := insertTestBuildingAction(t, conn)

	err := repo.Delete(context.Background(), tx, action.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
}

func TestIT_BuildingActionRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	err := repo.Delete(context.Background(), tx, nonExistingId)
	tx.Close(context.Background())

	assert.Nil(t, err)
}

func TestIT_BuildingActionRepository_DeleteForPlanet(t *testing.T) {
	repo, conn, tx := newTestBuildingActionRepositoryAndTransaction(t)
	defer conn.Close(context.Background())
	action, planet := insertTestBuildingAction(t, conn)

	err := repo.DeleteForPlanet(context.Background(), tx, planet.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
}

func TestIT_BuildingActionRepository_CreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestBuildingActionRepository(t)
	defer conn.Close(context.Background())

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building := insertTestBuilding(t, conn)

	action := persistence.BuildingAction{
		Id:           uuid.New(),
		Planet:       planet.Id,
		Building:     building.Id,
		CurrentLevel: 26,
		DesiredLevel: 27,
		CreatedAt:    time.Date(2024, 12, 7, 20, 26, 47, 0, time.UTC),
		CompletedAt:  time.Date(2024, 12, 7, 21, 26, 47, 0, time.UTC),
	}

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		_, err = repo.Create(context.Background(), tx, action)
		require.Nil(t, err)
	}()

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		actionFromDb, err := repo.Get(context.Background(), tx, action.Id)
		require.Nil(t, err)

		assert.True(t, eassert.EqualsIgnoringFields(action, actionFromDb))
	}()

	func() {
		tx, err := conn.BeginTx(context.Background())
		require.Nil(t, err)

		defer tx.Close(context.Background())

		err = repo.Delete(context.Background(), tx, action.Id)
		require.Nil(t, err)
	}()

	assertBuildingActionDoesNotExist(t, conn, action.Id)
}

func newTestBuildingActionRepository(t *testing.T) (BuildingActionRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewBuildingActionRepository(), conn
}

func newTestBuildingActionRepositoryAndTransaction(t *testing.T) (BuildingActionRepository, db.Connection, db.Transaction) {
	repo, conn := newTestBuildingActionRepository(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return repo, conn, tx
}

func insertTestBuildingAction(t *testing.T, conn db.Connection) (persistence.BuildingAction, persistence.Planet) {
	someTime := time.Date(2024, 12, 7, 20, 8, 48, 0, time.UTC)

	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building := insertTestBuilding(t, conn)

	action := persistence.BuildingAction{
		Id:           uuid.New(),
		Planet:       planet.Id,
		Building:     building.Id,
		CurrentLevel: 4,
		DesiredLevel: 5,
		CreatedAt:    someTime,
		CompletedAt:  someTime.Add(1*time.Hour + 2*time.Minute),
	}

	sqlQuery := `INSERT INTO building_action (id, planet, building, current_level, desired_level, created_at, completed_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		action.Id,
		action.Planet,
		action.Building,
		action.CurrentLevel,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
	)
	require.Nil(t, err)

	return action, planet
}

func assertBuildingActionExists(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertBuildingActionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Zero(t, value)
}
