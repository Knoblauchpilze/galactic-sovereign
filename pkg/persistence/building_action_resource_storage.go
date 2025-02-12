package persistence

import (
	"github.com/google/uuid"
)

type BuildingActionResourceStorage struct {
	Action   uuid.UUID
	Resource uuid.UUID
	Storage  int
}

func MergeWithPlanetResourceStorage(actionStorage BuildingActionResourceStorage, planetStorage PlanetResourceStorage) PlanetResourceStorage {
	out := planetStorage
	out.Storage = actionStorage.Storage
	return out
}

func ToPlanetResourceStorage(actionStorage BuildingActionResourceStorage, action BuildingAction) PlanetResourceStorage {
	out := PlanetResourceStorage{
		Planet:   action.Planet,
		Resource: actionStorage.Resource,
		Storage:  actionStorage.Storage,

		CreatedAt: action.CompletedAt,
		UpdatedAt: action.CompletedAt,

		Version: 0,
	}

	return out
}
