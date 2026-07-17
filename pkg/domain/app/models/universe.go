package models

import (
	"time"

	"github.com/google/uuid"
)

type Universe struct {
	Id       uuid.UUID
	Name     string
	Topology UniverseTopology

	CreatedAt time.Time

	Version int

	Resources []Resource
	Buildings []Building

	OccupancyMap OccupancyMap
}

type UniverseTopology struct {
	Galaxies     int
	SolarSystems int
	Orbits       int
}

func (u Universe) CreatePlanet(player uuid.UUID, homeworld bool) Planet {
	createdAt := time.Now()

	coordinate := u.OccupancyMap.PickPosition()
	fields := coordinate.Fields(homeworld)

	planetResources := make([]PlanetResource, 0, len(u.Resources))
	planetStorages := make([]PlanetResourceStorage, 0, len(u.Resources))
	planetProductions := make([]PlanetResourceProduction, 0, len(u.Resources))
	planetBuildings := make([]PlanetBuilding, 0, len(u.Buildings))

	for _, r := range u.Resources {
		pr := PlanetResource{
			Resource: r.Id,
			Amount:   float64(r.StartAmount),
		}
		planetResources = append(planetResources, pr)

		ps := PlanetResourceStorage{
			Resource: r.Id,
			Storage:  r.StartStorage,
		}
		planetStorages = append(planetStorages, ps)

		pp := PlanetResourceProduction{
			Resource:   r.Id,
			Production: r.StartProduction,
		}
		planetProductions = append(planetProductions, pp)
	}

	for _, b := range u.Buildings {
		pb := PlanetBuilding{
			Building: b.Id,
			Level:    0,
		}
		planetBuildings = append(planetBuildings, pb)
	}

	name := homeworldDefaultName
	if !homeworld {
		name = planetDefaultName
	}

	return Planet{
		Id:             uuid.New(),
		Player:         player,
		Name:           name,
		Homeworld:      homeworld,
		Coordinate:     coordinate,
		Fields:         fields,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
		Version:        0,
		Resources:      planetResources,
		Storages:       planetStorages,
		Productions:    planetProductions,
		Buildings:      planetBuildings,
		BuildingAction: nil,
	}
}
