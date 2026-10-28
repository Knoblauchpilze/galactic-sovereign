package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
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

func saveAndReloadPlanet(
	ctx context.Context,
	tx database.Transaction,
	planet models.Planet,
	expectedVersion int,
) (models.Planet, error) {
	if planet.Version == expectedVersion {
		// Voluntarily ignoring error to prioritize the domain error
		// nolint:errcheck
		tx.Rollback()
		return models.Planet{}, domainerrors.ErrMutationWithoutVersionBump
	}

	err := updatePlanetDetails(ctx, tx, planet, expectedVersion)
	if err != nil {
		// Voluntarily ignoring error to prioritize the domain error
		// nolint:errcheck
		tx.Rollback()
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
