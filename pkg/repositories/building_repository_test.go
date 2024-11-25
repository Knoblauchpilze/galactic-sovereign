package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func TestUnit__BuildingRepository_Transaction(t *testing.T) {
	dummyStr := ""

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"list": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewBuildingRepository()
					_, err := s.List(ctx, tx)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, name, created_at, updated_at FROM building`,
				},
			},
		},

		dbGetAllTestCases: map[string]dbTransactionGetAllTestCase{
			"list": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewBuildingRepository()
					_, err := repo.List(ctx, tx)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyStr,
						&time.Time{},
						&time.Time{},
					},
				},
			},
		},
	}

	suite.Run(t, &s)
}
