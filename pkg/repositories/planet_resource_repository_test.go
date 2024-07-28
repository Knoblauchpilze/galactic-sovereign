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
	Planet:    defaultPlanetId,
	Resource:  defaultPlayerId,
	Amount:    1234.567,
	CreatedAt: time.Date(2024, 7, 28, 10, 59, 41, 651387233, time.UTC),
	UpdatedAt: time.Date(2024, 7, 28, 10, 59, 42, 651387233, time.UTC),
	Version:   5,
}
var defaultUpdatedPlanetResource = persistence.PlanetResource{
	Planet:    defaultPlanetId,
	Resource:  defaultResourceId,
	Amount:    456.789,
	CreatedAt: time.Date(2024, 7, 28, 10, 59, 41, 651387233, time.UTC),
	UpdatedAt: time.Date(2024, 7, 28, 11, 23, 12, 651387233, time.UTC),
	Version:   6,
}
var defaultTransactionTime = time.Date(2024, 7, 28, 16, 32, 14, 651387233, time.UTC)

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
					`INSERT INTO planet_resource (planet, resource, amount, created_at) VALUES($1, $2, $3, $4)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetResource.Planet,
						defaultPlanetResource.Resource,
						defaultPlanetResource.Amount,
						defaultPlanetResource.CreatedAt,
					},
				},
			},
			"createForPlanet": {
				sqlMode: ExecBased,
				generateMock: func() db.Transaction {
					return &mockTransaction{
						timeStamp: defaultTransactionTime,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceRepository()
					return s.CreateForPlanet(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`
INSERT INTO
	planet_resource (planet, resource, amount, created_at)
SELECT
	$1,
	id,
	start_amount,
	$2
FROM
	resource
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
						defaultTransactionTime,
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
						affectedRows: 1,
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
	updated_at = $2,
	version = $3
WHERE
	planet = $4
	AND resource = $5
	AND version = $6
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUpdatedPlanetResource.Amount,
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
						affectedRows: 1,
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
						affectedRows: 0,
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
						affectedRows: 2,
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
