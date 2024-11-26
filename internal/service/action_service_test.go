package service

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var someTime = time.Date(2024, 8, 17, 14, 31, 13, 651387247, time.UTC)

func TestUnit__ActionService(t *testing.T) {
	s := ServicePoolTestSuite{
		generateRepositoriesMocks: generateActionServiceMocks,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"processActionsUntil_listActions": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(defaultPlanetId, m.listBeforeCompletionTimePlanet)
					assert.Equal(1, m.listBeforeCompletionTimeCalled)
					assert.Equal(someTime, m.listBeforeCompletionTime)
				},
			},
			"processActionsUntil_listActionsFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultBuildingAction.Planet, m.listForPlanetIds[0])
				},
			},
			"processActionsUntil_listPlanetResourcesFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"processActionsUntil_listPlanetResourceProductions": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultBuildingAction.Planet, m.listForPlanetIds[0])
				},
			},
			"processActionsUntil_listPlanetResourceProductionsFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"processActionsUntil_updatePlanetResources": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()

					resource := defaultPlanetResource
					resource.UpdatedAt = defaultBuildingAction.CompletedAt.Add(-2 * time.Minute)

					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: resource,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))

					actual := m.updatedPlanetResources[0]

					assert.Equal(defaultPlanetResource.Planet, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					expectedAmount := defaultPlanetResource.Amount + 2.0/60.0*float64(defaultPlanetResourceProduction.Production)
					assert.Equal(expectedAmount, actual.Amount)
					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					assert.Equal(defaultBuildingAction.CompletedAt, actual.UpdatedAt)
					assert.Equal(defaultPlanetResource.Version, actual.Version)
				},
			},
			"processActionsUntil_updatePlanetResourcesFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: defaultPlanetResource,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"processActionsUntil_updatePlanetResources_whenResourceIsNotProduced_expectNoUpdate": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()

					production := defaultPlanetResourceProduction
					production.Resource = crystalResourceId

					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						planetResourceProduction: production,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(0, m.updateCalled)
				},
			},
			"processActionsUntil_updatePlanetResources_whenStorageIsAlreadyFull_expectNoUpdate": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()

					resource := defaultPlanetResource
					resource.UpdatedAt = defaultBuildingAction.CompletedAt.Add(-2 * time.Minute)
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: resource,
					}

					storage := defaultPlanetResourceStorage
					storage.Storage = int(defaultPlanetResource.Amount) - 10
					repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
						planetResourceStorage: storage,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))
					actual := m.updatedPlanetResources[0]

					assert.Equal(defaultPlanetId, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					assert.Equal(defaultPlanetResource.Amount, actual.Amount)

					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					assert.Equal(defaultBuildingAction.CompletedAt, actual.UpdatedAt)
					assert.Equal(defaultPlanetResource.Version, actual.Version)
				},
			},
			"whenUpdatingPlanetResources_whenStorageIsNotEnoughToAbsorbAllProduction_expectPartialUpdate": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()

					resource := defaultPlanetResource
					resource.Amount = 60.0
					resource.UpdatedAt = defaultBuildingAction.CompletedAt.Add(-2 * time.Minute)
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: resource,
					}

					production := defaultPlanetResourceProduction
					production.Production = 60
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						planetResourceProduction: production,
					}

					storage := defaultPlanetResourceStorage
					storage.Storage = 62
					repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
						planetResourceStorage: storage,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))
					actual := m.updatedPlanetResources[0]

					assert.Equal(defaultPlanetId, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					assert.Equal(62.0, actual.Amount)

					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					assert.Equal(defaultBuildingAction.CompletedAt, actual.UpdatedAt)
					assert.Equal(defaultPlanetResource.Version, actual.Version)
				},
			},
			"processActionsUntil_getBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForPlanetAndBuildingCalled)
					assert.Equal(defaultBuildingAction.Planet, m.getForPlanetAndBuildingPlanet)
					assert.Equal(defaultBuildingAction.Building, m.getForPlanetAndBuildingBuilding)
				},
			},
			"processActionsUntil_getBuildingFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetBuilding = &mockPlanetBuildingRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetBuilding = &mockPlanetBuildingRepository{
						planetBuilding: defaultPlanetBuilding,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"processActionsUntil_getForPlanetAndBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForPlanetAndBuildingCalled)
					assert.Equal(defaultBuildingAction.Planet, m.getForPlanetAndBuildingPlanet)
					assert.Equal(&defaultBuildingAction.Building, m.getForPlanetAndBuildingBuilding)
				},
			},
			"processActionsUntil_getForPlanetAndBuildingFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						planetResourceProduction: defaultPlanetResourceProduction,
						errs: []error{
							nil,
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForPlanetAndBuildingCalled)
				},
			},
			"processActionsUntil_listProductionForAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.listForActionId)
				},
			},
			"processActionsUntil_listProductionForActionFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.BuildingActionResourceProduction = &mockBuildingActionResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForActionCalled)
				},
			},
			"processActionsUntil_getsExistingProductionForPlanetAndBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForPlanetAndBuildingCalled)
					assert.Equal(defaultBuildingAction.Planet, m.getForPlanetAndBuildingPlanet)
					assert.Equal(&defaultBuildingAction.Building, m.getForPlanetAndBuildingBuilding)
				},
			},
			"processActionsUntil_getsExistingProductionForPlanetAndBuildingFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						errs: []error{
							nil,
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getForPlanetAndBuildingCalled)
				},
			},
			"processActionsUntil_whenResourceIsNotProduced_createsResourceProduction": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						errs: []error{
							nil,
							errors.NewCode(db.NoMatchingSqlRows),
							nil,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, time.Now())
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)

					actual := m.createdPlanetResourceProduction

					assert.Equal(defaultBuildingAction.Planet, actual.Planet)
					assert.Equal(defaultBuildingActionResourceProduction.Resource, actual.Resource)
					assert.Equal(&defaultBuildingAction.Building, actual.Building)
					assert.Equal(defaultBuildingActionResourceProduction.Production, actual.Production)
					assert.Equal(defaultBuildingAction.CompletedAt, actual.CreatedAt)
					assert.Equal(defaultBuildingAction.CompletedAt, actual.UpdatedAt)
					assert.Equal(0, actual.Version)
				},
			},
			"processActionsUntil_whenResourceIsNotProducedAndResourceProductionCreationFails_expectError": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						errs: []error{
							nil,
							errors.NewCode(db.NoMatchingSqlRows),
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
				},
			},
			"processActionsUntil_whenResourceIsAlreadyProduced_updatesResourceProduction": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()

					production := defaultPlanetResourceProduction
					production.UpdatedAt = time.Now().Add(-1 * time.Minute)

					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						planetResourceProduction: production,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)

					actual := m.updatedPlanetResourceProductions[0]

					assert.Equal(defaultPlanetResourceProduction.Planet, actual.Planet)
					assert.Equal(defaultPlanetResourceProduction.Building, actual.Building)
					assert.Equal(defaultPlanetResourceProduction.Resource, actual.Resource)
					assert.Equal(defaultBuildingActionResourceProduction.Production, actual.Production)
					assert.Equal(defaultPlanetResourceProduction.CreatedAt, actual.CreatedAt)
					assert.Equal(defaultBuildingAction.CompletedAt, actual.UpdatedAt)
					assert.Equal(defaultPlanetResourceProduction.Version, actual.Version)
				},
			},
			"processActionsUntil_whenResourceIsAlreadyProducedAndResourceProductionUpdateFails_expectError": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()

					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						planetResourceProduction: defaultPlanetResourceProduction,
						updateErr:                errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"processActionsUntil_deleteActionResourceProduction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteForActionId)
				},
			},
			"processActionsUntil_deleteActionResourceProductionFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteForActionId)
				},
			},
			"processActionsUntil_deleteActionCostFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteId)
				},
			},
			"processActionsUntil_deleteActionFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(pool db.ConnectionPool, assert *require.Assertions) {
					m := assertConnectionPoolIsAMock(pool, assert)

					assert.Equal(1, len(m.txs))
				},
			},
			"processActionsUntil_whenFetchingActionsFail_expectSingleTransactionToBCreated": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateActionServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewActionService(pool, repos)
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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
					return s.ProcessActionsUntil(ctx, defaultPlanetId, someTime)
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

func generateActionServiceMocks() repositories.Repositories {
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
		PlanetResourceProduction: &mockPlanetResourceProductionRepository{
			planetResourceProduction: defaultPlanetResourceProduction,
		},
		PlanetResourceStorage: &mockPlanetResourceStorageRepository{
			planetResourceStorage: defaultPlanetResourceStorage,
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
