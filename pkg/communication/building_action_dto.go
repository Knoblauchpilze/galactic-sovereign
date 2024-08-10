package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type BuildingActionDtoRequest struct {
	Planet       uuid.UUID `json:"planet"`
	Building     uuid.UUID `json:"building"`
	CurrentLevel int       `json:"currentLevel"`
}

type BuildingActionDtoResponse struct {
	Id           uuid.UUID `json:"id"`
	Planet       uuid.UUID `json:"planet"`
	Building     uuid.UUID `json:"building"`
	CurrentLevel int       `json:"currentLevel"`
	DesiredLevel int       `json:"desiredLevel"`
	CreatedAt    time.Time `json:"createdAt"`
	CompletedAt  time.Time `json:"completedAt"`
}

func FromBuildingActionDtoRequest(action BuildingActionDtoRequest) persistence.BuildingAction {
	t := time.Now()
	return persistence.BuildingAction{
		Id:           uuid.New(),
		Planet:       action.Planet,
		Building:     action.Building,
		CurrentLevel: action.CurrentLevel,
		DesiredLevel: action.CurrentLevel + 1,
		CreatedAt:    t,
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
