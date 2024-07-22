package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultPlanetId = uuid.MustParse("0a05c1be-3235-48d6-b0e2-849b8664a515")
var defaultPlanetName = "my-planet"
var defaultPlanet = persistence.Planet{
	Id:        defaultPlanetId,
	Player:    defaultPlayerId,
	Name:      defaultPlanetName,
	Homeworld: true,
	CreatedAt: time.Date(2024, 7, 9, 20, 11, 21, 651387230, time.UTC),
	UpdatedAt: time.Date(2024, 7, 9, 20, 11, 21, 651387230, time.UTC),
}

func Test_PlanetRepository(t *testing.T) {
	dummyStr := ""
	dummyBool := false

	s := RepositoryPoolTestSuite{
		dbInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlanetRepository(pool)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
WHERE
	id = $1
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
					},
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlanetRepository(pool)
					_, err := s.List(ctx)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
`,
				},
			},
			"listForPlayer": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewPlanetRepository(pool)
					_, err := s.ListForPlayer(ctx, defaultPlayerId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	p.id,
	p.player,
	p.name,
	CASE
		WHEN h.planet IS NOT NULL THEN true
		ELSE false
	END AS homeworld,
	p.created_at,
	p.updated_at
FROM
	planet AS p
	LEFT JOIN homeworld AS h ON h.planet = p.id
WHERE
	p.player = $1
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlayerId,
					},
				},
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
						&dummyBool,
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
						&dummyBool,
						&time.Time{},
						&time.Time{},
					},
				},
			},
			"listForPlayer": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewPlanetRepository(pool)
					_, err := repo.ListForPlayer(ctx, defaultPlayerId)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyStr,
						&dummyBool,
						&time.Time{},
						&time.Time{},
					},
				},
			},
		},
	}

	suite.Run(t, &s)
}

func Test_PlanetRepository_Transaction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetRepository(&mockConnectionPool{})

					planet := persistence.Planet{
						Id:        defaultPlanet.Id,
						Player:    defaultPlanet.Player,
						Name:      defaultPlanet.Name,
						Homeworld: false,
						CreatedAt: defaultPlanet.CreatedAt,
						UpdatedAt: defaultPlanet.UpdatedAt,
					}

					_, err := s.Create(ctx, tx, planet)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO planet (id, player, name, created_at) VALUES($1, $2, $3, $4)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanet.Id,
						defaultPlanet.Player,
						defaultPlanet.Name,
						defaultPlanet.CreatedAt,
					},
				},
			},
			"create_homeworld": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetRepository(&mockConnectionPool{})
					_, err := s.Create(ctx, tx, defaultPlanet)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO planet (id, player, name, created_at) VALUES($1, $2, $3, $4)`,
					`INSERT INTO homeworld (player, planet) VALUES($1, $2)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanet.Id,
						defaultPlanet.Player,
						defaultPlanet.Name,
						defaultPlanet.CreatedAt,
					},
					{
						defaultPlanet.Player,
						defaultPlanet.Id,
					},
				},
			},

			"delete": {
				sqlMode: ExecBased,
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: 1,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM planet WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
					},
				},
			},
		},

		dbReturnTestCases: map[string]dbTransactionReturnTestCase{
			"create": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					s := NewPlanetRepository(&mockConnectionPool{})
					out, _ := s.Create(ctx, tx, defaultPlanet)
					return out
				},
				expectedContent: defaultPlanet,
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"delete_noRowsAffected": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: 0,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultPlanetId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
			"delete_moreThanOneRowAffected": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: 2,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultPlanetId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
		},
	}

	suite.Run(t, &s)
}
