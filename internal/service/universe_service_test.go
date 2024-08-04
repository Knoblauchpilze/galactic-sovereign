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
var defaultResourceName = "my-resource"
var defaultResource = persistence.Resource{
	Id:   defaultResourceId,
	Name: defaultResourceName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func Test_UniverseService(t *testing.T) {
	s := ServiceTestSuite{
		generateRepositoriesMock:      generateValidUniverseRepositoryMock,
		generateErrorRepositoriesMock: generateErrorUniverseRepositoryMock,

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
			"create_repositoryFails": {
				generateRepositoriesMock: generateErrorUniverseRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Create(ctx, defaultUniverseDtoRequest)
					return err
				},
				expectedError: errDefault,
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
			"get_universeRepositoryFails": {
				generateRepositoriesMock: generateErrorUniverseRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				expectedError: errDefault,
			},
			"get_resourceRepositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Resource: &mockResourceRepository{
							err: errDefault,
						},
						Universe: &mockUniverseRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				expectedError: errDefault,
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
			"list_repositoryFails": {
				generateRepositoriesMock: generateErrorUniverseRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
				expectedError: errDefault,
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
			"delete_repositoryFails": {
				generateRepositoriesMock: generateErrorUniverseRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					return s.Delete(ctx, defaultUniverseId)
				},
				expectedError: errDefault,
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

				expectedContent: communication.FullUniverseDtoResponse{
					UniverseDtoResponse: communication.UniverseDtoResponse{
						Id:   defaultUniverse.Id,
						Name: defaultUniverse.Name,

						CreatedAt: defaultUniverse.CreatedAt,
					},
					Resources: []communication.ResourceDtoResponse{
						{
							Id:   defaultResource.Id,
							Name: defaultResource.Name,

							CreatedAt: defaultResource.CreatedAt,
						},
					},
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
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
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
	}

	suite.Run(t, &s)
}

func generateValidUniverseRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Resource: &mockResourceRepository{
			resource: defaultResource,
		},
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
