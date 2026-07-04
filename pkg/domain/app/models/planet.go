package models

import (
	"math"
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

// AddBuildingAction adds a building action to the planet.
// The action will be added with a creation date equal to the UpdatedAt
// field of the planet. This means that prior to calling this function,
// callers are expected to trigger UpdateToTime to the desired time.
// The UpdatedAt field will not be updated.
func (p *Planet) AddBuildingAction(building Building) error {
	if p.BuildingAction != nil {
		return domainerrors.ErrActionAlreadyInProgress
	}

	pb, err := p.findBuildingById(building.Id)
	if err != nil {
		return err
	}

	action := building.CreateBuildingAction(pb.Level+1, p.UpdatedAt)

	if err := p.validateEnoughResources(action); err != nil {
		return err
	}

	p.deductResources(action)

	p.BuildingAction = &action

	p.Version++

	return nil
}

// CancelBuildingAction deletes a building action from the planet.
// In case there is no action running an error will be returned.
// The resources used up by the action will be credited back to the
// resources stored on the planet.
// This means that prior to calling this function, callers are
// expected to trigger UpdateToTime to the desired time.
// The UpdatedAt field will not be updated.
func (p *Planet) CancelBuildingAction() error {
	if p.BuildingAction == nil {
		return domainerrors.ErrNoActionInProgress
	}

	p.creditResources(*p.BuildingAction)

	p.BuildingAction = nil

	p.Version++

	return nil
}

func (p *Planet) UpdateToTime(moment time.Time) error {
	if p.UpdatedAt.After(moment) {
		return nil
	}

	if p.BuildingAction != nil && moment.After(p.BuildingAction.CompletedAt) {
		return domainerrors.ErrPlanetNotUpToDate
	}

	elapsed := moment.Sub(p.UpdatedAt)
	hours := elapsed.Hours()

	productions := make(map[uuid.UUID]float64)
	for _, pr := range p.Productions {
		existing := productions[pr.Resource]
		existing += float64(pr.Production)
		productions[pr.Resource] = existing
	}

	storages := make(map[uuid.UUID]float64)
	for _, s := range p.Storages {
		storages[s.Resource] = float64(s.Storage)
	}

	for id, r := range p.Resources {
		prod, ok := productions[r.Resource]
		if !ok {
			continue
		}

		storage, ok := storages[r.Resource]
		if !ok {
			continue
		}

		fullAmount := p.Resources[id].Amount
		if fullAmount >= storage {
			continue
		}

		fullAmount += prod * hours
		fullAmount = math.Min(fullAmount, storage)

		p.Resources[id].Amount = fullAmount
	}

	p.UpdatedAt = moment
	p.Version++

	return nil
}

func (p *Planet) ApplyAction() error {
	if p.BuildingAction == nil {
		return domainerrors.ErrNoActionInProgress
	}

	if p.BuildingAction.CompletedAt != p.UpdatedAt {
		return domainerrors.ErrActionNotCompleted
	}

	p.updateProductions()
	p.updateStorages()

	for id := range p.Buildings {
		if p.Buildings[id].Building == p.BuildingAction.Building {
			p.Buildings[id].Level = p.BuildingAction.DesiredLevel
		}
	}

	p.BuildingAction = nil

	p.Version++

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

func (p *Planet) updateProductions() {
	temp := make(map[uuid.UUID]int)
	for id, pr := range p.Productions {
		if p.Productions[id].Building == nil || *p.Productions[id].Building != p.BuildingAction.Building {
			continue
		}

		temp[pr.Resource] = id
	}

	for _, pp := range p.BuildingAction.Productions {
		id, ok := temp[pp.Resource]

		if ok {
			p.Productions[id].Production = pp.Production
		} else {
			newProd := PlanetResourceProduction{
				Resource:   pp.Resource,
				Building:   &p.BuildingAction.Building,
				Production: pp.Production,
			}
			p.Productions = append(p.Productions, newProd)
		}
	}
}

func (p *Planet) updateStorages() {
	temp := make(map[uuid.UUID]int)
	for id, s := range p.Storages {
		temp[s.Resource] = id
	}

	for _, s := range p.BuildingAction.Storages {
		id := temp[s.Resource]
		p.Storages[id].Storage = s.Storage
	}
}

func (p *Planet) findBuildingById(id uuid.UUID) (PlanetBuilding, error) {
	for _, b := range p.Buildings {
		if b.Building == id {
			return b, nil
		}
	}

	return PlanetBuilding{}, domainerrors.ErrBuildingNotFound
}
