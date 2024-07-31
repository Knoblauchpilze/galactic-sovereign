package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultPlayer = uuid.MustParse("efc01287-830f-4b95-8b26-3deff7135f2d")
var defaultPlanet = persistence.Planet{
	Id:        defaultUuid,
	Player:    defaultPlayer,
	Name:      "my-player",
	Homeworld: true,

	CreatedAt: someTime,
}
var defaultPlanetDtoResponse = PlanetDtoResponse{
	Id:        defaultUuid,
	Player:    defaultPlayer,
	Name:      "my-planet",
	Homeworld: true,
	CreatedAt: someTime,
}

func TestPlanetDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := PlanetDtoRequest{
		Player: defaultPlayer,
		Name:   "my-planet",
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"player":"efc01287-830f-4b95-8b26-3deff7135f2d","name":"my-planet"}`, string(out))
}

func TestFromPlanetDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := PlanetDtoRequest{
		Player: defaultPlayer,
		Name:   "my-planet",
	}

	actual := FromPlanetDtoRequest(dto)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal(defaultPlayer, actual.Player)
	assert.Equal("my-planet", actual.Name)
	assert.False(actual.Homeworld)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestToPlanetDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToPlanetDtoResponse(defaultPlanet)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal(defaultPlayer, actual.Player)
	assert.Equal("my-player", actual.Name)
	assert.True(actual.Homeworld)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestPlanetDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultPlanetDtoResponse)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","player":"efc01287-830f-4b95-8b26-3deff7135f2d","name":"my-planet","homeworld":true,"createdAt":"2024-05-05T20:50:18.651387237Z"}`, string(out))
}
