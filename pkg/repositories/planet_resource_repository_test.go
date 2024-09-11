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

var defaultPlanetResource = persistence.PlanetResource{
	Planet:     defaultPlanetId,
	Resource:   defaultResourceId,
	Amount:     1234.567,
	Production: 36,
	CreatedAt:  time.Date(2024, 7, 28, 10, 59, 41, 651387233, time.UTC),
	UpdatedAt:  time.Date(2024, 7, 28, 10, 59, 42, 651387233, time.UTC),
	Version:    5,
}
var defaultUpdatedPlanetResource = persistence.PlanetResource{
	Planet:    defaultPlanetId,
	Resource:  defaultResourceId,
	Amount:    456.789,
	CreatedAt: time.Date(2024, 7, 28, 10, 59, 41, 651387233, time.UTC),
	UpdatedAt: time.Date(2024, 7, 28, 11, 23, 12, 651387233, time.UTC),
	Version:   6,
}

func Test_PlanetResourceRepository_Transaction(t *testing.T) {
	var dummyFloat64 float64
	var dummyInt int

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceRepository()
					_, err := s.Create(ctx, tx, defaultPlanetResource)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO planet_resource (planet, resource, amount, production, created_at) VALUES($1, $2, $3, $4, $5)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetResource.Planet,
						defaultPlanetResource.Resource,
						defaultPlanetResource.Amount,
						defaultPlanetResource.Production,
						defaultPlanetResource.CreatedAt,
					},
				},
			},
			"listForPlanet": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceRepository()
					_, err := s.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	planet,
	resource,
	amount,
	production,
	created_at,
	updated_at,
	version
FROM
	planet_resource
WHERE
	planet = $1
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
					},
				},
			},
			"update": {
				sqlMode: ExecBased,
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{1},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResource)
					return err
				},
				expectedSqlQueries: []string{
					`
UPDATE
	planet_resource
SET
	amount = $1,
	production = $2,
	updated_at = $3,
	version = $4
WHERE
	planet = $5
	AND resource = $6
	AND version = $7
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUpdatedPlanetResource.Amount,
						defaultUpdatedPlanetResource.Production,
						defaultUpdatedPlanetResource.UpdatedAt,
						defaultUpdatedPlanetResource.Version + 1,
						defaultUpdatedPlanetResource.Planet,
						defaultUpdatedPlanetResource.Resource,
						defaultUpdatedPlanetResource.Version,
					},
				},
			},
			"deleteForPlanet": {
				sqlMode: ExecBased,
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{1},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceRepository()
					return s.DeleteForPlanet(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM planet_resource WHERE planet = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
					},
				},
			},
		},

		dbGetAllTestCases: map[string]dbTransactionGetAllTestCase{
			"listForPlanet": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewPlanetResourceRepository()
					_, err := repo.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyFloat64,
						&dummyInt,
						&time.Time{},
						&time.Time{},
						&dummyInt,
					},
				},
			},
		},

		dbReturnTestCases: map[string]dbTransactionReturnTestCase{
			"create": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					s := NewPlanetResourceRepository()
					out, _ := s.Create(ctx, tx, defaultPlanetResource)
					return out
				},
				expectedContent: defaultPlanetResource,
			},
			"update": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					s := NewPlanetResourceRepository()
					out, _ := s.Update(ctx, tx, defaultUpdatedPlanetResource)
					return out
				},
				expectedContent: defaultUpdatedPlanetResource,
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"update_optimisticLockException": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{0},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResource)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.OptimisticLockException))
				},
			},
			"update_moreThanOneRowAffected": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{2},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResource)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.MoreThanOneMatchingSqlRows))
				},
			},
		},
	}

	suite.Run(t, &s)
}
