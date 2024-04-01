package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user persistence.User) (persistence.User, error)
	Get(ctx context.Context, id uuid.UUID) (persistence.User, error)
	List(ctx context.Context) ([]uuid.UUID, error)
	Update(ctx context.Context, user persistence.User) (persistence.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepositoryImpl struct {
	conn db.Connection
}

func NewUserRepository(conn db.Connection) UserRepository {
	return &userRepositoryImpl{
		conn: conn,
	}
}

const createUserSqlTemplate = "INSERT INTO api_user (id, email, password, created_at) VALUES($1, $2, $3, $4)"

func (r *userRepositoryImpl) Create(ctx context.Context, user persistence.User) (persistence.User, error) {
	_, err := r.conn.Exec(ctx, createUserSqlTemplate, user.Id, user.Email, user.Password, user.CreatedAt)
	return user, err
}

const getUserSqlTemplate = "SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE id = $1"

func (r *userRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.User, error) {
	res := r.conn.Query(ctx, getUserSqlTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.User{}, err
	}

	var out persistence.User
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Email, &out.Password, &out.CreatedAt, &out.UpdatedAt, &out.Version)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.User{}, err
	}

	return out, nil
}

const listUserSqlTemplate = "SELECT id FROM api_user"

func (r *userRepositoryImpl) List(ctx context.Context) ([]uuid.UUID, error) {
	res := r.conn.Query(ctx, listUserSqlTemplate)
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

const updateUserSqlTemplate = "UPDATE api_user SET email = $1, password = $2, version = $3 WHERE id = $4 AND version = $5"

func (r *userRepositoryImpl) Update(ctx context.Context, user persistence.User) (persistence.User, error) {
	version := user.Version + 1
	affected, err := r.conn.Exec(ctx, updateUserSqlTemplate, user.Email, user.Password, version, user.Id, user.Version)
	if err != nil {
		return user, err
	}
	if affected != 1 {
		return user, errors.NewCode(db.OptimisticLockException)
	}

	user.Version = version

	return user, nil
}

const deleteUserSqlTemplate = "DELETE FROM api_user WHERE id = $1"

func (r *userRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	affected, err := r.conn.Exec(ctx, deleteUserSqlTemplate, id)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.NewCode(db.NoMatchingSqlRows)
	}
	return nil
}
