package request

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type UniverseCreationRequest struct {
	Name     string
	Topology TopologyRequest
}

type TopologyRequest struct {
	Galaxies     int
	SolarSystems int
	Orbits       int
}

func FromUniverseCreationRequest(universe UniverseCreationRequest) models.Universe {
	t := time.Now()
	return models.Universe{
		Id:   uuid.New(),
		Name: universe.Name,
		Topology: models.UniverseTopology{
			Galaxies:     universe.Topology.Galaxies,
			SolarSystems: universe.Topology.SolarSystems,
			Orbits:       universe.Topology.Orbits,
		},

		CreatedAt: t,

		Version: 0,
	}
}
