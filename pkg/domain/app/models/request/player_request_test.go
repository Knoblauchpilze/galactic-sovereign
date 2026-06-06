package request

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_FromPlayerCreationRequest(t *testing.T) {
	beforeConversion := time.Now()

	request := PlayerCreationRequest{
		ApiUser:  uuid.New(),
		Universe: uuid.New(),
		Name:     "my-universe",
	}

	actual := FromPlayerCreationRequest(request)

	assert.Equal(t, request.ApiUser, actual.ApiUser)
	assert.Equal(t, request.Universe, actual.Universe)
	assert.Equal(t, request.Name, actual.Name)
	assert.True(t, actual.CreatedAt.After(beforeConversion))
	assert.Zero(t, actual.Version)
}
