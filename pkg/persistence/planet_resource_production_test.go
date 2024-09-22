package persistence

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var someTime = time.Date(2024, 9, 22, 21, 9, 46, 651387252, time.UTC)

var metalProduction = PlanetResourceProduction{
	Planet:     planetId,
	Resource:   uuid.MustParse("6e6c7d14-1ae2-41e4-9fbd-43217e335f28"),
	Building:   nil,
	Production: 32,

	CreatedAt: someTime,
	UpdatedAt: someTime,

	Version: 12,
}

var crystalProduction = PlanetResourceProduction{
	Planet:     planetId,
	Resource:   uuid.MustParse("500de111-ef71-4e15-91e3-34b87bbad396"),
	Building:   &buildingId,
	Production: 26,

	CreatedAt: someTime,
	UpdatedAt: someTime,

	Version: 7,
}

func TestToPlanetResourceProductionMap(t *testing.T) {
	assert := assert.New(t)

	in := []PlanetResourceProduction{metalProduction, crystalProduction}

	actual := ToPlanetResourceProductionMap(in)

	assert.Equal(2, len(actual))

	metal, ok := actual[metalProduction.Resource]
	assert.True(ok)
	assert.Equal(metalProduction, metal)

	crystal, ok := actual[crystalProduction.Resource]
	assert.True(ok)
	assert.Equal(crystalProduction, crystal)
}
