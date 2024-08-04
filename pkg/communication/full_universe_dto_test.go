package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func TestToFullUniverseDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToFullUniverseDtoResponse(defaultUniverse, []persistence.Resource{defaultResource})

	assert.Equal(defaultUniverseId, actual.Id)
	assert.Equal("my-universe", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)

	assert.Equal(1, len(actual.Resources))
	assert.Equal(defaultResourceDtoResponse, actual.Resources[0])
}

func TestFullUniverseDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := FullUniverseDtoResponse{
		UniverseDtoResponse: defaultUniverseDtoResponse,
		Resources: []ResourceDtoResponse{
			defaultResourceDtoResponse,
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"06fedf46-80ed-4188-b94c-ed0a494ec7bd","name":"my-universe","createdAt":"2024-05-05T20:50:18.651387237Z","resources":[{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","name":"my-resource","createdAt":"2024-05-05T20:50:18.651387237Z"}]}`, string(out))
}

func TestFullUniverseDtoResponse_WhenResourcesAreEmpty_MarshalsToEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := FullUniverseDtoResponse{
		UniverseDtoResponse: defaultUniverseDtoResponse,
		Resources:           nil,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"06fedf46-80ed-4188-b94c-ed0a494ec7bd","name":"my-universe","createdAt":"2024-05-05T20:50:18.651387237Z","resources":[]}`, string(out))
}
