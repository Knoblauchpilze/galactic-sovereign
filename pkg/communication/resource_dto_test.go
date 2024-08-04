package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

var defaultResource = persistence.Resource{
	Id:   defaultUuid,
	Name: "my-resource",

	CreatedAt: someTime,
}
var defaultResourceDtoResponse = ResourceDtoResponse{
	Id:        defaultUuid,
	Name:      "my-resource",
	CreatedAt: someTime,
}

func TestToResourceDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToResourceDtoResponse(defaultResource)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal("my-resource", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestResourceDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultResourceDtoResponse)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","name":"my-resource","createdAt":"2024-05-05T20:50:18.651387237Z"}`, string(out))
}
