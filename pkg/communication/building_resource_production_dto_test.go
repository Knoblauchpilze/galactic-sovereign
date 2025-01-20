package communication

import (
	"encoding/json"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

var defaultBuildingResourceProduction = persistence.BuildingResourceProduction{
	Building: defaultBuildingId,
	Resource: defaultResourceId,
	Base:     54,
	Progress: 1.3,
}
var defaultBuildingResourceProductionDtoResponse = BuildingResourceProductionDtoResponse{
	Building: defaultBuildingId,
	Resource: defaultResourceId,
	Base:     54,
	Progress: 1.3,
}

func TestUnit_ToBuildingResourceProductionDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToBuildingResourceProductionDtoResponse(defaultBuildingResourceProduction)

	assert.Equal(defaultBuildingId, actual.Building)
	assert.Equal(defaultResourceId, actual.Resource)
	assert.Equal(54, actual.Base)
	assert.Equal(1.3, actual.Progress)
}

func TestUnit_BuildingResourceProductionDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultBuildingResourceProductionDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
		"base": 54,
		"progress": 1.3
	}`
	assert.JSONEq(expectedJson, string(out))
}
