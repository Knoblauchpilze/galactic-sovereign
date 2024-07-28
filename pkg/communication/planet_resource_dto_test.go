package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var planetUuid = uuid.MustParse("65801b9b-84e6-411d-805f-2eb89587c5a7")
var resourceUuid = uuid.MustParse("97ddca58-8eee-41af-8bda-f37a3080f618")
var someOtherTime = time.Date(2024, 07, 28, 10, 30, 02, 651387236, time.UTC)

func TestToPlanetResourceDtoResponse(t *testing.T) {
	assert := assert.New(t)

	p := persistence.PlanetResource{
		Planet:   planetUuid,
		Resource: resourceUuid,
		Amount:   1234.567,

		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
	}

	actual := ToPlanetResourceDtoResponse(p)

	assert.Equal(planetUuid, actual.Planet)
	assert.Equal(resourceUuid, actual.Resource)
	assert.Equal(1234.567, actual.Amount)
	assert.Equal(someTime, actual.CreatedAt)
	assert.Equal(someOtherTime, actual.UpdatedAt)
}

func TestPlanetResourceDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	u := PlanetResourceDtoResponse{
		Planet:    planetUuid,
		Resource:  resourceUuid,
		Amount:    1234.567,
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
	}

	out, err := json.Marshal(u)

	assert.Nil(err)
	assert.Equal(`{"planet":"65801b9b-84e6-411d-805f-2eb89587c5a7","resource":"97ddca58-8eee-41af-8bda-f37a3080f618","amount":1234.567,"createdAt":"2024-05-05T20:50:18.651387237Z","updatedAt":"2024-07-28T10:30:02.651387236Z"}`, string(out))
}
