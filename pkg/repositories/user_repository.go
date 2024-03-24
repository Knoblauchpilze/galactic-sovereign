package repositories

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user persistence.User) error
	Get(ctx context.Context, id uuid.UUID) (persistence.User, error)
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

const sqlCreateUserTemplate = "INSERT INTO api_user (id, email, password, created_at) VALUES($1, $2, $3, $4)"

func (r *userRepositoryImpl) Create(ctx context.Context, user persistence.User) error {
	_, err := r.conn.Exec(ctx, sqlCreateUserTemplate, user.Id, user.Email, user.Password, user.CreatedAt)
	return err
}

const sqlQueryUserTemplate = "SELECT id, email, password, created_at, updated_at FROM api_user WHERE id = $1"

func (r *userRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.User, error) {
	res := r.conn.Query(ctx, sqlQueryUserTemplate, id)
	if err := res.Err(); err != nil {
		return persistence.User{}, err
	}

	var out persistence.User
	parser := func(rows db.Scannable) error {
		return rows.Scan(&out.Id, &out.Email, &out.Password, &out.CreatedAt, &out.UpdatedAt)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return persistence.User{}, err
	}

	return out, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user persistence.User) (persistence.User, error) {
	return persistence.User{}, errors.NewCode(errors.NotImplementedCode)
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.NewCode(errors.NotImplementedCode)
}
