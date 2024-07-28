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
	"github.com/stretchr/testify/require"
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

	s := RepositoryPoolTestSuite{
		dbInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Create(ctx, defaultUser)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO api_user (id, email, password, created_at) VALUES($1, $2, $3, $4)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUser.Id,
						defaultUser.Email,
						defaultUser.Password,
						defaultUser.CreatedAt,
					},
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Get(ctx, defaultUserId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{defaultUserId},
				},
			},
			"getByEmail": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.GetByEmail(ctx, defaultUserEmail)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, email, password, created_at, updated_at, version FROM api_user WHERE email = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUserEmail,
					},
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.List(ctx)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id FROM api_user`,
				},
			},
			"update": {
				sqlMode: ExecBased,
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						affectedRows: 1,
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Update(ctx, defaultUpdatedUser)
					return err
				},
				expectedSqlQueries: []string{
					`UPDATE api_user SET email = $1, password = $2, version = $3 WHERE id = $4 AND version = $5`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUpdatedUser.Email,
						defaultUpdatedUser.Password,
						defaultUpdatedUser.Version + 1,
						defaultUpdatedUser.Id,
						defaultUpdatedUser.Version,
					},
				},
			},
		},

		dbSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
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
					_, err := repo.GetByEmail(ctx, defaultUserEmail)
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

		dbGetAllTestCases: map[string]dbPoolGetAllTestCase{
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

		dbReturnTestCases: map[string]dbPoolReturnTestCase{
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
					out, _ := s.Update(ctx, defaultUpdatedUser)
					return out
				},
				expectedContent: defaultUpdatedUser,
			},
		},

		dbErrorTestCases: map[string]dbPoolErrorTestCase{
			"create_duplicatedKey": {
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						execErr: fmt.Errorf(`duplicate key value violates unique constraint "api_user_email_key" (SQLSTATE 23505)`),
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Create(ctx, defaultUser)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey))
				},
			},
			"update_optimisticLockException": {
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						affectedRows: 0,
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Update(ctx, defaultUpdatedUser)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.OptimisticLockException))
				},
			},
			"update_moreThanOneRowAffected": {
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						affectedRows: 2,
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewUserRepository(pool)
					_, err := s.Update(ctx, defaultUpdatedUser)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.MoreThanOneMatchingSqlRows))
				},
			},
		},
	}

	suite.Run(t, &s)
}

func Test_UserRepository_Transaction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"delete": {
				sqlMode: ExecBased,
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: 1,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUserRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultUserId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM api_user WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUserId,
					},
				},
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{
			"delete_noRowsAffected": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: 0,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUserRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultUserId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
			"delete_moreThanOneRowAffected": {
				generateMock: func() db.Transaction {
					return &mockTransaction{
						affectedRows: 2,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUserRepository(&mockConnectionPool{})
					return s.Delete(ctx, tx, defaultUserId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
		},
	}

	suite.Run(t, &s)
}
