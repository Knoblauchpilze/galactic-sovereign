package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	eassert "github.com/KnoblauchPilze/easy-assert/assert"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/communication"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
		BuildingAction: &mockBuildingActionRepository{
			action: defaultBuildingAction,
		},
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
		BuildingAction: &mockBuildingActionRepository{},
		Planet:         &mockPlanetRepository{},
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

func TestIT_PlayerService_Create_ExpectHomeworldRegisteredForPlayer(t *testing.T) {
	conn := newTestConnection(t)
	defer conn.Close(context.Background())
	repos := repositories.Repositories{
		Planet: repositories.NewPlanetRepository(conn),
		Player: repositories.NewPlayerRepository(conn),
	}
	universe := insertTestUniverse(t, conn)

	service := NewPlayerService(conn, repos)

	var err error
	var playerResponse communication.PlayerDtoResponse
	func() {
		playerRequest := communication.PlayerDtoRequest{
			ApiUser:  uuid.New(),
			Universe: universe.Id,
			Name:     fmt.Sprintf("my-player-%s", uuid.NewString()),
		}

		playerResponse, err = service.Create(context.Background(), playerRequest)
		require.Nil(t, err)
	}()

	assertHomeworldExistsForPlayer(t, conn, playerResponse.Id)
}

func TestIT_PlayerService_Delete_ExpectBuildingActionToBeDeleted(t *testing.T) {
	conn := newTestConnection(t)
	defer conn.Close(context.Background())
	repos := repositories.Repositories{
		Planet:         repositories.NewPlanetRepository(conn),
		Player:         repositories.NewPlayerRepository(conn),
		BuildingAction: repositories.NewBuildingActionRepository(),
	}
	universe := insertTestUniverse(t, conn)

	service := NewPlayerService(conn, repos)

	playerRequest := communication.PlayerDtoRequest{
		ApiUser:  uuid.New(),
		Universe: universe.Id,
		Name:     fmt.Sprintf("my-player-%s", uuid.NewString()),
	}

	var err error
	var playerResponse communication.PlayerDtoResponse
	func() {
		playerResponse, err = service.Create(context.Background(), playerRequest)
		require.Nil(t, err)
	}()

	planet := getHomeworldForPlayer(t, conn, playerResponse.Id)
	action, _ := insertTestBuildingActionForPlanet(t, conn, planet)

	func() {
		err = service.Delete(context.Background(), playerResponse.Id)
		require.Nil(t, err)
	}()

	assertPlayerDoesNotExist(t, conn, playerResponse.Id)
	assertHomeworldDoesNotExistForPlayer(t, conn, playerResponse.Id)
	assertBuildingActionDoesNotExist(t, conn, action.Id)
}

func TestIT_PlayerService_CreationDeletionWorkflow(t *testing.T) {
	conn := newTestConnection(t)
	defer conn.Close(context.Background())
	repos := repositories.Repositories{
		Planet:         repositories.NewPlanetRepository(conn),
		Player:         repositories.NewPlayerRepository(conn),
		BuildingAction: repositories.NewBuildingActionRepository(),
	}
	universe := insertTestUniverse(t, conn)

	service := NewPlayerService(conn, repos)

	playerRequest := communication.PlayerDtoRequest{
		ApiUser:  uuid.New(),
		Universe: universe.Id,
		Name:     fmt.Sprintf("my-player-%s", uuid.NewString()),
	}

	beforeCreation := time.Now()
	// Make sure that there's a bit of time between the creation and this timestamp
	time.Sleep(100 * time.Millisecond)

	var err error
	var playerResponse communication.PlayerDtoResponse
	func() {
		playerResponse, err = service.Create(context.Background(), playerRequest)
		require.Nil(t, err)
	}()

	assertPlayerExists(t, conn, playerResponse.Id)
	expected := communication.PlayerDtoResponse{
		ApiUser:  playerRequest.ApiUser,
		Universe: playerRequest.Universe,
		Name:     playerRequest.Name,
	}
	assert.True(t, eassert.EqualsIgnoringFields(playerResponse, expected, "Id", "CreatedAt"))
	assert.True(t, playerResponse.CreatedAt.After(beforeCreation))

	func() {
		playerFromDb, err := service.Get(context.Background(), playerResponse.Id)
		require.Nil(t, err)

		assert.True(t, eassert.EqualsIgnoringFields(playerFromDb, playerResponse, "CreatedAt"))
		assert.True(t, eassert.AreTimeCloserThan(
			playerFromDb.CreatedAt,
			playerResponse.CreatedAt,
			1*time.Second,
		))
		assert.True(t, playerFromDb.CreatedAt.After(beforeCreation), "actual: %v, expected: %v", playerFromDb.CreatedAt, beforeCreation)
	}()

	func() {
		err = service.Delete(context.Background(), playerResponse.Id)
		require.Nil(t, err)
	}()

	assertPlayerDoesNotExist(t, conn, playerResponse.Id)
	assertHomeworldDoesNotExistForPlayer(t, conn, playerResponse.Id)
}

func assertPlayerExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	sqlQuery := `SELECT id FROM player WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertPlayerDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	sqlQuery := `SELECT COUNT(id) FROM player WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertHomeworldExistsForPlayer(t *testing.T, conn db.Connection, player uuid.UUID) {
	sqlQuery := `SELECT COUNT(id) FROM planet WHERE player = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, player)
	require.Nil(t, err)
	require.Equal(t, 1, value)

	sqlQuery = `SELECT COUNT(planet) FROM homeworld WHERE player = $1`
	value, err = db.QueryOne[int](context.Background(), conn, sqlQuery, player)
	require.Nil(t, err)
	require.Equal(t, 1, value)
}

func assertHomeworldDoesNotExistForPlayer(t *testing.T, conn db.Connection, player uuid.UUID) {
	sqlQuery := `SELECT COUNT(planet) FROM homeworld WHERE player = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, player)
	require.Nil(t, err)
	require.Zero(t, value)
}

func getHomeworldForPlayer(t *testing.T, conn db.Connection, player uuid.UUID) uuid.UUID {
	sqlQuery := `SELECT planet FROM homeworld WHERE player = $1`
	planet, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, player)
	require.Nil(t, err)
	return planet
}
