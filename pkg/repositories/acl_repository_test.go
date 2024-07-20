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

var defaultAclId = uuid.MustParse("fd538d90-cf52-41f9-8ed5-69b1afea5964")
var defaultAcl = persistence.Acl{
	Id:   defaultAclId,
	User: defaultUserId,

	Resource: "some-resource",
	Permissions: []string{
		"GET",
		"PATCH",
	},

	CreatedAt: time.Date(2024, 06, 22, 9, 58, 20, 651387237, time.UTC),
	UpdatedAt: time.Date(2024, 06, 22, 9, 58, 40, 651387237, time.UTC),
}

func Test_AclRepository(t *testing.T) {
	dummyStr := ""

	s := RepositoryTransactionTestSuiteNew{
		dbInteractionTestCases: map[string]dbTransactionInteractionTestCase{
			"create": {
				sqlMode: QueryBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewAclRepository()
					_, err := s.Create(ctx, tx, defaultAcl)
					return err
				},
				expectedSqlQueries: []string{
					`
INSERT INTO acl (id, api_user, resource)
	VALUES($1, $2, $3)
	ON CONFLICT (api_user, resource) DO NOTHING
	RETURNING
		acl.id
`,
					`
INSERT INTO acl_permission (acl, permission)
	VALUES($1, $2)
`,
					`
INSERT INTO acl_permission (acl, permission)
	VALUES($1, $2)
`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultAcl.Id,
						defaultAcl.User,
						defaultAcl.Resource,
					},
					{"GET"},
					{"PATCH"},
				},
			},
			"get": {
				sqlMode: QueryBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewAclRepository()
					_, err := s.Get(ctx, tx, defaultAclId)
					return err
				},
				expectedSqlQueries: []string{
					`
SELECT
	id,
	api_user,
	resource,
	created_at,
	updated_at
FROM
	acl
WHERE
	id = $1
`,
					`SELECT permission FROM acl_permission WHERE acl = $1`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultAclId,
					},
					{
						defaultAclId,
					},
				},
			},
			"getForUser": {
				sqlMode: QueryBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewAclRepository()
					_, err := s.GetForUser(ctx, tx, defaultUserId)
					return err
				},
				expectedSqlQueries: []string{
					`SELECT id FROM acl WHERE api_user = $1`,
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
					s := NewAclRepository()
					return s.Delete(ctx, tx, []uuid.UUID{defaultAclId})
				},
				expectedSqlQueries: []string{
					`DELETE FROM acl_permission WHERE acl in ($1)`,
					`DELETE FROM acl WHERE id IN ($1)`,
				},
				expectedArguments: [][]interface{}{
					{
						defaultAclId,
					},
					{
						defaultAclId,
					},
				},
			},
			"delete_multipleIds": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					ids := []uuid.UUID{
						uuid.MustParse("50714fb2-db52-4e3a-8315-cf8e4a8abcf8"),
						uuid.MustParse("9fc0def1-d51c-4af0-8db5-40310796d16d"),
					}
					s := NewAclRepository()
					return s.Delete(ctx, tx, ids)
				},
				expectedSqlQueries: []string{
					`DELETE FROM acl_permission WHERE acl in ($1,$2)`,
					`DELETE FROM acl WHERE id IN ($1,$2)`,
				},
				expectedArguments: [][]interface{}{
					{
						uuid.MustParse("50714fb2-db52-4e3a-8315-cf8e4a8abcf8"),
						uuid.MustParse("9fc0def1-d51c-4af0-8db5-40310796d16d"),
					},
					{
						uuid.MustParse("50714fb2-db52-4e3a-8315-cf8e4a8abcf8"),
						uuid.MustParse("9fc0def1-d51c-4af0-8db5-40310796d16d"),
					},
				},
			},
			"deleteForUser": {
				sqlMode: ExecBased,
				handler: func(ctx context.Context, tx db.Transaction) error {
					s := NewAclRepository()
					return s.DeleteForUser(ctx, tx, defaultUserId)
				},
				expectedSqlQueries: []string{
					`
DELETE FROM acl_permission
	WHERE acl
		IN (SELECT id FROM acl WHERE api_user = $1)
`,
					`DELETE FROM acl WHERE api_user = $1`,
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
					repo := NewAclRepository()
					_, err := repo.Create(ctx, tx, defaultAcl)
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
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewAclRepository()
					_, err := repo.Get(ctx, tx, defaultAclId)
					return err
				},
				expectedGetSingleValueCalls: 1,
				expectedScanCalls:           2,
				expectedScannedProps: [][]interface{}{
					{
						&uuid.UUID{},
						&uuid.UUID{},
						&dummyStr,
						&time.Time{},
						&time.Time{},
					},
					{
						&dummyStr,
					},
				},
			},
		},

		dbGetAllTestCases: map[string]dbTransactionGetAllTestCase{
			"getForUser": {
				handler: func(ctx context.Context, tx db.Transaction) error {
					repo := NewAclRepository()
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
					s := NewAclRepository()
					out, _ := s.Create(ctx, tx, defaultAcl)
					return out
				},
				expectedContent: defaultAcl,
			},
		},
	}

	suite.Run(t, &s)
}
