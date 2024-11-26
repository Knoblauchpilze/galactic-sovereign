package repositories

import (
	"context"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type ResourceRepository interface {
	Create(ctx context.Context, resource persistence.Resource) (persistence.Resource, error)
	Get(ctx context.Context, id uuid.UUID) (persistence.Resource, error)
	List(ctx context.Context, tx db.Transaction) ([]persistence.Resource, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type resourceRepositoryImpl struct {
	conn db.ConnectionPool
}

func NewResourceRepository(conn db.ConnectionPool) ResourceRepository {
	return &resourceRepositoryImpl{
		conn: conn,
	}
}

const createResourceSqlTemplate = "INSERT INTO resource (id, name, created_at) VALUES($1, $2, $3)"

func (r *resourceRepositoryImpl) Create(ctx context.Context, resource persistence.Resource) (persistence.Resource, error) {
	_, err := r.conn.Exec(ctx, createResourceSqlTemplate, resource.Id, resource.Name, resource.CreatedAt)
	if err != nil && duplicatedKeySqlErrorRegexp.MatchString(err.Error()) {
		return persistence.Resource{}, errors.NewCode(db.DuplicatedKeySqlKey)
	}

	return resource, err
}

const getResourceSqlTemplate = "SELECT id, name, created_at, updated_at FROM resource WHERE id = $1"

func (r *resourceRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.Resource, error) {
	res := r.conn.Query(ctx, getResourceSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.Resource{}, err
	}

	var out persistence.Resource
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Name, &out.CreatedAt, &out.UpdatedAt)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.Resource{}, err
	}

	return out, nil
}

const listResourceSqlTemplate = "SELECT id, name, created_at, updated_at FROM resource"

func (r *resourceRepositoryImpl) List(ctx context.Context, tx db.Transaction) ([]persistence.Resource, error) {
	res := tx.Query(ctx, listResourceSqlTemplate)
	if err := res.Err(); err != nil {
		return []persistence.Resource{}, err
	}

	var out []persistence.Resource
	parser := func(rows db.Scannable) error {
		var resource persistence.Resource
		err := rows.Scan(&resource.Id, &resource.Name, &resource.CreatedAt, &resource.UpdatedAt)
		if err != nil {
			return err
		}

		out = append(out, resource)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []persistence.Resource{}, err
	}

	return out, nil
}

const deleteResourceSqlTemplate = "DELETE FROM resource WHERE id = $1"

func (r *resourceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	affected, err := r.conn.Exec(ctx, deleteResourceSqlTemplate, id)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.NewCode(db.NoMatchingSqlRows)
	}
	return nil
}
