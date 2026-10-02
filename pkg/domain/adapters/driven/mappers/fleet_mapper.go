package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbFleet struct {
	Id     uuid.UUID
	Player uuid.UUID
	Source uuid.UUID

	Galaxy      int
	SolarSystem int
	Position    int
	Target      *uuid.UUID

	CreatedAt time.Time
	ArrivalAt time.Time
	ReturnAt  time.Time
	UpdatedAt time.Time
	Version   int
}

func (f DbFleet) ToDomain() models.Fleet {
	return models.Fleet{
		Id:     f.Id,
		Player: f.Player,
		Source: f.Source,
		Destination: models.FleetDestination{
			Coordinate: models.Coordinate{
				Galaxy:      f.Galaxy,
				SolarSystem: f.SolarSystem,
				Position:    f.Position,
			},
			Target: f.Target,
		},
		CreatedAt: f.CreatedAt,
		ArrivalAt: f.ArrivalAt,
		ReturnAt:  f.ReturnAt,
		UpdatedAt: f.UpdatedAt,
		Version:   f.Version,
	}
}
