package drivenadapters

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
)

type PlanetDeleter struct {
	conn database.Connection
}

func NewPlanetDeleter(conn database.Connection) *PlanetDeleter {
	return &PlanetDeleter{
		conn: conn,
	}
}

func (m *PlanetDeleter) Delete(
	ctx context.Context,
	id uuid.UUID,
	deleter drivenports.PlanetDeleter,
) error {
	tx, err := m.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	actual, err := db.QueryOneTx[uuid.UUID](ctx, tx, lockPlanetForUpdateQuery, id)
	if err != nil {
		return parseDbError(err)
	}
	if actual != id {
		return domainerrors.ErrNotFound
	}

	planet, err := loadPlanetAndDetails(ctx, tx, id)
	if err != nil {
		return err
	}

	err = deleter(&planet)
	if err != nil {
		// There's no point in checking the error here: it is not logged
		// and there's already an error pending.
		// nolint:errcheck
		tx.Rollback()

		return err
	}

	return deletePlanetAndDetails(ctx, tx, id)
}
