package repositories

import (
	"context"
	"fmt"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error)
	Get(ctx context.Context, id uuid.UUID) (persistence.ApiKey, error)
	GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error)
	GetForUser(ctx context.Context, user uuid.UUID) ([]uuid.UUID, error)
	GetForUserTx(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error)
	DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error
}

type apiKeyRepositoryImpl struct {
	conn db.ConnectionPool
}

func NewApiKeyRepository(conn db.ConnectionPool) ApiKeyRepository {
	return &apiKeyRepositoryImpl{
		conn: conn,
	}
}

const createApiKeySqlTemplate = `
INSERT INTO api_key (id, key, api_user, valid_until)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (api_user) DO UPDATE
	SET
		valid_until = excluded.valid_until
	WHERE
		api_key.api_user = excluded.api_user
	RETURNING
		api_key.key
`

func (r *apiKeyRepositoryImpl) Create(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error) {
	res := r.conn.Query(ctx, createApiKeySqlTemplate, apiKey.Id, apiKey.Key, apiKey.ApiUser, apiKey.ValidUntil)
	if err := res.Err(); err != nil {
		return persistence.ApiKey{}, err
	}

	parser := func(rows db.Scannable) error {
		return rows.Scan(&apiKey.Key)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.ApiKey{}, err
	}

	return apiKey, nil
}

const getApiKeySqlTemplate = "SELECT id, key, api_user, valid_until FROM api_key WHERE id = $1"

func (r *apiKeyRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.ApiKey, error) {
	res := r.conn.Query(ctx, getApiKeySqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.ApiKey{}, err
	}

	var out persistence.ApiKey
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Key, &out.ApiUser, &out.ValidUntil)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.ApiKey{}, err
	}

	return out, nil
}

const getApiKeyForKeySqlTemplate = "SELECT id, key, api_user, valid_until FROM api_key WHERE key = $1"

func (r *apiKeyRepositoryImpl) GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error) {
	res := r.conn.Query(ctx, getApiKeyForKeySqlTemplate, apiKey)
	if err := res.Err(); err != nil {
		return persistence.ApiKey{}, err
	}

	var out persistence.ApiKey
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Key, &out.ApiUser, &out.ValidUntil)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.ApiKey{}, err
	}

	return out, nil
}

const getApiKeyForUserSqlTemplate = "SELECT id FROM api_key WHERE api_user = $1"

func (r *apiKeyRepositoryImpl) GetForUser(ctx context.Context, user uuid.UUID) ([]uuid.UUID, error) {
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

func (r *apiKeyRepositoryImpl) GetForUserTx(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	res := tx.Query(ctx, getApiKeyForUserSqlTemplate, user)
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

const deleteApiKeysSqlTemplate = "DELETE FROM api_key WHERE id IN (%s)"

func (r *apiKeyRepositoryImpl) Delete(ctx context.Context, ids []uuid.UUID) error {
	in := db.ToSliceInterface(ids)
	sqlQuery := fmt.Sprintf(deleteApiKeysSqlTemplate, db.GenerateInClauseForArgs(len(ids)))

	_, err := r.conn.Exec(ctx, sqlQuery, in...)
	return err
}

const deleteApiKeyForUserSqlTemplate = `DELETE FROM api_key WHERE api_user = $1`

func (r *apiKeyRepositoryImpl) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteApiKeyForUserSqlTemplate, user)
	if err != nil {
		return err
	}

	return nil
}
