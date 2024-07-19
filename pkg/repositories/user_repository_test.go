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
var defaultUpdatedUser = persistence.User{
	Id:        defaultUserId,
	Email:     "my-updated-email",
	Password:  "my-updated-password",
	CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	UpdatedAt: time.Date(2009, 11, 17, 20, 34, 59, 651387237, time.UTC),
	Version:   4,
}

func Test_UserRepository(t *testing.T) {
	dummyStr := ""
	dummyInt := 0

	s := RepositoryTestSuite{
		dbPoolInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Create(ctx, defaultUser)
					return err
				},
				expectedSql: `INSERT INTO api_user (id, email, password, created_at) VALUES($1, $2, $3, $4)`,
				expectedArguments: []interface{}{
					defaultUser.Id,
					defaultUser.Email,
					defaultUser.Password,
					defaultUser.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Get(ctx, defaultUserId)
					return err
				},
				expectedSql: `SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE id = $1`,
				expectedArguments: []interface{}{
					defaultUserId,
				},
			},
			"getByEmail": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.GetByEmail(ctx, defaultUserEmail)
					return err
				},
				expectedSql: `SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE email = $1`,
				expectedArguments: []interface{}{
					defaultUserEmail,
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.List(ctx)
					return err
				},
				expectedSql: `SELECT id FROM api_user`,
			},
			"update": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Update(ctx, defaultUpdatedUser)
					return err
				},
				expectedSql: `UPDATE api_user SET email = $1, password = $2, version = $3 WHERE id = $4 AND version = $5`,
				expectedArguments: []interface{}{
					defaultUpdatedUser.Email,
					defaultUpdatedUser.Password,
					defaultUpdatedUser.Version + 1,
					defaultUpdatedUser.Id,
					defaultUpdatedUser.Version,
				},
			},
		},

		dbPoolSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewUserRepository(pool)
					_, err := repo.Get(ctx, defaultUserId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyStr,
						&dummyStr,
						&time.Time{},
						&time.Time{},
						&dummyInt,
					},
				},
			},
			"getByEmail": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewUserRepository(pool)
					_, err := repo.Get(ctx, defaultUserId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyStr,
						&dummyStr,
						&time.Time{},
						&time.Time{},
						&dummyInt,
					},
				},
			},
		},

		dbPoolGetAllTestCases: map[string]dbPoolGetAllTestCase{
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewUserRepository(pool)
					_, err := repo.List(ctx)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
					},
				},
			},
		},

		dbPoolReturnTestCases: map[string]dbPoolReturnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool) interface{} {
					s := NewUserRepository(pool)
					out, _ := s.Create(ctx, defaultUser)
					return out
				},
				expectedContent: defaultUser,
			},
			"update": {
				handler: func(ctx context.Context, pool db.ConnectionPool) interface{} {
					s := NewUserRepository(pool)
					out, _ := s.Update(ctx, defaultUser)
					return out
				},
				expectedContent: defaultUser,
			},
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

func TestUserRepository_Delete_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserRepository(&mockConnectionPool{})
			return repo.Delete(ctx, tx, defaultUserId)
		},
		expectedSql: []string{`DELETE FROM api_user WHERE id = $1`},
		expectedArguments: [][]interface{}{
			{
				defaultUserId,
			},
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
