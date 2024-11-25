package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultUserId = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
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

func TestUnit__PlayerRepository(t *testing.T) {
	dummyStr := ""
	dummyInt := 0

	s := RepositoryPoolTestSuite{
		dbInteractionTestCases: map[string]dbPoolInteractionTestCase{
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
			"listForApiUser": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlayerRepository(pool)
					_, err := s.ListForApiUser(ctx, defaultUserId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, api_user, universe, name, created_at, updated_at, version FROM player where api_user = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUserId,
					},
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
			"listForApiUser": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewPlayerRepository(pool)
					_, err := repo.ListForApiUser(ctx, defaultUserId)
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
	}

	suite.Run(t, &s)
}

func TestUnit__PlayerRepository_Transaction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlayerRepository(&mockConnectionPool{})
					_, err := s.Create(ctx, tx, defaultPlayer)
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
			"delete": {
				sqlMode: ExecBased,
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{1},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlayerRepository(&mockConnectionPool{})
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

		dbReturnTestCases: map[string]dbTransactionReturnTestCase{
			"create": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					s := NewPlayerRepository(&mockConnectionPool{})
					out, _ := s.Create(ctx, tx, defaultPlayer)
					return out
				},
				expectedContent: defaultPlayer,
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"create_duplicatedKey": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						execErrs: []error{
							fmt.Errorf(`duplicate key value violates unique constraint "player_universe_name_key" (SQLSTATE 23505)`),
						},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlayerRepository(&mockConnectionPool{})
					_, err := s.Create(ctx, tx, defaultPlayer)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey))
				},
			},
			"delete_noRowsAffected": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{0},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlayerRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultPlayerId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
			"delete_moreThanOneRowAffected": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{2},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlayerRepository(&mockConnectionPool{})
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
