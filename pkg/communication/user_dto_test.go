package communication

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var defaultApiKey = uuid.MustParse("cc1742fa-77b4-4f5f-ac92-058c2e47a5d6")
var someTime = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

func TestToUserDtoResponse(t *testing.T) {
	assert := assert.New(t)

	u := persistence.User{
		Id:       defaultUuid,
		Email:    "email",
		Password: "password",

		CreatedAt: someTime,
	}

	actual := ToUserDtoResponse(u, []uuid.UUID{defaultApiKey})

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal("email", actual.Email)
	assert.Equal("password", actual.Password)
	assert.Equal([]uuid.UUID{defaultApiKey}, actual.ApiKeys)
	assert.Equal(someTime, actual.CreatedAt)
}

func TestToUserDtoResponse_WhenApiKeysIsNull_OutptusAnEmptySlice(t *testing.T) {
	assert := assert.New(t)

	u := persistence.User{
		Id:       defaultUuid,
		Email:    "email",
		Password: "password",

		CreatedAt: someTime,
	}

	actual := ToUserDtoResponse(u, nil)

	assert.Equal([]uuid.UUID{}, actual.ApiKeys)
}

func TestFromUserDtoRequest(t *testing.T) {
	assert := assert.New(t)

	beforeConversion := time.Now()

	u := UserDtoRequest{
		Email:    "email",
		Password: "password",
	}

	actual := FromUserDtoRequest(u)

	assert.Nil(uuid.Validate(actual.Id.String()))
	assert.Equal("email", actual.Email)
	assert.Equal("password", actual.Password)
	assert.True(actual.CreatedAt.After(beforeConversion))
	assert.Equal(actual.CreatedAt, actual.UpdatedAt)
}
