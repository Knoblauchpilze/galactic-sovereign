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

func TestIT_Server_PlayerBuildingActionLifecycle(t *testing.T) {
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
	assert.Equal(t, []uuid.UUID{player.Homeworld}, player.Planets)

	// Get the homeworld and assert basic properties
	homeworld := doGet[dtos.PlanetDtoResponse](
		t, urlFor(conf, "planets", player.Homeworld.String()),
	)
	assert.True(t, homeworld.Homeworld)
	assert.Equal(t, "homeworld", homeworld.Name)
	assert.Equal(t, player.Id, homeworld.Player)
	assert.Len(t, homeworld.Resources, 2)
	assert.Len(t, homeworld.Buildings, 4)
	assert.Nil(t, homeworld.BuildingAction)

	// Create a building action on the planet
	actionReq := dtos.BuildingActionDtoRequest{
		Building: metalMineId,
	}
	action := doPost[dtos.BuildingActionDtoResponse](
		t, urlFor(conf, "planets", homeworld.Id.String(), "actions"), actionReq,
	)
	assert.Equal(t, metalMineId, action.Building)
	assert.Len(t, action.Costs, 2)
	assert.Len(t, action.Productions, 1)
	assert.Empty(t, action.Storages)

	homeworld = doGet[dtos.PlanetDtoResponse](
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
}

func TestIT_Server_PlayerDeletionRemovesPlanets(t *testing.T) {
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

	// Create a second planet
	planet := doPost[dtos.PlanetDtoResponse](
		t, urlFor(conf, "players", player.Id.String(), "planets"), nil,
	)

	homeworld := doGet[dtos.PlanetDtoResponse](
		t, urlFor(conf, "planets", player.Homeworld.String()),
	)
	secondPlanet := doGet[dtos.PlanetDtoResponse](
		t, urlFor(conf, "planets", planet.Id.String()),
	)

	assert.Equal(t, player.Id, homeworld.Player)
	assert.Equal(t, player.Id, secondPlanet.Player)
	assert.True(t, homeworld.Homeworld)
	assert.False(t, secondPlanet.Homeworld)

	// Delete the player
	doDelete(t, urlFor(conf, "players", player.Id.String()))

	assertGetStatus(t, urlFor(conf, "planets", homeworld.Id.String()), http.StatusNotFound)
	assertGetStatus(t, urlFor(conf, "planets", secondPlanet.Id.String()), http.StatusNotFound)
	assertGetStatus(t, urlFor(conf, "players", player.Id.String()), http.StatusNotFound)
}

func assertGetStatus(t *testing.T, url string, expectedStatus int) {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err, "GET %s: %v", url, err)
	require.Equal(t, expectedStatus, resp.StatusCode, "GET %s returned %d", url, resp.StatusCode)
}
