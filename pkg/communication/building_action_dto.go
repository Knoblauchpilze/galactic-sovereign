package communication

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionDtoRequest struct {
	Planet   uuid.UUID `json:"planet" format:"uuid"`
	Building uuid.UUID `json:"building" format:"uuid"`
}

type BuildingActionDtoResponse struct {
	Id           uuid.UUID `json:"id" format:"uuid"`
	Planet       uuid.UUID `json:"planet" format:"uuid"`
	Building     uuid.UUID `json:"building" format:"uuid"`
	CurrentLevel int       `json:"currentLevel"`
	DesiredLevel int       `json:"desiredLevel"`
	CreatedAt    time.Time `json:"createdAt" format:"date-time"`
	CompletedAt  time.Time `json:"completedAt" format:"date-time"`
}

func FromBuildingActionDtoRequest(action BuildingActionDtoRequest) persistence.BuildingAction {
	t := time.Now()
	return persistence.BuildingAction{
		Id:        uuid.New(),
		Planet:    action.Planet,
		Building:  action.Building,
		CreatedAt: t,
	}
}

func ToBuildingActionDtoResponse(action persistence.BuildingAction) BuildingActionDtoResponse {
	return BuildingActionDtoResponse{
		Id:           action.Id,
		Planet:       action.Planet,
		Building:     action.Building,
		CurrentLevel: action.CurrentLevel,
		DesiredLevel: action.DesiredLevel,
		CreatedAt:    action.CreatedAt,
		CompletedAt:  action.CompletedAt,
	}
}
