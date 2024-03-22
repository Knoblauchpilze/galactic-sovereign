package communication

import (
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultUuid = uuid.MustParse("08ce96a3-3430-48a8-a3b2-b1c987a207ca")
var someTime = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

func TestFromUser(t *testing.T) {
	assert := assert.New(t)

	u := persistence.User{
		Id:       defaultUuid,
		Email:    "email",
		Password: "password",

		CreatedAt: someTime,
	}

	actual := FromUser(u)

	assert.Equal(defaultUuid, actual.Id)
	assert.Equal("email", actual.Email)
	assert.Equal("password", actual.Password)
	assert.Equal(someTime, actual.CreatedAt)
}
