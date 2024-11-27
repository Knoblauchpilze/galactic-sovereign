package repositories

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

var defaultBuildingId = uuid.MustParse("9c5a9f5c-b53e-4f5c-af6f-cb04e47abd95")

func TestUnit_BuildingCostRepository_Transaction(t *testing.T) {
	var dummyInt int
	var dummyFloat64 float64

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"listForBuilding": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingCostRepository()
					_, err := s.ListForBuilding(ctx, tx, defaultBuildingId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	building,
	resource,
	cost,
	progress
FROM
	building_cost
WHERE
	building = $1
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultBuildingId,
					},
				},
			},
		},

		dbGetAllTestCases: map[string]dbTransactionGetAllTestCase{
			"listForBuilding": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewBuildingCostRepository()
					_, err := repo.ListForBuilding(ctx, tx, defaultBuildingId)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyInt,
						&dummyFloat64,
					},
				},
			},
		},
	}

	suite.Run(t, &s)
}
