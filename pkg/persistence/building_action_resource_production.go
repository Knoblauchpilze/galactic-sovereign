package persistence

import (
	"github.com/google/uuid"
)

type BuildingActionResourceProduction struct {
	Action     uuid.UUID
	Resource   uuid.UUID
	Production int
}

func ToPlanetResource(production BuildingActionResourceProduction, resource PlanetResource) PlanetResource {
	out := resource
	out.Production = production.Production
	return out
}
