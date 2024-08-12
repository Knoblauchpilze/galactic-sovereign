package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultBuildinActionId = uuid.MustParse("1ec83c50-0f25-4918-b49a-52b7816189b9")
var defaultBuildingAction = persistence.BuildingAction{
	Id:           defaultBuildinActionId,
	Planet:       defaultPlanetId,
	Building:     defaultBuildingId,
	CurrentLevel: 56,
	DesiredLevel: 61,
	CreatedAt:    time.Date(2024, 8, 11, 21, 40, 51, 651387244, time.UTC),
	CompletedAt:  time.Date(2024, 7, 11, 21, 40, 54, 651387244, time.UTC),
}

func Test_BuildingActionRepository(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{},
	}

	suite.Run(t, &s)
}

func Test_BuildingActionRepository_Transaction(t *testing.T) {
	var dummyInt int

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionRepository()
					_, err := s.Create(ctx, tx, defaultBuildingAction)
					return err
				},
				expectedSqlQueries: []string{
					`
INSERT INTO
	building_action (id, planet, building, current_level, desired_level, created_at, completed_at)
	VALUES($1, $2, $3, $4, $5, $6, $7)
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultBuildingAction.Id,
						defaultBuildingAction.Planet,
						defaultBuildingAction.Building,
						defaultBuildingAction.CurrentLevel,
						defaultBuildingAction.DesiredLevel,
						defaultBuildingAction.CreatedAt,
						defaultBuildingAction.CompletedAt,
					},
				},
			},
			"listForPlanet": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionRepository()
					_, err := s.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	id,
	planet,
	building,
	current_level,
	desired_level,
	created_at,
	completed_at
FROM
	building_action
WHERE
	planet = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultPlanetId,
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
					s := NewBuildingActionRepository()
					return s.DeleteForPlanet(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM building_action WHERE planet = $1`,
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
					repo := NewBuildingActionRepository()
					_, err := repo.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyInt,
						&dummyInt,
						&time.Time{},
						&time.Time{},
					},
				},
			},
		},

		dbReturnTestCases: map[string]dbTransactionReturnTestCase{
			"create": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					s := NewBuildingActionRepository()
					out, _ := s.Create(ctx, tx, defaultBuildingAction)
					return out
				},
				expectedContent: defaultBuildingAction,
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"create_duplicatedKey": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						execErrs: []error{
							fmt.Errorf(`duplicate key value violates unique constraint "universe_name_key" (SQLSTATE 23505)`),
						},
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingActionRepository()
					_, err := s.Create(ctx, tx, defaultBuildingAction)
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