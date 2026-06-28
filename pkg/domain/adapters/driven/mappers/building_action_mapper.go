package mappers

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type DbBuildingAction struct {
	Id           uuid.UUID
	Planet       uuid.UUID
	Building     uuid.UUID
	DesiredLevel int

	CreatedAt   time.Time
	CompletedAt time.Time
}

func (a DbBuildingAction) ToDomain() models.BuildingAction {
	return models.BuildingAction{
		Id:           a.Id,
		Planet:       a.Planet,
		Building:     a.Building,
		DesiredLevel: a.DesiredLevel,

		CreatedAt:   a.CreatedAt,
		CompletedAt: a.CompletedAt,
	}
}
