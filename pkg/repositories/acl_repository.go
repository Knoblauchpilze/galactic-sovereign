package repositories

import (
	"context"
	"fmt"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type AclRepository interface {
	Create(ctx context.Context, tx db.Transaction, acl persistence.Acl) (persistence.Acl, error)
	Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Acl, error)
	GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error)
	Delete(ctx context.Context, tx db.Transaction, ids []uuid.UUID) error
	DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error
}

type aclRepositoryImpl struct{}

func NewAclRepository() AclRepository {
	return &aclRepositoryImpl{}
}

const createAclSqlTemplate = `
INSERT INTO acl (id, api_user, resource)
	VALUES($1, $2, $3)
	ON CONFLICT (api_user, resource) DO NOTHING
	RETURNING
		acl.id
`
const createAclPermissionsSqlTemplate = `
INSERT INTO acl_permission (acl, permission)
	VALUES($1, $2)
`

func (r *aclRepositoryImpl) Create(ctx context.Context, tx db.Transaction, acl persistence.Acl) (persistence.Acl, error) {
	res := tx.Query(ctx, createAclSqlTemplate, acl.Id, acl.User, acl.Resource)
	if err := res.Err(); err != nil {
		return persistence.Acl{}, err
	}

	parser := func(rows db.Scannable) error {
		return rows.Scan(&acl.Id)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.Acl{}, err
	}

	for _, permission := range acl.Permissions {
		res = tx.Query(ctx, createAclPermissionsSqlTemplate, permission)
		if err := res.Err(); err != nil {
			return persistence.Acl{}, err
		}
	}

	return acl, nil
}

const getAclSqlTemplate = `
SELECT
	id,
	api_user,
	resource,
	created_at,
	updated_at
FROM
	acl
WHERE
	id = $1
`
const getAclPermissionsSqlTemplate = `SELECT permission FROM acl_permission WHERE acl = $1`

func (r *aclRepositoryImpl) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Acl, error) {
	res := tx.Query(ctx, getAclSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.Acl{}, err
	}

	var out persistence.Acl
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.User, &out.Resource, &out.CreatedAt, &out.UpdatedAt)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.Acl{}, err
	}

	res = tx.Query(ctx, getAclPermissionsSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.Acl{}, err
	}

	parser = func(rows db.Scannable) error {
		var permission string
		err := rows.Scan(&permission)
		if err == nil {
			out.Permissions = append(out.Permissions, permission)
		}
		return err
	}

	if err := res.GetAll(parser); err != nil {
		return persistence.Acl{}, err
	}

	return out, nil
}

const getAclForUserSqlTemplate = "SELECT id FROM acl WHERE api_user = $1"

func (r *aclRepositoryImpl) GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	res := tx.Query(ctx, getAclForUserSqlTemplate, user)
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

const deleteAclSqlTemplate = "DELETE FROM acl WHERE id IN (%s)"
const deleteAclPermissionsSqlTemplate = "DELETE FROM acl_permission WHERE acl in (%s)"

func (r *aclRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, ids []uuid.UUID) error {
	in := db.ToSliceInterface(ids)
	permissionsSqlQuery := fmt.Sprintf(deleteAclPermissionsSqlTemplate, db.GenerateInClauseForArgs(len(ids)))
	_, err := tx.Exec(ctx, permissionsSqlQuery, in...)
	if err != nil {
		return err
	}

	aclSqlQuery := fmt.Sprintf(deleteAclSqlTemplate, db.GenerateInClauseForArgs(len(ids)))

	_, err = tx.Exec(ctx, aclSqlQuery, in...)
	return err
}

const deleteAclForUserSqlTemplate = "DELETE FROM acl WHERE api_user = $1"
const deleteAclPermissionsForUserSqlTemplate = `
DELETE FROM acl_permission
	WHERE acl
		IN (SELECT id FROM acl WHERE api_user = $1)
`

func (r *aclRepositoryImpl) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteAclPermissionsForUserSqlTemplate, user)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, deleteAclForUserSqlTemplate, user)
	return err
}
