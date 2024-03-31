package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey persistence.ApiKey) error
	Get(ctx context.Context, id uuid.UUID) (persistence.ApiKey, error)
	List(ctx context.Context) ([]uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type apiUserRepositoryImpl struct {
	conn db.Connection
}

func NewApiKeyRepository(conn db.Connection) ApiKeyRepository {
	return &apiUserRepositoryImpl{
		conn: conn,
	}
}

const createApiKeySqlTemplate = "INSERT INTO api_key (id, key, api_user) VALUES($1, $2, $3)"

func (r *apiUserRepositoryImpl) Create(ctx context.Context, apiKey persistence.ApiKey) error {
	_, err := r.conn.Exec(ctx, createApiKeySqlTemplate, apiKey.Id, apiKey.Key, apiKey.ApiUser)
	return err
}

const getApiKeySqlTemplate = "SELECT id, key, api_user FROM api_key WHERE id = $1"

func (r *apiUserRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.ApiKey, error) {
	res := r.conn.Query(ctx, getApiKeySqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.ApiKey{}, err
	}

	var out persistence.ApiKey
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Key, &out.ApiUser)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.ApiKey{}, err
	}

	return out, nil
}

const listApiKeySqlTemplate = "SELECT id FROM api_key"

func (r *apiUserRepositoryImpl) List(ctx context.Context) ([]uuid.UUID, error) {
	res := r.conn.Query(ctx, listApiKeySqlTemplate)
	if err := res.Err(); err != nil {
		return []uuid.UUID{}, err
	}

	var out []uuid.UUID
	parser := func(rows db.Scannable) error {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return err
		}

		out = append(out, id)
		return nil
	}

	if err := res.GetAll(parser); err != nil {
		return []uuid.UUID{}, err
	}

	return out, nil
}

const deleteApiKeySqlTemplate = "DELETE FROM api_key WHERE id = $1"

func (r *apiUserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	affected, err := r.conn.Exec(ctx, deleteApiKeySqlTemplate, id)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.NewCode(db.NoMatchingSqlRows)
	}
	return nil
}
