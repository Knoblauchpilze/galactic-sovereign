package request

import (
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/stretchr/testify/assert"
)

func TestUnit_FromUniverseCreationRequest(t *testing.T) {
	beforeConversion := time.Now()

	request := UniverseCreationRequest{
		Name: "my-universe",
		Topology: TopologyRequest{
			Galaxies:     36,
			SolarSystems: 487,
			Orbits:       8,
		},
	}

	actual := FromUniverseCreationRequest(request)

	assert.Equal(t, request.Name, actual.Name)
	expectedTopology := models.UniverseTopology{
		Galaxies:     request.Topology.Galaxies,
		SolarSystems: request.Topology.SolarSystems,
		Orbits:       request.Topology.Orbits,
	}
	assert.Equal(t, expectedTopology, actual.Topology)
	assert.True(t, actual.CreatedAt.After(beforeConversion))
	assert.Zero(t, actual.Version)
}
