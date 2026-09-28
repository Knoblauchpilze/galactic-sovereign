package request

import (
	"testing"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/stretchr/testify/assert"
)

func TestUnit_FleetDestinationRequest_ToCoordinates(t *testing.T) {
	request := FleetDestinationRequest{
		Galaxy:      2,
		SolarSystem: 142,
		Position:    6,
	}

	actual := request.ToCoordinates()

	expected := models.Coordinate{
		Galaxy:      request.Galaxy,
		SolarSystem: request.SolarSystem,
		Position:    request.Position,
	}
	assert.Equal(t, expected, actual)
}
