package communication

import (
	"encoding/json"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

var defaultBuildingResourceStorage = persistence.BuildingResourceStorage{
	Building: defaultBuildingId,
	Resource: defaultResourceId,
	Base:     74,
	Scale:    1.08,
	Progress: 2.97,
}
var defaultBuildingResourceStorageDtoResponse = BuildingResourceStorageDtoResponse{
	Building: defaultBuildingId,
	Resource: defaultResourceId,
	Base:     74,
	Scale:    1.08,
	Progress: 2.97,
}

func TestUnit_ToBuildingResourceStorageDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToBuildingResourceStorageDtoResponse(defaultBuildingResourceStorage)

	assert.Equal(defaultBuildingId, actual.Building)
	assert.Equal(defaultResourceId, actual.Resource)
	assert.Equal(74, actual.Base)
	assert.Equal(1.08, actual.Scale)
	assert.Equal(2.97, actual.Progress)
}

func TestUnit_BuildingResourceStorageDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultBuildingResourceStorageDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
		"base": 74,
		"scale": 1.08,
		"progress": 2.97
	}`
	assert.JSONEq(expectedJson, string(out))
}
