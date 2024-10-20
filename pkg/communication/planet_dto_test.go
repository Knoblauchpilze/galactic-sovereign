package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultPlanetId = uuid.MustParse("65801b9b-84e6-411d-805f-2eb89587c5a7")
var defaultPlanet = persistence.Planet{
	Id:        defaultPlanetId,
	Player:    defaultPlayerId,
	Name:      "my-player",
	Homeworld: true,

	CreatedAt: someTime,
}
var defaultPlanetDtoResponse = PlanetDtoResponse{
	Id:        defaultPlanetId,
	Player:    defaultPlayerId,
	Name:      "my-planet",
	Homeworld: true,
	CreatedAt: someTime,
}

func TestPlanetDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := PlanetDtoRequest{
		Player: defaultPlayerId,
		Name:   "my-planet",
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"player": "efc01287-830f-4b95-8b26-3deff7135f2d",
		"name": "my-planet"
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFromPlanetDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := PlanetDtoRequest{
		Player: defaultPlayerId,
		Name:   "my-planet",
	}

	actual := FromPlanetDtoRequest(dto)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal(defaultPlayerId, actual.Player)
	assert.Equal("my-planet", actual.Name)
	assert.False(actual.Homeworld)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestToPlanetDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToPlanetDtoResponse(defaultPlanet)

	assert.Equal(defaultPlanetId, actual.Id)
	assert.Equal(defaultPlayerId, actual.Player)
	assert.Equal("my-player", actual.Name)
	assert.True(actual.Homeworld)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestPlanetDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultPlanetDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"player": "efc01287-830f-4b95-8b26-3deff7135f2d",
		"name": "my-planet",
		"homeworld": true,
		"createdAt": "2024-05-05T20:50:18.651387237Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
