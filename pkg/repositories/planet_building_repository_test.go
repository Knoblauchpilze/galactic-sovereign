package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

var defaultUpdatedPlanetBuilding = persistence.PlanetBuilding{
	Planet:    defaultPlanetId,
	Building:  defaultBuildingId,
	Level:     36,
	CreatedAt: time.Date(2024, 8, 16, 14, 9, 21, 651387245, time.UTC),
	UpdatedAt: time.Date(2024, 8, 16, 14, 9, 22, 651387245, time.UTC),
	Version:   4,
}

func Test_PlanetBuildingRepository_Transaction(t *testing.T) {
	var dummyInt int

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"listForPlanet": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewPlanetBuildingRepository()
					_, err := s.ListForPlanet(ctx, tx, defaultPlanetId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	planet,
	building,
	level,
	created_at,
	updated_at,
	version
FROM
	planet_building
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
					s := NewPlanetBuildingRepository()
					_, err := s.Update(ctx, tx, defaultUpdatedPlanetBuilding)
					return err
				},
				expectedSqlQueries: []string{
					`
UPDATE
	planet_building
SET
	level = $1,
	version = $2
WHERE
	planet = $3
	AND building = $4
	AND version = $5
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUpdatedPlanetBuilding.Level,
						defaultUpdatedPlanetBuilding.Version + 1,
						defaultUpdatedPlanetBuilding.Planet,
						defaultUpdatedPlanetBuilding.Building,
						defaultUpdatedPlanetBuilding.Version,
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
					s := NewPlanetBuildingRepository()
					return s.DeleteForPlanet(ctx, tx, defaultPlanetId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM planet_building WHERE planet = $1`,
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
					repo := NewPlanetBuildingRepository()
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
	}

	suite.Run(t, &s)
}
