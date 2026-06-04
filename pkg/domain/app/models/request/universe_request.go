package request

import (
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
)

type UniverseCreationRequest struct {
	Name string
}

func FromUniverseCreationRequest(universe UniverseCreationRequest) models.Universe {
	t := time.Now()
	return models.Universe{
		Id:   uuid.New(),
		Name: universe.Name,

		CreatedAt: t,

		Version: 0,
	}
}
