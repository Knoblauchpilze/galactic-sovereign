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

func Test_UserLimitRepository(t *testing.T) {
	dummyStr := ""

	s := RepositoryTransactionTestSuiteNew{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: QueryBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUserLimitRepository()
					_, err := s.Create(ctx, tx, defaultUserLimit)
					return err
				},
				expectedSqlQueries: []string{
					`
INSERT INTO user_limit (id, name, api_user)
	VALUES($1, $2, $3)
	ON CONFLICT (name, api_user) DO NOTHING
	RETURNING
		user_limit.id
`,
					`
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
`,
					`
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
`,
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
			},
			"get": {
				sqlMode: QueryBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUserLimitRepository()
					_, err := s.Get(ctx, tx, defaultUserLimitId)
					return err
				},
				expectedSqlQueries: []string{
					`
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
`,
					`
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
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUserLimitId,
					},
					{
						defaultUserLimitId,
					},
				},
			},
			"getForUser": {
				sqlMode: QueryBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUserLimitRepository()
					_, err := s.GetForUser(ctx, tx, defaultUserId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id FROM user_limit WHERE api_user = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUserId,
					},
				},
			},
			"delete_singleId": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUserLimitRepository()
					return s.Delete(ctx, tx, []uuid.UUID{defaultUserLimitId})
				},
				expectedSqlQueries: []string{
					`DELETE FROM limits WHERE user_limit in ($1)`,
					`DELETE FROM user_limit WHERE id IN ($1)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUserLimitId,
					},
					{
						defaultUserLimitId,
					},
				},
			},
			"delete_multipleIds": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					ids := []uuid.UUID{
						uuid.MustParse("4daffd0f-c71b-4afb-bcee-730149c0ecad"),
						uuid.MustParse("c6dea18c-f8ed-4b50-9ba2-d2ec10660d9f"),
					}
					s := NewUserLimitRepository()
					return s.Delete(ctx, tx, ids)
				},
				expectedSqlQueries: []string{
					`DELETE FROM limits WHERE user_limit in ($1,$2)`,
					`DELETE FROM user_limit WHERE id IN ($1,$2)`,
				},
				expectedArguments: [][]interface{}{
					{
						uuid.MustParse("4daffd0f-c71b-4afb-bcee-730149c0ecad"),
						uuid.MustParse("c6dea18c-f8ed-4b50-9ba2-d2ec10660d9f"),
					},
					{
						uuid.MustParse("4daffd0f-c71b-4afb-bcee-730149c0ecad"),
						uuid.MustParse("c6dea18c-f8ed-4b50-9ba2-d2ec10660d9f"),
					},
				},
			},
			"deleteForUser": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewUserLimitRepository()
					return s.DeleteForUser(ctx, tx, defaultUserId)
				},
				expectedSqlQueries: []string{
					`
DELETE FROM limits
	WHERE user_limit
		IN (SELECT id FROM user_limit WHERE api_user = $1)
`,
					`DELETE FROM user_limit WHERE api_user = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultUserId,
					},
					{
						defaultUserId,
					},
				},
			},
		},

		dbSingleValueTestCases: map[string]dbTransactionSingleValueTestCase{
			"create": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewUserLimitRepository()
					_, err := repo.Create(ctx, tx, defaultUserLimit)
					return err
				},
				expectedGetSingleValueCalls: 3,
				expectedScanCalls:           3,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
					},
					{
						&uuid.UUID{},
					},
					{
						&uuid.UUID{},
					},
				},
			},
			"get": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewUserLimitRepository()
					_, err := repo.Get(ctx, tx, defaultAclId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           2,
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
			},
		},

		dbGetAllTestCases: map[string]dbTransactionGetAllTestCase{
			"getForUser": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewUserLimitRepository()
					_, err := repo.GetForUser(ctx, tx, defaultUserId)
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

		dbReturnTestCases: map[string]dbTransactionReturnTestCase{
			"create": {
				handler: func(ctx context.Context, tx db.Transaction) interface{} {
					repo := NewUserLimitRepository()
					out, _ := repo.Create(ctx, tx, defaultUserLimit)
					return out
				},
				expectedContent: defaultUserLimit,
			},
		},

		dbErrorTestCases: map[string]dbTransactionErrorTestCase{},
	}

	suite.Run(t, &s)
}
