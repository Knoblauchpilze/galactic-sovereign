package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/communication"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultUserId = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultPlayerId = uuid.MustParse("f16455b7-514c-44b1-847f-ba2cf99c765b")
var defaultPlayerName = "my-player"

var defaultPlayerDtoRequest = communication.PlayerDtoRequest{
	ApiUser:  defaultUserId,
	Universe: defaultUniverseId,
	Name:     defaultPlayerName,
}
var defaultPlayer = persistence.Player{
	Id:       defaultPlayerId,
	ApiUser:  defaultUserId,
	Universe: defaultUniverseId,
	Name:     defaultPlayerName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func TestUnit_PlayerService(t *testing.T) {
	s := ServicePoolTestSuite{
		generateRepositoriesMocks:      generatePlayerServiceMocks,
		generateErrorRepositoriesMocks: generateErrorPlayerServiceMocks,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.Create(ctx, defaultPlayerDtoRequest)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlayerRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultPlayerDtoRequest.ApiUser, m.createdPlayer.ApiUser)
					assert.Equal(defaultPlayerDtoRequest.Universe, m.createdPlayer.Universe)
					assert.Equal(defaultPlayerDtoRequest.Name, m.createdPlayer.Name)
				},
			},
			"create_createPlanet": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.Create(ctx, defaultPlayerDtoRequest)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(defaultPlayer.Id, m.createdPlanet.Player)
					assert.Equal("homeworld", m.createdPlanet.Name)
					assert.True(m.createdPlanet.Homeworld)
				},
			},
			"create_playerRepositoryFails": {
				generateRepositoriesMocks: generateErrorPlayerServiceMocks,
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.Create(ctx, defaultPlayerDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"create_planetRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					return repositories.Repositories{
						Planet: &mockPlanetRepository{
							err: errDefault,
						},
						Player: &mockPlayerRepository{},
					}
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.Create(ctx, defaultPlayerDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"get": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.Get(ctx, defaultPlayerId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlayerRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultPlayerId, m.getId)
				},
			},
			"get_repositoryFails": {
				generateRepositoriesMocks: generateErrorPlayerServiceMocks,
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.Get(ctx, defaultPlayerId)
					return err
				},
				expectedError: errDefault,
			},
			"list": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.List(ctx)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlayerRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"list_repositoryFails": {
				generateRepositoriesMocks: generateErrorPlayerServiceMocks,
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.List(ctx)
					return err
				},
				expectedError: errDefault,
			},
			"listForApiUser": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.ListForApiUser(ctx, defaultUserId)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlayerRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForApiUserCalled)
					assert.Equal(defaultUserId, m.listForApiUserId)
				},
			},
			"listForApiUser_repositoryFails": {
				generateRepositoriesMocks: generateErrorPlayerServiceMocks,
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.ListForApiUser(ctx, defaultUserId)
					return err
				},
				expectedError: errDefault,
			},
			"delete": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					return s.Delete(ctx, defaultPlayerId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlayerRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultPlayerId, m.deleteId)
				},
			},
			"delete_repositoryFails": {
				generateRepositoriesMocks: generateErrorPlayerServiceMocks,
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					return s.Delete(ctx, defaultPlayerId)
				},
				expectedError: errDefault,
			},
		},

		returnTestCases: map[string]returnTestCase{
			"create": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) interface{} {
					s := NewPlayerService(conn, repos)
					out, _ := s.Create(ctx, defaultPlayerDtoRequest)
					return out
				},

				expectedContent: communication.PlayerDtoResponse{
					Id:       defaultPlayer.Id,
					ApiUser:  defaultPlayer.ApiUser,
					Universe: defaultPlayer.Universe,
					Name:     defaultPlayer.Name,

					CreatedAt: defaultPlayer.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) interface{} {
					s := NewPlayerService(conn, repos)
					out, _ := s.Get(ctx, defaultPlayerId)
					return out
				},

				expectedContent: communication.PlayerDtoResponse{
					Id:       defaultPlayer.Id,
					ApiUser:  defaultPlayer.ApiUser,
					Universe: defaultPlayer.Universe,
					Name:     defaultPlayer.Name,

					CreatedAt: defaultPlayer.CreatedAt,
				},
			},
			"list": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) interface{} {
					s := NewPlayerService(conn, repos)
					out, _ := s.List(ctx)
					return out
				},

				expectedContent: []communication.PlayerDtoResponse{
					{
						Id:       defaultPlayer.Id,
						ApiUser:  defaultPlayer.ApiUser,
						Universe: defaultPlayer.Universe,
						Name:     defaultPlayer.Name,

						CreatedAt: defaultPlayer.CreatedAt,
					},
				},
			},
			"listForApiUser": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) interface{} {
					s := NewPlayerService(conn, repos)
					out, _ := s.ListForApiUser(ctx, defaultUserId)
					return out
				},

				expectedContent: []communication.PlayerDtoResponse{
					{
						Id:       defaultPlayer.Id,
						ApiUser:  defaultPlayer.ApiUser,
						Universe: defaultPlayer.Universe,
						Name:     defaultPlayer.Name,

						CreatedAt: defaultPlayer.CreatedAt,
					},
				},
			},
		},

		transactionTestCases: map[string]transactionTestCase{
			"create": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					_, err := s.Create(ctx, defaultPlayerDtoRequest)
					return err
				},
			},
			"delete": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlayerService(conn, repos)
					return s.Delete(ctx, defaultPlayerId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generatePlayerServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		Planet: &mockPlanetRepository{
			planet: defaultPlanet,
		},
		Player: &mockPlayerRepository{
			player: defaultPlayer,
		},
	}
}

func generateErrorPlayerServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		Player: &mockPlayerRepository{
			err: errDefault,
		},
	}
}

func assertPlayerRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlayerRepository {
	m, ok := repos.Player.(*mockPlayerRepository)
	if !ok {
		assert.Fail("Provided player repository is not a mock")
	}
	return m
}
