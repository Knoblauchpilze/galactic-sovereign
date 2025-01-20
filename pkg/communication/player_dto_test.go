package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUser = uuid.MustParse("c74a22da-8a05-43a9-a8b9-717e422b0af4")
var defaultPlayerId = uuid.MustParse("efc01287-830f-4b95-8b26-3deff7135f2d")
var someTime = time.Date(2024, 05, 05, 20, 50, 18, 651387237, time.UTC)

func TestUnit_PlayerDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := PlayerDtoRequest{
		ApiUser:  defaultUser,
		Universe: defaultUniverseId,
		Name:     "my-player",
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"api_user": "c74a22da-8a05-43a9-a8b9-717e422b0af4",
		"universe": "06fedf46-80ed-4188-b94c-ed0a494ec7bd",
		"name": "my-player"
	}`

	assert.JSONEq(expectedJson, string(out))
}

func TestUnit_FromPlayerDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := PlayerDtoRequest{
		ApiUser:  defaultUser,
		Universe: defaultUniverseId,
		Name:     "my-player",
	}

	actual := FromPlayerDtoRequest(dto)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal(defaultUser, actual.ApiUser)
	assert.Equal(defaultUniverseId, actual.Universe)
	assert.Equal("my-player", actual.Name)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestUnit_ToPlayerDtoResponse(t *testing.T) {
	assert := assert.New(t)

	entity := persistence.Player{
		Id:       defaultPlayerId,
		ApiUser:  defaultUser,
		Universe: defaultUniverseId,
		Name:     "my-player",

		CreatedAt: someTime,
	}

	actual := ToPlayerDtoResponse(entity)

	assert.Equal(defaultPlayerId, actual.Id)
	assert.Equal(defaultUser, actual.ApiUser)
	assert.Equal(defaultUniverseId, actual.Universe)
	assert.Equal("my-player", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestUnit_PlayerDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := PlayerDtoResponse{
		Id:       defaultPlayerId,
		ApiUser:  defaultUser,
		Universe: defaultUniverseId,
		Name:     "my-player",

		CreatedAt: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "efc01287-830f-4b95-8b26-3deff7135f2d",
		"api_user": "c74a22da-8a05-43a9-a8b9-717e422b0af4",
		"universe": "06fedf46-80ed-4188-b94c-ed0a494ec7bd",
		"name": "my-player",
		"createdAt": "2024-05-05T20:50:18.651387237Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
