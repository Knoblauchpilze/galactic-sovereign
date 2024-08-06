package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var buildingUuid = uuid.MustParse("04a197b0-63fc-4a4b-9ce4-0f76cb7b545f")
var defaultPlanetBuilding = persistence.PlanetBuilding{
	Planet:   planetUuid,
	Building: buildingUuid,
	Level:    37,

	CreatedAt: someTime,
	UpdatedAt: someOtherTime,
}
var defaultPlanetBuildingDtoResponse = PlanetBuildingDtoResponse{
	Planet:    planetUuid,
	Building:  buildingUuid,
	Level:     37,
	CreatedAt: someTime,
	UpdatedAt: someOtherTime,
}

func TestToPlanetBuildingDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToPlanetBuildingDtoResponse(defaultPlanetBuilding)

	assert.Equal(planetUuid, actual.Planet)
	assert.Equal(buildingUuid, actual.Building)
	assert.Equal(37, actual.Level)
	assert.Equal(someTime, actual.CreatedAt)
	assert.Equal(someOtherTime, actual.UpdatedAt)
}

func TestPlanetBuildingDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultPlanetBuildingDtoResponse)

	assert.Nil(err)
	assert.Equal(`{"planet":"65801b9b-84e6-411d-805f-2eb89587c5a7","building":"04a197b0-63fc-4a4b-9ce4-0f76cb7b545f","level":37,"createdAt":"2024-05-05T20:50:18.651387237Z","updatedAt":"2024-07-28T10:30:02.651387236Z"}`, string(out))
}
