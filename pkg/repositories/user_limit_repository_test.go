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

var defaultUserLimitId = uuid.MustParse("75481254-9eb9-4a07-8d66-1882b80a8421")
var defaultUserLimit = persistence.UserLimit{
	Id:   defaultUserLimitId,
	Name: "my-limit",
	User: defaultUserId,

	Limits: []persistence.Limit{
		{
			Id:    uuid.MustParse("1f686a74-5ecf-4aa5-bf24-75e744068909"),
			Name:  "limit-1",
			Value: "2",
		},
		{
			Id:    uuid.MustParse("a4cd5f53-3f35-4690-b4bc-f434acffbff9"),
			Name:  "limit-2",
			Value: "my-value",
		},
	},

	CreatedAt: time.Date(2024, 06, 22, 16, 5, 20, 651387237, time.UTC),
	UpdatedAt: time.Date(2024, 06, 22, 16, 5, 40, 651387237, time.UTC),
}

func TestUserLimitRepository_Create_DbInteraction(t *testing.T) {
	expectedSql := `
INSERT INTO user_limit (id, name, api_user)
	VALUES($1, $2, $3)
	ON CONFLICT (name, api_user) DO NOTHING
	RETURNING
		user_limit.id
`
	expectedLimitsSql := `
INSERT INTO limits (id, name, value, user_limit)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (name, user_limit) DO UPDATE
	SET
		value = excluded.value
	WHERE
		limits.name = excluded.name
		AND limits.user_limit = excluded.user_limit
	RETURNING
		limits.id
`

	s := RepositoryTransactionTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			_, err := repo.Create(context.Background(), tx, defaultUserLimit)
			return err
		},
		expectedSql: []string{
			expectedSql,
			expectedLimitsSql,
			expectedLimitsSql,
		},
		expectedArguments: [][]interface{}{
			{
				defaultUserLimit.Id,
				defaultUserLimit.Name,
				defaultAcl.User,
			},
			{
				defaultUserLimit.Limits[0].Id,
				defaultUserLimit.Limits[0].Name,
				defaultUserLimit.Limits[0].Value,
				defaultUserLimit.Id,
			},
			{
				defaultUserLimit.Limits[1].Id,
				defaultUserLimit.Limits[1].Name,
				defaultUserLimit.Limits[1].Value,
				defaultUserLimit.Id,
			},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_Create_RetrievesGeneratedUserLimit(t *testing.T) {
	s := RepositorySingleValueTransactionTestSuite{
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			_, err := repo.Create(ctx, tx, defaultUserLimit)
			return err
		},
		expectedSingleValueCalls: 3,
		expectedScanCalls:        3,
		expectedScannedProps: [][]interface{}{
			{&uuid.UUID{}},
			{&uuid.UUID{}},
			{&uuid.UUID{}},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_Create_ReturnsInputUserLimit(t *testing.T) {
	assert := assert.New(t)

	repo := NewUserLimitRepository()
	mt := &mockTransaction{}

	actual, err := repo.Create(context.Background(), mt, defaultUserLimit)

	assert.Nil(err)
	assert.Equal(defaultUserLimit, actual)
}

func TestUserLimitRepository_Get_DbInteraction(t *testing.T) {
	expectedSql := `
SELECT
	id,
	name,
	api_user,
	created_at,
	updated_at
FROM
	user_limit
WHERE
	id = $1
`
	expectedLimitsSql := `
SELECT
	id,
	name,
	value,
	created_at,
	updated_at
FROM
	limits
WHERE
	user_limit = $1
`

	s := RepositoryTransactionTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			_, err := repo.Get(context.Background(), tx, defaultUserLimitId)
			return err
		},
		expectedSql: []string{
			expectedSql,
			expectedLimitsSql,
		},
		expectedArguments: [][]interface{}{
			{defaultUserLimitId},
			{defaultUserLimitId},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_Get_InterpretDbData(t *testing.T) {
	dummyStr := ""

	s := RepositorySingleValueTransactionTestSuite{
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			_, err := repo.Get(ctx, tx, defaultUserLimitId)
			return err
		},
		expectedSingleValueCalls: 1,
		expectedScanCalls:        2,
		expectedScannedProps: [][]interface{}{
			{
				&uuid.UUID{},
				&dummyStr,
				&uuid.UUID{},
				&time.Time{},
				&time.Time{},
			},
			{
				&uuid.UUID{},
				&dummyStr,
				&dummyStr,
				&time.Time{},
				&time.Time{},
			},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_GetForUser_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			_, err := repo.GetForUser(context.Background(), tx, defaultUserId)
			return err
		},
		expectedSql: []string{`SELECT id FROM user_limit WHERE api_user = $1`},
		expectedArguments: [][]interface{}{
			{
				defaultUserId,
			},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_GetForUser_InterpretDbData(t *testing.T) {
	s := RepositoryGetAllTransactionTestSuite{
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			_, err := repo.GetForUser(ctx, tx, defaultUserId)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{&uuid.UUID{}},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_Delete_SingleId_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			return repo.Delete(context.Background(), tx, []uuid.UUID{defaultUserLimitId})
		},
		expectedSql: []string{
			`DELETE FROM limits WHERE user_limit in ($1)`,
			`DELETE FROM user_limit WHERE id IN ($1)`,
		},
		expectedArguments: [][]interface{}{
			{defaultUserLimitId},
			{defaultUserLimitId},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_Delete_MultipleIds_DbInteraction(t *testing.T) {
	ids := []uuid.UUID{
		uuid.MustParse("50714fb2-db52-4e3a-8315-cf8e4a8abcf8"),
		uuid.MustParse("9fc0def1-d51c-4af0-8db5-40310796d16d"),
	}

	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			return repo.Delete(context.Background(), tx, ids)
		},
		expectedSql: []string{
			`DELETE FROM limits WHERE user_limit in ($1,$2)`,
			`DELETE FROM user_limit WHERE id IN ($1,$2)`,
		},
		expectedArguments: [][]interface{}{
			{ids[0], ids[1]},
			{ids[0], ids[1]},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_Delete_NominalCase(t *testing.T) {
	assert := assert.New(t)

	repo := NewUserLimitRepository()
	mt := &mockTransaction{}

	err := repo.Delete(context.Background(), mt, []uuid.UUID{defaultUserLimitId})

	assert.Nil(err)
}

func TestUserLimitRepository_DeleteForUser_DbInteraction(t *testing.T) {
	expectedLimitSqlQuery := `
DELETE FROM limits
	WHERE user_limit
		IN (SELECT id FROM user_limit WHERE api_user = $1)
`

	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewUserLimitRepository()
			return repo.DeleteForUser(context.Background(), tx, defaultUserId)
		},
		expectedSql: []string{
			expectedLimitSqlQuery,
			`DELETE FROM user_limit WHERE api_user = $1`,
		},
		expectedArguments: [][]interface{}{
			{defaultUserId},
			{defaultUserId},
		},
	}

	suite.Run(t, &s)
}

func TestUserLimitRepository_DeleteForUser_NominalCase(t *testing.T) {
	assert := assert.New(t)

	repo := NewUserLimitRepository()
	mt := &mockTransaction{}

	err := repo.DeleteForUser(context.Background(), mt, defaultUserId)

	assert.Nil(err)
}
