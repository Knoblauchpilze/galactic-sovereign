package request

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type PlanetCreationRequest struct {
	Player uuid.UUID `json:"player" format:"uuid"`
	Name   string    `json:"name" form:"name"`
}

func FromPlanetCreationRequest(planet PlanetCreationRequest) models.Planet {
	t := time.Now()
	return models.Planet{
		Id:          uuid.New(),
		Player:      planet.Player,
		Name:        planet.Name,
		Homeworld:   false,
		CreatedAt:   t,
		UpdatedAt:   t,
		Version:     0,
		Resources:   []models.PlanetResource{},
		Storages:    []models.PlanetResourceStorage{},
		Productions: []models.PlanetResourceProduction{},
	}
}
