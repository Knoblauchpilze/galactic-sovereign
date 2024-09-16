package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

var defaultPlanetResourceProduction = persistence.PlanetResourceProduction{
	Planet:     defaultPlanetId,
	Building:   &defaultBuildingId,
	Resource:   defaultResourceId,
	Production: 12,

	CreatedAt: someTime,
	UpdatedAt: someOtherTime,
}
var defaultPlanetResourceProductionDtoResponse = PlanetResourceProductionDtoResponse{
	Planet:     defaultPlanetId,
	Building:   &defaultBuildingId,
	Resource:   defaultResourceId,
	Production: 12,
}

func TestToPlanetResourceProductionDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToPlanetResourceProductionDtoResponse(defaultPlanetResourceProduction)

	assert.Equal(defaultPlanetId, actual.Planet)
	assert.Equal(defaultResourceId, actual.Resource)
	assert.Equal(&defaultBuildingId, actual.Building)
	assert.Equal(12, actual.Production)
}

func TestPlanetResourceProductionDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultPlanetResourceProductionDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
		"production": 12
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestPlanetResourceProductionDtoResponse_whenBuildingIsNotSet_expectBuildingIsOmitted(t *testing.T) {
	assert := assert.New(t)

	withoutBuilding := PlanetResourceProductionDtoResponse{
		Planet:     defaultPlanetId,
		Resource:   defaultResourceId,
		Production: 12,
	}

	out, err := json.Marshal(withoutBuilding)

	assert.Nil(err)
	expectedJson := `
	{
		"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
		"production": 12
	}`
	assert.JSONEq(expectedJson, string(out))
}
