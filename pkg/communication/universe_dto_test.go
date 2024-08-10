package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUniverseId = uuid.MustParse("06fedf46-80ed-4188-b94c-ed0a494ec7bd")
var defaultUniverse = persistence.Universe{
	Id:   defaultUniverseId,
	Name: "my-universe",

	CreatedAt: someTime,

	Version: 9,
}
var defaultUniverseDtoResponse = UniverseDtoResponse{
	Id:        defaultUniverseId,
	Name:      "my-universe",
	CreatedAt: someTime,
}

func TestUniverseDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := UniverseDtoRequest{
		Name: "my-universe",
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"name": "my-universe"
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFromUniverseDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := UniverseDtoRequest{
		Name: "my-universe",
	}

	actual := FromUniverseDtoRequest(dto)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal("my-universe", actual.Name)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestToUniverseDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToUniverseDtoResponse(defaultUniverse)

	assert.Equal(defaultUniverseId, actual.Id)
	assert.Equal("my-universe", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestUniverseDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultUniverseDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "06fedf46-80ed-4188-b94c-ed0a494ec7bd",
		"name": "my-universe",
		"createdAt": "2024-05-05T20:50:18.651387237Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
