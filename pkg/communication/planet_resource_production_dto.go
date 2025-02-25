package communication

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceProductionDtoResponse struct {
	Planet     uuid.UUID  `json:"planet"`
	Building   *uuid.UUID `json:"building,omitempty"`
	Resource   uuid.UUID  `json:"resource"`
	Production int        `json:"production"`
}

func ToPlanetResourceProductionDtoResponse(production persistence.PlanetResourceProduction) PlanetResourceProductionDtoResponse {
	return PlanetResourceProductionDtoResponse{
		Planet:     production.Planet,
		Resource:   production.Resource,
		Building:   production.Building,
		Production: production.Production,
	}
}
