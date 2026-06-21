package models

import (
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
)

type Planet struct {
	Id        uuid.UUID
	Player    uuid.UUID
	Name      string
	Homeworld bool

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int

	Resources   []PlanetResource
	Storages    []PlanetResourceStorage
	Productions []PlanetResourceProduction

	Buildings []PlanetBuilding

	BuildingAction *BuildingAction
}

type PlanetResource struct {
	Resource uuid.UUID
	Amount   float64
}

type PlanetResourceStorage struct {
	Resource uuid.UUID
	Storage  int
}

type PlanetResourceProduction struct {
	Resource   uuid.UUID
	Building   *uuid.UUID
	Production int
}

type PlanetBuilding struct {
	Building uuid.UUID
	Level    int
}

func (p *Planet) AddBuildingAction(building Building) error {
	if p.BuildingAction != nil {
		return domainerrors.ErrActionAlreadyInProgress
	}

	pb, err := p.findBuildingById(building.Id)
	if err != nil {
		return err
	}

	action := building.CreateBuildingAction(p.Id, pb.Level+1)

	if err := p.validateEnoughResources(action); err != nil {
		return err
	}

	p.deductResources(action)

	p.BuildingAction = &action

	p.bumpVersion(action.CreatedAt)

	return nil
}

func (p *Planet) CancelBuildingAction() error {
	if p.BuildingAction == nil {
		return domainerrors.ErrNoActionInProgress
	}

	p.creditResources(*p.BuildingAction)

	p.BuildingAction = nil

	p.bumpVersion(time.Now().UTC())

	return nil
}

func (p *Planet) validateEnoughResources(
	action BuildingAction,
) error {
	temp := make(map[uuid.UUID]PlanetResource)
	for _, resource := range p.Resources {
		temp[resource.Resource] = resource
	}

	for _, cost := range action.Costs {
		actual, ok := temp[cost.Resource]
		if !ok || actual.Amount < float64(cost.Amount) {
			return domainerrors.ErrNotEnoughResources
		}
	}

	return nil
}

func (p *Planet) deductResources(
	action BuildingAction,
) {
	temp := make(map[uuid.UUID]BuildingActionCost)
	for _, cost := range action.Costs {
		temp[cost.Resource] = cost
	}

	for id, resource := range p.Resources {
		cost, ok := temp[resource.Resource]
		if ok {
			p.Resources[id].Amount -= float64(cost.Amount)
		}
	}
}

func (p *Planet) creditResources(
	action BuildingAction,
) {
	temp := make(map[uuid.UUID]int)
	for id, pr := range p.Resources {
		temp[pr.Resource] = id
	}

	for _, c := range action.Costs {
		id, ok := temp[c.Resource]
		if ok {
			p.Resources[id].Amount += float64(c.Amount)
		} else {
			pr := PlanetResource{
				Resource: c.Resource,
				Amount:   float64(c.Amount),
			}
			p.Resources = append(p.Resources, pr)
		}
	}
}

func (p *Planet) bumpVersion(timestamp time.Time) {
	p.UpdatedAt = timestamp
	p.Version++
}

func (p *Planet) findBuildingById(id uuid.UUID) (PlanetBuilding, error) {
	for _, b := range p.Buildings {
		if b.Building == id {
			return b, nil
		}
	}

	return PlanetBuilding{}, domainerrors.ErrBuildingNotFound
}
