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
	CurrentLevel int
	DesiredLevel int

	CreatedAt   time.Time
	CompletedAt time.Time

	Version int
}

func (a DbBuildingAction) ToDomain() models.BuildingAction {
	return models.BuildingAction{
		Id:           a.Id,
		Planet:       a.Planet,
		Building:     a.Building,
		CurrentLevel: a.CurrentLevel,
		DesiredLevel: a.DesiredLevel,

		CreatedAt:   a.CreatedAt,
		CompletedAt: a.CompletedAt,

		Version: a.Version,
	}
}
