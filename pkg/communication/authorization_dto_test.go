package communication

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultAcls = []AclDtoResponse{
	{
		Id:   defaultUuid,
		User: defaultUser,

		Resource:    "my-resource",
		Permissions: []string{"GET", "DELETE"},
		CreatedAt:   someTime,
	},
}
var defaultLimits = []LimitDtoResponse{
	{
		Name:  "my-limit-1",
		Value: "value-1",
	},
	{
		Name:  "my-limit-2",
		Value: "value-2",
	},
}

func TestAuthorizationDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := AuthorizationDtoResponse{
		Acls:   defaultAcls,
		Limits: defaultLimits,
	}
	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"acls": [
			{
				"id": "08ce96a3-3430-48a8-a3b2-b1c987a207ca",
				"user": "c74a22da-8a05-43a9-a8b9-717e422b0af4",
				"resource": "my-resource",
				"permissions": [
					"GET",
					"DELETE"
				],
				"createdAt": "2024-05-05T20:50:18.651387237Z"
			}
		],
		"limits": [
			{
				"name": "my-limit-1",
				"value": "value-1"
			},
			{
				"name": "my-limit-2",
				"value": "value-2"
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestAuthorizationDtoResponse_WhenAclsAreNil_OutputIsEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := AuthorizationDtoResponse{
		Acls:   nil,
		Limits: defaultLimits,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"acls": [],
		"limits": [
			{
				"name": "my-limit-1",
				"value": "value-1"
			},
			{
				"name": "my-limit-2",
				"value": "value-2"
			}
		]
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestAuthorizationDtoResponse_WhenLimitsIsNil_OutputIsEmptyArray(t *testing.T) {
	assert := assert.New(t)

	dto := AuthorizationDtoResponse{
		Acls:   defaultAcls,
		Limits: nil,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"acls": [
			{
				"id": "08ce96a3-3430-48a8-a3b2-b1c987a207ca",
				"user": "c74a22da-8a05-43a9-a8b9-717e422b0af4",
				"resource": "my-resource",
				"permissions": [
					"GET",
					"DELETE"
				],
				"createdAt": "2024-05-05T20:50:18.651387237Z"
			}
		],
		"limits": []
	}`
	assert.JSONEq(expectedJson, string(out))
}
