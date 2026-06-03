package driven

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

const (
	createBuildingActionQuery = `
INSERT INTO
	building_action (id, planet, building, current_level, desired_level, created_at, completed_at, version)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
`

	getBuildingActionQuery = `
SELECT
	id,
	planet,
	building,
	current_level,
	desired_level,
	created_at,
	completed_at,
	version
FROM
	building_action
WHERE
	id = $1`

	listBuildingActionForPlanetQuery = `
SELECT
	id,
	planet,
	building,
	current_level,
	desired_level,
	created_at,
	completed_at,
	version
FROM
	building_action
WHERE
	planet = $1`

	listBuildingActionBeforeCompletionTimeQuery = `
SELECT
	id,
	planet,
	building,
	current_level,
	desired_level,
	created_at,
	completed_at,
	version
FROM
	building_action
WHERE
	completed_at <= $1
	AND planet = $2`

	deleteBuildingActionQuery = `DELETE FROM building_action WHERE id = $1`

	deleteBuildingActionForPlanetQuery = `DELETE FROM building_action WHERE planet = $1`

	deleteBuildingActionForPlayerQuery = `
DELETE FROM
	building_action AS bad
USING
	building_action AS ba
	LEFT JOIN planet AS p ON p.id = ba.planet
WHERE
	bad.id = ba.id
	AND p.player = $1`
)

type buildingActionRepositoryImpl struct {
	conn db.Connection
}

func NewBuildingActionRepository(conn db.Connection) driven.ForManagingBuildingActions {
	return &buildingActionRepositoryImpl{
		conn: conn,
	}
}

func (r *buildingActionRepositoryImpl) Create(
	ctx context.Context,
	action models.BuildingAction,
) error {
	_, err := r.conn.Exec(
		ctx,
		createBuildingActionQuery,
		action.Id,
		action.Planet,
		action.Building,
		action.CurrentLevel,
		action.DesiredLevel,
		action.CreatedAt,
		action.CompletedAt,
		action.Version,
	)
	return err
}

func (r *buildingActionRepositoryImpl) Get(
	ctx context.Context,
	id uuid.UUID,
) (models.BuildingAction, error) {
	return db.QueryOne[models.BuildingAction](
		ctx,
		r.conn,
		getBuildingActionQuery,
		id,
	)
}

func (r *buildingActionRepositoryImpl) ListForPlanet(
	ctx context.Context,
	planet uuid.UUID,
) ([]models.BuildingAction, error) {
	return db.QueryAll[models.BuildingAction](
		ctx,
		r.conn,
		listBuildingActionForPlanetQuery,
		planet,
	)
}

func (r *buildingActionRepositoryImpl) ListBeforeCompletionTime(
	ctx context.Context,
	planet uuid.UUID,
	until time.Time,
) ([]models.BuildingAction, error) {
	return db.QueryAll[models.BuildingAction](
		ctx,
		r.conn,
		listBuildingActionBeforeCompletionTimeQuery,
		until,
		planet,
	)
}

func (r *buildingActionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deleteBuildingActionQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *buildingActionRepositoryImpl) DeleteForPlanet(ctx context.Context, planet uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deleteBuildingActionForPlanetQuery, planet)
	if err != nil {
		return err
	}

	return nil
}

func (r *buildingActionRepositoryImpl) DeleteForPlayer(ctx context.Context, player uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	_, err = tx.Exec(ctx, deleteBuildingActionForPlayerQuery, player)
	if err != nil {
		return err
	}

	return nil
}
