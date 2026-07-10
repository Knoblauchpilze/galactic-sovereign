package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

const (
	lockPlanetForUpdateQuery = `
SELECT
	id
FROM
	planet
WHERE
	id = $1
FOR UPDATE
`
)

type PlanetMutator struct {
	conn db.Connection
}

func NewPlanetMutator(conn db.Connection) *PlanetMutator {
	return &PlanetMutator{
		conn: conn,
	}
}

func (m *PlanetMutator) Mutate(
	ctx context.Context,
	id uuid.UUID,
	mutator drivenports.PlanetMutator,
) (models.PlanetMutationResult, error) {
	tx, err := m.conn.BeginTx(ctx)
	if err != nil {
		return models.PlanetMutationResult{}, err
	}
	defer tx.Close(ctx)

	actual, err := db.QueryOneTx[uuid.UUID](ctx, tx, lockPlanetForUpdateQuery, id)
	if err != nil {
		return models.PlanetMutationResult{}, parseDbError(err)
	}
	if actual != id {
		return models.PlanetMutationResult{}, domainerrors.ErrNotFound
	}

	planet, err := loadPlanetAndDetails(ctx, tx, id)
	if err != nil {
		return models.PlanetMutationResult{}, err
	}

	expectedVersion := planet.Version

	deleted, err := mutator(&planet)
	if err != nil {
		// There's no point in checking the error here: it is not logged
		// and there's already an error pending.
		// nolint:errcheck
		tx.Rollback()

		return models.PlanetMutationResult{}, err
	}

	out := models.PlanetMutationResult{Deleted: deleted}

	if deleted {
		err = deletePlanetAndDetails(ctx, tx, id)
		return out, err
	}

	out.Planet, err = saveAndReloadPlanet(ctx, tx, planet, expectedVersion)
	if err != nil {
		return out, err
	}

	return out, nil
}

func saveAndReloadPlanet(
	ctx context.Context,
	tx db.Transaction,
	planet models.Planet,
	expectedVersion int,
) (models.Planet, error) {
	if planet.Version == expectedVersion {
		return models.Planet{}, domainerrors.ErrMutationWithoutVersionBump
	}

	err := updatePlanetDetails(ctx, tx, planet, expectedVersion)
	if err != nil {
		return models.Planet{}, err
	}

	// This second load is to make sure that the returned planet value
	// corresponds to the up to date data stored in the DB. There is
	// nothing preventing the mutation function to perform updates to
	// the planet which are not reflected by the SQL queries (such as
	// update to the identifier, etc.).
	// A better solution would be to constrain the shape of the mutation
	// function so that it can only perform valid modifications but this
	// solution (reload the data form the DB) is acceptable.
	out, err := loadPlanetAndDetails(ctx, tx, planet.Id)
	if err != nil {
		return models.Planet{}, err
	}

	return out, nil
}
