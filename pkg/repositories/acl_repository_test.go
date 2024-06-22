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

func TestAclRepository_Create_DbInteraction(t *testing.T) {
	expectedSql := `
INSERT INTO acl (id, api_user, resource)
	VALUES($1, $2, $3)
	ON CONFLICT (api_user, resource) DO NOTHING
	RETURNING
		acl.id
`
	expectedPermissionSql := `
INSERT INTO acl_permissions (acl, permission)
	VALUES($1, $2)
`

	s := RepositoryTransactionTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewAclRepository(&mockConnectionPool{})
			_, err := repo.Create(context.Background(), tx, defaultAcl)
			return err
		},
		expectedSql: []string{
			expectedSql,
			expectedPermissionSql,
			expectedPermissionSql,
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
	}

	suite.Run(t, &s)
}

func TestAclRepository_Create_RetrievesGeneratedAcl(t *testing.T) {
	s := RepositorySingleValueTransactionTestSuite{
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewAclRepository(&mockConnectionPool{})
			_, err := repo.Create(ctx, tx, defaultAcl)
			return err
		},
		expectedScanCalls: 1,
		expectedScannedProps: [][]interface{}{
			{&uuid.UUID{}},
		},
	}

	suite.Run(t, &s)
}

func TestAclRepository_Create_ReturnsInputAcl(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewAclRepository(mc)
	mt := &mockTransaction{}

	actual, err := repo.Create(context.Background(), mt, defaultAcl)

	assert.Nil(err)
	assert.Equal(defaultAcl, actual)
}

func TestAclRepository_Get_DbInteraction(t *testing.T) {
	expectedSql := `
SELECT
	id,
	api_user,
	resource,
	created_at,
	updared_at
FROM
	acl
WHERE
	id = $1
`

	s := RepositoryTransactionTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewAclRepository(&mockConnectionPool{})
			_, err := repo.Get(context.Background(), tx, defaultAclId)
			return err
		},
		expectedSql: []string{
			expectedSql,
			`SELECT permission FROM acl_permissions WHERE acl = $1`,
		},
		expectedArguments: [][]interface{}{
			{defaultAclId},
			{defaultAclId},
		},
	}

	suite.Run(t, &s)
}

func TestAclRepository_Get_InterpretDbData(t *testing.T) {
	dummyStr := ""

	s := RepositorySingleValueTransactionTestSuite{
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewAclRepository(&mockConnectionPool{})
			_, err := repo.Get(ctx, tx, defaultAclId)
			return err
		},
		expectedScanCalls: 2,
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
	}

	suite.Run(t, &s)
}

func TestAclRepository_GetForUser_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: QueryBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewAclRepository(&mockConnectionPool{})
			_, err := repo.GetForUser(context.Background(), tx, defaultUserId)
			return err
		},
		expectedSql: []string{`SELECT id FROM acl WHERE api_user = $1`},
		expectedArguments: [][]interface{}{
			{
				defaultUserId,
			},
		},
	}

	suite.Run(t, &s)
}

func TestAclRepository_GetForUser_InterpretDbData(t *testing.T) {
	s := RepositoryGetAllTransactionTestSuite{
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewAclRepository(&mockConnectionPool{})
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

func TestAclRepository_Delete_SingleId_DbInteraction(t *testing.T) {
	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewAclRepository(&mockConnectionPool{})
			return repo.Delete(context.Background(), tx, []uuid.UUID{defaultAclId})
		},
		expectedSql: []string{
			`DELETE FROM acl_permissions WHERE acl in ($1)`,
			`DELETE FROM acl WHERE id IN ($1)`,
		},
		expectedArguments: [][]interface{}{
			{defaultAclId},
			{defaultAclId},
		},
	}

	suite.Run(t, &s)
}

func TestAclRepository_Delete_MultipleIds_DbInteraction(t *testing.T) {
	ids := []uuid.UUID{
		uuid.MustParse("50714fb2-db52-4e3a-8315-cf8e4a8abcf8"),
		uuid.MustParse("9fc0def1-d51c-4af0-8db5-40310796d16d"),
	}

	s := RepositoryTransactionTestSuite{
		sqlMode: ExecBased,
		testFunc: func(ctx context.Context, tx db.Transaction) error {
			repo := NewAclRepository(&mockConnectionPool{})
			return repo.Delete(context.Background(), tx, ids)
		},
		expectedSql: []string{
			`DELETE FROM acl_permissions WHERE acl in ($1,$2)`,
			`DELETE FROM acl WHERE id IN ($1,$2)`,
		},
		expectedArguments: [][]interface{}{
			{ids[0], ids[1]},
			{ids[0], ids[1]},
		},
	}

	suite.Run(t, &s)
}

func TestAclRepository_Delete_NominalCase(t *testing.T) {
	assert := assert.New(t)

	mc := &mockConnectionPool{}
	repo := NewAclRepository(mc)
	mt := &mockTransaction{}

	err := repo.Delete(context.Background(), mt, []uuid.UUID{defaultAclId})

	assert.Nil(err)
}
