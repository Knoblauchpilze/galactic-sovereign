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

func TestUniverseRepository_Create_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUniverseRepository(pool)
			_, err := repo.Create(ctx, defaultUniverse)
			return err
		},
		expectedSql: `INSERT INTO universe (id, name, created_at) VALUES($1, $2, $3)`,
		expectedArguments: []interface{}{
			defaultUniverse.Id,
			defaultUniverse.Name,
			defaultUniverse.CreatedAt,
		},
	}

	suite.Run(t, &s)
}

func TestUniverseRepository_Create_WhenQueryIndicatesDuplicatedKey_ReturnsDuplicatedKey(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		execErr: fmt.Errorf(`duplicate key value violates unique constraint "api_user_email_key" (SQLSTATE 23505)`),
	}
	repo := NewUniverseRepository(mc)

	_, err := repo.Create(context.Background(), defaultUniverse)

	assert.True(errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey))
}

func TestUniverseRepository_Create_ReturnsInputUniverse(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUniverseRepository(mc)

	actual, err := repo.Create(context.Background(), defaultUniverse)

	assert.Nil(err)
	assert.Equal(defaultUniverse, actual)
}

func TestUniverseRepository_Get_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUniverseRepository(pool)
			_, err := repo.Get(ctx, defaultUniverseId)
			return err
		},
		expectedSql: `SELECT id, name, created_at, updated_at, version FROM universe WHERE id = $1`,
		expectedArguments: []interface{}{
			defaultUniverseId,
		},
	}

	suite.Run(t, &s)
}

func TestUniverseRepository_Get_InterpretDbData(t *testing.T) {
	dummyStr := ""
	dummyInt := 0

	s := RepositorySingleValueTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUniverseRepository(pool)
			_, err := repo.Get(ctx, defaultUniverseId)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
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

func TestUniverseRepository_List_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUniverseRepository(pool)
			_, err := repo.List(ctx)
			return err
		},
		expectedSql: `SELECT id, name, created_at, updated_at, version FROM universe`,
	}

	suite.Run(t, &s)
}

func TestUniverseRepository_List_InterpretDbData(t *testing.T) {
	dummyStr := ""
	dummyInt := 0

	s := RepositoryGetAllTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUniverseRepository(pool)
			_, err := repo.List(ctx)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
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
