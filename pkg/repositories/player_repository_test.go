package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var defaultPlayerId = uuid.MustParse("bfd8259a-a8fd-4fba-b5db-f627d6dc055c")
var defaultPlayerName = "my-player"
var defaultPlayer = persistence.Player{
	Id:        defaultPlayerId,
	ApiUser:   defaultUserId,
	Universe:  defaultUniverseId,
	Name:      defaultPlayerName,
	CreatedAt: time.Date(2024, 7, 8, 22, 16, 0, 651387231, time.UTC),
	UpdatedAt: time.Date(2024, 7, 8, 22, 16, 0, 651387231, time.UTC),
	Version:   5,
}

func Test_PlayerRepository(t *testing.T) {
	dummyStr := ""
	dummyInt := 0

	s := RepositoryTestSuite{
		dbPoolInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlayerRepository(pool)
					_, err := s.Create(ctx, defaultPlayer)
					return err
				},
				expectedSql: `INSERT INTO player (id, api_user, universe, name, created_at) VALUES($1, $2, $3, $4, $5)`,
				expectedArguments: []interface{}{
					defaultPlayer.Id,
					defaultPlayer.ApiUser,
					defaultPlayer.Universe,
					defaultPlayer.Name,
					defaultPlayer.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlayerRepository(pool)
					_, err := s.Get(ctx, defaultPlayerId)
					return err
				},
				expectedSql: `SELECT id, api_user, universe, name, created_at, updated_at, version FROM player WHERE id = $1`,
				expectedArguments: []interface{}{
					defaultPlayerId,
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlayerRepository(pool)
					_, err := s.List(ctx)
					return err
				},
				expectedSql: `SELECT id, api_user, universe, name, created_at, updated_at, version FROM player`,
			},
		},

		dbPoolSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewPlayerRepository(pool)
					_, err := repo.Get(ctx, defaultPlayerId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyStr,
						&time.Time{},
						&time.Time{},
						&dummyInt,
					},
				},
			},
		},

		dbPoolReturnTestCases: map[string]dbPoolReturnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool) interface{} {
					s := NewPlayerRepository(pool)
					out, _ := s.Create(ctx, defaultPlayer)
					return out
				},
				expectedContent: defaultPlayer,
			},
		},
	}

	suite.Run(t, &s)
}

func TestPlayerRepository_Create_WhenQueryIndicatesDuplicatedKey_ReturnsDuplicatedKey(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		execErr: fmt.Errorf(`duplicate key value violates unique constraint "player_universe_name_key" (SQLSTATE 23505)`),
	}
	repo := NewPlayerRepository(mc)

	_, err := repo.Create(context.Background(), defaultPlayer)

	assert.True(errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey))
}

func TestPlayerRepository_List_InterpretDbData(t *testing.T) {
	dummyStr := ""
	dummyInt := 0

	s := RepositoryGetAllTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewPlayerRepository(pool)
			_, err := repo.List(ctx)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
				&uuid.UUID{},
				&uuid.UUID{},
				&uuid.UUID{},
				&dummyStr,
				&time.Time{},
				&time.Time{},
				&dummyInt,
			},
		},
	}

	suite.Run(t, &s)
}

func TestPlayerRepository_Delete_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewPlayerRepository(&mockConnectionPool{})
			return repo.Delete(ctx, tx, defaultPlayerId)
		},
		expectedSql: []string{`DELETE FROM player WHERE id = $1`},
		expectedArguments: [][]interface{}{
			{
				defaultPlayerId,
			},
		},
	}

	suite.Run(t, &s)
}

func TestPlayerRepository_Delete_WhenAffectedRowsIsZero_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewPlayerRepository(mc)
	mt := &mockTransaction{
		affectedRows: 0,
	}

	err := repo.Delete(context.Background(), mt, defaultPlayerId)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestPlayerRepository_Delete_WhenAffectedRowsIsGreaterThanOne_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewPlayerRepository(mc)
	mt := &mockTransaction{
		affectedRows: 2,
	}

	err := repo.Delete(context.Background(), mt, defaultPlayerId)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestPlayerRepository_Delete_WhenAffectedRowsIsOne_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewPlayerRepository(mc)
	mt := &mockTransaction{
		affectedRows: 1,
	}

	err := repo.Delete(context.Background(), mt, defaultPlayerId)

	assert.Nil(err)
}
