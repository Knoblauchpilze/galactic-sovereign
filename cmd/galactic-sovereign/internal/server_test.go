package internal

import (
	"log/slog"
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

	s := CreateGameServer(testServerConfig, conn, slog.Default())
	done := asyncStartServer(t, s)

	// Create a player
	playerReq := dtos.PlayerDtoRequest{
		ApiUser:  uuid.New(),
		Universe: oberonUniverseId,
		Name:     "test-player",
	}
	player := doPost[dtos.PlayerDtoResponse](
		t, urlFor("players"), playerReq,
	)
	assert.Equal(t, oberonUniverseId, player.Universe)
	assert.Equal(t, "test-player", player.Name)
	assert.Equal(t, []uuid.UUID{player.Homeworld}, player.Planets)

	// Get the homeworld and assert basic properties
	homeworld := doGet[dtos.PlanetDtoResponse](
		t, urlFor("planets", player.Homeworld.String()),
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
		t, urlFor("planets", homeworld.Id.String(), "actions"), actionReq,
	)
	assert.Equal(t, metalMineId, action.Building)
	assert.Len(t, action.Costs, 2)
	assert.Len(t, action.Productions, 1)
	assert.Empty(t, action.Storages)
	assert.Equal(t, homeworld.Id, action.Planet)

	homeworld = doGet[dtos.PlanetDtoResponse](
		t, urlFor("planets", player.Homeworld.String()),
	)
	require.NotNil(t, homeworld.BuildingAction)
	assert.Equal(t, action.Id, *homeworld.BuildingAction)

	// Cancel the building action
	doDelete(t, urlFor("actions", action.Id.String()))

	homeworld = doGet[dtos.PlanetDtoResponse](
		t, urlFor("planets", player.Homeworld.String()),
	)
	assert.Nil(t, homeworld.BuildingAction)

	err := s.Stop()
	require.NoError(t, err, "Actual err: %v", err)
	<-done
}
