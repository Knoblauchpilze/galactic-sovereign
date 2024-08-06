package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func TestToFullPlanetDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToFullPlanetDtoResponse(defaultPlanet, []persistence.PlanetResource{defaultPlanetResource}, []persistence.PlanetBuilding{defaultPlanetBuilding})

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal(defaultPlayer, actual.Player)
	assert.Equal("my-player", actual.Name)
	assert.True(actual.Homeworld)
	assert.Equal(someTime, actual.CreatedAt)

	assert.Equal(1, len(actual.Resources))
	assert.Equal(defaultPlanetResourceDtoResponse, actual.Resources[0])

	assert.Equal(1, len(actual.Buildings))
	assert.Equal(defaultPlanetBuildingDtoResponse, actual.Buildings[0])
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
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","player":"efc01287-830f-4b95-8b26-3deff7135f2d","name":"my-planet","homeworld":true,"createdAt":"2024-05-05T20:50:18.651387237Z","resources":[{"planet":"65801b9b-84e6-411d-805f-2eb89587c5a7","resource":"97ddca58-8eee-41af-8bda-f37a3080f618","amount":1234.567,"createdAt":"2024-05-05T20:50:18.651387237Z","updatedAt":"2024-07-28T10:30:02.651387236Z"}],"buildings":[{"planet":"65801b9b-84e6-411d-805f-2eb89587c5a7","building":"04a197b0-63fc-4a4b-9ce4-0f76cb7b545f","level":37,"createdAt":"2024-05-05T20:50:18.651387237Z","updatedAt":"2024-07-28T10:30:02.651387236Z"}]}`, string(out))
}

func TestFullPlanetDtoResponse_WhenResourcesAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullPlanetDtoResponse{
		PlanetDtoResponse: defaultPlanetDtoResponse,
		Resources:         nil,
		Buildings: []PlanetBuildingDtoResponse{
			defaultPlanetBuildingDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","player":"efc01287-830f-4b95-8b26-3deff7135f2d","name":"my-planet","homeworld":true,"createdAt":"2024-05-05T20:50:18.651387237Z","resources":[],"buildings":[{"planet":"65801b9b-84e6-411d-805f-2eb89587c5a7","building":"04a197b0-63fc-4a4b-9ce4-0f76cb7b545f","level":37,"createdAt":"2024-05-05T20:50:18.651387237Z","updatedAt":"2024-07-28T10:30:02.651387236Z"}]}`, string(out))
}

func TestFullPlanetDtoResponse_WhenBuildingsAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullPlanetDtoResponse{
		PlanetDtoResponse: defaultPlanetDtoResponse,
		Resources: []PlanetResourceDtoResponse{
			defaultPlanetResourceDtoResponse,
		},
		Buildings: nil,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","player":"efc01287-830f-4b95-8b26-3deff7135f2d","name":"my-planet","homeworld":true,"createdAt":"2024-05-05T20:50:18.651387237Z","resources":[{"planet":"65801b9b-84e6-411d-805f-2eb89587c5a7","resource":"97ddca58-8eee-41af-8bda-f37a3080f618","amount":1234.567,"createdAt":"2024-05-05T20:50:18.651387237Z","updatedAt":"2024-07-28T10:30:02.651387236Z"}],"buildings":[]}`, string(out))
}
