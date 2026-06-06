package request

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_FromPlanetCreationRequest(t *testing.T) {
	beforeConversion := time.Now()

	request := PlanetCreationRequest{
		Player: uuid.New(),
		Name:   "my-planet",
	}

	actual := FromPlanetCreationRequest(request)

	assert.Equal(t, request.Player, actual.Player)
	assert.Equal(t, request.Name, actual.Name)
	assert.False(t, actual.Homeworld)
	assert.True(t, actual.CreatedAt.After(beforeConversion))
	assert.Equal(t, actual.CreatedAt, actual.UpdatedAt)
	assert.Zero(t, actual.Version)
	assert.Equal(t, []models.PlanetResource{}, actual.Resources)
	assert.Equal(t, []models.PlanetResourceStorage{}, actual.Storages)
	assert.Equal(t, []models.PlanetResourceProduction{}, actual.Productions)
}
