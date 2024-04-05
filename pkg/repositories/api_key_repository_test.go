package repositories

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultApiKeyId = uuid.MustParse("cc1742fa-77b4-4f5f-ac92-058c2e47a5d6")
var defaultApiKeyValue = uuid.MustParse("b01b9b1f-b651-4702-9b58-905b19584d69")
var defaultApiKey = persistence.ApiKey{
	Id:      defaultApiKeyId,
	Key:     defaultApiKeyValue,
	ApiUser: defaultUserId,
}

func TestApiKeyRepository_Create_UsesTransactionToExec(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	repo.Create(context.Background(), mt, defaultApiKey)

	assert.Equal(0, mc.execCalled)
	assert.Equal(1, mt.execCalled)
}

func TestApiKeyRepository_Create_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	repo.Create(context.Background(), mt, defaultApiKey)

	assert.Equal("INSERT INTO api_key (id, key, api_user) VALUES($1, $2, $3)", mt.sqlQuery)
	assert.Equal(3, len(mt.args))
	assert.Equal(defaultApiKey.Id, mt.args[0])
	assert.Equal(defaultApiKey.Key, mt.args[1])
	assert.Equal(defaultApiKey.ApiUser, mt.args[2])
}

func TestApiKeyRepository_Create_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{
		execErr: errDefault,
	}

	_, err := repo.Create(context.Background(), mt, defaultApiKey)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_Create_ReturnsInputApiKey(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	actual, err := repo.Create(context.Background(), mt, defaultApiKey)

	assert.Nil(err)
	assert.Equal(defaultApiKey, actual)
}

func TestApiKeyRepository_Get_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	repo.Get(context.Background(), uuid.UUID{})

	assert.Equal(1, mc.queryCalled)
}

func TestApiKeyRepository_Get_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	repo.Get(context.Background(), defaultApiKeyId)

	assert.Equal("SELECT id, key, api_user FROM api_key WHERE id = $1", mc.sqlQuery)
	assert.Equal(1, len(mc.args))
	assert.Equal(defaultApiKeyId, mc.args[0])
}

func TestApiKeyRepository_Get_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			err: errDefault,
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.Get(context.Background(), defaultApiKeyId)

	assert.Equal(errDefault, err)
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
	assert.Equal(3, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
	assert.IsType(&uuid.UUID{}, props[1])
	assert.IsType(&uuid.UUID{}, props[2])
}

func TestApiKeyRepository_GetForKey_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	repo.GetForKey(context.Background(), defaultApiKeyValue)

	assert.Equal(1, mc.queryCalled)
}

func TestApiKeyRepository_GetForKey_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)

	repo.GetForKey(context.Background(), defaultApiKeyValue)

	assert.Equal("SELECT id, key, api_user, enabled FROM api_key WHERE key = $1", mc.sqlQuery)
	assert.Equal(1, len(mc.args))
	assert.Equal(defaultApiKeyValue, mc.args[0])
}

func TestApiKeyRepository_GetForKey_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{
		rows: mockRows{
			err: errDefault,
		},
	}
	repo := NewApiKeyRepository(mc)

	_, err := repo.GetForKey(context.Background(), defaultApiKeyValue)

	assert.Equal(errDefault, err)
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
	var b bool
	assert.IsType(&b, props[3])
}

func TestApiKeyRepository_GetForUser_UsesConnectionToQuery(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	repo.GetForUser(context.Background(), mt, defaultUserId)

	assert.Equal(0, mc.queryCalled)
	assert.Equal(1, mt.queryCalled)
}

func TestApiKeyRepository_GetForUser_GeneratesValidSql(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	repo.GetForUser(context.Background(), mt, defaultUserId)

	assert.Equal("SELECT id FROM api_key WHERE api_user = $1", mt.sqlQuery)
	assert.Equal(1, len(mt.args))
	assert.Equal(defaultUserId, mt.args[0])
}

func TestApiKeyRepository_GetForUser_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{
		rows: mockRows{
			err: errDefault,
		},
	}

	_, err := repo.GetForUser(context.Background(), mt, defaultUserId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForUser_CallsGetAll(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	repo.GetForUser(context.Background(), mt, defaultUserId)

	assert.Equal(1, mt.rows.allCalled)
}

func TestApiKeyRepository_GetForUser_WhenResultReturnsError_Fails(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{
		rows: mockRows{
			allErr: errDefault,
		},
	}

	_, err := repo.GetForUser(context.Background(), mt, defaultUserId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForUser_WhenResultSucceeds_Success(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	_, err := repo.GetForUser(context.Background(), mt, defaultUserId)

	assert.Nil(err)
}

func TestApiKeyRepository_GetForUser_PropagatesScanErrors(t *testing.T) {
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

	_, err := repo.GetForUser(context.Background(), mt, defaultUserId)

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_GetForUser_ScansApiKeyProperties(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{
		rows: mockRows{
			scanner: &mockScannable{},
		},
	}

	_, err := repo.GetForUser(context.Background(), mt, defaultUserId)

	assert.Nil(err)

	props := mt.rows.scanner.props
	assert.Equal(1, mt.rows.scanner.scannCalled)
	assert.Equal(1, len(props))
	assert.IsType(&uuid.UUID{}, props[0])
}

func TestApiKeyRepository_Delete_UsesTransactionToExec(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	repo.Delete(context.Background(), mt, []uuid.UUID{defaultApiKeyId})

	assert.Equal(0, mc.execCalled)
	assert.Equal(1, mt.execCalled)
}

func TestApiKeyRepository_Delete_GeneratesValidSql_ForSingleId(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	repo.Delete(context.Background(), mt, []uuid.UUID{defaultApiKeyId})

	assert.Equal("DELETE FROM api_key WHERE id IN ($1)", mt.sqlQuery)
	assert.Equal(1, len(mt.args))
	assert.Equal(defaultApiKeyId, mt.args[0])
}

func TestApiKeyRepository_Delete_GeneratesValidSql_ForMultipleId(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	ids := []uuid.UUID{
		uuid.MustParse("50714fb2-db52-4e3a-8315-cf8e4a8abcf8"),
		uuid.MustParse("9fc0def1-d51c-4af0-8db5-40310796d16d"),
	}

	repo.Delete(context.Background(), mt, ids)

	assert.Equal("DELETE FROM api_key WHERE id IN ($1,$2)", mt.sqlQuery)
	assert.Equal(2, len(mt.args))
	assert.Equal(ids[0], mt.args[0])
	assert.Equal(ids[1], mt.args[1])
}

func TestApiKeyRepository_Delete_PropagatesQueryFailure(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{
		execErr: errDefault,
	}

	err := repo.Delete(context.Background(), mt, []uuid.UUID{defaultApiKeyId})

	assert.Equal(errDefault, err)
}

func TestApiKeyRepository_Delete_NominalCase(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewApiKeyRepository(mc)
	mt := &mockTransaction{}

	err := repo.Delete(context.Background(), mt, []uuid.UUID{defaultApiKeyId})

	assert.Nil(err)
}
