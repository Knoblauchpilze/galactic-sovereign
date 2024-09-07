package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

var defaultBuildingCost = persistence.BuildingCost{
	Building: defaultBuildingId,
	Resource: defaultResourceId,
	Cost:     54,
	Progress: 1.3,
}
var defaultBuildingCostDtoResponse = BuildingCostDtoResponse{
	Building: defaultBuildingId,
	Resource: defaultResourceId,
	Cost:     54,
	Progress: 1.3,
}

func TestToBuildingCostDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToBuildingCostDtoResponse(defaultBuildingCost)

	assert.Equal(defaultBuildingId, actual.Building)
	assert.Equal(defaultResourceId, actual.Resource)
	assert.Equal(54, actual.Cost)
	assert.Equal(1.3, actual.Progress)
}

func TestBuildingCostDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultBuildingCostDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
		"cost": 54,
		"progress": 1.3
	}`
	assert.JSONEq(expectedJson, string(out))
}
