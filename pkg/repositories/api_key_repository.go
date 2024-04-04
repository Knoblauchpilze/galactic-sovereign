package repositories

import (
	"context"
	"fmt"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, tx db.Transaction, apiKey persistence.ApiKey) (persistence.ApiKey, error)
	Get(ctx context.Context, id uuid.UUID) (persistence.ApiKey, error)
	GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error)
	List(ctx context.Context) ([]uuid.UUID, error)
	Update(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error)
	Delete(ctx context.Context, tx db.Transaction, ids []uuid.UUID) error
}

type apiUserRepositoryImpl struct {
	conn db.ConnectionPool
}

func NewApiKeyRepository(conn db.ConnectionPool) ApiKeyRepository {
	return &apiUserRepositoryImpl{
		conn: conn,
	}
}

const createApiKeySqlTemplate = "INSERT INTO api_key (id, key, api_user) VALUES($1, $2, $3)"

func (r *apiUserRepositoryImpl) Create(ctx context.Context, tx db.Transaction, apiKey persistence.ApiKey) (persistence.ApiKey, error) {
	_, err := tx.Exec(ctx, createApiKeySqlTemplate, apiKey.Id, apiKey.Key, apiKey.ApiUser)
	return apiKey, err
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

const getApiKeyForUserSqlTemplate = "SELECT id FROM api_key WHERE api_user = $1"

func (r *apiUserRepositoryImpl) GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	res := r.conn.Query(ctx, getApiKeyForUserSqlTemplate, user)
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

const updateApiKeySqlTemplate = "UPDATE api_key SET enabled = $1, version = $2 WHERE id = $3 AND version = $4"

func (r *apiUserRepositoryImpl) Update(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error) {
	version := apiKey.Version + 1
	affected, err := r.conn.Exec(ctx, updateApiKeySqlTemplate, apiKey.Enabled, version, apiKey.Id, apiKey.Version)
	if err != nil {
		return apiKey, err
	}
	if affected == 0 {
		return apiKey, errors.NewCode(db.OptimisticLockException)
	} else if affected != 1 {
		return apiKey, errors.NewCode(db.MoreThanOneMatchingSqlRows)
	}

	apiKey.Version = version

	return apiKey, nil
}

const deleteApiKeysSqlTemplate = "DELETE FROM api_key WHERE id IN (%s)"

func (r *apiUserRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, ids []uuid.UUID) error {
	in := db.ToSliceInterface(ids)
	sqlQuery := fmt.Sprintf(deleteApiKeysSqlTemplate, db.GenerateInClauseForArgs(len(ids)))

	_, err := tx.Exec(ctx, sqlQuery, in...)
	return err
}
