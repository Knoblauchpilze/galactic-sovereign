package models

import (
	"math/rand"

	"github.com/google/uuid"
)

type OccupancyMap struct {
	Topology  UniverseTopology
	UsedSlots map[Coordinate]uuid.UUID
}

// PickPosition picks a random unoccupied position in the topology and attach it
// to the planet provided as input. This function can take a long time to complete
// or even stall in case all slots have been used.
// It's a rare enough event that it is not accounted for yet.
func (m *OccupancyMap) PickPosition(planet uuid.UUID) Coordinate {
	used := true

	if m.UsedSlots == nil {
		m.UsedSlots = make(map[Coordinate]uuid.UUID)
	}

	var out Coordinate

	for used {
		out = Coordinate{
			Galaxy:      rand.Intn(m.Topology.Galaxies),
			SolarSystem: rand.Intn(m.Topology.SolarSystems),
			Position:    rand.Intn(m.Topology.Orbits),
		}

		_, used = m.UsedSlots[out]
	}

	m.UsedSlots[out] = planet

	return out
}
