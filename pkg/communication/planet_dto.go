package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type PlanetDtoRequest struct {
	Player uuid.UUID `json:"player"`
	Name   string    `json:"name" form:"name"`
}

type PlanetDtoResponse struct {
	Id     uuid.UUID `json:"id"`
	Player uuid.UUID `json:"player"`
	Name   string    `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
}

func FromPlanetDtoRequest(planet PlanetDtoRequest) persistence.Planet {
	t := time.Now()
	return persistence.Planet{
		Id:     uuid.New(),
		Player: planet.Player,
		Name:   planet.Name,

		CreatedAt: t,
		UpdatedAt: t,
	}
}

func ToPlanetDtoResponse(planet persistence.Planet) PlanetDtoResponse {
	return PlanetDtoResponse{
		Id:     planet.Id,
		Player: planet.Player,
		Name:   planet.Name,

		CreatedAt: planet.CreatedAt,
	}
}
