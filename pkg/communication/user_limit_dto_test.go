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

func TestLimitDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	u := LimitDtoResponse{
		Name:  "my-name",
		Value: "my-value",
	}

	out, err := json.Marshal(u)

	assert.Nil(err)
	assert.Equal(`{"name":"my-name","value":"my-value"}`, string(out))
}

func TestUserLimitDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	u := UserLimitDtoRequest{
		Name: "my-name",
		User: uuid.MustParse("7aaa145a-5ad8-4d63-be87-177d6abcf1b5"),
		Limits: []LimitDtoRequest{
			{
				Name:  "limit",
				Value: "test",
			},
		},
	}

	out, err := json.Marshal(u)

	assert.Nil(err)
	assert.Equal(`{"name":"my-name","user":"7aaa145a-5ad8-4d63-be87-177d6abcf1b5","limits":[{"name":"limit","value":"test"}]}`, string(out))
}

func TestFromUserLimitDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	u := UserLimitDtoRequest{
		Name: "my-name",
		User: uuid.MustParse("7aaa145a-5ad8-4d63-be87-177d6abcf1b5"),
		Limits: []LimitDtoRequest{
			{
				Name:  "limit",
				Value: "test",
			},
		},
	}

	actual := FromUserLimitDtoRequest(u)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal("my-name", actual.Name)
	assert.Equal("7aaa145a-5ad8-4d63-be87-177d6abcf1b5", actual.User.String())
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)

	assert.Equal(1, len(actual.Limits))
	assert.Nil(uuid.Validate(actual.Limits[0].Id.String()))
	assert.Equal("limit", actual.Limits[0].Name)
	assert.Equal("test", actual.Limits[0].Value)
	assert.True(actual.Limits[0].CreatedAt.After(beforeConversion))
	assert.Equal(actual.Limits[0].CreatedAt, actual.Limits[0].UpdatedAt)
}

func TestToUserLimitDtoResponse(t *testing.T) {
	assert := assert.New(t)

	u := persistence.UserLimit{
		Id:   uuid.MustParse("2f3b7c63-5b4a-422a-bd9d-7da0f78b6294"),
		Name: "my-limit",
		User: uuid.MustParse("3657b088-ba88-497a-a158-9d6c7faae94f"),
		Limits: []persistence.Limit{
			{
				Id:    defaultUuid,
				Name:  "my-name",
				Value: "my-value",

				CreatedAt: someTime,
				UpdatedAt: someTime,
			},
		},
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	actual := ToUserLimitDtoResponse(u)

	assert.Equal(u.Id, actual.Id)
	assert.Equal(u.Name, actual.Name)
	assert.Equal(u.User, actual.User)
	assert.Equal(u.CreatedAt, actual.CreatedAt)

	assert.Equal(len(u.Limits), len(actual.Limits))
	assert.Equal(u.Limits[0].Name, actual.Limits[0].Name)
	assert.Equal(u.Limits[0].Value, actual.Limits[0].Value)
}

func TestUserLimitDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	u := UserLimitDtoResponse{
		Id:        uuid.MustParse("2f3b7c63-5b4a-422a-bd9d-7da0f78b6294"),
		Name:      "my-limit",
		User:      uuid.MustParse("3657b088-ba88-497a-a158-9d6c7faae94f"),
		CreatedAt: someTime,
		Limits: []LimitDtoResponse{
			{
				Name:  "limit-1",
				Value: "my-value",
			},
		},
	}

	out, err := json.Marshal(u)

	assert.Nil(err)
	assert.Equal(`{"id":"2f3b7c63-5b4a-422a-bd9d-7da0f78b6294","name":"my-limit","user":"3657b088-ba88-497a-a158-9d6c7faae94f","limits":[{"name":"limit-1","value":"my-value"}],"createdAt":"2024-05-05T20:50:18.651387237Z"}`, string(out))
}
