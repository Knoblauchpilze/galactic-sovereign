package communication

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultKey = uuid.MustParse("872e9e40-ce61-497e-b606-c7a08a4faa14")
var defaultUser = uuid.MustParse("c74a22da-8a05-43a9-a8b9-717e422b0af4")

func TestApiKeyDtoResponse_MarshalsToCamelCase(t *testing.T) {
	assert := assert.New(t)

	ak := ApiKeyDtoResponse{
		User:       defaultUser,
		Key:        defaultKey,
		ValidUntil: someTime,
	}

	out, err := json.Marshal(ak)

	assert.Nil(err)
	assert.Equal(`{"user":"c74a22da-8a05-43a9-a8b9-717e422b0af4","key":"872e9e40-ce61-497e-b606-c7a08a4faa14","validUntil":"2024-05-05T20:50:18.651387237Z"}`, string(out))
}

func TestToApiKeyDtoResponse(t *testing.T) {
	assert := assert.New(t)

	k := persistence.ApiKey{
		Id:         defaultUuid,
		Key:        defaultKey,
		ApiUser:    defaultUser,
		ValidUntil: someTime.Add(2 * time.Hour),
	}

	actual := ToApiKeyDtoResponse(k)

	assert.Equal(defaultUser, actual.User)
	assert.Equal(defaultKey, actual.Key)
	assert.Equal(k.ValidUntil, actual.ValidUntil)
}
