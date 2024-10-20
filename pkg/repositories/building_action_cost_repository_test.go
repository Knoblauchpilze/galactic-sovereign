package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultBuildingActionCost = persistence.BuildingActionCost{
	Action:   defaultBuildinActionId,
	Resource: defaultResourceId,
	Amount:   26,
}

func Test_BuildingActionCostRepository_Transaction(t *testing.T) {
	dummyInt := 0

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionCostRepository()
					_, err := s.Create(ctx, tx, defaultBuildingActionCost)
					return err
				},
				expectedSqlQueries: []string{
					`
INSERT INTO
	building_action_cost (action, resource, amount)
	VALUES($1, $2, $3)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultBuildingActionCost.Action,
						defaultBuildingActionCost.Resource,
						defaultBuildingActionCost.Amount,
					},
				},
			},
			"listForAction": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionCostRepository()
					_, err := s.ListForAction(ctx, tx, defaultBuildinActionId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	action,
	resource,
	amount
FROM
	building_action_cost
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
					s := NewBuildingActionCostRepository()
					return s.DeleteForAction(ctx, tx, defaultBuildinActionId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM building_action_cost WHERE action = $1`,
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
					s := NewBuildingActionCostRepository()
					return s.DeleteForPlanet(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`
DELETE FROM
	building_action_cost
USING
	building_action_cost AS bac
	LEFT JOIN building_action AS ba ON ba.id = bac.action
WHERE
	building_action_cost.action = bac.action
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
					repo := NewBuildingActionCostRepository()
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
					s := NewBuildingActionCostRepository()
					out, _ := s.Create(ctx, tx, defaultBuildingActionCost)
					return out
				},
				expectedContent: defaultBuildingActionCost,
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"create_duplicatedKey": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						execErrs: []error{
							fmt.Errorf(`duplicate key value violates unique constraint "building_action_cost_action_resource_key" (SQLSTATE 23505)`),
						},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionCostRepository()
					_, err := s.Create(ctx, tx, defaultBuildingActionCost)
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
