package models

import (
	"math"
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
)

type Planet struct {
	Id         uuid.UUID
	Player     uuid.UUID
	Name       string
	Homeworld  bool
	Coordinate Coordinate
	Fields     int

	CreatedAt time.Time
	UpdatedAt time.Time

	Version int

	Resources   []PlanetResource
	Storages    []PlanetResourceStorage
	Productions []PlanetResourceProduction
	Buildings   []PlanetBuilding
	Ships       []PlanetShip

	BuildingAction *BuildingAction
	ShipActions    []ShipAction
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

type PlanetShip struct {
	Ship  uuid.UUID
	Count int
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

	if !p.fieldsAvailable() {
		return domainerrors.ErrAllFieldsUsed
	}

	action := building.CreateBuildingAction(pb.Level+1, p.UpdatedAt)

	if err := p.validateEnoughResourcesForBuildingAction(action); err != nil {
		return err
	}

	p.deductBuildingActionResources(action)

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

// AddShipAction adds a ship action to the planet.
// The action will be added with a creation date equal to the UpdatedAt
// field of the planet. This means that prior to calling this function,
// callers are expected to trigger UpdateToTime to the desired time.
// The UpdatedAt field will not be updated.
func (p *Planet) AddShipAction(ship Ship, count int) error {
	err := p.validateShipExists(ship.Id)
	if err != nil {
		return err
	}

	nextActionStartTime := p.determineShipActionStartTime()
	action := ship.CreateShipAction(count, nextActionStartTime)

	if err := p.validateEnoughResourcesForShipAction(action); err != nil {
		return err
	}

	p.deductShipActionResources(action)

	p.ShipActions = append(p.ShipActions, action)

	p.Version++

	return nil
}

func (p *Planet) UpdateToTime(moment time.Time) error {
	if p.UpdatedAt.Compare(moment) >= 0 {
		return nil
	}

	err := p.checkUpToDate(moment)
	if err != nil {
		return err
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

func (p *Planet) ApplyBuildingAction() error {
	if p.BuildingAction == nil {
		return domainerrors.ErrNoActionInProgress
	}

	if p.BuildingAction.CompletedAt != p.UpdatedAt {
		return domainerrors.ErrBuildingActionNotCompleted
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

func (p *Planet) ApplyShipAction() error {
	if len(p.ShipActions) == 0 {
		return domainerrors.ErrNoActionInProgress
	}

	action := &p.ShipActions[0]
	if action.NextCompletionAt != p.UpdatedAt {
		return domainerrors.ErrShipActionNotCompleted
	}

	for id := range p.Ships {
		if p.Ships[id].Ship == action.Ship {
			p.Ships[id].Count++
		}
	}

	err := action.CompleteOne()
	if err != nil {
		return err
	}

	if action.Completed() {
		p.ShipActions = p.ShipActions[1:]
	}

	p.Version++

	return nil
}

func (p *Planet) findBuildingById(id uuid.UUID) (PlanetBuilding, error) {
	for _, b := range p.Buildings {
		if b.Building == id {
			return b, nil
		}
	}

	return PlanetBuilding{}, domainerrors.ErrBuildingNotFound
}

func (p *Planet) fieldsAvailable() bool {
	used := 0
	for _, b := range p.Buildings {
		used += b.Level
	}

	return used < p.Fields
}

func (p *Planet) validateEnoughResourcesForBuildingAction(
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

func (p *Planet) deductBuildingActionResources(
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

func (p *Planet) validateShipExists(id uuid.UUID) error {
	for _, s := range p.Ships {
		if s.Ship == id {
			return nil
		}
	}

	return domainerrors.ErrShipNotFound
}

func (p *Planet) determineShipActionStartTime() time.Time {
	if len(p.ShipActions) == 0 {
		return p.UpdatedAt
	}

	earliest := p.UpdatedAt
	for _, action := range p.ShipActions {
		completedAt := action.CompletionTime()
		if completedAt.After(earliest) {
			earliest = completedAt
		}
	}

	return earliest
}

func (p *Planet) validateEnoughResourcesForShipAction(
	action ShipAction,
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

func (p *Planet) deductShipActionResources(
	action ShipAction,
) {
	temp := make(map[uuid.UUID]ShipActionCost)
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

func (p *Planet) checkUpToDate(moment time.Time) error {
	if p.BuildingAction != nil && p.BuildingAction.CompletedAt.Compare(moment) < 0 {
		return domainerrors.ErrPlanetNotUpToDate
	}

	for _, action := range p.ShipActions {
		if action.NextCompletionAt.Compare(moment) < 0 {
			return domainerrors.ErrPlanetNotUpToDate
		}
	}

	return nil
}
