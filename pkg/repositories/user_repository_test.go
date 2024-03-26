package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockConnection struct {
	queryCalled int
	execCalled  int

	affectedRows int
	execErr      error

	sqlQuery string
	args     []interface{}

	rows mockRows
}

type mockRows struct {
	err            error
	singleValueErr error

	singleValueCalled int
	allCalled         int
}

var errDefault = fmt.Errorf("some error")
var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultUser = persistence.User{
	Id:        defaultUuid,
	Email:     "e.mail@domain.com",
	Password:  "password",
	CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	UpdatedAt: time.Date(2009, 11, 17, 20, 34, 59, 651387237, time.UTC),
}

func TestUserRepository_Create_UsesConnectionToExec(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Create(context.Background(), defaultUser)

	assert.Equal(1, mc.execCalled)
}

func TestUserRepository_Create_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Create(context.Background(), defaultUser)

	assert.Equal("INSERT INTO api_user (id, email, password, created_at) VALUES($1, $2, $3, $4)", mc.sqlQuery)
	assert.Equal(4, len(mc.args))
	assert.Equal(defaultUser.Id, mc.args[0])
	assert.Equal(defaultUser.Email, mc.args[1])
	assert.Equal(defaultUser.Password, mc.args[2])
	assert.Equal(defaultUser.CreatedAt, mc.args[3])
}

func TestUserRepository_Create_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		execErr: errDefault,
	}
	repo := NewUserRepository(mc)

	err := repo.Create(context.Background(), defaultUser)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Get_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Get(context.Background(), uuid.UUID{})

	assert.Equal(1, mc.queryCalled)
}

func TestUserRepository_Get_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Get(context.Background(), defaultUuid)

	assert.Equal("SELECT id, email, password, created_at, updated_at FROM api_user WHERE id = $1", mc.sqlQuery)
	assert.Equal(1, len(mc.args))
	assert.Equal(defaultUuid, mc.args[0])
}

func TestUserRepository_Get_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		rows: mockRows{
			err: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.Get(context.Background(), defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Get_CallsGetSingleValue(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Get(context.Background(), defaultUuid)

	assert.Equal(1, mc.rows.singleValueCalled)
}

func TestUserRepository_Get_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		rows: mockRows{
			singleValueErr: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.Get(context.Background(), defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Get_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	_, err := repo.Get(context.Background(), defaultUuid)

	assert.Nil(err)
}

func TestUserRepository_Update_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Update(context.Background(), defaultUser)

	assert.Equal(1, mc.execCalled)
}

func TestUserRepository_Update_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Update(context.Background(), defaultUser)

	assert.Equal("UPDATE api_user SET email = $1, password = $2 WHERE id = $3", mc.sqlQuery)
	assert.Equal(3, len(mc.args))
	assert.Equal(defaultUser.Email, mc.args[0])
	assert.Equal(defaultUser.Password, mc.args[1])
	assert.Equal(defaultUser.Id, mc.args[2])
}

func TestUserRepository_Update_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		execErr: errDefault,
	}
	repo := NewUserRepository(mc)

	_, err := repo.Update(context.Background(), defaultUser)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Update_WhenAffectedRowsIsNotOne_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		affectedRows: 2,
	}
	repo := NewUserRepository(mc)

	_, err := repo.Update(context.Background(), defaultUser)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestUserRepository_Update_WhenAffectedRowsIsOne_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		affectedRows: 1,
	}
	repo := NewUserRepository(mc)

	_, err := repo.Update(context.Background(), defaultUser)

	assert.Nil(err)
}

func TestUserRepository_Update_ReturnsUpdatedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		affectedRows: 1,
	}
	repo := NewUserRepository(mc)

	actual, _ := repo.Update(context.Background(), defaultUser)

	assert.Equal(defaultUser, actual)
}

func TestUserRepository_Delete_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Delete(context.Background(), uuid.UUID{})

	assert.Equal(1, mc.execCalled)
}

func TestUserRepository_Delete_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Delete(context.Background(), defaultUuid)

	assert.Equal("DELETE FROM api_user WHERE id = $1", mc.sqlQuery)
	assert.Equal(1, len(mc.args))
	assert.Equal(defaultUuid, mc.args[0])
}

func TestUserRepository_Delete_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		execErr: errDefault,
	}
	repo := NewUserRepository(mc)

	err := repo.Delete(context.Background(), defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Delete_WhenAffectedRowsIsNotOne_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		affectedRows: 2,
	}
	repo := NewUserRepository(mc)

	err := repo.Delete(context.Background(), defaultUuid)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestUserRepository_Delete_WhenAffectedRowsIsOne_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		affectedRows: 1,
	}
	repo := NewUserRepository(mc)

	err := repo.Delete(context.Background(), defaultUuid)

	assert.Nil(err)
}

func (m *mockConnection) Connect() error { return nil }

func (m *mockConnection) Close() {}

func (m *mockConnection) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.queryCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return &m.rows
}

func (m *mockConnection) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	m.execCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return m.affectedRows, m.execErr
}

func (m *mockRows) Err() error { return m.err }

func (m *mockRows) Empty() bool { return false }

func (m *mockRows) Close() {}

func (m *mockRows) GetSingleValue(parser db.RowParser) error {
	m.singleValueCalled++
	return m.singleValueErr
}

func (m *mockRows) GetAll(parser db.RowParser) error {
	m.allCalled++
	return nil
}
