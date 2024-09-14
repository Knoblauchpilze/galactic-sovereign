package persistence

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var resourceId = uuid.MustParse("e0e56162-b462-4c05-9bda-828a373037a7")

func TestToPlanetResource(t *testing.T) {
	assert := assert.New(t)

	production := BuildingActionResourceProduction{
		Action:     actionId,
		Resource:   resourceId,
		Production: 54,
	}
	resource := PlanetResource{
		Planet:     planetId,
		Resource:   resourceId,
		Amount:     27.0,
		Production: 10,

		CreatedAt: time.Date(2024, 9, 14, 20, 40, 05, 651387249, time.UTC),
		UpdatedAt: time.Date(2024, 9, 14, 20, 40, 06, 651387249, time.UTC),

		Version: 1,
	}

	actual := ToPlanetResource(production, resource)

	assert.Equal(planetId, actual.Planet)
	assert.Equal(resourceId, actual.Resource)
	assert.Equal(resource.Amount, actual.Amount)
	assert.Equal(production.Production, actual.Production)
	assert.Equal(resource.CreatedAt, actual.CreatedAt)
	assert.Equal(resource.UpdatedAt, actual.UpdatedAt)
	assert.Equal(resource.Version, actual.Version)
}
