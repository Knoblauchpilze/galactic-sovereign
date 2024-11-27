package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultBuildingActionResourceProduction = persistence.BuildingActionResourceProduction{
	Action:     defaultBuildinActionId,
	Resource:   defaultResourceId,
	Production: 18,
}

func TestUnit_BuildingActionResourceProductionRepository_Transaction(t *testing.T) {
	dummyInt := 0

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionResourceProductionRepository()
					_, err := s.Create(ctx, tx, defaultBuildingActionResourceProduction)
					return err
				},
				expectedSqlQueries: []string{
					`
INSERT INTO
	building_action_resource_production (action, resource, production)
	VALUES($1, $2, $3)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultBuildingActionResourceProduction.Action,
						defaultBuildingActionResourceProduction.Resource,
						defaultBuildingActionResourceProduction.Production,
					},
				},
			},
			"listForAction": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionResourceProductionRepository()
					_, err := s.ListForAction(ctx, tx, defaultBuildinActionId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	action,
	resource,
	production
FROM
	building_action_resource_production
WHERE
	action = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultBuildinActionId,
					},
				},
			},
			"deleteForAction": {
				sqlMode: ExecBased,
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: []int{1},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionResourceProductionRepository()
					return s.DeleteForAction(ctx, tx, defaultBuildinActionId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM building_action_resource_production WHERE action = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultBuildinActionId,
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
					s := NewBuildingActionResourceProductionRepository()
					return s.DeleteForPlanet(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`
DELETE FROM
	building_action_resource_production
USING
	building_action_resource_production AS barp
	LEFT JOIN building_action AS ba ON ba.id = barp.action
WHERE
	building_action_resource_production.action = barp.action
	AND ba.planet = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
					},
				},
			},
		},

		dbGetAllTestCases: map[string]dbTransactionGetAllTestCase{
			"listForAction": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewBuildingActionResourceProductionRepository()
					_, err := repo.ListForAction(ctx, tx, defaultBuildinActionId)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyInt,
					},
				},
			},
		},

		dbReturnTestCases: map[string]dbTransactionReturnTestCase{
			"create": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					s := NewBuildingActionResourceProductionRepository()
					out, _ := s.Create(ctx, tx, defaultBuildingActionResourceProduction)
					return out
				},
				expectedContent: defaultBuildingActionResourceProduction,
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"create_duplicatedKey": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						execErrs: []error{
							fmt.Errorf(`duplicate key value violates unique constraint "building_action_resource_production_action_resource_key" (SQLSTATE 23505)`),
						},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionResourceProductionRepository()
					_, err := s.Create(ctx, tx, defaultBuildingActionResourceProduction)
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
