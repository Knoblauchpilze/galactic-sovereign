package communication

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var someTime = time.Date(2024, 05, 05, 20, 50, 18, 651387237, time.UTC)

func TestToUserDtoResponse(t *testing.T) {
	assert := assert.New(t)

	u := persistence.User{
		Id:       defaultUuid,
		Email:    "email",
		Password: "password",

		CreatedAt: someTime,
	}

	actual := ToUserDtoResponse(u)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal("email", actual.Email)
	assert.Equal("password", actual.Password)
	assert.Equal(someTime, actual.CreatedAt)
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
