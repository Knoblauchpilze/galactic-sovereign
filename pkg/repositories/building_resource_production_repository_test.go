package repositories

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func Test_BuildingResourceProductionRepository_Transaction(t *testing.T) {
	var dummyInt int
	var dummyFloat64 float64

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"listForBuilding": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingResourceProductionRepository()
					_, err := s.ListForBuilding(ctx, tx, defaultBuildingId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	building,
	resource,
	base,
	progress
FROM
	building_resource_production
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
					repo := NewBuildingResourceProductionRepository()
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
