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
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listBeforeCompletionTimeCalled)
				},
				expectedError: errDefault,
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
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForPlanetAndBuildingCalled)
				},
				expectedError: errDefault,
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
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
				expectedError: errDefault,
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
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
				},
				expectedError: errDefault,
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
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
				},
				expectedError: errDefault,
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
		PlanetBuilding: &mockPlanetBuildingRepository{
			planetBuilding: defaultPlanetBuilding,
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
