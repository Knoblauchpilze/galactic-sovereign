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

type mockConnectionPool struct {
	db.ConnectionPool

	queryCalled int
	execCalled  int

	affectedRows int
	execErr      error

	sqlQuery string
	args     []interface{}

	rows mockRows
}

type mockTransaction struct {
	db.Transaction

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
	allErr         error

	singleValueCalled int
	allCalled         int
	scanner           *mockScannable
}

type mockScannable struct {
	err error

	scannCalled int
	props       []interface{}
}

var errDefault = fmt.Errorf("some error")
var defaultUserId = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultUserEmail = "e.mail@domain.com"
var defaultUser = persistence.User{
	Id:        defaultUserId,
	Email:     defaultUserEmail,
	Password:  "password",
	CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	UpdatedAt: time.Date(2009, 11, 17, 20, 34, 59, 651387237, time.UTC),
	Version:   4,
}

func TestUserRepository_Create_UsesConnectionToExec(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.Create(context.Background(), defaultUser)

	assert.Equal(1, mc.execCalled)
}

func TestUserRepository_Create_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
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

	mc := &mockConnectionPool{
		execErr: errDefault,
	}
	repo := NewUserRepository(mc)

	_, err := repo.Create(context.Background(), defaultUser)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Create_WhenQueryIndicatesDuplicatedKey_ReturnsDuplicatedKey(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		execErr: fmt.Errorf(`duplicate key value violates unique constraint "api_user_email_key" (SQLSTATE 23505)`),
	}
	repo := NewUserRepository(mc)

	_, err := repo.Create(context.Background(), defaultUser)

	assert.True(errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey))
}

func TestUserRepository_Create_ReturnsInputUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	actual, err := repo.Create(context.Background(), defaultUser)

	assert.Nil(err)
	assert.Equal(defaultUser, actual)
}

func TestUserRepository_Get_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.Get(context.Background(), uuid.UUID{})

	assert.Equal(1, mc.queryCalled)
}

func TestUserRepository_Get_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.Get(context.Background(), defaultUserId)

	assert.Equal("SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE id = $1", mc.sqlQuery)
	assert.Equal(1, len(mc.args))
	assert.Equal(defaultUserId, mc.args[0])
}

func TestUserRepository_Get_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			err: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.Get(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Get_CallsGetSingleValue(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.Get(context.Background(), defaultUserId)

	assert.Equal(1, mc.rows.singleValueCalled)
}

func TestUserRepository_Get_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			singleValueErr: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.Get(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Get_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	_, err := repo.Get(context.Background(), defaultUserId)

	assert.Nil(err)
}

func TestUserRepository_Get_PropagatesScanErrors(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.Get(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Get_ScansUserProperties(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.Get(context.Background(), defaultUserId)

	assert.Nil(err)

	props := mc.rows.scanner.props
	assert.Equal(1, mc.rows.scanner.scannCalled)
	assert.Equal(6, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
	var str string
	assert.IsType(&str, props[1])
	assert.IsType(&str, props[2])
	assert.IsType(&time.Time{}, props[3])
	assert.IsType(&time.Time{}, props[4])
	var itg int
	assert.IsType(&itg, props[5])
}

func TestUserRepository_GetByEmail_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.GetByEmail(context.Background(), defaultUserEmail)

	assert.Equal(1, mc.queryCalled)
}

func TestUserRepository_GetByEmail_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.GetByEmail(context.Background(), defaultUserEmail)

	assert.Equal("SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE email = $1", mc.sqlQuery)
	assert.Equal(1, len(mc.args))
	assert.Equal(defaultUserEmail, mc.args[0])
}

func TestUserRepository_GetByEmail_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			err: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.GetByEmail(context.Background(), defaultUserEmail)

	assert.Equal(errDefault, err)
}

func TestUserRepository_GetByEmail_CallsGetSingleValue(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.GetByEmail(context.Background(), defaultUserEmail)

	assert.Equal(1, mc.rows.singleValueCalled)
}

func TestUserRepository_GetByEmail_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			singleValueErr: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.GetByEmail(context.Background(), defaultUserEmail)

	assert.Equal(errDefault, err)
}

func TestUserRepository_GetByEmail_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	_, err := repo.GetByEmail(context.Background(), defaultUserEmail)

	assert.Nil(err)
}

func TestUserRepository_GetByEmail_PropagatesScanErrors(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.GetByEmail(context.Background(), defaultUserEmail)

	assert.Equal(errDefault, err)
}

func TestUserRepository_GetByEmail_ScansUserProperties(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.GetByEmail(context.Background(), defaultUserEmail)

	assert.Nil(err)

	props := mc.rows.scanner.props
	assert.Equal(1, mc.rows.scanner.scannCalled)
	assert.Equal(6, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
	var str string
	assert.IsType(&str, props[1])
	assert.IsType(&str, props[2])
	assert.IsType(&time.Time{}, props[3])
	assert.IsType(&time.Time{}, props[4])
	var itg int
	assert.IsType(&itg, props[5])
}

func TestUserRepository_List_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.List(context.Background())

	assert.Equal(1, mc.queryCalled)
}

func TestUserRepository_List_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.List(context.Background())

	assert.Equal("SELECT id FROM api_user", mc.sqlQuery)
	assert.Equal(0, len(mc.args))
}

func TestUserRepository_List_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			err: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.List(context.Background())

	assert.Equal(errDefault, err)
}

func TestUserRepository_List_CallsGetAll(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.List(context.Background())

	assert.Equal(1, mc.rows.allCalled)
}

func TestUserRepository_List_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			allErr: errDefault,
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.List(context.Background())

	assert.Equal(errDefault, err)
}

func TestUserRepository_List_PropagatesScanErrors(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.List(context.Background())

	assert.Equal(errDefault, err)
}

func TestUserRepository_List_ScansIds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}
	repo := NewUserRepository(mc)

	_, err := repo.List(context.Background())

	assert.Nil(err)

	props := mc.rows.scanner.props
	assert.Equal(1, mc.rows.scanner.scannCalled)
	assert.Equal(1, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
}

func TestUserRepository_Update_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.Update(context.Background(), defaultUser)

	assert.Equal(1, mc.execCalled)
}

func TestUserRepository_Update_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)

	repo.Update(context.Background(), defaultUser)

	assert.Equal("UPDATE api_user SET email = $1, password = $2, version = $3 WHERE id = $4 AND version = $5", mc.sqlQuery)
	assert.Equal(5, len(mc.args))
	assert.Equal(defaultUser.Email, mc.args[0])
	assert.Equal(defaultUser.Password, mc.args[1])
	assert.Equal(defaultUser.Version+1, mc.args[2])
	assert.Equal(defaultUser.Id, mc.args[3])
	assert.Equal(defaultUser.Version, mc.args[4])
}

func TestUserRepository_Update_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		execErr: errDefault,
	}
	repo := NewUserRepository(mc)

	_, err := repo.Update(context.Background(), defaultUser)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Update_WhenAffectedRowsIsZero_ReturnsOptimisticLockException(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		affectedRows: 0,
	}
	repo := NewUserRepository(mc)

	_, err := repo.Update(context.Background(), defaultUser)

	assert.True(errors.IsErrorWithCode(err, db.OptimisticLockException))
}

func TestUserRepository_Update_WhenAffectedRowsIsGreaterThanOne_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		affectedRows: 2,
	}
	repo := NewUserRepository(mc)

	_, err := repo.Update(context.Background(), defaultUser)

	assert.True(errors.IsErrorWithCode(err, db.MoreThanOneMatchingSqlRows))
}

func TestUserRepository_Update_WhenAffectedRowsIsOne_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		affectedRows: 1,
	}
	repo := NewUserRepository(mc)

	_, err := repo.Update(context.Background(), defaultUser)

	assert.Nil(err)
}

func TestUserRepository_Update_ReturnsUpdatedUser(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		affectedRows: 1,
	}
	repo := NewUserRepository(mc)

	actual, _ := repo.Update(context.Background(), defaultUser)

	expected := persistence.User{
		Id:       defaultUser.Id,
		Email:    defaultUser.Email,
		Password: defaultUser.Password,

		CreatedAt: defaultUser.CreatedAt,
		UpdatedAt: defaultUser.UpdatedAt,

		Version: defaultUser.Version + 1,
	}

	assert.Equal(expected, actual)
}

func TestUserRepository_Delete_UsesTransactionToExec(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)
	mt := &mockTransaction{}

	repo.Delete(context.Background(), mt, defaultUserId)

	assert.Equal(0, mc.execCalled)
	assert.Equal(1, mt.execCalled)
}

func TestUserRepository_Delete_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)
	mt := &mockTransaction{}

	repo.Delete(context.Background(), mt, defaultUserId)

	assert.Equal("DELETE FROM api_user WHERE id = $1", mt.sqlQuery)
	assert.Equal(1, len(mt.args))
	assert.Equal(defaultUserId, mt.args[0])
}

func TestUserRepository_Delete_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)
	mt := &mockTransaction{
		execErr: errDefault,
	}

	err := repo.Delete(context.Background(), mt, defaultUserId)

	assert.Equal(errDefault, err)
}

func TestUserRepository_Delete_WhenAffectedRowsIsNotOne_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)
	mt := &mockTransaction{
		affectedRows: 2,
	}

	err := repo.Delete(context.Background(), mt, defaultUserId)

	assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
}

func TestUserRepository_Delete_WhenAffectedRowsIsOne_Succeeds(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewUserRepository(mc)
	mt := &mockTransaction{
		affectedRows: 1,
	}

	err := repo.Delete(context.Background(), mt, defaultUserId)

	assert.Nil(err)
}

func (m *mockConnectionPool) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.queryCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return &m.rows
}

func (m *mockConnectionPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	m.execCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return m.affectedRows, m.execErr
}

func (m *mockTransaction) Close(ctx context.Context) {}

func (m *mockTransaction) Query(ctx context.Context, sql string, arguments ...interface{}) db.Rows {
	m.queryCalled++
	m.sqlQuery = sql
	m.args = append(m.args, arguments...)
	return &m.rows
}

func (m *mockTransaction) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
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
	if m.scanner != nil {
		return parser(m.scanner)
	}
	return m.singleValueErr
}

func (m *mockRows) GetAll(parser db.RowParser) error {
	m.allCalled++
	if m.scanner != nil {
		return parser(m.scanner)
	}
	return m.allErr
}

func (m *mockScannable) Scan(dest ...interface{}) error {
	m.scannCalled++
	m.props = append(m.props, dest...)
	return m.err
}
