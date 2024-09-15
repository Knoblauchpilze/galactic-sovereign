package persistence

import (
	"github.com/google/uuid"
)

type BuildingActionResourceProduction struct {
	Action     uuid.UUID
	Resource   uuid.UUID
	Production int
}

func ToPlanetResourceProduction(actionProduction BuildingActionResourceProduction, planetProduction PlanetResourceProduction) PlanetResourceProduction {
	out := planetProduction
	out.Production = actionProduction.Production
	return out
}
