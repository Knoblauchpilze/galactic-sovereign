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
	Id:   metalResourceId,
	Name: defaultResourceName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}
var defaultBuildingId = uuid.MustParse("5ec0f2cb-adc9-4f09-bb77-61d0ccdbcc52")
var defaultBuildingName = "my-building"
var defaultBuilding = persistence.Building{
	Id:   defaultBuildingId,
	Name: defaultBuildingName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}
var defaultBuildingCost = persistence.BuildingCost{
	Building: defaultBuildingId,
	Resource: metalResourceId,
	Cost:     250,
	Progress: 1.5,
}
var defaultBuildingResourceProduction = persistence.BuildingResourceProduction{
	Building: defaultBuildingId,
	Resource: metalResourceId,
	Base:     30,
	Progress: 1.1,
}

func Test_UniverseService(t *testing.T) {
	s := ServicePoolTestSuite{
		generateRepositoriesMocks:      generateValidUniverseServiceMocks,
		generateErrorRepositoriesMocks: generateErrorUniverseRepositoryMock,

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
				generateRepositoriesMocks: generateErrorUniverseRepositoryMock,
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
				generateRepositoriesMocks: generateErrorUniverseRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				expectedError: errDefault,
			},
			"get_resourceRepository": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"get_resourceRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateValidUniverseServiceMocks()
					repos.Resource = &mockResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				expectedError: errDefault,
			},
			"get_buildingRepository": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"get_buildingRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateValidUniverseServiceMocks()
					repos.Building = &mockBuildingRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				expectedError: errDefault,
			},
			"get_buildingCostRepository": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForBuildingCalled)
				},
			},
			"get_buildingCostRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateValidUniverseServiceMocks()
					repos.BuildingCost = &mockBuildingCostRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				expectedError: errDefault,
			},
			"get_buildingResourceProductionRepository": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewUniverseService(pool, repos)
					_, err := s.Get(ctx, defaultUniverseId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForBuildingCalled)
				},
			},
			"get_buildingResourceProductionRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateValidUniverseServiceMocks()
					repos.BuildingResourceProduction = &mockBuildingResourceProductionRepository{
						err: errDefault,
					}

					return repos
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
				generateRepositoriesMocks: generateErrorUniverseRepositoryMock,
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
				generateRepositoriesMocks: generateErrorUniverseRepositoryMock,
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
					Buildings: []communication.FullBuildingDtoResponse{
						{
							BuildingDtoResponse: communication.BuildingDtoResponse{
								Id:   defaultBuilding.Id,
								Name: defaultBuilding.Name,

								CreatedAt: defaultBuilding.CreatedAt,
							},
							Costs: []communication.BuildingCostDtoResponse{
								{
									Building: defaultBuilding.Id,
									Resource: metalResourceId,
									Cost:     250,
									Progress: 1.5,
								},
							},
							Productions: []communication.BuildingResourceProductionDtoResponse{
								{
									Building: defaultBuilding.Id,
									Resource: metalResourceId,
									Base:     30,
									Progress: 1.1,
								},
							},
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

func generateValidUniverseServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		Building: &mockBuildingRepository{
			building: defaultBuilding,
		},
		BuildingCost: &mockBuildingCostRepository{
			buildingCost: defaultBuildingCost,
		},
		BuildingResourceProduction: &mockBuildingResourceProductionRepository{
			buildingResourceProduction: defaultBuildingResourceProduction,
		},
		Resource: &mockResourceRepository{
			resources: []persistence.Resource{defaultResource},
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

func assertResourceRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockResourceRepository {
	m, ok := repos.Resource.(*mockResourceRepository)
	if !ok {
		assert.Fail("Provided resource repository is not a mock")
	}
	return m
}

func assertBuildingRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockBuildingRepository {
	m, ok := repos.Building.(*mockBuildingRepository)
	if !ok {
		assert.Fail("Provided building repository is not a mock")
	}
	return m
}

func assertBuildingCostRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockBuildingCostRepository {
	m, ok := repos.BuildingCost.(*mockBuildingCostRepository)
	if !ok {
		assert.Fail("Provided building cost repository is not a mock")
	}
	return m
}

func assertBuildingResourceProductionRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockBuildingResourceProductionRepository {
	m, ok := repos.BuildingResourceProduction.(*mockBuildingResourceProductionRepository)
	if !ok {
		assert.Fail("Provided building resiyrce production repository is not a mock")
	}
	return m
}
