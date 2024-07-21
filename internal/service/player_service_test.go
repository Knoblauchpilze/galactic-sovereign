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

func Test_PlayerService(t *testing.T) {
	s := ServiceTestSuite{
		generateRepositoriesMock:      generateValidPlayerRepositoryMock,
		generateErrorRepositoriesMock: generateErrorPlayerRepositoryMock,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
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
			"create_repositoryFails": {
				generateRepositoriesMock: generateErrorPlayerRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
					_, err := s.Create(ctx, defaultPlayerDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
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
				generateRepositoriesMock: generateErrorPlayerRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
					_, err := s.Get(ctx, defaultPlayerId)
					return err
				},
				expectedError: errDefault,
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
					_, err := s.List(ctx)
					return err
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlayerRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"list_repositoryFails": {
				generateRepositoriesMock: generateErrorPlayerRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
				expectedError: errDefault,
			},
			"listForApiUser": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
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
				generateRepositoriesMock: generateErrorPlayerRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
					_, err := s.ListForApiUser(ctx, defaultUserId)
					return err
				},
				expectedError: errDefault,
			},
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
					return s.Delete(ctx, defaultPlayerId)
				},

				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlayerRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultPlayerId, m.deleteId)
				},
			},
			"delete_repositoryFails": {
				generateRepositoriesMock: generateErrorPlayerRepositoryMock,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
					return s.Delete(ctx, defaultPlayerId)
				},
				expectedError: errDefault,
			},
		},

		returnTestCases: map[string]returnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlayerService(pool, repos)
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
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlayerService(pool, repos)
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
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlayerService(pool, repos)
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
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlayerService(pool, repos)
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
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlayerService(pool, repos)
					return s.Delete(ctx, defaultPlayerId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidPlayerRepositoryMock() repositories.Repositories {
	return repositories.Repositories{
		Player: &mockPlayerRepository{
			player: defaultPlayer,
		},
	}
}

func generateErrorPlayerRepositoryMock() repositories.Repositories {
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
