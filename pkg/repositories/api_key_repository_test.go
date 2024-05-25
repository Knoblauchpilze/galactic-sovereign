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

func TestApiKeyRepository_Create(t *testing.T) {
	expectedSql := `
INSERT INTO api_key (id, key, api_user, valid_until)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (api_user) DO UPDATE
	SET
		valid_until = excluded.valid_until
	WHERE
		api_key.api_user = excluded.api_user
	RETURNING
		api_key.key
`

	s := RepositoryPoolTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.Create(context.Background(), defaultApiKey)
			return err
		},
		expectedSql: expectedSql,
		expectedArguments: []interface{}{
			defaultApiKey.Id,
			defaultApiKey.Key,
			defaultApiKey.ApiUser,
			defaultApiKey.ValidUntil,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_Create_GetsReturnedValue(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{},
	}
	repo := NewApiKeyRepository(mc)

	repo.Create(context.Background(), defaultApiKey)

	assert.Equal(1, mc.rows.singleValueCalled)
}

func TestApiKeyRepository_Create_WhenReturnValueFails_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			singleValueErr: errDefault,
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.Create(context.Background(), defaultApiKey)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_Create_ReturnsInputApiKey(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	actual, err := repo.Create(context.Background(), defaultApiKey)

	assert.Nil(err)
	assert.Equal(defaultApiKey, actual)
}

func TestApiKeyRepository_Get(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.Get(context.Background(), defaultApiKeyId)
			return err
		},
		expectedSql: `SELECT id, key, api_user, valid_until FROM api_key WHERE id = $1`,
		expectedArguments: []interface{}{
			defaultApiKeyId,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_Get_CallsGetSingleValue(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	repo.Get(context.Background(), defaultApiKeyId)

	assert.Equal(1, mc.rows.singleValueCalled)
}

func TestApiKeyRepository_Get_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			singleValueErr: errDefault,
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.Get(context.Background(), defaultApiKeyId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_Get_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	_, err := repo.Get(context.Background(), defaultApiKeyId)

	assert.Nil(err)
}

func TestApiKeyRepository_Get_PropagatesScanErrors(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.Get(context.Background(), defaultApiKeyId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_Get_ScansApiKeyProperties(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.Get(context.Background(), defaultApiKeyId)

	assert.Nil(err)

	props := mc.rows.scanner.props
	assert.Equal(1, mc.rows.scanner.scannCalled)
	assert.Equal(4, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
	assert.IsType(&uuid.UUID{}, props[1])
	assert.IsType(&uuid.UUID{}, props[2])
	assert.IsType(&time.Time{}, props[3])
}

func TestApiKeyRepository_GetForKey(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.GetForKey(context.Background(), defaultApiKeyValue)
			return err
		},
		expectedSql: `SELECT id, key, api_user, valid_until FROM api_key WHERE key = $1`,
		expectedArguments: []interface{}{
			defaultApiKeyValue,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForKey_CallsGetSingleValue(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	repo.GetForKey(context.Background(), defaultApiKeyValue)

	assert.Equal(1, mc.rows.singleValueCalled)
}

func TestApiKeyRepository_GetForKey_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			singleValueErr: errDefault,
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForKey(context.Background(), defaultApiKeyValue)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForKey_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForKey(context.Background(), defaultApiKeyValue)

	assert.Nil(err)
}

func TestApiKeyRepository_GetForKey_PropagatesScanErrors(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForKey(context.Background(), defaultApiKeyValue)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForKey_ScansApiKeyProperties(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForKey(context.Background(), defaultApiKeyValue)

	assert.Nil(err)

	props := mc.rows.scanner.props
	assert.Equal(1, mc.rows.scanner.scannCalled)
	assert.Equal(4, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
	assert.IsType(&uuid.UUID{}, props[1])
	assert.IsType(&uuid.UUID{}, props[2])
	assert.IsType(&time.Time{}, props[3])
}

func TestApiKeyRepository_GetForUser(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			_, err := repo.GetForUser(context.Background(), defaultUserId)
			return err
		},
		expectedSql: `SELECT id FROM api_key WHERE api_user = $1`,
		expectedArguments: []interface{}{
			defaultUserId,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForUser_CallsGetAll(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	repo.GetForUser(context.Background(), defaultUserId)

	assert.Equal(1, mc.rows.allCalled)
}

func TestApiKeyRepository_GetForUser_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			allErr: errDefault,
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForUser(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForUser_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForUser(context.Background(), defaultUserId)

	assert.Nil(err)
}

func TestApiKeyRepository_GetForUser_PropagatesScanErrors(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForUser(context.Background(), defaultUserId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForUser_ScansApiKeyProperties(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForUser(context.Background(), defaultUserId)

	assert.Nil(err)

	props := mc.rows.scanner.props
	assert.Equal(1, mc.rows.scanner.scannCalled)
	assert.Equal(1, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
}

func TestApiKeyRepository_GetForUserTx(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewApiKeyRepository(&mockConnectionPool{})
			_, err := repo.GetForUserTx(context.Background(), tx, defaultUserId)
			return err
		},
		expectedSql: `SELECT id FROM api_key WHERE api_user = $1`,
		expectedArguments: []interface{}{
			defaultUserId,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_GetForUserTx_CallsGetAll(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	repo.GetForUserTx(context.Background(), mt, defaultUserId)

	assert.Equal(1, mt.rows.allCalled)
}

func TestApiKeyRepository_GetForUserTx_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{
		rows: mockRows{
			allErr: errDefault,
		},
	}

	_, err := repo.GetForUserTx(context.Background(), mt, defaultUserId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForUserTx_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	_, err := repo.GetForUserTx(context.Background(), mt, defaultUserId)

	assert.Nil(err)
}

func TestApiKeyRepository_GetForUserTx_PropagatesScanErrors(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{
		rows: mockRows{
			scanner: &mockScannable{
				err: errDefault,
			},
		},
	}

	_, err := repo.GetForUserTx(context.Background(), mt, defaultUserId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForUserTx_ScansApiKeyProperties(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}

	_, err := repo.GetForUserTx(context.Background(), mt, defaultUserId)

	assert.Nil(err)

	props := mt.rows.scanner.props
	assert.Equal(1, mt.rows.scanner.scannCalled)
	assert.Equal(1, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
}

func TestApiKeyRepository_Delete_SingleId(t *testing.T) {
	s := RepositoryPoolTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			return repo.Delete(context.Background(), []uuid.UUID{defaultApiKeyId})
		},
		expectedSql: `DELETE FROM api_key WHERE id IN ($1)`,
		expectedArguments: []interface{}{
			defaultApiKeyId,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_Delete_MultipleIds(t *testing.T) {
	ids := []uuid.UUID{
		uuid.MustParse("50714fb2-db52-4e3a-8315-cf8e4a8abcf8"),
		uuid.MustParse("9fc0def1-d51c-4af0-8db5-40310796d16d"),
	}

	s := RepositoryPoolTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, pool db.ConnectionPool) error {
			repo := NewApiKeyRepository(pool)
			return repo.Delete(context.Background(), ids)
		},
		expectedSql: `DELETE FROM api_key WHERE id IN ($1,$2)`,
		expectedArguments: []interface{}{
			ids[0],
			ids[1],
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_Delete_NominalCase(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	err := repo.Delete(context.Background(), []uuid.UUID{defaultApiKeyId})

	assert.Nil(err)
}

func TestApiKeyRepository_DeleteTx_SingleId(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewApiKeyRepository(&mockConnectionPool{})
			return repo.DeleteTx(context.Background(), tx, []uuid.UUID{defaultApiKeyId})
		},
		expectedSql: `DELETE FROM api_key WHERE id IN ($1)`,
		expectedArguments: []interface{}{
			defaultApiKeyId,
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_DeleteTx_MultipleIds(t *testing.T) {
	ids := []uuid.UUID{
		uuid.MustParse("50714fb2-db52-4e3a-8315-cf8e4a8abcf8"),
		uuid.MustParse("9fc0def1-d51c-4af0-8db5-40310796d16d"),
	}

	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewApiKeyRepository(&mockConnectionPool{})
			return repo.DeleteTx(context.Background(), tx, ids)
		},
		expectedSql: `DELETE FROM api_key WHERE id IN ($1,$2)`,
		expectedArguments: []interface{}{
			ids[0],
			ids[1],
		},
	}

	suite.Run(t, &s)
}

func TestApiKeyRepository_DeleteTx_NominalCase(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	err := repo.DeleteTx(context.Background(), mt, []uuid.UUID{defaultApiKeyId})

	assert.Nil(err)
}
