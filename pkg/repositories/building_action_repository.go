package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionRepository interface {
	ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.BuildingAction, error)
	DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error
}

type buildingActionRepositoryImpl struct{}

func NewBuildingActionRepository() BuildingActionRepository {
	return &buildingActionRepositoryImpl{}
}

const listBuildingActionForPlanetSqlTemplate = `
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
	planet = $1`

func (r *buildingActionRepositoryImpl) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.BuildingAction, error) {
	res := tx.Query(ctx, listBuildingActionForPlanetSqlTemplate, planet)
	if err := res.Err(); err != nil {
		return []persistence.BuildingAction{}, err
	}

	var out []persistence.BuildingAction
	parser := func(rows db.Scannable) error {
		var action persistence.BuildingAction
		err := rows.Scan(&action.Id, &action.Planet, &action.Building, &action.CurrentLevel, &action.DesiredLevel, &action.CreatedAt, &action.CompletedAt)
		if err != nil {
			return err
		}

		out = append(out, action)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.BuildingAction{}, err
	}

	return out, nil
}

const deleteBuildingActionSqlTemplate = "DELETE FROM building_action WHERE planet = $1"

func (r *buildingActionRepositoryImpl) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteBuildingActionSqlTemplate, planet)
	return err
}
