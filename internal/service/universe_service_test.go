package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultUniverseId = uuid.MustParse("3e7fde5c-ac70-4e5d-bd09-73029725048d")
var defaultUniverseName = "my-universe"

var defaultUniverseDtoRequest = communication.UniverseDtoRequest{
	Name: defaultUniverseName,
}
var defaultUniverse = persistence.Universe{
	Id:   defaultUniverseId,
	Name: defaultUniverseName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func Test_UniverseService(t *testing.T) {
	s := ServiceTestSuite{
		generateRepositoriesMock:      generateValidUniverseRepositoryMock,
		generateErrorRepositoriesMock: generateErrorUniverseRepositoryMock,

		errorTestCases: map[string]errorTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Create(ctx, defaultUniverseDtoRequest)
					return err
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
			},
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					return s.Delete(ctx, defaultUniverseId)
				},
			},
		},

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Create(ctx, defaultUniverseDtoRequest)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUniverseRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultUniverseDtoRequest.Name, m.createdUniverse.Name)
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUniverseRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultUniverseId, m.getId)
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.List(ctx)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUniverseRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					return s.Delete(ctx, defaultUniverseId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertUniverseRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultUniverseId, m.deleteId)
				},
			},
		},

		returnTestCases: map[string]returnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewUniverseService(pool, repos)
					out, _ := s.Create(ctx, defaultUniverseDtoRequest)
					return out
				},

				expectedContent: communication.UniverseDtoResponse{
					Id:   defaultUniverse.Id,
					Name: defaultUniverse.Name,

					CreatedAt: defaultUniverse.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewUniverseService(pool, repos)
					out, _ := s.Get(ctx, defaultUniverseId)
					return out
				},

				expectedContent: communication.UniverseDtoResponse{
					Id:   defaultUniverse.Id,
					Name: defaultUniverse.Name,

					CreatedAt: defaultUniverse.CreatedAt,
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewUniverseService(pool, repos)
					out, _ := s.List(ctx)
					return out
				},

				expectedContent: []communication.UniverseDtoResponse{
					{
						Id:   defaultUniverse.Id,
						Name: defaultUniverse.Name,

						CreatedAt: defaultUniverse.CreatedAt,
					},
				},
			},
		},

		transactionTestCases: map[string]transactionTestCase{
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					return s.Delete(ctx, defaultUniverseId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidUniverseRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Universe: &mockUniverseRepository{
			universe: defaultUniverse,
		},
	}
}

func generateErrorUniverseRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Universe: &mockUniverseRepository{
			err: errDefault,
		},
	}
}

func assertUniverseRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockUniverseRepository {
	m, ok := repos.Universe.(*mockUniverseRepository)
	if !ok {
		assert.Fail("Provided universe repository is not a mock")
	}
	return m
}
