package persistence

import (
	"time"

	"github.com/google/uuid"
)

type PlanetResourceProduction struct {
	Planet     uuid.UUID
	Resource   uuid.UUID
	Building   *uuid.UUID
	Production int

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int
}

func ToPlanetResourceProductionMap(in []PlanetResourceProduction) map[uuid.UUID]PlanetResourceProduction {
	out := make(map[uuid.UUID]PlanetResourceProduction)

	for _, production := range in {
		out[production.Resource] = production
	}

	return out
}
