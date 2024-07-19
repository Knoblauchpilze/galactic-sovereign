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

func Test_PlanetRepository(t *testing.T) {
	dummyStr := ""

	s := RepositoryPoolTestSuite{
		dbInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlanetRepository(pool)
					_, err := s.Create(ctx, defaultPlanet)
					return err
				},
				expectedSql: `INSERT INTO planet (id, player, name, created_at) VALUES($1, $2, $3, $4)`,
				expectedArguments: []interface{}{
					defaultPlanet.Id,
					defaultPlanet.Player,
					defaultPlanet.Name,
					defaultPlanet.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlanetRepository(pool)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedSql: `SELECT id, player, name, created_at, updated_at FROM planet WHERE id = $1`,
				expectedArguments: []interface{}{
					defaultPlanetId,
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlanetRepository(pool)
					_, err := s.List(ctx)
					return err
				},
				expectedSql: `SELECT id, player, name, created_at, updated_at FROM planet`,
			},
		},

		dbSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewPlanetRepository(pool)
					_, err := repo.Get(ctx, defaultPlanetId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyStr,
						&time.Time{},
						&time.Time{},
					},
				},
			},
		},

		dbGetAllTestCases: map[string]dbPoolGetAllTestCase{
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewPlanetRepository(pool)
					_, err := repo.List(ctx)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyStr,
						&time.Time{},
						&time.Time{},
					},
				},
			},
		},

		dbReturnTestCases: map[string]dbPoolReturnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool) interface{} {
					s := NewPlanetRepository(pool)
					out, _ := s.Create(ctx, defaultPlanet)
					return out
				},
				expectedContent: defaultPlanet,
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
