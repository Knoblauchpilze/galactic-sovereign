package persistence

import (
	"github.com/google/uuid"
)

type BuildingActionResourceProduction struct {
	Action     uuid.UUID
	Resource   uuid.UUID
	Production int
}

func MergeWithPlanetResourceProduction(actionProduction BuildingActionResourceProduction, planetProduction PlanetResourceProduction) PlanetResourceProduction {
	out := planetProduction
	out.Production = actionProduction.Production
	return out
}

func ToPlanetResourceProduction(actionProduction BuildingActionResourceProduction, action BuildingAction) PlanetResourceProduction {
	out := PlanetResourceProduction{
		Planet:     action.Planet,
		Resource:   actionProduction.Resource,
		Building:   &action.Building,
		Production: actionProduction.Production,

		CreatedAt: action.CompletedAt,
		UpdatedAt: action.CompletedAt,

		Version: 0,
	}

	return out
}
