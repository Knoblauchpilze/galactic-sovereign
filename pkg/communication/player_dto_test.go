package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUniverse = uuid.MustParse("06fedf46-80ed-4188-b94c-ed0a494ec7bd")

func TestPlayerDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	p := PlayerDtoRequest{
		ApiUser:  defaultUser,
		Universe: defaultUniverse,
		Name:     "my-player",
	}

	out, err := json.Marshal(p)

	assert.Nil(err)
	assert.Equal(`{"api_user":"c74a22da-8a05-43a9-a8b9-717e422b0af4","universe":"06fedf46-80ed-4188-b94c-ed0a494ec7bd","name":"my-player"}`, string(out))
}

func TestFromPlayerDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	p := PlayerDtoRequest{
		ApiUser:  defaultUser,
		Universe: defaultUniverse,
		Name:     "my-player",
	}

	actual := FromPlayerDtoRequest(p)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal(defaultUser, actual.ApiUser)
	assert.Equal(defaultUniverse, actual.Universe)
	assert.Equal("my-player", actual.Name)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestToPlayerDtoResponse(t *testing.T) {
	assert := assert.New(t)

	p := persistence.Player{
		Id:       defaultUuid,
		ApiUser:  defaultUser,
		Universe: defaultUniverse,
		Name:     "my-player",

		CreatedAt: someTime,
	}

	actual := ToPlayerDtoResponse(p)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal(defaultUser, actual.ApiUser)
	assert.Equal(defaultUniverse, actual.Universe)
	assert.Equal("my-player", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestPlayerDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	p := PlayerDtoResponse{
		Id:       defaultUuid,
		ApiUser:  defaultUser,
		Universe: defaultUniverse,
		Name:     "my-player",

		CreatedAt: someTime,
	}

	out, err := json.Marshal(p)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","api_user":"c74a22da-8a05-43a9-a8b9-717e422b0af4","universe":"06fedf46-80ed-4188-b94c-ed0a494ec7bd","name":"my-player","createdAt":"2024-05-05T20:50:18.651387237Z"}`, string(out))
}
