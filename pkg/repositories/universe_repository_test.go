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
				expectedSql: `INSERT INTO universe (id, name, created_at) VALUES($1, $2, $3)`,
				expectedArguments: []interface{}{
					defaultUniverse.Id,
					defaultUniverse.Name,
					defaultUniverse.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUniverseRepository(pool)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				expectedSql: `SELECT id, name, created_at, updated_at, version FROM universe WHERE id = $1`,
				expectedArguments: []interface{}{
					defaultUniverseId,
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUniverseRepository(pool)
					_, err := s.List(ctx)
					return err
				},
				expectedSql: `SELECT id, name, created_at, updated_at, version FROM universe`,
			},
		},

		dbSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewUniverseRepository(pool)
					_, err := repo.Get(ctx, defaultUniverseId)
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
	}

	suite.Run(t, &s)
}

func TestUniverseRepository_Create_WhenQueryIndicatesDuplicatedKey_ReturnsDuplicatedKey(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		execErr: fmt.Errorf(`duplicate key value violates unique constraint "universe_name_key" (SQLSTATE 23505)`),
	}
	repo := NewUniverseRepository(mc)

	_, err := repo.Create(context.Background(), defaultUniverse)

	assert.True(errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey))
}

func TestUniverseRepository_Delete_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUniverseRepository(&mockConnectionPool{})
			return repo.Delete(ctx, tx, defaultUniverseId)
		},
		expectedSql: []string{`DELETE FROM universe WHERE id = $1`},
		expectedArguments: [][]interface{}{
			{
				defaultUniverseId,
			},
		},
	}

	suite.Run(t, &s)
}

func TestUniverseRepository_Delete_WhenAffectedRowsIsZero_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUniverseRepository(mc)
	mt := &mockTransaction{
		affectedRows: 0,
	}

	err := repo.Delete(context.Background(), mt, defaultUniverseId)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestUniverseRepository_Delete_WhenAffectedRowsIsGreaterThanOne_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUniverseRepository(mc)
	mt := &mockTransaction{
		affectedRows: 2,
	}

	err := repo.Delete(context.Background(), mt, defaultUniverseId)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestUniverseRepository_Delete_WhenAffectedRowsIsOne_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUniverseRepository(mc)
	mt := &mockTransaction{
		affectedRows: 1,
	}

	err := repo.Delete(context.Background(), mt, defaultUniverseId)

	assert.Nil(err)
}
