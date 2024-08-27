package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultBuildingActionCost = persistence.BuildingActionCost{
	Action:   defaultBuildinActionId,
	Resource: defaultResourceId,
	Amount:   26,
}

func Test_BuildingActionCostRepository_Transaction(t *testing.T) {
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
	VALUES($1, $2, $3)
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultBuildingActionCost.Action,
						defaultBuildingActionCost.Resource,
						defaultBuildingActionCost.Amount,
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
