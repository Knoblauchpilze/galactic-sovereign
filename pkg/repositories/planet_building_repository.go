package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetBuildingRepository interface {
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetBuilding, error)
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
}

type planetBuildingRepositoryImpl struct{}

func NewPlanetBuildingRepository() PlanetBuildingRepository {
	return &planetBuildingRepositoryImpl{}
}

const listPlanetBuildingForPlanetSqlTemplate = `
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
`

func (r *planetBuildingRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetBuilding, error) {
	res := tx.Query(ctx, listPlanetBuildingForPlanetSqlTemplate, planet)
	if err := res.Err(); err != nil {
		return []persistence.PlanetBuilding{}, err
	}

	var out []persistence.PlanetBuilding
	parser := func(rows db.Scannable) error {
		var building persistence.PlanetBuilding
		err := rows.Scan(&building.Planet, &building.Building, &building.Level, &building.CreatedAt, &building.UpdatedAt, &building.Version)
		if err != nil {
			return err
		}

		out = append(out, building)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.PlanetBuilding{}, err
	}

	return out, nil
}

const deletePlanetBuildingSqlTemplate = "DELETE FROM planet_building WHERE planet = $1"

func (r *planetBuildingRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetBuildingSqlTemplate, planet)
	return err
}
