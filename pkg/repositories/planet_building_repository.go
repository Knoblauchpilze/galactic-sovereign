package repositories

import (
	"context"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetBuildingRepository interface {
	GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building uuid.UUID) (persistence.PlanetBuilding, error)
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetBuilding, error)
	Update(ctx context.Context, tx db.Transaction, building persistence.PlanetBuilding) (persistence.PlanetBuilding, error)
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
}

type planetBuildingRepositoryImpl struct{}

func NewPlanetBuildingRepository() PlanetBuildingRepository {
	return &planetBuildingRepositoryImpl{}
}

const getPlanetBuildingForPlanetAndBuildingSqlTemplate = `
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
	AND building = $2`

func (r *planetBuildingRepositoryImpl) GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building uuid.UUID) (persistence.PlanetBuilding, error) {
	res := tx.Query(ctx, getPlanetBuildingForPlanetAndBuildingSqlTemplate, planet, building)
	if err := res.Err(); err != nil {
		return persistence.PlanetBuilding{}, err
	}

	var out persistence.PlanetBuilding
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Planet, &out.Building, &out.Level, &out.CreatedAt, &out.UpdatedAt, &out.Version)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.PlanetBuilding{}, err
	}

	return out, nil
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

const updatePlanetBuildingSqlTemplate = `
UPDATE
	planet_building
SET
	level = $1,
	version = $2
WHERE
	planet = $3
	AND building = $4
	AND version = $5
`

func (r *planetBuildingRepositoryImpl) Update(ctx context.Context, tx db.Transaction, building persistence.PlanetBuilding) (persistence.PlanetBuilding, error) {
	version := building.Version + 1
	affected, err := tx.Exec(ctx, updatePlanetBuildingSqlTemplate, building.Level, version, building.Planet, building.Building, building.Version)
	if err != nil {
		return building, err
	}
	if affected == 0 {
		return building, errors.NewCode(db.OptimisticLockException)
	} else if affected != 1 {
		return building, errors.NewCode(db.MoreThanOneMatchingSqlRows)
	}

	building.Version = version

	return building, nil
}

const deletePlanetBuildingSqlTemplate = "DELETE FROM planet_building WHERE planet = $1"

func (r *planetBuildingRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetBuildingSqlTemplate, planet)
	return err
}
