package repositories

import (
	"context"
	"fmt"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UserLimitRepository interface {
	Create(ctx context.Context, tx db.Transaction, userLimit persistence.UserLimit) (persistence.UserLimit, error)
	Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.UserLimit, error)
	GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error)
	Delete(ctx context.Context, tx db.Transaction, ids []uuid.UUID) error
	DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error
}

type userLimitRepositoryImpl struct{}

func NewUserLimitRepository() UserLimitRepository {
	return &userLimitRepositoryImpl{}
}

const createUserLimitSqlTemplate = `
INSERT INTO user_limit (id, name, api_user)
	VALUES($1, $2, $3)
	ON CONFLICT (name, api_user) DO NOTHING
	RETURNING
		user_limit.id
`
const createLimitsSqlTemplate = `
INSERT INTO limits (id, name, value, user_limit)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (name, user_limit) DO UPDATE
	SET
		value = excluded.value
	WHERE
		limits.name = excluded.name
		AND limits.user_limit = excluded.user_limit
	RETURNING
		limits.id
`

func (r *userLimitRepositoryImpl) Create(ctx context.Context, tx db.Transaction, userLimit persistence.UserLimit) (persistence.UserLimit, error) {
	res := tx.Query(ctx, createUserLimitSqlTemplate, userLimit.Id, userLimit.Name, userLimit.User)
	if err := res.Err(); err != nil {
		return persistence.UserLimit{}, err
	}

	parser := func(rows db.Scannable) error {
		return rows.Scan(&userLimit.Id)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.UserLimit{}, err
	}

	var newLimits []persistence.Limit

	for _, limit := range userLimit.Limits {
		res = tx.Query(ctx, createLimitsSqlTemplate, limit.Id, limit.Name, limit.Value, userLimit.Id)
		if err := res.Err(); err != nil {
			return persistence.UserLimit{}, err
		}

		newLimit := limit
		parser = func(rows db.Scannable) error {
			return rows.Scan(&newLimit.Id)
		}

		if err := res.GetSingleValue(parser); err != nil {
			return persistence.UserLimit{}, err
		}

		newLimits = append(newLimits, newLimit)
	}

	userLimit.Limits = newLimits

	return userLimit, nil
}

const getUserLimitSqlTemplate = `
SELECT
	id,
	name,
	api_user,
	created_at,
	updated_at
FROM
	user_limit
WHERE
	id = $1
`
const getLimitsSqlTemplate = `
SELECT
	id,
	name,
	value,
	created_at,
	updated_at
FROM
	limits
WHERE
	user_limit = $1
`

func (r *userLimitRepositoryImpl) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.UserLimit, error) {
	res := tx.Query(ctx, getUserLimitSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.UserLimit{}, err
	}

	var out persistence.UserLimit
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Name, &out.User, &out.CreatedAt, &out.UpdatedAt)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.UserLimit{}, err
	}

	res = tx.Query(ctx, getLimitsSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.UserLimit{}, err
	}

	parser = func(rows db.Scannable) error {
		var limit persistence.Limit
		err := rows.Scan(&limit.Id, &limit.Name, &limit.Value, &limit.CreatedAt, &limit.UpdatedAt)
		if err == nil {
			out.Limits = append(out.Limits, limit)
		}
		return err
	}

	if err := res.GetAll(parser); err != nil {
		return persistence.UserLimit{}, err
	}

	return out, nil
}

const getUserLimitForUserSqlTemplate = "SELECT id FROM user_limit WHERE api_user = $1"

func (r *userLimitRepositoryImpl) GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	res := tx.Query(ctx, getUserLimitForUserSqlTemplate, user)
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

const deleteUserLimitSqlTemplate = "DELETE FROM user_limit WHERE id IN (%s)"
const deleteLimitsSqlTemplate = "DELETE FROM limits WHERE user_limit in (%s)"

func (r *userLimitRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, ids []uuid.UUID) error {
	in := db.ToSliceInterface(ids)
	limitsSqlQuery := fmt.Sprintf(deleteLimitsSqlTemplate, db.GenerateInClauseForArgs(len(ids)))
	_, err := tx.Exec(ctx, limitsSqlQuery, in...)
	if err != nil {
		return err
	}

	userLimitSqlQuery := fmt.Sprintf(deleteUserLimitSqlTemplate, db.GenerateInClauseForArgs(len(ids)))

	_, err = tx.Exec(ctx, userLimitSqlQuery, in...)
	return err
}

const deleteUserLimitForUserSqlTemplate = "DELETE FROM user_limit WHERE api_user = $1"
const deleteLimitsForUserSqlTemplate = `
DELETE FROM limits
	WHERE user_limit
		IN (SELECT id FROM user_limit WHERE api_user = $1)
`

func (r *userLimitRepositoryImpl) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteLimitsForUserSqlTemplate, user)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteUserLimitForUserSqlTemplate, user)
	return err
}
