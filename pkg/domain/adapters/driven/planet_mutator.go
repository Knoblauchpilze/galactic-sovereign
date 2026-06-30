package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
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
) (models.Planet, error) {
	tx, err := m.conn.BeginTx(ctx)
	if err != nil {
		return models.Planet{}, err
	}
	defer tx.Close(ctx)

	planet, err := loadPlanetAndDetails(ctx, tx, id)
	if err != nil {
		return models.Planet{}, err
	}

	err = mutator(&planet, tx.TimeStamp())
	if err != nil {
		// There's no point in checking the error here: it is not logged
		// and there's already an error pending.
		// nolint:errcheck
		tx.Rollback()

		return models.Planet{}, err
	}

	err = updatePlanetDetails(ctx, tx, planet)
	if err != nil {
		return models.Planet{}, err
	}

	return planet, nil
}
