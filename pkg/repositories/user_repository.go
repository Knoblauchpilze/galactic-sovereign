package repositories

import (
	"fmt"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user persistence.User) error
	Get(id uuid.UUID) (persistence.User, error)
	Update(user persistence.User) (persistence.User, error)
	Delete(id uuid.UUID) error
}

type userRepositoryImpl struct {
	conn db.Connection
}

func NewUserRepository(conn db.Connection) UserRepository {
	return &userRepositoryImpl{
		conn: conn,
	}
}

func (r *userRepositoryImpl) Create(user persistence.User) error {
	return errors.NewCode(errors.NotImplementedCode)
}

const sqlQueryUserTemplate = "select id, email, password, created_at, updated_at from api_user where id = '%s'"

func (r *userRepositoryImpl) Get(id uuid.UUID) (persistence.User, error) {
	query := fmt.Sprintf(sqlQueryUserTemplate, id)
	res := r.conn.Query(query)
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

func (r *userRepositoryImpl) Update(user persistence.User) (persistence.User, error) {
	return persistence.User{}, errors.NewCode(errors.NotImplementedCode)
}

func (r *userRepositoryImpl) Delete(id uuid.UUID) error {
	return errors.NewCode(errors.NotImplementedCode)
}
