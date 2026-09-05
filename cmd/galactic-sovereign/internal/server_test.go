package internal

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driving/dtos"
	integrationdb "github.com/Knoblauchpilze/galactic-sovereign/pkg/testing/integrationdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_Server(t *testing.T) {
	t.Run("create a player and assert homeworld properties", func(t *testing.T) {
		dbContainer := integrationdb.NewDatabaseSharedContainer(t)
		conn := dbContainer.NewTestConnection(t)
		conf := newTestServerConfig()

		s := CreateGameServer(conf, conn, slog.Default())
		asyncStartServer(t, s)

		// Create a player
		playerReq := dtos.PlayerDtoRequest{
			ApiUser:  uuid.New(),
			Universe: oberonUniverseId,
			Name:     "test-player",
		}
		player := doPost[dtos.PlayerDtoResponse](
			t, urlFor(conf, "players"), playerReq,
		)
		assert.Equal(t, oberonUniverseId, player.Universe)
		assert.Equal(t, "test-player", player.Name)
		assert.Len(t, player.Planets, 1)
		assert.Equal(t, player.Homeworld, player.Planets[0].Id)
		assert.Equal(t, "homeworld", player.Planets[0].Name)

		// Get the homeworld and assert basic properties
		homeworld := doGet[dtos.PlanetDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String()),
		)
		assert.True(t, homeworld.Homeworld)
		assert.Equal(t, "homeworld", homeworld.Name)
		assert.Equal(t, player.Id, homeworld.Player)
		assert.Len(t, homeworld.Resources, 3)
		assert.Len(t, homeworld.Buildings, 7)
		assert.Nil(t, homeworld.BuildingAction)
	})

	t.Run("create a building action and cancel it", func(t *testing.T) {
		dbContainer := integrationdb.NewDatabaseSharedContainer(t)
		conn := dbContainer.NewTestConnection(t)
		conf := newTestServerConfig()

		s := CreateGameServer(conf, conn, slog.Default())
		asyncStartServer(t, s)

		// Create a player
		playerReq := dtos.PlayerDtoRequest{
			ApiUser:  uuid.New(),
			Universe: oberonUniverseId,
			Name:     "test-player",
		}
		player := doPost[dtos.PlayerDtoResponse](
			t, urlFor(conf, "players"), playerReq,
		)
		assert.Equal(t, oberonUniverseId, player.Universe)
		assert.Equal(t, "test-player", player.Name)
		assert.Len(t, player.Planets, 1)
		assert.Equal(t, player.Homeworld, player.Planets[0].Id)
		assert.Equal(t, "homeworld", player.Planets[0].Name)

		// Create a building action on the planet
		actionReq := dtos.BuildingActionDtoRequest{
			Building: metalMineId,
		}
		action := doPost[dtos.BuildingActionDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String(), "actions"), actionReq,
		)
		assert.Equal(t, metalMineId, action.Building)
		assert.Len(t, action.Costs, 2)
		assert.Len(t, action.Productions, 1)
		assert.Empty(t, action.Storages)

		homeworld := doGet[dtos.PlanetDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String()),
		)
		require.NotNil(t, homeworld.BuildingAction)
		assert.Equal(t, action, *homeworld.BuildingAction)

		// Cancel the building action
		doDelete(t, urlFor(conf, "planets", homeworld.Id.String(), "actions"))

		homeworld = doGet[dtos.PlanetDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String()),
		)
		assert.Nil(t, homeworld.BuildingAction)
	})

	t.Run("create a player and a building action and delete the player", func(t *testing.T) {
		dbContainer := integrationdb.NewDatabaseSharedContainer(t)
		conn := dbContainer.NewTestConnection(t)
		conf := newTestServerConfig()

		s := CreateGameServer(conf, conn, slog.Default())
		asyncStartServer(t, s)

		// Create a player
		playerReq := dtos.PlayerDtoRequest{
			ApiUser:  uuid.New(),
			Universe: oberonUniverseId,
			Name:     "test-player-b",
		}
		player := doPost[dtos.PlayerDtoResponse](
			t, urlFor(conf, "players"), playerReq,
		)

		// Create a building action
		actionReq := dtos.BuildingActionDtoRequest{Building: metalMineId}
		action := doPost[dtos.BuildingActionDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String(), "actions"), actionReq,
		)

		homeworld := doGet[dtos.PlanetDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String()),
		)

		assert.Equal(t, player.Id, homeworld.Player)
		assert.True(t, homeworld.Homeworld)
		require.NotNil(t, homeworld.BuildingAction)
		assert.Equal(t, action.Id, homeworld.BuildingAction.Id)
		assert.Equal(t, metalMineId, homeworld.BuildingAction.Building)
		assert.Equal(t, 1, homeworld.BuildingAction.DesiredLevel)

		// Delete the player
		doDelete(t, urlFor(conf, "players", player.Id.String()))

		assertGetStatus(t, urlFor(conf, "planets", homeworld.Id.String()), http.StatusNotFound)
		assertGetStatus(t, urlFor(conf, "players", player.Id.String()), http.StatusNotFound)
	})

	t.Run("create a player and a ship action and delete the player", func(t *testing.T) {
		dbContainer := integrationdb.NewDatabaseSharedContainer(t)
		conn := dbContainer.NewTestConnection(t)
		conf := newTestServerConfig()

		s := CreateGameServer(conf, conn, slog.Default())
		asyncStartServer(t, s)

		// Create a player
		playerReq := dtos.PlayerDtoRequest{
			ApiUser:  uuid.New(),
			Universe: oberonUniverseId,
			Name:     "test-player",
		}
		player := doPost[dtos.PlayerDtoResponse](
			t, urlFor(conf, "players"), playerReq,
		)

		// Fetch the universe and pick a ship from it
		universe := doGet[dtos.UniverseDtoResponse](
			t, urlFor(conf, "universes", oberonUniverseId.String()),
		)
		require.NotEmpty(t, universe.Ships)
		ship := universe.Ships[0]

		// Credit the homeworld with enough resources to afford the ship
		for _, cost := range ship.Costs {
			addPlanetResources(t, conn, player.Homeworld, cost.Resource, cost.Cost)
		}

		// Create a ship action on the planet
		actionReq := dtos.ShipActionDtoRequest{
			Ship:  ship.Id,
			Count: 1,
		}
		action := doPost[dtos.ShipActionDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String(), "ships"), actionReq,
		)
		assert.Equal(t, ship.Id, action.Ship)
		assert.Equal(t, 1, action.Count)

		homeworld := doGet[dtos.PlanetDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String()),
		)
		require.Len(t, homeworld.ShipActions, 1)
		assert.Equal(t, action, homeworld.ShipActions[0])

		// Delete the player
		doDelete(t, urlFor(conf, "players", player.Id.String()))

		assertGetStatus(t, urlFor(conf, "planets", player.Homeworld.String()), http.StatusNotFound)
		assertGetStatus(t, urlFor(conf, "players", player.Id.String()), http.StatusNotFound)
	})

	t.Run("create a player and a ship action and delete the planet", func(t *testing.T) {
		dbContainer := integrationdb.NewDatabaseSharedContainer(t)
		conn := dbContainer.NewTestConnection(t)
		conf := newTestServerConfig()

		s := CreateGameServer(conf, conn, slog.Default())
		asyncStartServer(t, s)

		// Create a player
		playerReq := dtos.PlayerDtoRequest{
			ApiUser:  uuid.New(),
			Universe: oberonUniverseId,
			Name:     "test-player",
		}
		player := doPost[dtos.PlayerDtoResponse](
			t, urlFor(conf, "players"), playerReq,
		)

		// Fetch the universe and pick a ship from it
		universe := doGet[dtos.UniverseDtoResponse](
			t, urlFor(conf, "universes", oberonUniverseId.String()),
		)
		require.NotEmpty(t, universe.Ships)
		ship := universe.Ships[0]

		// Credit the homeworld with enough resources to afford the ship
		for _, cost := range ship.Costs {
			addPlanetResources(t, conn, player.Homeworld, cost.Resource, cost.Cost)
		}

		// Create a ship action on the planet
		actionReq := dtos.ShipActionDtoRequest{
			Ship:  ship.Id,
			Count: 1,
		}
		action := doPost[dtos.ShipActionDtoResponse](
			t, urlFor(conf, "planets", player.Homeworld.String(), "ships"), actionReq,
		)
		assert.Equal(t, ship.Id, action.Ship)
		assert.Equal(t, 1, action.Count)

		// Delete the planet: it should fail
		assertDeleteStatus(t, urlFor(conf, "planets", player.Homeworld.String()), http.StatusConflict)
	})
}

func assertGetStatus(t *testing.T, url string, expectedStatus int) {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err, "GET %s: %v", url, err)
	require.Equal(t, expectedStatus, resp.StatusCode, "GET %s returned %d", url, resp.StatusCode)
}

func assertDeleteStatus(t *testing.T, url string, expectedStatus int) {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err, "Actual err: %v", err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "DELETE %s: %v", url, err)
	require.Equal(t, expectedStatus, resp.StatusCode, "DELETE %s returned %d", url, resp.StatusCode)
}
