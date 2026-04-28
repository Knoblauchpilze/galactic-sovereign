package communication

import (
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetResourceProductionDtoResponse struct {
	Planet     uuid.UUID  `json:"planet" format:"uuid"`
	Building   *uuid.UUID `json:"building,omitempty" format:"uuid"`
	Resource   uuid.UUID  `json:"resource" format:"uuid"`
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
