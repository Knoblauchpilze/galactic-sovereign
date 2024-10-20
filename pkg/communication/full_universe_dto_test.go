package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultFullBuildingDtoResponse = FullBuildingDtoResponse{
	BuildingDtoResponse: defaultBuildingDtoResponse,
	Costs: []BuildingCostDtoResponse{
		defaultBuildingCostDtoResponse,
	},
	Productions: []BuildingResourceProductionDtoResponse{
		defaultBuildingResourceProductionDtoResponse,
	},
}

func TestToFullUniverseDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToFullUniverseDtoResponse(defaultUniverse,
		[]persistence.Resource{defaultResource},
		[]persistence.Building{defaultBuilding},
		map[uuid.UUID][]persistence.BuildingCost{defaultBuilding.Id: {defaultBuildingCost}},
		map[uuid.UUID][]persistence.BuildingResourceProduction{defaultBuilding.Id: {defaultBuildingResourceProduction}})

	assert.Equal(defaultUniverseId, actual.Id)
	assert.Equal("my-universe", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)

	assert.Equal(1, len(actual.Resources))
	assert.Equal(defaultResourceDtoResponse, actual.Resources[0])

	assert.Equal(1, len(actual.Buildings))
	assert.Equal(defaultFullBuildingDtoResponse, actual.Buildings[0])
}

func TestFullUniverseDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := FullUniverseDtoResponse{
		UniverseDtoResponse: defaultUniverseDtoResponse,
		Resources: []ResourceDtoResponse{
			defaultResourceDtoResponse,
		},
		Buildings: []FullBuildingDtoResponse{
			defaultFullBuildingDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "06fedf46-80ed-4188-b94c-ed0a494ec7bd",
		"name": "my-universe",
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"resources": [
			{
				"id": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"name": "my-resource",
				"createdAt": "2024-05-05T20:50:18.651387237Z"
			}
		],
		"buildings": [
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
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFullUniverseDtoResponse_WhenResourcesAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullUniverseDtoResponse{
		UniverseDtoResponse: defaultUniverseDtoResponse,
		Resources:           nil,
		Buildings: []FullBuildingDtoResponse{
			defaultFullBuildingDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "06fedf46-80ed-4188-b94c-ed0a494ec7bd",
		"name": "my-universe",
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"resources": [],
		"buildings": [
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
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFullUniverseDtoResponse_WhenBuildingsAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullUniverseDtoResponse{
		UniverseDtoResponse: defaultUniverseDtoResponse,
		Resources: []ResourceDtoResponse{
			defaultResourceDtoResponse,
		},
		Buildings: nil,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "06fedf46-80ed-4188-b94c-ed0a494ec7bd",
		"name": "my-universe",
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"resources": [
			{
				"id": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"name": "my-resource",
				"createdAt": "2024-05-05T20:50:18.651387237Z"
			}
		],
		"buildings": []
	}`
	assert.JSONEq(expectedJson, string(out))
}
