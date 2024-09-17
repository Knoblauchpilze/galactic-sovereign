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
	Amount:   250,
}
var defaultBuildingActionResourceProduction = persistence.BuildingActionResourceProduction{
	Action:     defaultBuildingActionId,
	Resource:   metalResourceId,
	Production: 380,
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
			"create_listsResourcesOnPlanet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal([]uuid.UUID{defaultPlanetId}, m.listForPlanetIds)
				},
			},
			"create_listsResource": {
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
			"create_listsBuildingsOnPlanet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultBuildingActionDtoRequest.Planet, m.listForPlanetId)
				},
			},
			"create_listsCostsForBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForBuildingCalled)
					assert.Equal(defaultBuildingActionDtoRequest.Building, m.listForBuildingId)
				},
			},
			"create_listsProductionsForBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForBuildingCalled)
					assert.Equal(defaultBuildingActionDtoRequest.Building, m.listForBuildingId)
				},
			},
			"create_updatesResourcesOnPlanet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))
					actual := m.updatedPlanetResources[0]
					assert.Equal(defaultBuildingActionDtoRequest.Planet, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					expectedCost := 562.0
					expectedAmount := defaultPlanetResource.Amount - expectedCost
					assert.Equal(expectedAmount, actual.Amount)
					assert.Equal(defaultPlanetResource.Version, actual.Version)
				},
			},
			"create_registersActionCosts": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultBuildingCost.Resource, m.createdBuildingActionCost.Resource)
					assert.Equal(562, m.createdBuildingActionCost.Amount)
				},
			},
			"create_registersActionProductions": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultBuildingActionResourceProduction.Resource, m.createdBuildingActionResourceProduction.Resource)
					assert.Equal(36, m.createdBuildingActionResourceProduction.Production)
				},
			},
			"create_createsAction": {
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
					expectedDuration, err := time.ParseDuration("13m29s280ms")
					assert.Nil(err)
					expectedCompletionTime := actual.CreatedAt.Add(expectedDuration)
					assert.Equal(expectedCompletionTime, actual.CompletedAt)
				},
			},
			"create_failure_listResourcesOnPlanet": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
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
			"create_failure_listResource": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.Resource = &mockResourceRepository{
						err: errDefault,
					}

					return repos
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
			"create_failure_listBuildingsOnPlanet": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.PlanetBuilding = &mockPlanetBuildingRepository{
						err: errDefault,
					}

					return repos
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
			"create_failure_listCostsForBuilding": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingCost = &mockBuildingCostRepository{
						err: errDefault,
					}

					return repos
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
			"create_failure_listProductionsForBuilding": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingResourceProduction = &mockBuildingResourceProductionRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForBuildingCalled)
				},
			},
			"create_failure_consolidation": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					consolidator := func(action persistence.BuildingAction, _ []persistence.Resource, _ []persistence.BuildingActionCost) (persistence.BuildingAction, error) {
						return action, errDefault
					}
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.PlanetBuilding, _ []persistence.BuildingActionCost) error {
						return nil
					}
					s := newBuildingActionService(pool, repos, consolidator, validator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"create_failure_validation": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					consolidator := func(action persistence.BuildingAction, _ []persistence.Resource, _ []persistence.BuildingActionCost) (persistence.BuildingAction, error) {
						return action, nil
					}
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.PlanetBuilding, _ []persistence.BuildingActionCost) error {
						return errDefault
					}
					s := newBuildingActionService(pool, repos, consolidator, validator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"create_failure_updateResourcesOnPlanet": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: defaultPlanetResource,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"create_failure_registerActionCosts": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingActionCost = &mockBuildingActionCostRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
				},
			},
			"create_failure_registerActionProductions": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingActionResourceProduction = &mockBuildingActionResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
				},
			},
			"create_failure_createsAction": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
				},
			},
			"delete_getsAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultBuildingAction.Id, m.getId)
				},
			},
			"delete_listsCostsForAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.listForActionId)
				},
			},
			"delete_listsResourcesOnPlanet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal([]uuid.UUID{defaultBuildingAction.Planet}, m.listForPlanetIds)
				},
			},
			"delete_updatesResourcesOnPlanet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))
					actual := m.updatedPlanetResources[0]
					assert.Equal(defaultBuildingActionDtoRequest.Planet, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					expectedAmount := defaultPlanetResource.Amount + float64(defaultBuildingActionCost.Amount)
					assert.Equal(expectedAmount, actual.Amount)
					assert.Equal(defaultPlanetResource.Version, actual.Version)

				},
			},
			"delete_deletesResourceProductionsForAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteForActionId)
				},
			},
			"delete_deletesCostsForAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteForActionId)
				},
			},
			"delete_deletesAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteId)
				},
			},
			"delete_failure_getAction": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
				},
			},
			"delete_failure_listCostsForAction": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingActionCost = &mockBuildingActionCostRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForActionCalled)
				},
			},
			"delete_failure_listResourcesOnPlanet": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"delete_failure_updateResourcesOnPlanet": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: defaultPlanetResource,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"delete_failure_deleteResourceProductionsForAction": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingActionResourceProduction = &mockBuildingActionResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
				},
			},
			"delete_failure_deleteCostsForAction": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingActionCost = &mockBuildingActionCostRepository{
						actionCost: defaultBuildingActionCost,
						errs: []error{
							nil,
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForActionCalled)
				},
			},
			"delete_failure_deleteAction": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidBuildingActionRepositoryMock()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{
							nil,
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewBuildingActionService(pool, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteId)
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
		PlanetResource: &mockPlanetResourceRepository{
			planetResource: defaultPlanetResource,
		},
		Resource: &mockResourceRepository{
			resources: defaultResources,
		},
		PlanetBuilding: &mockPlanetBuildingRepository{
			planetBuilding: defaultPlanetBuilding,
		},
		BuildingCost: &mockBuildingCostRepository{
			buildingCost: defaultBuildingCost,
		},
		BuildingResourceProduction: &mockBuildingResourceProductionRepository{
			buildingResourceProduction: defaultBuildingResourceProduction,
		},
		BuildingAction: &mockBuildingActionRepository{
			action: defaultBuildingAction,
		},
		BuildingActionCost: &mockBuildingActionCostRepository{
			actionCost: defaultBuildingActionCost,
		},
		BuildingActionResourceProduction: &mockBuildingActionResourceProductionRepository{
			actionResourceProduction: defaultBuildingActionResourceProduction,
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

func assertBuildingActionResourceProductionRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockBuildingActionResourceProductionRepository {
	m, ok := repos.BuildingActionResourceProduction.(*mockBuildingActionResourceProductionRepository)
	if !ok {
		assert.Fail("Provided building action resource production repository is not a mock")
	}
	return m
}
