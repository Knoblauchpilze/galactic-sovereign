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

var defaultPlanetId = uuid.MustParse("5b0efd85-8817-4454-b8f3-7af5d93253a1")
var defaultPlanetName = "my-planet"
var defaultResourceId = uuid.MustParse("3e0aaf91-8b81-403f-b967-8bdba748594d")
var defaultBuildingActionId = uuid.MustParse("38a739bd-79db-453e-ab03-44f9f300c3c8")

var defaultPlanetDtoRequest = communication.PlanetDtoRequest{
	Player: defaultPlayerId,
	Name:   defaultPlanetName,
}
var defaultPlanet = persistence.Planet{
	Id:     defaultPlanetId,
	Player: defaultPlayerId,
	Name:   defaultPlanetName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}
var defaultPlanetResource = persistence.PlanetResource{
	Planet:    defaultPlanetId,
	Resource:  defaultResourceId,
	Amount:    321.0987,
	CreatedAt: testDate,
	UpdatedAt: testDate,
}
var defaultPlanetBuilding = persistence.PlanetBuilding{
	Planet:    defaultPlanetId,
	Building:  defaultBuildingId,
	Level:     38,
	CreatedAt: testDate,
	UpdatedAt: testDate,
}
var defaultBuildingAction = persistence.BuildingAction{
	Id:           defaultBuildingActionId,
	Planet:       defaultPlanetId,
	Building:     defaultBuildingId,
	CurrentLevel: 27,
	DesiredLevel: 33,
	CreatedAt:    testDate,
	CompletedAt:  testDate,
}

func Test_PlanetService(t *testing.T) {
	s := ServiceTestSuite{
		generateRepositoriesMock:      generateValidPlanetRepositoryMock,
		generateErrorRepositoriesMock: generateErrorPlanetRepositoryMock,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Create(ctx, defaultPlanetDtoRequest)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultPlanetDtoRequest.Player, m.createdPlanet.Player)
					assert.Equal(defaultPlanetDtoRequest.Name, m.createdPlanet.Name)
				},
			},
			"create_repositoryFails": {
				generateRepositoriesMock: generateErrorPlanetRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Create(ctx, defaultPlanetDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"get_planet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultPlanetId, m.getId)
				},
			},
			"get_planetRepositoryFails": {
				generateRepositoriesMock: generateErrorPlanetRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"get_planetResource": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultPlanetId, m.listForPlanetId)
				},
			},
			"get_planetResourceRepositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Planet: &mockPlanetRepository{},
						PlanetResource: &mockPlanetResourceRepository{
							err: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"get_planetBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultPlanetId, m.listForPlanetId)
				},
			},
			"get_planetBuildingRepositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						Planet: &mockPlanetRepository{},
						PlanetBuilding: &mockPlanetBuildingRepository{
							err: errDefault,
						},
						PlanetResource: &mockPlanetResourceRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"get_buildingAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultPlanetId, m.listForPlanetId)
				},
			},
			"get_buildingActionRepositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						BuildingAction: &mockBuildingActionRepository{
							err: errDefault,
						},
						Planet:         &mockPlanetRepository{},
						PlanetBuilding: &mockPlanetBuildingRepository{},
						PlanetResource: &mockPlanetResourceRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.List(ctx)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"list_planetRepositoryFails": {
				generateRepositoriesMock: generateErrorPlanetRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
				expectedError: errDefault,
			},
			"listForPlayer": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.ListForPlayer(ctx, defaultPlayerId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlayerCalled)
					assert.Equal(defaultPlayerId, m.listForPlayerId)
				},
			},
			"listForPlayer_planetRepositoryFails": {
				generateRepositoriesMock: generateErrorPlanetRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.ListForPlayer(ctx, defaultPlayerId)
					return err
				},
				expectedError: errDefault,
			},
			"delete_planet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultPlanetId, m.deleteId)
				},
			},
			"delete_planetResource": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_planetBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_buildingAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_planetRepositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						BuildingAction: &mockBuildingActionRepository{},
						Planet: &mockPlanetRepository{
							err: errDefault,
						},
						PlanetBuilding: &mockPlanetBuildingRepository{},
						PlanetResource: &mockPlanetResourceRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
			},
			"delete_planetResourceRepositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						BuildingAction: &mockBuildingActionRepository{},
						Planet:         &mockPlanetRepository{},
						PlanetBuilding: &mockPlanetBuildingRepository{},
						PlanetResource: &mockPlanetResourceRepository{
							err: errDefault,
						},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
			},
			"delete_planetBuildingRepositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						BuildingAction: &mockBuildingActionRepository{},
						Planet:         &mockPlanetRepository{},
						PlanetBuilding: &mockPlanetBuildingRepository{
							err: errDefault,
						},
						PlanetResource: &mockPlanetResourceRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
			},
			"delete_buildingActionRepositoryFails": {
				generateRepositoriesMock: func() repositories.Repositories {
					return repositories.Repositories{
						BuildingAction: &mockBuildingActionRepository{
							err: errDefault,
						},
						Planet:         &mockPlanetRepository{},
						PlanetBuilding: &mockPlanetBuildingRepository{},
						PlanetResource: &mockPlanetResourceRepository{},
					}
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
			},
		},

		returnTestCases: map[string]returnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlanetService(pool, repos)
					out, _ := s.Create(ctx, defaultPlanetDtoRequest)
					return out
				},

				expectedContent: communication.PlanetDtoResponse{
					Id:     defaultPlanet.Id,
					Player: defaultPlanet.Player,
					Name:   defaultPlanet.Name,

					CreatedAt: defaultPlanet.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlanetService(pool, repos)
					out, _ := s.Get(ctx, defaultPlanetId)
					return out
				},

				expectedContent: communication.FullPlanetDtoResponse{
					PlanetDtoResponse: communication.PlanetDtoResponse{
						Id:     defaultPlanet.Id,
						Player: defaultPlanet.Player,
						Name:   defaultPlanet.Name,

						CreatedAt: defaultPlanet.CreatedAt,
					},
					Resources: []communication.PlanetResourceDtoResponse{
						{
							Planet:    defaultPlanet.Id,
							Resource:  defaultResourceId,
							Amount:    321.0987,
							CreatedAt: testDate,
							UpdatedAt: testDate,
						},
					},
					Buildings: []communication.PlanetBuildingDtoResponse{
						{
							Planet:    defaultPlanet.Id,
							Building:  defaultBuildingId,
							Level:     38,
							CreatedAt: testDate,
							UpdatedAt: testDate,
						},
					},
					BuildingActions: []communication.BuildingActionDtoResponse{
						{
							Id:           defaultBuildingAction.Id,
							Planet:       defaultBuildingAction.Planet,
							Building:     defaultBuildingAction.Building,
							CurrentLevel: 27,
							DesiredLevel: 33,
							CreatedAt:    testDate,
							CompletedAt:  testDate,
						},
					},
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlanetService(pool, repos)
					out, _ := s.List(ctx)
					return out
				},

				expectedContent: []communication.PlanetDtoResponse{
					{
						Id:     defaultPlanet.Id,
						Player: defaultPlanet.Player,
						Name:   defaultPlanet.Name,

						CreatedAt: defaultPlanet.CreatedAt,
					},
				},
			},
			"listForPlayer": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlanetService(pool, repos)
					out, _ := s.ListForPlayer(ctx, defaultPlayerId)
					return out
				},

				expectedContent: []communication.PlanetDtoResponse{
					{
						Id:     defaultPlanet.Id,
						Player: defaultPlanet.Player,
						Name:   defaultPlanet.Name,

						CreatedAt: defaultPlanet.CreatedAt,
					},
				},
			},
		},

		transactionTestCases: map[string]transactionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Create(ctx, defaultPlanetDtoRequest)
					return err
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
			},
			"listForPlayer": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.ListForPlayer(ctx, defaultPlayerId)
					return err
				},
			},
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidPlanetRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		BuildingAction: &mockBuildingActionRepository{
			action: defaultBuildingAction,
		},
		Planet: &mockPlanetRepository{
			planet: defaultPlanet,
		},
		PlanetBuilding: &mockPlanetBuildingRepository{
			planetBuilding: defaultPlanetBuilding,
		},
		PlanetResource: &mockPlanetResourceRepository{
			planetResource: defaultPlanetResource,
		},
	}
}

func generateErrorPlanetRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Planet: &mockPlanetRepository{
			err: errDefault,
		},
	}
}

func assertPlanetRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlanetRepository {
	m, ok := repos.Planet.(*mockPlanetRepository)
	if !ok {
		assert.Fail("Provided planet repository is not a mock")
	}
	return m
}

func assertPlanetResourceRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlanetResourceRepository {
	m, ok := repos.PlanetResource.(*mockPlanetResourceRepository)
	if !ok {
		assert.Fail("Provided planet resource repository is not a mock")
	}
	return m
}

func assertPlanetBuildingRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlanetBuildingRepository {
	m, ok := repos.PlanetBuilding.(*mockPlanetBuildingRepository)
	if !ok {
		assert.Fail("Provided planet building repository is not a mock")
	}
	return m
}

func assertBuildingActionRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockBuildingActionRepository {
	m, ok := repos.BuildingAction.(*mockBuildingActionRepository)
	if !ok {
		assert.Fail("Provided building action repository is not a mock")
	}
	return m
}
