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
	Planet:       defaultPlanetId,
	Building:     defaultBuildingId,
	CurrentLevel: defaultPlanetBuilding.Level,
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
					assert.Equal(defaultBuildingAction.CurrentLevel, actual.CurrentLevel)
					assert.Equal(defaultBuildingAction.DesiredLevel, actual.DesiredLevel)
					assert.True(beforeTestSuite.Before(actual.CreatedAt))
					// TODO: Improve this.
					// assert.False(actual.CompletedAt.IsZero())
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
			"create_validationFails": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					validator := func(_ persistence.BuildingAction, _ []persistence.PlanetResource, _ []persistence.BuildingCost, _ []persistence.PlanetBuilding) error {
						return errDefault
					}
					s := newBuildingActionService(pool, repos, validator)
					_, err := s.Create(ctx, defaultBuildingActionDtoRequest)
					return err
				},
				expectedError: errDefault,
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
		},
	}

	suite.Run(t, &s)
}

func generateValidBuildingActionRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
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
	}
}

func generateErrorBuildingActionRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		PlanetResource: &mockPlanetResourceRepository{
			err: errDefault,
		},
	}
}
