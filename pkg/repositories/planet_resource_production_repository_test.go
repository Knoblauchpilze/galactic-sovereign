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

var defaultPlanetResourceProduction = persistence.PlanetResourceProduction{
	Planet:     defaultPlanetId,
	Building:   &defaultBuildingId,
	Resource:   defaultResourceId,
	Production: 89,
	CreatedAt:  time.Date(2024, 9, 15, 11, 57, 02, 651387250, time.UTC),
	UpdatedAt:  time.Date(2024, 9, 15, 11, 57, 03, 651387250, time.UTC),
	Version:    5,
}
var defaultUpdatedPlanetResourceProduction = persistence.PlanetResourceProduction{
	Planet:     defaultPlanetId,
	Building:   &defaultBuildingId,
	Resource:   defaultResourceId,
	Production: 321,
	CreatedAt:  time.Date(2024, 9, 15, 11, 57, 02, 651387250, time.UTC),
	UpdatedAt:  time.Date(2024, 9, 15, 11, 58, 53, 651387250, time.UTC),
	Version:    6,
}

func TestUnit__PlanetResourceProductionRepository_Transaction(t *testing.T) {
	var dummyUuid *uuid.UUID
	var dummyInt int

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceProductionRepository()
					_, err := s.Create(ctx, tx, defaultPlanetResourceProduction)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO planet_resource_production (planet, building, resource, production, created_at) VALUES($1, $2, $3, $4, $5)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetResourceProduction.Planet,
						defaultPlanetResourceProduction.Building,
						defaultPlanetResourceProduction.Resource,
						defaultPlanetResourceProduction.Production,
						defaultPlanetResourceProduction.CreatedAt,
					},
				},
			},
			"getForPlanetAndBuilding": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceProductionRepository()
					_, err := s.GetForPlanetAndBuilding(ctx, tx, defaultPlanetId, &defaultBuildingId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	planet,
	building,
	resource,
	production,
	created_at,
	updated_at,
	version
FROM
	planet_resource_production
WHERE
	planet = $1
	AND building = $2
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
						&defaultBuildingId,
					},
				},
			},
			"getForPlanetAndBuilding_buildingIsNil": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceProductionRepository()
					_, err := s.GetForPlanetAndBuilding(ctx, tx, defaultPlanetId, nil)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	planet,
	building,
	resource,
	production,
	created_at,
	updated_at,
	version
FROM
	planet_resource_production
WHERE
	planet = $1
	AND building = $2
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
						dummyUuid,
					},
				},
			},
			"listForPlanet": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetResourceProductionRepository()
					_, err := s.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	planet,
	building,
	resource,
	production,
	created_at,
	updated_at,
	version
FROM
	planet_resource_production
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
					s := NewPlanetResourceProductionRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResourceProduction)
					return err
				},
				expectedSqlQueries: []string{
					`
UPDATE
	planet_resource_production
SET
	production = $1,
	updated_at = $2,
	version = $3
WHERE
	planet = $4
	AND building = $5
	AND resource = $6
	AND version = $7
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUpdatedPlanetResourceProduction.Production,
						defaultUpdatedPlanetResourceProduction.UpdatedAt,
						defaultUpdatedPlanetResourceProduction.Version + 1,
						defaultUpdatedPlanetResourceProduction.Planet,
						defaultUpdatedPlanetResourceProduction.Building,
						defaultUpdatedPlanetResourceProduction.Resource,
						defaultUpdatedPlanetResourceProduction.Version,
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
					s := NewPlanetResourceProductionRepository()
					return s.DeleteForPlanet(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM planet_resource_production WHERE planet = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
					},
				},
			},
		},

		dbSingleValueTestCases: map[string]dbTransactionSingleValueTestCase{
			"getForPlanetAndBuilding": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewPlanetResourceProductionRepository()
					_, err := repo.GetForPlanetAndBuilding(ctx, tx, defaultPlanetId, &defaultBuildingId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyUuid,
						&uuid.UUID{},
						&dummyInt,
						&time.Time{},
						&time.Time{},
						&dummyInt,
					},
				},
			},
			"getForPlanetAndBuilding_buildingIsNil": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewPlanetResourceProductionRepository()
					_, err := repo.GetForPlanetAndBuilding(ctx, tx, defaultPlanetId, nil)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyUuid,
						&uuid.UUID{},
						&dummyInt,
						&time.Time{},
						&time.Time{},
						&dummyInt,
					},
				},
			},
		},

		dbGetAllTestCases: map[string]dbTransactionGetAllTestCase{
			"listForPlanet": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewPlanetResourceProductionRepository()
					_, err := repo.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyUuid,
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
					s := NewPlanetResourceProductionRepository()
					out, _ := s.Create(ctx, tx, defaultPlanetResourceProduction)
					return out
				},
				expectedContent: defaultPlanetResourceProduction,
			},
			"update": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					s := NewPlanetResourceProductionRepository()
					out, _ := s.Update(ctx, tx, defaultUpdatedPlanetResourceProduction)
					return out
				},
				expectedContent: defaultUpdatedPlanetResourceProduction,
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
					s := NewPlanetResourceProductionRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResourceProduction)
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
					s := NewPlanetResourceProductionRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetResourceProduction)
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
