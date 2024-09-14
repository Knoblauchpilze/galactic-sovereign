package service

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var someTime = time.Date(2024, 8, 17, 14, 31, 13, 651387247, time.UTC)

func Test_ActionService(t *testing.T) {
	s := ServiceTestSuite{
		generateRepositoriesMock: generateValidActionServiceMocks,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"processActionsUntil_listActions": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listBeforeCompletionTimeCalled)
					assert.Equal(someTime, m.listBeforeCompletionTime)
				},
			},
			"processActionsUntil_listActionsFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listBeforeCompletionTimeCalled)
				},
			},
			"processActionsUntil_listPlanetResources": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(2, m.listForPlanetCalled)
					assert.Equal(defaultBuildingAction.Planet, m.listForPlanetIds[0])
					assert.Equal(defaultBuildingAction.Planet, m.listForPlanetIds[1])
				},
			},
			"processActionsUntil_listPlanetResourcesFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"processActionsUntil_updatePlanetResources": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()

					resource := defaultPlanetResource
					resource.UpdatedAt = time.Now().Add(-1 * time.Minute)

					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: resource,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, time.Now())
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(2, m.updateCalled)
					assert.Equal(2, len(m.updatedPlanetResources))

					actual := m.updatedPlanetResources[0]

					assert.Equal(defaultPlanetResource.Planet, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					// Not asserting amount because it does not get updated in the mock
					assert.Equal(defaultPlanetResource.Production, actual.Production)
					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					// TODO: Improve this logic by allowing to specify the pool creation
					assert.True(actual.UpdatedAt.After(defaultPlanetResource.UpdatedAt))
					assert.Equal(defaultPlanetResource.Version, actual.Version)

					actual = m.updatedPlanetResources[1]
					assert.Equal(defaultPlanetResource.Planet, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					// Not asserting amount because it does not get updated in the mock
					// The production gets updated the value defined in the action
					assert.Equal(defaultBuildingActionResourceProduction.Production, actual.Production)
					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					// TODO: Improve the assertion logic for UpdatedAt
					assert.Equal(defaultPlanetResource.Version, actual.Version)

				},
			},
			"processActionsUntil_updatePlanetResourcesFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						updateErr: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"processActionsUntil_getBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForPlanetAndBuildingCalled)
					assert.Equal(defaultBuildingAction.Planet, m.getForPlanetAndBuildingPlanet)
					assert.Equal(defaultBuildingAction.Building, m.getForPlanetAndBuildingBuilding)
				},
			},
			"processActionsUntil_getBuildingFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.PlanetBuilding = &mockPlanetBuildingRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForPlanetAndBuildingCalled)
				},
			},
			"processActionsUntil_updateBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					expectedBuilding := persistence.PlanetBuilding{
						Planet:    defaultPlanetId,
						Building:  defaultBuildingId,
						Level:     defaultBuildingAction.DesiredLevel,
						CreatedAt: defaultBuilding.CreatedAt,
						UpdatedAt: defaultBuildingAction.CompletedAt,
					}
					assert.Equal(expectedBuilding, m.updateBuilding)
				},
			},
			"processActionsUntil_updateBuildingFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.PlanetBuilding = &mockPlanetBuildingRepository{
						planetBuilding: defaultPlanetBuilding,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"processActionsUntil_deleteActionResourceProduction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteForActionId)
				},
			},
			"processActionsUntil_deleteActionResourceProductionFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.BuildingActionResourceProduction = &mockBuildingActionResourceProductionRepository{
						actionResourceProduction: defaultBuildingActionResourceProduction,
						errs: []error{
							nil,
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
				},
			},
			"processActionsUntil_deleteActionCost": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteForActionId)
				},
			},
			"processActionsUntil_deleteActionCostFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.BuildingActionCost = &mockBuildingActionCostRepository{
						actionCost: defaultBuildingActionCost,
						errs: []error{
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
				},
			},
			"processActionsUntil_deleteAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteId)
				},
			},
			"processActionsUntil_deleteActionFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						action: defaultBuildingAction,
						errs: []error{
							nil,
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
				},
			},
		},

		transactionInteractionTestCases: map[string]transactionInteractionTestCase{
			"processActionsUntil_createsTwoTransactionAndClosesThem": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				verifyInteractions: func(pool db.ConnectionPool, assert *require.Assertions) {
					m := assertConnectionPoolIsAMock(pool, assert)

					assert.Equal(2, len(m.txs))
					for _, tx := range m.txs {
						assert.Equal(1, tx.closeCalled)
					}
				},
			},
			"processActionsUntil_whenFirstTransactionFails_returnsError": {
				generateConnectionPoolMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						errs: []error{
							errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(pool db.ConnectionPool, assert *require.Assertions) {
					m := assertConnectionPoolIsAMock(pool, assert)

					assert.Equal(1, len(m.txs))
				},
			},
			"processActionsUntil_whenFetchingActionsFail_expectSingleTransactionToBCreated": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidActionServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(pool db.ConnectionPool, assert *require.Assertions) {
					m := assertConnectionPoolIsAMock(pool, assert)

					assert.Equal(1, len(m.txs))
				},
			},
			"processActionsUntil_whenFailureToCreateTransactionForAction_expectPlanetResourcesNotUpdated": {
				generateConnectionPoolMock: func() db.ConnectionPool {
					return &mockConnectionPool{
						errs: []error{
							nil,
							errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, someTime)
				},
				expectedError: errDefault,
				verifyMockInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(0, m.listForPlanetCalled)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidActionServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		BuildingAction: &mockBuildingActionRepository{
			action: defaultBuildingAction,
		},
		BuildingActionCost: &mockBuildingActionCostRepository{
			actionCost: defaultBuildingActionCost,
		},
		BuildingActionResourceProduction: &mockBuildingActionResourceProductionRepository{
			actionResourceProduction: defaultBuildingActionResourceProduction,
		},
		PlanetBuilding: &mockPlanetBuildingRepository{
			planetBuilding: defaultPlanetBuilding,
		},
		PlanetResource: &mockPlanetResourceRepository{
			planetResource: defaultPlanetResource,
		},
	}
}

func assertConnectionPoolIsAMock(pool db.ConnectionPool, assert *require.Assertions) *mockConnectionPool {
	m, ok := pool.(*mockConnectionPool)
	if !ok {
		assert.Fail("Provided connection pool is not a mock")
	}
	return m
}
