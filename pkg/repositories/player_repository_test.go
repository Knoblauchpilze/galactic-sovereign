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
	"github.com/stretchr/testify/require"
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

	s := RepositoryPoolTestSuite{
		dbInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlayerRepository(pool)
					_, err := s.Create(ctx, defaultPlayer)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO player (id, api_user, universe, name, created_at) VALUES($1, $2, $3, $4, $5)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlayer.Id,
						defaultPlayer.ApiUser,
						defaultPlayer.Universe,
						defaultPlayer.Name,
						defaultPlayer.CreatedAt,
					},
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlayerRepository(pool)
					_, err := s.Get(ctx, defaultPlayerId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, api_user, universe, name, created_at, updated_at, version FROM player WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{defaultPlayerId},
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlayerRepository(pool)
					_, err := s.List(ctx)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, api_user, universe, name, created_at, updated_at, version FROM player`,
				},
			},
		},

		dbSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
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

		dbGetAllTestCases: map[string]dbPoolGetAllTestCase{
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewPlayerRepository(pool)
					_, err := repo.List(ctx)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
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

		dbReturnTestCases: map[string]dbPoolReturnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool) interface{} {
					s := NewPlayerRepository(pool)
					out, _ := s.Create(ctx, defaultPlayer)
					return out
				},
				expectedContent: defaultPlayer,
			},
		},

		dbErrorTestCases: map[string]dbPoolErrorTestCase{
			"create_duplicatedKey": {
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPoolNew{
						execErr: fmt.Errorf(`duplicate key value violates unique constraint "player_universe_name_key" (SQLSTATE 23505)`),
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlayerRepository(pool)
					_, err := s.Create(ctx, defaultPlayer)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey))
				},
			},
		},
	}

	suite.Run(t, &s)
}

func Test_PlayerRepository_Transaction(t *testing.T) {
	s := RepositoryTransactionTestSuiteNew{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"delete": {
				sqlMode: ExecBased,
				generateMock: func() db.Transaction {
					return &mockTransactionNew{
						affectedRows: 1,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlayerRepository(&mockConnectionPoolNew{})
					return s.Delete(ctx, tx, defaultPlayerId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM player WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlayerId,
					},
				},
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"delete_noRowsAffected": {
				generateMock: func() db.Transaction {
					return &mockTransactionNew{
						affectedRows: 0,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlayerRepository(&mockConnectionPoolNew{})
					return s.Delete(ctx, tx, defaultPlayerId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
			"delete_moreThanOneRowAffected": {
				generateMock: func() db.Transaction {
					return &mockTransactionNew{
						affectedRows: 2,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlayerRepository(&mockConnectionPoolNew{})
					return s.Delete(ctx, tx, defaultPlayerId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
		},
	}

	suite.Run(t, &s)
}
