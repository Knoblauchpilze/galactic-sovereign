package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultBuildingId = uuid.MustParse("461ba465-86e6-4234-94b8-fc8fab03fa74")
var defaultBuilding = persistence.Building{
	Id:   defaultBuildingId,
	Name: "my-building",

	CreatedAt: someTime,
}
var defaultBuildingDtoResponse = BuildingDtoResponse{
	Id:        defaultBuildingId,
	Name:      "my-building",
	CreatedAt: someTime,
}

func TestToBuildingDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToBuildingDtoResponse(defaultBuilding)

	assert.Equal(defaultBuildingId, actual.Id)
	assert.Equal("my-building", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestBuildingDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultBuildingDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"name": "my-building",
		"createdAt": "2024-05-05T20:50:18.651387237Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
