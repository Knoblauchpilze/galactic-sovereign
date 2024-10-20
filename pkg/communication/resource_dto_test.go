package communication

import (
	"encoding/json"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultResourceId = uuid.MustParse("97ddca58-8eee-41af-8bda-f37a3080f618")
var defaultResource = persistence.Resource{
	Id:   defaultResourceId,
	Name: "my-resource",

	CreatedAt: someTime,
}
var defaultResourceDtoResponse = ResourceDtoResponse{
	Id:        defaultResourceId,
	Name:      "my-resource",
	CreatedAt: someTime,
}

func TestToResourceDtoResponse(t *testing.T) {
	assert := assert.New(t)

	actual := ToResourceDtoResponse(defaultResource)

	assert.Equal(defaultResourceId, actual.Id)
	assert.Equal("my-resource", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestResourceDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	out, err := json.Marshal(defaultResourceDtoResponse)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "97ddca58-8eee-41af-8bda-f37a3080f618",
		"name": "my-resource",
		"createdAt": "2024-05-05T20:50:18.651387237Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
