package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

var defaultPlanetResourceStorage = persistence.PlanetResourceStorage{
	Planet:   defaultPlanetId,
	Resource: defaultResourceId,
	Storage:  20,

	CreatedAt: someTime,
	UpdatedAt: someOtherTime,
}
var defaultPlanetResourceStorageDtoResponse = PlanetResourceStorageDtoResponse{
	Planet:   defaultPlanetId,
	Resource: defaultResourceId,
	Storage:  20,
}

func TestToPlanetResourceStorageDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToPlanetResourceStorageDtoResponse(defaultPlanetResourceStorage)

	assert.Equal(defaultPlanetId, actual.Planet)
	assert.Equal(defaultResourceId, actual.Resource)
	assert.Equal(20, actual.Storage)
}

func TestPlanetResourceStorageDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultPlanetResourceStorageDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
		"storage": 20
	}`
	assert.JSONEq(expectedJson, string(out))
}
