package service

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestUnit_BuildingActionService(t *testing.T) {
	beforeTestSuite := time.Now()

	s := ServicePoolTestSuite{
		generateRepositoriesMocks:      generateBuildingActionServiceMocks,
		generateErrorRepositoriesMocks: generateErrorBuildingActionServiceMocks,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"create_listsResourcesOnPlanet": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"create_listsBuildingsOnPlanet": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultBuildingActionResourceProduction.Resource, m.createdBuildingActionResourceProduction.Resource)
					assert.Equal(72, m.createdBuildingActionResourceProduction.Production)
				},
			},
			"create_createsAction": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.Resource = &mockResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.PlanetBuilding = &mockPlanetBuildingRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.BuildingCost = &mockBuildingCostRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.BuildingResourceProduction = &mockBuildingResourceProductionRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					consolidator := func(action persistence.BuildingAction, _ []persistence.Resource, _ []persistence.BuildingActionCost) (persistence.BuildingAction, error) {
						return action, errDefault
					}
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.PlanetBuilding, _ []persistence.BuildingActionCost) error {
						return nil
					}
					s := newBuildingActionService(conn, repos, consolidator, validator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"create_failure_validation": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					consolidator := func(action persistence.BuildingAction, _ []persistence.Resource, _ []persistence.BuildingActionCost) (persistence.BuildingAction, error) {
						return action, nil
					}
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.PlanetBuilding, _ []persistence.BuildingActionCost) error {
						return errDefault
					}
					s := newBuildingActionService(conn, repos, consolidator, validator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"create_failure_updateResourcesOnPlanet": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: defaultPlanetResource,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.BuildingActionCost = &mockBuildingActionCostRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.BuildingActionResourceProduction = &mockBuildingActionResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultBuildingAction.Id, m.getId)
				},
			},
			"delete_listsCostsForAction": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForActionCalled)
					assert.Equal(defaultBuildingAction.Id, m.listForActionId)
				},
			},
			"delete_listsResourcesOnPlanet": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal([]uuid.UUID{defaultBuildingAction.Planet}, m.listForPlanetIds)
				},
			},
			"delete_updatesResourcesOnPlanet": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
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
			"delete_deletesAction": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteId)
				},
			},
			"delete_failure_getAction": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
				},
			},
			"delete_failure_listCostsForAction": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.BuildingActionCost = &mockBuildingActionCostRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForActionCalled)
				},
			},
			"delete_failure_listResourcesOnPlanet": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"delete_failure_updateResourcesOnPlanet": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: defaultPlanetResource,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"delete_failure_deleteAction": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generateBuildingActionServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{
							nil,
							errDefault,
						},
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultBuildingAction.Id, m.deleteId)
				},
			},
			"delete_failure_actionCompletedInThePast": {
				generateConnectionMock: func() db.Connection {
					return &mockConnection{
						timeStamp: defaultBuildingAction.CompletedAt.Add(2 * time.Minute),
					}
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingAction.Id)
				},
				verifyError: func(err error, assert *require.Assertions) {
					assert.True(errors.IsErrorWithCode(err, ActionAlreadyCompleted))
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultBuildingAction.Id, m.getId)
				},
			},
		},

		returnTestCases: map[string]returnTestCase{
			"create": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) interface{} {
					s := NewBuildingActionService(conn, repos)
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
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
			},
			"delete": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewBuildingActionService(conn, repos)
					return s.Delete(ctx, defaultBuildingActionId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateBuildingActionServiceMocks() repositories.Repositories {
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

func generateErrorBuildingActionServiceMocks() repositories.Repositories {
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

func TestIT_BuildingActionService_CreationDeletionWorkflow(t *testing.T) {
	conn := newTestConnection(t)
	repos := repositories.Repositories{
		Resource:                         repositories.NewResourceRepository(),
		PlanetResource:                   repositories.NewPlanetResourceRepository(),
		PlanetBuilding:                   repositories.NewPlanetBuildingRepository(),
		BuildingCost:                     repositories.NewBuildingCostRepository(),
		BuildingResourceProduction:       repositories.NewBuildingResourceProductionRepository(),
		BuildingAction:                   repositories.NewBuildingActionRepository(),
		BuildingActionCost:               repositories.NewBuildingActionCostRepository(),
		BuildingActionResourceProduction: repositories.NewBuildingActionResourceProductionRepository(),
	}
	planet, _, _ := insertTestPlanetForPlayer(t, conn)
	building, _ := insertTestPlanetBuildingForPlanet(t, conn, planet.Id)
	_, resource := insertTestBuildingCost(t, conn, building.Building)
	insertTestPlanetResourceForResource(t, conn, planet.Id, resource.Id)

	var generatedCreatedAt time.Time
	var returnedCompletionTime time.Time
	completionTimeFunc := func(action persistence.BuildingAction, resources []persistence.Resource, costs []persistence.BuildingActionCost) (persistence.BuildingAction, error) {
		generatedCreatedAt = action.CreatedAt
		returnedCompletionTime = time.Now().Add(1 * time.Hour)
		action.CompletedAt = returnedCompletionTime
		return action, nil
	}

	service := newBuildingActionServiceWithCompletionTime(conn, repos, completionTimeFunc)

	actionRequest := communication.BuildingActionDtoRequest{
		Planet:   planet.Id,
		Building: building.Building,
	}

	var err error
	var actionResponse communication.BuildingActionDtoResponse
	func() {
		actionResponse, err = service.Create(context.Background(), actionRequest)
		require.Nil(t, err)
	}()

	assertBuildingActionExists(t, conn, actionResponse.Id)
	expected := communication.BuildingActionDtoResponse{
		Planet:       actionRequest.Planet,
		Building:     actionRequest.Building,
		CurrentLevel: 4,
		DesiredLevel: 5,
		CreatedAt:    generatedCreatedAt,
		CompletedAt:  returnedCompletionTime,
	}
	assert.True(t, eassert.EqualsIgnoringFields(actionResponse, expected, "Id"))

	func() {
		err = service.Delete(context.Background(), actionResponse.Id)
		require.Nil(t, err)
	}()

	assertBuildingActionDoesNotExist(t, conn, actionResponse.Id)
}

func assertBuildingActionExists(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertBuildingActionDoesNotExist(t *testing.T, conn db.Connection, action uuid.UUID) {
	sqlQuery := `SELECT COUNT(*) FROM building_action WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, action)
	require.Nil(t, err)
	require.Zero(t, value)
}
