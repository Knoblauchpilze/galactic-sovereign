package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUniverseDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := UniverseDtoRequest{
		Name: "my-universe",
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"name":"my-universe"}`, string(out))
}

func TestFromUniverseDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := UniverseDtoRequest{
		Name: "my-universe",
	}

	actual := FromUniverseDtoRequest(dto)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal("my-universe", actual.Name)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestToUniverseDtoResponse(t *testing.T) {
	assert := assert.New(t)

	entity := persistence.Universe{
		Id:   defaultUuid,
		Name: "my-universe",

		CreatedAt: someTime,
	}

	actual := ToUniverseDtoResponse(entity)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal("my-universe", actual.Name)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestUniverseDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := UniverseDtoResponse{
		Id:        defaultUuid,
		Name:      "my-universe",
		CreatedAt: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	assert.Equal(`{"id":"08ce96a3-3430-48a8-a3b2-b1c987a207ca","name":"my-universe","createdAt":"2024-05-05T20:50:18.651387237Z"}`, string(out))
}
