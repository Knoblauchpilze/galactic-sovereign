package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	homeworldDefaultName string = "homeworld"
	planetDefaultName    string = "colony"
)

type Player struct {
	Id       uuid.UUID
	ApiUser  uuid.UUID
	Universe uuid.UUID
	Name     string

	CreatedAt time.Time

	Version int

	Homeworld uuid.UUID
	Planets   []uuid.UUID
}

func (p *Player) CreateHomeworld(
	resources []Resource,
	buildings []Building,
) Planet {
	createdAt := time.Now()

	planet := createPlanet(p.Id, createdAt, true, resources, buildings)

	p.Homeworld = planet.Id
	p.Planets = []uuid.UUID{planet.Id}

	return planet
}

func (p *Player) Colonize(
	resources []Resource,
	buildings []Building,
) Planet {
	createdAt := time.Now()

	planet := createPlanet(p.Id, createdAt, false, resources, buildings)

	p.Planets = append(p.Planets, planet.Id)

	return planet
}

func createPlanet(
	player uuid.UUID,
	createdAt time.Time,
	homeworld bool,
	resources []Resource,
	buildings []Building,
) Planet {
	planetResources := make([]PlanetResource, 0, len(resources))
	planetStorages := make([]PlanetResourceStorage, 0, len(resources))
	planetProductions := make([]PlanetResourceProduction, 0, len(resources))
	planetBuildings := make([]PlanetBuilding, 0, len(buildings))

	for _, r := range resources {
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

	for _, b := range buildings {
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
