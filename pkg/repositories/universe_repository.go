package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type UniverseRepository interface {
	Create(ctx context.Context, universe persistence.Universe) (persistence.Universe, error)
	Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Universe, error)
	List(ctx context.Context) ([]persistence.Universe, error)
	Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error
}

type universeRepositoryImpl struct {
	conn db.ConnectionPool
}

func NewUniverseRepository(conn db.ConnectionPool) UniverseRepository {
	return &universeRepositoryImpl{
		conn: conn,
	}
}

const createUniverseSqlTemplate = "INSERT INTO universe (id, name, created_at) VALUES($1, $2, $3)"

func (r *universeRepositoryImpl) Create(ctx context.Context, universe persistence.Universe) (persistence.Universe, error) {
	_, err := r.conn.Exec(ctx, createUniverseSqlTemplate, universe.Id, universe.Name, universe.CreatedAt)
	if err != nil && duplicatedKeySqlErrorRegexp.MatchString(err.Error()) {
		return persistence.Universe{}, errors.NewCode(db.DuplicatedKeySqlKey)
	}

	return universe, err
}

const getUniverseSqlTemplate = "SELECT id, name, created_at, updated_at, version FROM universe WHERE id = $1"

func (r *universeRepositoryImpl) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Universe, error) {
	res := tx.Query(ctx, getUniverseSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.Universe{}, err
	}

	var out persistence.Universe
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Name, &out.CreatedAt, &out.UpdatedAt, &out.Version)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.Universe{}, err
	}

	return out, nil
}

const listUniverseSqlTemplate = "SELECT id, name, created_at, updated_at, version FROM universe"

func (r *universeRepositoryImpl) List(ctx context.Context) ([]persistence.Universe, error) {
	res := r.conn.Query(ctx, listUniverseSqlTemplate)
	if err := res.Err(); err != nil {
		return []persistence.Universe{}, err
	}

	var out []persistence.Universe
	parser := func(rows db.Scannable) error {
		var universe persistence.Universe
		err := rows.Scan(&universe.Id, &universe.Name, &universe.CreatedAt, &universe.UpdatedAt, &universe.Version)
		if err != nil {
			return err
		}

		out = append(out, universe)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.Universe{}, err
	}

	return out, nil
}

const deleteUniverseSqlTemplate = "DELETE FROM universe WHERE id = $1"

func (r *universeRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	affected, err := tx.Exec(ctx, deleteUniverseSqlTemplate, id)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.NewCode(db.NoMatchingSqlRows)
	}
	return nil
}
