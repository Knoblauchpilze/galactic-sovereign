package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlayerDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := PlayerDtoRequest{
		ApiUser:  defaultUser,
		Universe: defaultUniverseId,
		Name:     "my-player",
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"api_user":"c74a22da-8a05-43a9-a8b9-717e422b0af4","universe":"06fedf46-80ed-4188-b94c-ed0a494ec7bd","name":"my-player"}`, string(out))
}

func TestFromPlayerDtoRequest(t *testing.T) {
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

func TestToPlayerDtoResponse(t *testing.T) {
	assert := assert.New(t)

	entity := persistence.Player{
		Id:       defaultUuid,
		ApiUser:  defaultUser,
		Universe: defaultUniverseId,
		Name:     "my-player",

		CreatedAt: someTime,
	}

	actual := ToPlayerDtoResponse(entity)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal(defaultUser, actual.ApiUser)
	assert.Equal(defaultUniverseId, actual.Universe)
	assert.Equal("my-player", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestPlayerDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := PlayerDtoResponse{
		Id:       defaultUuid,
		ApiUser:  defaultUser,
		Universe: defaultUniverseId,
		Name:     "my-player",

		CreatedAt: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","api_user":"c74a22da-8a05-43a9-a8b9-717e422b0af4","universe":"06fedf46-80ed-4188-b94c-ed0a494ec7bd","name":"my-player","createdAt":"2024-05-05T20:50:18.651387237Z"}`, string(out))
}
