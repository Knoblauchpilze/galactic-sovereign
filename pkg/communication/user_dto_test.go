package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var someTime = time.Date(2024, 05, 05, 20, 50, 18, 651387237, time.UTC)

func TestUserDtoRequest_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := UserDtoRequest{
		Email:    "some@e.mail",
		Password: "secret",
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"email": "some@e.mail",
		"password": "secret"
	}`
	assert.JSONEq(expectedJson, string(out))
}

func TestFromUserDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	dto := UserDtoRequest{
		Email:    "email",
		Password: "password",
	}

	actual := FromUserDtoRequest(dto)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal("email", actual.Email)
	assert.Equal("password", actual.Password)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}

func TestToUserDtoResponse(t *testing.T) {
	assert := assert.New(t)

	entity := persistence.User{
		Id:       defaultUuid,
		Email:    "email",
		Password: "password",

		CreatedAt: someTime,
	}

	actual := ToUserDtoResponse(entity)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal("email", actual.Email)
	assert.Equal("password", actual.Password)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestUserDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	dto := UserDtoResponse{
		Id:        defaultUuid,
		Email:     "some@e.mail",
		Password:  "secret",
		CreatedAt: someTime,
	}

	out, err := json.Marshal(dto)

	assert.Nil(err)
	expectedJson := `
	{
		"id": "08ce96a3-3430-48a8-a3b2-b1c987a207ca",
		"email": "some@e.mail",
		"password": "secret",
		"createdAt": "2024-05-05T20:50:18.651387237Z"
	}`
	assert.JSONEq(expectedJson, string(out))
}
