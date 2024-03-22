package repositories

import (
	"fmt"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
)

type User struct {
	id       uuid.UUID
	email    string
	password string

	createdAt  time.Time
	updated_at time.Time
}

type UserRepository interface {
	Create(user User) error
	Get(id uuid.UUID) (User, error)
	Update(user User) (User, error)
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

func (r *userRepositoryImpl) Create(user User) error {
	return errors.NewCode(errors.NotImplementedCode)
}

const sqlQueryUserTemplate = "select id, email, password created_at, updated_at from api_user where id = '%s'"

func (r *userRepositoryImpl) Get(id uuid.UUID) (User, error) {
	query := fmt.Sprintf(sqlQueryUserTemplate, id)
	res := r.conn.Query(query)
	if err := res.Err(); err != nil {
		return User{}, err
	}

	var out User
	parser := func(rows db.Scannable) error {
		return rows.Scan(out.id, out.email, out.password, out.createdAt, out.updated_at)
	}

	if err := res.GetSingleValue(parser); err != nil {
		return User{}, err
	}

	return out, nil
}

func (r *userRepositoryImpl) Update(user User) (User, error) {
	return User{}, errors.NewCode(errors.NotImplementedCode)
}

func (r *userRepositoryImpl) Delete(id uuid.UUID) error {
	return errors.NewCode(errors.NotImplementedCode)
}
