package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultPlanetResourceStorage = persistence.PlanetResourceStorage{
	Planet:    defaultPlanetId,
	Resource:  defaultResourceId,
	Storage:   31,
	CreatedAt: time.Date(2024, 9, 28, 11, 04, 31, 651387252, time.UTC),
	UpdatedAt: time.Date(2024, 9, 28, 11, 04, 32, 651387252, time.UTC),
	Version:   2,
}
var defaultUpdatedPlanetResourceStorage = persistence.PlanetResourceStorage{
	Planet:    defaultPlanetId,
	Resource:  defaultResourceId,
	Storage:   321,
	CreatedAt: time.Date(2024, 9, 28, 11, 04, 52, 651387252, time.UTC),
	UpdatedAt: time.Date(2024, 9, 28, 11, 04, 53, 651387252, time.UTC),
	Version:   3,
}

func TestUnit__PlanetResourceStorageRepository_Transaction(t *testing.T) {
	var dummyInt int

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceStorageRepository()
					_, err := s.Create(ctx, tx, defaultPlanetResourceStorage)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO planet_resource_storage (planet, resource, storage, created_at) VALUES($1, $2, $3, $4)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetResourceStorage.Planet,
						defaultPlanetResourceStorage.Resource,
						defaultPlanetResourceStorage.Storage,
						defaultPlanetResourceStorage.CreatedAt,
					},
				},
			},
			"listForPlanet": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceStorageRepository()
					_, err := s.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	planet,
	resource,
	storage,
	created_at,
	updated_at,
	version
FROM
	planet_resource_storage
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
					s := NewPlanetResourceStorageRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResourceStorage)
					return err
				},
				expectedSqlQueries: []string{
					`
UPDATE
	planet_resource_storage
SET
	storage = $1,
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
						defaultUpdatedPlanetResourceStorage.Storage,
						defaultUpdatedPlanetResourceStorage.UpdatedAt,
						defaultUpdatedPlanetResourceStorage.Version + 1,
						defaultUpdatedPlanetResourceStorage.Planet,
						defaultUpdatedPlanetResourceStorage.Resource,
						defaultUpdatedPlanetResourceStorage.Version,
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
					s := NewPlanetResourceStorageRepository()
					return s.DeleteForPlanet(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM planet_resource_storage WHERE planet = $1`,
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
					repo := NewPlanetResourceStorageRepository()
					_, err := repo.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
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
					s := NewPlanetResourceStorageRepository()
					out, _ := s.Create(ctx, tx, defaultPlanetResourceStorage)
					return out
				},
				expectedContent: defaultPlanetResourceStorage,
			},
			"update": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					s := NewPlanetResourceStorageRepository()
					out, _ := s.Update(ctx, tx, defaultUpdatedPlanetResourceStorage)
					return out
				},
				expectedContent: defaultUpdatedPlanetResourceStorage,
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
					s := NewPlanetResourceStorageRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResourceStorage)
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
					s := NewPlanetResourceStorageRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResourceStorage)
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
