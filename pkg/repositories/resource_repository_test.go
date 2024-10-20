package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultResourceId = uuid.MustParse("4b60c8fe-39df-492b-95e8-71474bbac5ea")
var defaultResourceName = "my-resource"
var defaultResource = persistence.Resource{
	Id:        defaultResourceId,
	Name:      defaultResourceName,
	CreatedAt: time.Date(2024, 7, 28, 9, 43, 8, 651387234, time.UTC),
	UpdatedAt: time.Date(2024, 7, 28, 9, 43, 8, 651387234, time.UTC),
}

func Test_ResourceRepository(t *testing.T) {
	dummyStr := ""

	s := RepositoryPoolTestSuite{
		dbInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewResourceRepository(pool)
					_, err := s.Create(ctx, defaultResource)
					return err
				},
				expectedSqlQueries: []string{
					`INSERT INTO resource (id, name, created_at) VALUES($1, $2, $3)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultResource.Id,
						defaultResource.Name,
						defaultResource.CreatedAt,
					},
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewResourceRepository(pool)
					_, err := s.Get(ctx, defaultResourceId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, name, created_at, updated_at FROM resource WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultResourceId,
					},
				},
			},
			"delete": {
				sqlMode: ExecBased,
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						affectedRows: 1,
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewResourceRepository(pool)
					return s.Delete(ctx, defaultResourceId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM resource WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultResourceId,
					},
				},
			},
		},

		dbSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewResourceRepository(pool)
					_, err := repo.Get(ctx, defaultResourceId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyStr,
						&time.Time{},
						&time.Time{},
					},
				},
			},
		},

		dbReturnTestCases: map[string]dbPoolReturnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool) interface{} {
					s := NewResourceRepository(pool)
					out, _ := s.Create(ctx, defaultResource)
					return out
				},
				expectedContent: defaultResource,
			},
		},

		dbErrorTestCases: map[string]dbPoolErrorTestCase{
			"create_duplicatedKey": {
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						execErr: fmt.Errorf(`duplicate key value violates unique constraint "resource_name_key" (SQLSTATE 23505)`),
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewResourceRepository(pool)
					_, err := s.Create(ctx, defaultResource)
					return err
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey))
				},
			},
			"delete_noRowsAffected": {
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						affectedRows: 0,
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewResourceRepository(pool)
					return s.Delete(ctx, defaultResourceId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
			"delete_moreThanOneRowAffected": {
				generateMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						affectedRows: 2,
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewResourceRepository(pool)
					return s.Delete(ctx, defaultResourceId)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, db.NoMatchingSqlRows))
				},
			},
		},
	}

	suite.Run(t, &s)
}

func Test_ResourceRepository_Transaction(t *testing.T) {
	dummyStr := ""

	s := RepositoryTransactionTestSuite{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"list": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewResourceRepository(&mockConnectionPool{})
					_, err := s.List(ctx, tx)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, name, created_at, updated_at FROM resource`,
				},
			},
		},

		dbGetAllTestCases: map[string]dbTransactionGetAllTestCase{
			"list": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewResourceRepository(&mockConnectionPool{})
					_, err := repo.List(ctx, tx)
					return err
				},
				expectedGetAllCalls: 1,
				expectedScanCalls:   1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&dummyStr,
						&time.Time{},
						&time.Time{},
					},
				},
			},
		},
	}

	suite.Run(t, &s)
}
