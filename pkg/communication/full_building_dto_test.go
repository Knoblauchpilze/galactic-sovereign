package communication

import (
	"encoding/json"
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func TestUnit_ToFullBuildingDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToFullBuildingDtoResponse(defaultBuilding, []persistence.BuildingCost{defaultBuildingCost}, []persistence.BuildingResourceProduction{defaultBuildingResourceProduction})

	assert.Equal(defaultBuildingId, actual.Id)
	assert.Equal("my-building", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)

	assert.Equal(1, len(actual.Costs))
	assert.Equal(defaultBuildingCostDtoResponse, actual.Costs[0])

	assert.Equal(1, len(actual.Productions))
	assert.Equal(defaultBuildingResourceProductionDtoResponse, actual.Productions[0])
}

func TestUnit_FullBuildingDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := FullBuildingDtoResponse{
		BuildingDtoResponse: defaultBuildingDtoResponse,
		Costs: []BuildingCostDtoResponse{
			defaultBuildingCostDtoResponse,
		},
		Productions: []BuildingResourceProductionDtoResponse{
			defaultBuildingResourceProductionDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"name": "my-building",
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"costs": [
			{
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"cost": 54,
				"progress": 1.3
			}
		],
		"productions": [
			{
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"base": 54,
				"progress": 1.3
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestUnit_FullBuildingDtoResponse_WhenCostsAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullBuildingDtoResponse{
		BuildingDtoResponse: defaultBuildingDtoResponse,
		Costs:               nil,
		Productions: []BuildingResourceProductionDtoResponse{
			defaultBuildingResourceProductionDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"name": "my-building",
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"costs": [],
		"productions": [
			{
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"base": 54,
				"progress": 1.3
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestUnit_FullBuildingDtoResponse_WhenProductionsAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullBuildingDtoResponse{
		BuildingDtoResponse: defaultBuildingDtoResponse,
		Costs: []BuildingCostDtoResponse{
			defaultBuildingCostDtoResponse,
		},
		Productions: nil,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "461ba465-86e6-4234-94b8-fc8fab03fa74",
		"name": "my-building",
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"costs": [
			{
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"cost": 54,
				"progress": 1.3
			}
		],
		"productions": []
	}`
	assert.JSONEq(expectedJson, string(out))
}
