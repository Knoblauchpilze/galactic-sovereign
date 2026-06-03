package request

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnit_FromUniverseCreationRequest(t *testing.T) {
	beforeConversion := time.Now()

	request := UniverseCreationRequest{
		Name: "my-universe",
	}

	actual := FromUniverseCreationRequest(request)

	assert.Equal(t, request.Name, actual.Name)
	assert.True(t, actual.CreatedAt.After(beforeConversion))
	assert.Zero(t, actual.Version)
}
