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
	"github.com/stretchr/testify/suite"
)

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

func TestUserRepository_Create(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUserRepository(pool)
			_, err := repo.Create(ctx, defaultUser)
			return err
		},
		expectedSql: `INSERT INTO api_user (id, email, password, created_at) VALUES($1, $2, $3, $4)`,
		expectedArguments: []interface{}{
			defaultUser.Id,
			defaultUser.Email,
			defaultUser.Password,
			defaultUser.CreatedAt,
		},
	}

	suite.Run(t, &s)
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

func TestUserRepository_Get(t *testing.T) {
	s := RepositoryPoolTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUserRepository(pool)
			_, err := repo.Get(ctx, defaultUserId)
			return err
		},
		expectedSql: `SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE id = $1`,
		expectedArguments: []interface{}{
			defaultUserId,
		},
	}

	suite.Run(t, &s)
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

func TestUserRepository_GetByEmail(t *testing.T) {
	s := RepositoryPoolTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUserRepository(pool)
			_, err := repo.GetByEmail(ctx, defaultUserEmail)
			return err
		},
		expectedSql: `SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE email = $1`,
		expectedArguments: []interface{}{
			defaultUserEmail,
		},
	}

	suite.Run(t, &s)
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

func TestUserRepository_List(t *testing.T) {
	s := RepositoryPoolTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUserRepository(pool)
			_, err := repo.List(ctx)
			return err
		},
		expectedSql: `SELECT id FROM api_user`,
	}

	suite.Run(t, &s)
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

func TestUserRepository_Update(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewUserRepository(pool)
			_, err := repo.Update(context.Background(), defaultUser)
			return err
		},
		expectedSql: `UPDATE api_user SET email = $1, password = $2, version = $3 WHERE id = $4 AND version = $5`,
		expectedArguments: []interface{}{
			defaultUser.Email,
			defaultUser.Password,
			defaultUser.Version + 1,
			defaultUser.Id,
			defaultUser.Version,
		},
	}

	suite.Run(t, &s)
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

func TestUserRepository_Delete(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserRepository(&mockConnectionPool{})
			return repo.Delete(ctx, tx, defaultUserId)
		},
		expectedSql: `DELETE FROM api_user WHERE id = $1`,
		expectedArguments: []interface{}{
			defaultUserId,
		},
	}

	suite.Run(t, &s)
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
