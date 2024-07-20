package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

var defaultApiKeyId = uuid.MustParse("cc1742fa-77b4-4f5f-ac92-058c2e47a5d6")
var defaultApiKeyValue = uuid.MustParse("b01b9b1f-b651-4702-9b58-905b19584d69")
var defaultApiKey = persistence.ApiKey{
	Id:      defaultApiKeyId,
	Key:     defaultApiKeyValue,
	ApiUser: defaultUserId,
}

func Test_ApiKeyRepository(t *testing.T) {
	s := RepositoryPoolTestSuite{
		dbInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewApiKeyRepository(pool)
					_, err := s.Create(ctx, defaultApiKey)
					return err
				},
				expectedSqlQueries: []string{`
INSERT INTO api_key (id, key, api_user, valid_until)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (api_user) DO UPDATE
	SET
		valid_until = excluded.valid_until
	WHERE
		api_key.api_user = excluded.api_user
	RETURNING
		api_key.key
`,
				},
				expectedArguments: [][]interface{}{
					{defaultApiKey.Id,
						defaultApiKey.Key,
						defaultApiKey.ApiUser,
						defaultApiKey.ValidUntil,
					},
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewApiKeyRepository(pool)
					_, err := s.Get(ctx, defaultApiKeyId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, key, api_user, valid_until FROM api_key WHERE id = $1`,
				},
				expectedArguments: [][]interface{}{
					{defaultApiKeyId},
				},
			},
			"getForKey": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewApiKeyRepository(pool)
					_, err := s.GetForKey(ctx, defaultApiKeyValue)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id, key, api_user, valid_until FROM api_key WHERE key = $1`,
				},
				expectedArguments: [][]interface{}{
					{defaultApiKeyValue},
				},
			},
			"getForUser": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewApiKeyRepository(pool)
					_, err := s.GetForUser(ctx, defaultUserId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id FROM api_key WHERE api_user = $1`,
				},
				expectedArguments: [][]interface{}{
					{defaultUserId},
				},
			},
		},

		dbSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewApiKeyRepository(pool)
					_, err := repo.Create(ctx, defaultApiKey)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
					},
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewApiKeyRepository(pool)
					_, err := repo.Get(ctx, defaultApiKeyId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&uuid.UUID{},
						&time.Time{},
					},
				},
			},
			"getForKey": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewApiKeyRepository(pool)
					_, err := repo.GetForKey(ctx, defaultApiKeyId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           1,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&uuid.UUID{},
						&time.Time{},
					},
				},
			},
		},

		dbGetAllTestCases: map[string]dbPoolGetAllTestCase{
			"getForUser": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					repo := NewApiKeyRepository(pool)
					_, err := repo.GetForUser(ctx, defaultUserId)
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
					s := NewApiKeyRepository(pool)
					out, _ := s.Create(ctx, defaultApiKey)
					return out
				},
				expectedContent: defaultApiKey,
			},
		},
	}

	suite.Run(t, &s)
}

func Test_ApiKeyRepository_Transaction(t *testing.T) {
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
					s := NewApiKeyRepository(&mockConnectionPool{})
					return s.DeleteForUser(ctx, tx, defaultUserId)
				},
				expectedSqlQueries: []string{
					`DELETE FROM api_key WHERE api_user = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUserId,
					},
				},
			},
		},
	}

	suite.Run(t, &s)
}
