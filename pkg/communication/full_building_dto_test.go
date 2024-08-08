package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func TestToFullBuildingDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToFullBuildingDtoResponse(defaultBuilding, []persistence.BuildingCost{defaultBuildingCost}, []persistence.Building{defaultBuilding})

	assert.Equal(defaultBuildingId, actual.Id)
	assert.Equal("my-building", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)

	assert.Equal(1, len(actual.Costs))
	assert.Equal(defaultBuildingCostDtoResponse, actual.Costs[0])
}

func TestFullBuildingDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := FullBuildingDtoResponse{
		BuildingDtoResponse: defaultBuildingDtoResponse,
		Costs: []BuildingCostDtoResponse{
			defaultBuildingCostDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"461ba465-86e6-4234-94b8-fc8fab03fa74","name":"my-building","createdAt":"2024-05-05T20:50:18.651387237Z","costs":[{"building":"461ba465-86e6-4234-94b8-fc8fab03fa74","resource":"97ddca58-8eee-41af-8bda-f37a3080f618","cost":54}]}`, string(out))
}

func TestFullBuildingDtoResponse_WhenCostsAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullBuildingDtoResponse{
		BuildingDtoResponse: defaultBuildingDtoResponse,
		Costs:               nil,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"461ba465-86e6-4234-94b8-fc8fab03fa74","name":"my-building","createdAt":"2024-05-05T20:50:18.651387237Z","costs":[]}`, string(out))
}
