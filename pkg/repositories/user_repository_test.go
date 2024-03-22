package repositories

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockConnection struct {
	queryCalled int

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

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")

func TestUserRepository_Create_NotImplemented(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	err := repo.Create(User{})
	assert.True(errors.IsErrorWithCode(err, errors.NotImplementedCode))
}

func TestUserRepository_Get_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Get(uuid.UUID{})

	assert.Equal(1, mc.queryCalled)
}

func TestUserRepository_Get_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Get(defaultUuid)

	assert.Equal("select id, email, password created_at, updated_at from api_user where id = '08ce96a3-3430-48a8-a3b2-b1c987a207ca'", mc.sqlQuery)
	assert.Equal(0, len(mc.args))
}

var errDefault = fmt.Errorf("some error")

func TestUserRepository_Get_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{
		rows: mockRows{
			err: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.Get(defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Get_CallsGetSingleValue(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	repo.Get(defaultUuid)

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

	_, err := repo.Get(defaultUuid)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Get_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	_, err := repo.Get(defaultUuid)

	assert.Nil(err)
}

func TestUserRepository_Update_NotImplemented(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	_, err := repo.Update(User{})
	assert.True(errors.IsErrorWithCode(err, errors.NotImplementedCode))
}

func TestUserRepository_Delete_NotImplemented(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnection{}
	repo := NewUserRepository(mc)

	err := repo.Delete(uuid.UUID{})
	assert.True(errors.IsErrorWithCode(err, errors.NotImplementedCode))
}

func (m *mockConnection) Connect() error { return nil }

func (m *mockConnection) Close() {}

func (m *mockConnection) Query(sql string, arguments ...interface{}) db.Rows {
	m.queryCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return &m.rows
}

func (m *mockConnection) Exec(sql string, arguments ...interface{}) (string, error) { return "", nil }

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
