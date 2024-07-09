package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var defaultPlanetId = uuid.MustParse("0a05c1be-3235-48d6-b0e2-849b8664a515")
var defaultPlanetName = "my-planet"
var defaultPlanet = persistence.Planet{
	Id:        defaultPlanetId,
	Player:    defaultPlayerId,
	Name:      defaultPlanetName,
	CreatedAt: time.Date(2024, 7, 9, 20, 11, 21, 651387230, time.UTC),
	UpdatedAt: time.Date(2024, 7, 9, 20, 11, 21, 651387230, time.UTC),
}

func TestPlanetRepository_Create_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewPlanetRepository(pool)
			_, err := repo.Create(ctx, defaultPlanet)
			return err
		},
		expectedSql: `INSERT INTO planet (id, player, name, created_at) VALUES($1, $2, $3, $4)`,
		expectedArguments: []interface{}{
			defaultPlanet.Id,
			defaultPlanet.Player,
			defaultPlanet.Name,
			defaultPlanet.CreatedAt,
		},
	}

	suite.Run(t, &s)
}

func TestPlanetRepository_Create_ReturnsInputPlanet(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewPlanetRepository(mc)

	actual, err := repo.Create(context.Background(), defaultPlanet)

	assert.Nil(err)
	assert.Equal(defaultPlanet, actual)
}

func TestPlanetRepository_Get_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewPlanetRepository(pool)
			_, err := repo.Get(ctx, defaultPlanetId)
			return err
		},
		expectedSql: `SELECT id, player, name, created_at, updated_at FROM planet WHERE id = $1`,
		expectedArguments: []interface{}{
			defaultPlanetId,
		},
	}

	suite.Run(t, &s)
}

func TestPlanetRepository_Get_InterpretDbData(t *testing.T) {
	dummyStr := ""

	s := RepositorySingleValueTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewPlanetRepository(pool)
			_, err := repo.Get(ctx, defaultPlanetId)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
				&uuid.UUID{},
				&uuid.UUID{},
				&dummyStr,
				&time.Time{},
				&time.Time{},
			},
		},
	}

	suite.Run(t, &s)
}

func TestPlanetRepository_List_DbInteraction(t *testing.T) {
	s := RepositoryPoolTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewPlanetRepository(pool)
			_, err := repo.List(ctx)
			return err
		},
		expectedSql: `SELECT id, player, name, created_at, updated_at FROM planet`,
	}

	suite.Run(t, &s)
}

func TestPlanetRepository_List_InterpretDbData(t *testing.T) {
	dummyStr := ""

	s := RepositoryGetAllTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewPlanetRepository(pool)
			_, err := repo.List(ctx)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
				&uuid.UUID{},
				&uuid.UUID{},
				&dummyStr,
				&time.Time{},
				&time.Time{},
			},
		},
	}

	suite.Run(t, &s)
}

func TestPlanetRepository_Delete_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewPlanetRepository(&mockConnectionPool{})
			return repo.Delete(ctx, tx, defaultPlanetId)
		},
		expectedSql: []string{`DELETE FROM planet WHERE id = $1`},
		expectedArguments: [][]interface{}{
			{
				defaultPlanetId,
			},
		},
	}

	suite.Run(t, &s)
}

func TestPlanetRepository_Delete_WhenAffectedRowsIsZero_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewPlanetRepository(mc)
	mt := &mockTransaction{
		affectedRows: 0,
	}

	err := repo.Delete(context.Background(), mt, defaultPlanetId)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestPlanetRepository_Delete_WhenAffectedRowsIsGreaterThanOne_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewPlanetRepository(mc)
	mt := &mockTransaction{
		affectedRows: 2,
	}

	err := repo.Delete(context.Background(), mt, defaultPlanetId)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestPlanetRepository_Delete_WhenAffectedRowsIsOne_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewPlanetRepository(mc)
	mt := &mockTransaction{
		affectedRows: 1,
	}

	err := repo.Delete(context.Background(), mt, defaultPlanetId)

	assert.Nil(err)
}
