package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

var defaultPlanetBuilding = persistence.PlanetBuilding{
	Planet:   defaultPlanetId,
	Building: defaultBuildingId,
	Level:    37,

	CreatedAt: someTime,
	UpdatedAt: someOtherTime,
}
var defaultPlanetBuildingDtoResponse = PlanetBuildingDtoResponse{
	Planet:    defaultPlanetId,
	Building:  defaultBuildingId,
	Level:     37,
	CreatedAt: someTime,
	UpdatedAt: someOtherTime,
}

func TestToPlanetBuildingDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToPlanetBuildingDtoResponse(defaultPlanetBuilding)

	assert.Equal(defaultPlanetId, actual.Planet)
	assert.Equal(defaultBuildingId, actual.Building)
	assert.Equal(37, actual.Level)
	assert.Equal(someTime, actual.CreatedAt)
	assert.Equal(someOtherTime, actual.UpdatedAt)
}

func TestPlanetBuildingDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultPlanetBuildingDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"level": 37,
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"updatedAt": "2024-07-28T10:30:02.651387236Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
