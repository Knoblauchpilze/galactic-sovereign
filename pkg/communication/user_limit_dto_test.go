package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLimitDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	u := LimitDtoRequest{
		Name:  "my-name",
		Value: "my-value",
	}

	out, err := json.Marshal(u)

	assert.Nil(err)
	assert.Equal(`{"name":"my-name","value":"my-value"}`, string(out))
}

func TestFromLimitDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	l := LimitDtoRequest{
		Name:  "my-name",
		Value: "my-value",
	}

	actual := FromLimitDtoRequest(l)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal("my-name", actual.Name)
	assert.Equal("my-value", actual.Value)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestToLimitDtoResponse(t *testing.T) {
	assert := assert.New(t)

	l := persistence.Limit{
		Id:    defaultUuid,
		Name:  "my-name",
		Value: "my-value",

		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	actual := ToLimitDtoResponse(l)

	assert.Equal("my-name", actual.Name)
	assert.Equal("my-value", actual.Value)
}

func TestUserLimitResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	u := LimitDtoResponse{
		Name:  "my-name",
		Value: "my-value",
	}

	out, err := json.Marshal(u)

	assert.Nil(err)
	assert.Equal(`{"name":"my-name","value":"my-value"}`, string(out))
}
