package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAclDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := AclDtoRequest{
		User: defaultUser,

		Resource: "my-resource",
		Permissions: []string{
			"permission-1",
			"permission-2",
		},
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	// https://stackoverflow.com/questions/32408890/how-to-compare-two-json-requests
	expectedJson := `
	{
		"user": "c74a22da-8a05-43a9-a8b9-717e422b0af4",
		"resource": "my-resource",
		"permissions": [
			"permission-1",
			"permission-2"
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFromAclDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := AclDtoRequest{
		User: defaultUser,

		Resource: "my-resource",
		Permissions: []string{
			"permission-1",
			"permission-2",
		},
	}

	actual := FromAclDtoRequest(dto)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal(defaultUser, actual.User)
	assert.Equal("my-resource", actual.Resource)
	assert.Equal([]string{"permission-1", "permission-2"}, actual.Permissions)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestToAclDtoResponse(t *testing.T) {
	assert := assert.New(t)

	entity := persistence.Acl{
		Id:   defaultUuid,
		User: defaultUser,

		Resource: "some-resource",
		Permissions: []string{
			"my-permission",
			"another-permission",
		},

		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	actual := ToAclDtoResponse(entity)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal(defaultUser, actual.User)

	assert.Equal("some-resource", actual.Resource)
	assert.Equal([]string{"my-permission", "another-permission"}, actual.Permissions)

	assert.Equal(someTime, actual.CreatedAt)
}

func TestAclDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := AclDtoResponse{
		Id:   defaultUuid,
		User: defaultUser,

		Resource: "some-resource",
		Permissions: []string{
			"my-permission",
			"another-permission",
		},

		CreatedAt: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "08ce96a3-3430-48a8-a3b2-b1c987a207ca",
		"user": "c74a22da-8a05-43a9-a8b9-717e422b0af4",
		"resource": "some-resource",
		"permissions": [
			"my-permission",
			"another-permission"
		],
		"createdAt": "2024-05-05T20:50:18.651387237Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
