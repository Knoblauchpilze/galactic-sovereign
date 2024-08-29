package service

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultBuildingActionDtoRequest = communication.BuildingActionDtoRequest{
	Planet:   defaultPlanetId,
	Building: defaultBuildingId,
}
var defaultBuildingAction = persistence.BuildingAction{
	Id:           defaultBuildingActionId,
	Planet:       defaultPlanetId,
	Building:     defaultBuildingId,
	CurrentLevel: defaultPlanetBuilding.Level,
	DesiredLevel: defaultPlanetBuilding.Level + 1,
	CreatedAt:    testDate,
	CompletedAt:  testDate,
}
var defaultBuildingActionCost = persistence.BuildingActionCost{
	Action:   defaultBuildingActionId,
	Resource: metalResourceId,
	Amount:   36,
}

var metalResourceId = uuid.MustParse("8ed8d1f2-f39a-404b-96e1-9805ae6cd175")
var crystalResourceId = uuid.MustParse("5caf0c30-3417-49d3-94ac-8476aaf460c2")
var defaultResources = []persistence.Resource{
	{
		Id:   metalResourceId,
		Name: "metal",
	},
	{
		Id:   crystalResourceId,
		Name: "crystal",
	},
}

func Test_BuildingActionService(t *testing.T) {
	beforeTestSuite := time.Now()

	s := ServiceTestSuite{
		generateRepositoriesMock:      generateValidBuildingActionRepositoryMock,
		generateErrorRepositoriesMock: generateErrorBuildingActionRepositoryMock,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					actual := m.createdBuildingAction

					assert.Equal(1, m.createCalled)
					assert.Nil(uuid.Validate(actual.Id.String()))
					assert.Equal(defaultPlanetId, actual.Planet)
					assert.Equal(defaultBuildingId, actual.Building)
					assert.Equal(defaultPlanetBuilding.Level, actual.CurrentLevel)
					assert.Equal(defaultPlanetBuilding.Level+1, actual.DesiredLevel)
					assert.True(beforeTestSuite.Before(actual.CreatedAt))
					expectedCompletionTime := actual.CreatedAt.Add(6 * time.Minute)
					assert.Equal(expectedCompletionTime, actual.CompletedAt)
				},
			},
			"create_resource": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"create_resourceFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Resource: &mockResourceRepository{
							err: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"create_planetResource": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultPlanetId, m.listForPlanetId)
				},
			},
			"create_planetResourceFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Resource: &mockResourceRepository{
							resources: defaultResources,
						},
						PlanetResource: &mockPlanetResourceRepository{
							err: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"create_planetBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultPlanetId, m.listForPlanetId)
				},
			},
			"create_planetBuildingFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Resource: &mockResourceRepository{
							resources: defaultResources,
						},
						PlanetResource: &mockPlanetResourceRepository{
							planetResource: defaultPlanetResource,
						},
						PlanetBuilding: &mockPlanetBuildingRepository{
							err: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"create_buildingCost": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForBuildingCalled)
					assert.Equal(defaultBuildingId, m.listForBuildingId)
				},
			},
			"create_buildingCostFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Resource: &mockResourceRepository{
							resources: defaultResources,
						},
						PlanetResource: &mockPlanetResourceRepository{
							planetResource: defaultPlanetResource,
						},
						PlanetBuilding: &mockPlanetBuildingRepository{
							planetBuilding: defaultPlanetBuilding,
						},
						BuildingCost: &mockBuildingCostRepository{
							err: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForBuildingCalled)
				},
			},
			"create_buildingActionCost": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.PlanetBuilding, _ []persistence.BuildingCost) error {
						return nil
					}
					consolidator := func(action persistence.BuildingAction, _ []persistence.PlanetBuilding, _ []persistence.Resource, _ []persistence.BuildingCost) (persistence.BuildingAction, []persistence.BuildingActionCost, error) {
						return action, []persistence.BuildingActionCost{defaultBuildingActionCost}, nil
					}
					s := newBuildingActionService(pool, repos, validator, consolidator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultBuildingActionCost, m.createdBuildingActionCost)
				},
			},
			"create_buildingActionCostFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Resource: &mockResourceRepository{
							resources: defaultResources,
						},
						PlanetResource: &mockPlanetResourceRepository{
							planetResource: defaultPlanetResource,
						},
						PlanetBuilding: &mockPlanetBuildingRepository{
							planetBuilding: defaultPlanetBuilding,
						},
						BuildingCost: &mockBuildingCostRepository{
							buildingCost: defaultBuildingCost,
						},
						BuildingAction: &mockBuildingActionRepository{
							action: defaultBuildingAction,
						},
						BuildingActionCost: &mockBuildingActionCostRepository{
							err: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.PlanetBuilding, _ []persistence.BuildingCost) error {
						return nil
					}
					consolidator := func(action persistence.BuildingAction, _ []persistence.PlanetBuilding, _ []persistence.Resource, _ []persistence.BuildingCost) (persistence.BuildingAction, []persistence.BuildingActionCost, error) {
						return action, []persistence.BuildingActionCost{defaultBuildingActionCost}, nil
					}
					s := newBuildingActionService(pool, repos, validator, consolidator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
				},
			},
			"create_validationFails": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.PlanetBuilding, _ []persistence.BuildingCost) error {
						return errDefault
					}
					consolidator := func(action persistence.BuildingAction, _ []persistence.PlanetBuilding, _ []persistence.Resource, _ []persistence.BuildingCost) (persistence.BuildingAction, []persistence.BuildingActionCost, error) {
						return action, []persistence.BuildingActionCost{}, nil
					}
					s := newBuildingActionService(pool, repos, validator, consolidator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"create_consolidationFails": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.PlanetBuilding, _ []persistence.BuildingCost) error {
						return nil
					}
					consolidator := func(action persistence.BuildingAction, _ []persistence.PlanetBuilding, _ []persistence.Resource, _ []persistence.BuildingCost) (persistence.BuildingAction, []persistence.BuildingActionCost, error) {
						return action, []persistence.BuildingActionCost{}, errDefault
					}
					s := newBuildingActionService(pool, repos, validator, consolidator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"create_repositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Resource: &mockResourceRepository{
							resources: defaultResources,
						},
						PlanetResource: &mockPlanetResourceRepository{
							planetResource: defaultPlanetResource,
						},
						PlanetBuilding: &mockPlanetBuildingRepository{
							planetBuilding: defaultPlanetBuilding,
						},
						BuildingCost: &mockBuildingCostRepository{
							buildingCost: defaultBuildingCost,
						},
						BuildingAction: &mockBuildingActionRepository{
							errs: []error{errDefault},
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingActionId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultBuildingActionId, m.deleteId)
				},
			},
		},

		returnTestCases: map[string]returnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewBuildingActionService(pool, repos)
					out, _ := s.Create(ctx, defaultBuildingActionDtoRequest)
					return out
				},
				expectedContent: communication.BuildingActionDtoResponse{
					Id:           defaultBuildingAction.Id,
					Planet:       defaultPlanetId,
					Building:     defaultBuildingId,
					CurrentLevel: defaultBuildingAction.CurrentLevel,
					DesiredLevel: defaultBuildingAction.DesiredLevel,
					CreatedAt:    defaultBuildingAction.CreatedAt,
					CompletedAt:  defaultBuildingAction.CompletedAt,
				},
			},
		},

		transactionTestCases: map[string]transactionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
			},
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingActionId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidBuildingActionRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Resource: &mockResourceRepository{
			resources: defaultResources,
		},
		PlanetResource: &mockPlanetResourceRepository{
			planetResource: defaultPlanetResource,
		},
		PlanetBuilding: &mockPlanetBuildingRepository{
			planetBuilding: defaultPlanetBuilding,
		},
		BuildingCost: &mockBuildingCostRepository{
			buildingCost: defaultBuildingCost,
		},
		BuildingAction: &mockBuildingActionRepository{
			action: defaultBuildingAction,
		},
		BuildingActionCost: &mockBuildingActionCostRepository{
			actionCost: defaultBuildingActionCost,
		},
	}
}

func generateErrorBuildingActionRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Resource: &mockResourceRepository{
			err: errDefault,
		},
	}
}

func assertBuildingActionCostRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockBuildingActionCostRepository {
	m, ok := repos.BuildingActionCost.(*mockBuildingActionCostRepository)
	if !ok {
		assert.Fail("Provided building action cost repository is not a mock")
	}
	return m
}
