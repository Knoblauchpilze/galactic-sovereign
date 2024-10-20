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

var defaultUniverseId = uuid.MustParse("80a18e0d-bf9c-4021-8986-4685f8a62601")
var defaultUniverseName = "my-universe"
var defaultUniverse = persistence.Universe{
	Id:        defaultUniverseId,
	Name:      defaultUniverseName,
	CreatedAt: time.Date(2024, 7, 8, 21, 20, 15, 651387239, time.UTC),
	UpdatedAt: time.Date(2024, 7, 8, 21, 20, 15, 651387239, time.UTC),
	Version:   5,
}

func Test_UniverseRepository(t *testing.T) {
	dummyStr := ""
	dummyInt := 0

	s := RepositoryPoolTestSuite{
		dbInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUniverseRepository(pool)
					_, err := s.Create(ctx, defaultUniverse)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO universe (id, name, created_at) VALUES($1, $2, $3)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUniverse.Id,
						defaultUniverse.Name,
						defaultUniverse.CreatedAt,
					},
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUniverseRepository(pool)
					_, err := s.List(ctx)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, name, created_at, updated_at, version FROM universe`,
				},
			},
		},

		dbGetAllTestCases: map[string]dbPoolGetAllTestCase{
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewUniverseRepository(pool)
					_, err := repo.List(ctx)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
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
					s := NewUniverseRepository(pool)
					out, _ := s.Create(ctx, defaultUniverse)
					return out
				},
				expectedContent: defaultUniverse,
			},
		},

		dbErrorTestCases: map[string]dbPoolErrorTestCase{
			"create_duplicatedKey": {
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						execErr: fmt.Errorf(`duplicate key value violates unique constraint "universe_name_key" (SQLSTATE 23505)`),
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUniverseRepository(pool)
					_, err := s.Create(ctx, defaultUniverse)
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

func Test_UniverseRepository_Transaction(t *testing.T) {
	dummyStr := ""
	dummyInt := 0

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"get": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUniverseRepository(&mockConnectionPool{})
					_, err := s.Get(ctx, tx, defaultUniverseId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, name, created_at, updated_at, version FROM universe WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUniverseId,
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
					s := NewUniverseRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultUniverseId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM universe WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUniverseId,
					},
				},
			},
		},

		dbSingleValueTestCases: map[string]dbTransactionSingleValueTestCase{
			"get": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewUniverseRepository(&mockConnectionPool{})
					_, err := repo.Get(ctx, tx, defaultUniverseId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyStr,
						&time.Time{},
						&time.Time{},
						&dummyInt,
					},
				},
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"delete_noRowsAffected": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{0},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUniverseRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultUniverseId)
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
					s := NewUniverseRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultUniverseId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
		},
	}

	suite.Run(t, &s)
}
