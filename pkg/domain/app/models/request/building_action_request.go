package request

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type BuildingActionCreationRequest struct {
	Planet   uuid.UUID `json:"planet" format:"uuid"`
	Building uuid.UUID `json:"building" format:"uuid"`
}

func FromBuildingActionCreationRequest(action BuildingActionCreationRequest) models.BuildingAction {
	t := time.Now()
	return models.BuildingAction{
		Id:          uuid.New(),
		Planet:      action.Planet,
		Building:    action.Building,
		CreatedAt:   t,
		Version:     0,
		Costs:       []models.BuildingActionCost{},
		Storages:    []models.BuildingActionResourceStorage{},
		Productions: []models.BuildingActionResourceProduction{},
	}
}
