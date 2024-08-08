package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

var someOtherTime = time.Date(2024, 07, 28, 10, 30, 02, 651387236, time.UTC)
var defaultPlanetResource = persistence.PlanetResource{
	Planet:   defaultPlanetId,
	Resource: defaultResourceId,
	Amount:   1234.567,

	CreatedAt: someTime,
	UpdatedAt: someOtherTime,
}
var defaultPlanetResourceDtoResponse = PlanetResourceDtoResponse{
	Planet:    defaultPlanetId,
	Resource:  defaultResourceId,
	Amount:    1234.567,
	CreatedAt: someTime,
	UpdatedAt: someOtherTime,
}

func TestToPlanetResourceDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToPlanetResourceDtoResponse(defaultPlanetResource)

	assert.Equal(defaultPlanetId, actual.Planet)
	assert.Equal(defaultResourceId, actual.Resource)
	assert.Equal(1234.567, actual.Amount)
	assert.Equal(someTime, actual.CreatedAt)
	assert.Equal(someOtherTime, actual.UpdatedAt)
}

func TestPlanetResourceDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultPlanetResourceDtoResponse)

	assert.Nil(err)
	assert.Equal(`{"planet":"65801b9b-84e6-411d-805f-2eb89587c5a7","resource":"97ddca58-8eee-41af-8bda-f37a3080f618","amount":1234.567,"createdAt":"2024-05-05T20:50:18.651387237Z","updatedAt":"2024-07-28T10:30:02.651387236Z"}`, string(out))
}
