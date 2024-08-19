package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func TestToFullPlanetDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToFullPlanetDtoResponse(defaultPlanet, []persistence.PlanetResource{defaultPlanetResource}, []persistence.PlanetBuilding{defaultPlanetBuilding}, []persistence.BuildingAction{defaultBuildignAction})

	assert.Equal(defaultPlanetId, actual.Id)
	assert.Equal(defaultPlayerId, actual.Player)
	assert.Equal("my-player", actual.Name)
	assert.True(actual.Homeworld)
	assert.Equal(someTime, actual.CreatedAt)

	assert.Equal(1, len(actual.Resources))
	assert.Equal(defaultPlanetResourceDtoResponse, actual.Resources[0])

	assert.Equal(1, len(actual.Buildings))
	assert.Equal(defaultPlanetBuildingDtoResponse, actual.Buildings[0])

	assert.Equal(1, len(actual.BuildingActions))
	assert.Equal(defaultBuildingActionDtoResponse, actual.BuildingActions[0])
}

func TestFullPlanetDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := FullPlanetDtoResponse{
		PlanetDtoResponse: defaultPlanetDtoResponse,
		Resources: []PlanetResourceDtoResponse{
			defaultPlanetResourceDtoResponse,
		},
		Buildings: []PlanetBuildingDtoResponse{
			defaultPlanetBuildingDtoResponse,
		},
		BuildingActions: []BuildingActionDtoResponse{
			defaultBuildingActionDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"player": "efc01287-830f-4b95-8b26-3deff7135f2d",
		"name": "my-planet",
		"homeworld": true,
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"resources": [
			{
				"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
				"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"amount": 1234.567,
				"createdAt": "2024-05-05T20:50:18.651387237Z",
				"updatedAt": "2024-07-28T10:30:02.651387236Z"
			}
		],
		"buildings": [
			{
				"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"level": 37,
				"createdAt": "2024-05-05T20:50:18.651387237Z",
				"updatedAt": "2024-07-28T10:30:02.651387236Z"
			}
		],
		"buildingActions": [
			{
			"id": "91336067-9884-4280-bb37-411124561e73",
			"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
			"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
			"currentLevel": 37,
			"desiredLevel": 38,
			"createdAt": "2024-05-05T20:50:18.651387237Z",
			"completedAt": "2024-07-28T10:30:02.651387236Z"
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFullPlanetDtoResponse_WhenResourcesAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullPlanetDtoResponse{
		PlanetDtoResponse: defaultPlanetDtoResponse,
		Resources:         nil,
		Buildings: []PlanetBuildingDtoResponse{
			defaultPlanetBuildingDtoResponse,
		},
		BuildingActions: []BuildingActionDtoResponse{
			defaultBuildingActionDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"player": "efc01287-830f-4b95-8b26-3deff7135f2d",
		"name": "my-planet",
		"homeworld": true,
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"resources": [],
		"buildings": [
			{
				"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"level": 37,
				"createdAt": "2024-05-05T20:50:18.651387237Z",
				"updatedAt": "2024-07-28T10:30:02.651387236Z"
			}
		],
		"buildingActions": [
			{
				"id": "91336067-9884-4280-bb37-411124561e73",
				"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"currentLevel": 37,
				"desiredLevel": 38,
				"createdAt": "2024-05-05T20:50:18.651387237Z",
				"completedAt": "2024-07-28T10:30:02.651387236Z"
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFullPlanetDtoResponse_WhenBuildingsAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullPlanetDtoResponse{
		PlanetDtoResponse: defaultPlanetDtoResponse,
		Resources: []PlanetResourceDtoResponse{
			defaultPlanetResourceDtoResponse,
		},
		Buildings: nil,
		BuildingActions: []BuildingActionDtoResponse{
			defaultBuildingActionDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"player": "efc01287-830f-4b95-8b26-3deff7135f2d",
		"name": "my-planet",
		"homeworld": true,
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"resources": [
			{
				"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
				"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"amount": 1234.567,
				"createdAt": "2024-05-05T20:50:18.651387237Z",
				"updatedAt": "2024-07-28T10:30:02.651387236Z"
			}
		],
		"buildings": [],
		"buildingActions": [
			{
				"id": "91336067-9884-4280-bb37-411124561e73",
				"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"currentLevel": 37,
				"desiredLevel": 38,
				"createdAt": "2024-05-05T20:50:18.651387237Z",
				"completedAt": "2024-07-28T10:30:02.651387236Z"
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFullPlanetDtoResponse_WhenBuildingActionsAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullPlanetDtoResponse{
		PlanetDtoResponse: defaultPlanetDtoResponse,
		Resources: []PlanetResourceDtoResponse{
			defaultPlanetResourceDtoResponse,
		},
		Buildings: []PlanetBuildingDtoResponse{
			defaultPlanetBuildingDtoResponse,
		},
		BuildingActions: nil,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "65801b9b-84e6-411d-805f-2eb89587c5a7",
		"player": "efc01287-830f-4b95-8b26-3deff7135f2d",
		"name": "my-planet",
		"homeworld": true,
		"createdAt": "2024-05-05T20:50:18.651387237Z",
		"resources": [
			{
				"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
				"resource": "97ddca58-8eee-41af-8bda-f37a3080f618",
				"amount": 1234.567,
				"createdAt": "2024-05-05T20:50:18.651387237Z",
				"updatedAt": "2024-07-28T10:30:02.651387236Z"
			}
		],
		"buildings": [
			{
				"planet": "65801b9b-84e6-411d-805f-2eb89587c5a7",
				"building": "461ba465-86e6-4234-94b8-fc8fab03fa74",
				"level": 37,
				"createdAt": "2024-05-05T20:50:18.651387237Z",
				"updatedAt": "2024-07-28T10:30:02.651387236Z"
			}
		],
		"buildingActions": []
	}`
	assert.JSONEq(expectedJson, string(out))
}
