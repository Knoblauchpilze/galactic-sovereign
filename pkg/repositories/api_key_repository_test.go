package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	s := RepositoryTestSuite{
		dbPoolInteractionTestCases: map[string]dbPoolInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewApiKeyRepository(pool)
					_, err := s.Create(ctx, defaultApiKey)
					return err
				},
				expectedSql: `
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
				expectedArguments: []interface{}{
					defaultApiKey.Id,
					defaultApiKey.Key,
					defaultApiKey.ApiUser,
					defaultApiKey.ValidUntil,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewApiKeyRepository(pool)
					_, err := s.Get(ctx, defaultApiKeyId)
					return err
				},
				expectedSql: `SELECT id, key, api_user, valid_until FROM api_key WHERE id = $1`,
				expectedArguments: []interface{}{
					defaultApiKeyId,
				},
			},
			"getForKey": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewApiKeyRepository(pool)
					_, err := s.GetForKey(ctx, defaultApiKeyValue)
					return err
				},
				expectedSql: `SELECT id, key, api_user, valid_until FROM api_key WHERE key = $1`,
				expectedArguments: []interface{}{
					defaultApiKeyValue,
				},
			},
			"getForUser": {
				handler: func(ctx context.Context, pool db.ConnectionPool) error {
					s := NewApiKeyRepository(pool)
					_, err := s.GetForUser(ctx, defaultUserId)
					return err
				},
				expectedSql: `SELECT id FROM api_key WHERE api_user = $1`,
				expectedArguments: []interface{}{
					defaultUserId,
				},
			},
		},

		dbPoolSingleValueTestCases: map[string]dbPoolSingleValueTestCase{
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
		},

		dbPoolReturnTestCases: map[string]dbPoolReturnTestCase{
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

func TestApiKeyRepository_Create_RetrievesGeneratedApiKey(t *testing.T) {
	s := RepositorySingleValueTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.Create(ctx, defaultApiKey)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{&uuid.UUID{}},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_Get_InterpretDbData(t *testing.T) {
	s := RepositorySingleValueTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.Get(ctx, defaultApiKeyId)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
				&uuid.UUID{},
				&uuid.UUID{},
				&uuid.UUID{},
				&time.Time{},
			},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForKey_InterpretDbData(t *testing.T) {
	s := RepositorySingleValueTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.GetForKey(ctx, defaultApiKeyValue)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{
				&uuid.UUID{},
				&uuid.UUID{},
				&uuid.UUID{},
				&time.Time{},
			},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForUser_InterpretDbData(t *testing.T) {
	s := RepositoryGetAllTestSuite{
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.GetForUser(ctx, defaultUserId)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{&uuid.UUID{}},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_DeleteForUser_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewApiKeyRepository(&mockConnectionPool{})
			return repo.DeleteForUser(context.Background(), tx, defaultUserId)
		},
		expectedSql: []string{
			`DELETE FROM api_key WHERE api_user = $1`,
		},
		expectedArguments: [][]interface{}{
			{defaultUserId},
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_DeleteForUser_NominalCase(t *testing.T) {
	assert := assert.New(t)

	repo := NewAclRepository()
	mt := &mockTransaction{}

	err := repo.DeleteForUser(context.Background(), mt, defaultUserId)

	assert.Nil(err)
}
