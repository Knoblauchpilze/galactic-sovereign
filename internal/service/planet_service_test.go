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
			"get": {
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
			"get_repositoryFails": {
				generateRepositoriesMock: generateErrorPlanetRepositoryMock,
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
			"list_repositoryFails": {
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
			"listForPlayer_repositoryFails": {
				generateRepositoriesMock: generateErrorPlanetRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.ListForPlayer(ctx, defaultPlayerId)
					return err
				},
				expectedError: errDefault,
			},
			"delete": {
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
			"delete_repositoryFails": {
				generateRepositoriesMock: generateErrorPlanetRepositoryMock,
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

				expectedContent: communication.PlanetDtoResponse{
					Id:     defaultPlanet.Id,
					Player: defaultPlanet.Player,
					Name:   defaultPlanet.Name,

					CreatedAt: defaultPlanet.CreatedAt,
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
		Planet: &mockPlanetRepository{
			planet: defaultPlanet,
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
