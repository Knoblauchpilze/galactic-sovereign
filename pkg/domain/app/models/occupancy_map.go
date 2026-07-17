package models

import "math/rand"

type OccupancyMap struct {
	Topology  UniverseTopology
	UsedSlots map[Coordinate]struct{}
}

// PickPosition picks a random unoccupied position in the topology. This function
// can take a long time to complete or even stall in case all slots have been used.
// It's a rare enough event that it is not accounted for yet.
func (m *OccupancyMap) PickPosition() Coordinate {
	used := true

	if m.UsedSlots == nil {
		m.UsedSlots = make(map[Coordinate]struct{})
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

	m.UsedSlots[out] = struct{}{}

	return out
}
