package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
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
	return db.QueryOneTx[persistence.PlanetBuilding](ctx, tx, getPlanetBuildingForPlanetAndBuildingSqlTemplate, planet, building)
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
	planet = $1`

func (r *planetBuildingRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetBuilding, error) {
	return db.QueryAllTx[persistence.PlanetBuilding](ctx, tx, listPlanetBuildingForPlanetSqlTemplate, planet)
}

const updatePlanetBuildingSqlTemplate = `
UPDATE
	planet_building
SET
	level = $1,
	updated_at = $2,
	version = $3
WHERE
	planet = $4
	AND building = $5
	AND version = $6`

func (r *planetBuildingRepositoryImpl) Update(ctx context.Context, tx db.Transaction, building persistence.PlanetBuilding) (persistence.PlanetBuilding, error) {
	version := building.Version + 1
	affectedRows, err := tx.Exec(
		ctx,
		updatePlanetBuildingSqlTemplate,
		building.Level,
		building.UpdatedAt,
		version,
		building.Planet,
		building.Building,
		building.Version,
	)
	if err != nil {
		return building, err
	}
	if affectedRows != 1 {
		return building, errors.NewCode(OptimisticLockException)
	}

	building.Version = version

	return building, nil
}

const deletePlanetBuildingSqlTemplate = `DELETE FROM planet_building WHERE planet = $1`

func (r *planetBuildingRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePlanetBuildingSqlTemplate, planet)
	return err
}
