package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func TestToResourceDtoResponse(t *testing.T) {
	assert := assert.New(t)

	entity := persistence.Resource{
		Id:   defaultUuid,
		Name: "my-resource",

		CreatedAt: someTime,
	}

	actual := ToResourceDtoResponse(entity)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal("my-resource", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestResourceDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := ResourceDtoResponse{
		Id:        defaultUuid,
		Name:      "my-resource",
		CreatedAt: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","name":"my-resource","createdAt":"2024-05-05T20:50:18.651387237Z"}`, string(out))
}
