package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
)

var errDefault = fmt.Errorf("some error")
var testDate = time.Date(2024, 04, 01, 11, 8, 47, 651387237, time.UTC)

type mockBuildingActionRepository struct {
	repositories.BuildingActionRepository

	action persistence.BuildingAction
	errs   []error
	calls  int

	createCalled                   int
	createdBuildingAction          persistence.BuildingAction
	getCalled                      int
	getId                          uuid.UUID
	listForPlanetId                uuid.UUID
	listForPlanetCalled            int
	listBeforeCompletionTimeCalled int
	listBeforeCompletionTimePlanet uuid.UUID
	listBeforeCompletionTime       time.Time
	deleteCalled                   int
	deleteId                       uuid.UUID
	deleteForPlanetCalled          int
	deleteForPlanetId              uuid.UUID
}

func (m *mockBuildingActionRepository) Create(ctx context.Context, tx db.Transaction, action persistence.BuildingAction) (persistence.BuildingAction, error) {
	m.createCalled++
	m.createdBuildingAction = action

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.action, *err
}

func (m *mockBuildingActionRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.BuildingAction, error) {
	m.getCalled++
	m.getId = id

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.action, *err
}

func (m *mockBuildingActionRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.BuildingAction, error) {
	m.listForPlanetCalled++
	m.listForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.BuildingAction{m.action}, *err
}

func (m *mockBuildingActionRepository) ListBeforeCompletionTime(ctx context.Context, tx db.Transaction, planet uuid.UUID, until time.Time) ([]persistence.BuildingAction, error) {
	m.listBeforeCompletionTimeCalled++
	m.listBeforeCompletionTimePlanet = planet
	m.listBeforeCompletionTime = until

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.BuildingAction{m.action}, *err
}

func (m *mockBuildingActionRepository) Delete(ctx context.Context, tx db.Transaction, action uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = action

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

func (m *mockBuildingActionRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

func (m *mockBuildingActionRepository) DeleteForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) error {
	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

type mockPlanetBuildingRepository struct {
	repositories.PlanetBuildingRepository

	planetBuilding persistence.PlanetBuilding
	err            error
	updateErr      error

	getForPlanetAndBuildingCalled   int
	getForPlanetAndBuildingPlanet   uuid.UUID
	getForPlanetAndBuildingBuilding uuid.UUID
	listForPlanetCalled             int
	listForPlanetId                 uuid.UUID
	updateCalled                    int
	updateBuilding                  persistence.PlanetBuilding
}

func (m *mockPlanetBuildingRepository) GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building uuid.UUID) (persistence.PlanetBuilding, error) {
	m.getForPlanetAndBuildingCalled++
	m.getForPlanetAndBuildingPlanet = planet
	m.getForPlanetAndBuildingBuilding = building
	return m.planetBuilding, m.err
}

func (m *mockPlanetBuildingRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetBuilding, error) {
	m.listForPlanetCalled++
	m.listForPlanetId = planet
	return []persistence.PlanetBuilding{m.planetBuilding}, m.err
}

func (m *mockPlanetBuildingRepository) Update(ctx context.Context, tx db.Transaction, building persistence.PlanetBuilding) (persistence.PlanetBuilding, error) {
	m.updateCalled++
	m.updateBuilding = building
	return m.updateBuilding, m.updateErr
}

type mockPlanetRepository struct {
	repositories.PlanetRepository

	planet persistence.Planet
	err    error

	createCalled          int
	createdPlanet         persistence.Planet
	getCalled             int
	getId                 uuid.UUID
	listCalled            int
	listForPlayerId       uuid.UUID
	listForPlayerCalled   int
	deleteCalled          int
	deleteId              uuid.UUID
	deleteForPlayerCalled int
	deleteForPlayerId     uuid.UUID
}

func (m *mockPlanetRepository) Create(ctx context.Context, tx db.Transaction, planet persistence.Planet) (persistence.Planet, error) {
	m.createCalled++
	m.createdPlanet = planet
	return m.planet, m.err
}

func (m *mockPlanetRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Planet, error) {
	m.getCalled++
	m.getId = id
	return m.planet, m.err
}

func (m *mockPlanetRepository) List(ctx context.Context, tx db.Transaction) ([]persistence.Planet, error) {
	m.listCalled++
	return []persistence.Planet{m.planet}, m.err
}

func (m *mockPlanetRepository) ListForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) ([]persistence.Planet, error) {
	m.listForPlayerCalled++
	m.listForPlayerId = player
	return []persistence.Planet{m.planet}, m.err
}

func (m *mockPlanetRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

func (m *mockPlanetRepository) DeleteForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) error {
	m.deleteForPlayerCalled++
	m.deleteForPlayerId = player
	return m.err
}

type mockPlanetResourceRepository struct {
	repositories.PlanetResourceRepository

	planetResource persistence.PlanetResource
	err            error
	updateErr      error

	createCalled           int
	createdPlanetResource  persistence.PlanetResource
	listForPlanetIds       []uuid.UUID
	listForPlanetCalled    int
	updateCalled           int
	updatedPlanetResources []persistence.PlanetResource
}

func (m *mockPlanetResourceRepository) Create(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error) {
	m.createCalled++
	m.createdPlanetResource = resource
	return m.planetResource, m.err
}

func (m *mockPlanetResourceRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResource, error) {
	m.listForPlanetCalled++
	m.listForPlanetIds = append(m.listForPlanetIds, planet)
	return []persistence.PlanetResource{m.planetResource}, m.err
}

func (m *mockPlanetResourceRepository) Update(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error) {
	m.updateCalled++
	m.updatedPlanetResources = append(m.updatedPlanetResources, resource)
	return resource, m.updateErr
}

type mockPlanetResourceProductionRepository struct {
	repositories.PlanetResourceProductionRepository

	planetResourceProduction persistence.PlanetResourceProduction
	errs                     []error
	calls                    int
	updateErr                error

	createCalled                     int
	createdPlanetResourceProduction  persistence.PlanetResourceProduction
	getForPlanetAndBuildingCalled    int
	getForPlanetAndBuildingPlanet    uuid.UUID
	getForPlanetAndBuildingBuilding  *uuid.UUID
	listForPlanetIds                 []uuid.UUID
	listForPlanetCalled              int
	updateCalled                     int
	updatedPlanetResourceProductions []persistence.PlanetResourceProduction
}

func (m *mockPlanetResourceProductionRepository) Create(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error) {
	m.createCalled++
	m.createdPlanetResourceProduction = production

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.planetResourceProduction, *err
}

func (m *mockPlanetResourceProductionRepository) GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building *uuid.UUID) (persistence.PlanetResourceProduction, error) {
	m.getForPlanetAndBuildingCalled++
	m.getForPlanetAndBuildingPlanet = planet
	m.getForPlanetAndBuildingBuilding = building

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.planetResourceProduction, *err
}

func (m *mockPlanetResourceProductionRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceProduction, error) {
	m.listForPlanetCalled++
	m.listForPlanetIds = append(m.listForPlanetIds, planet)

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.PlanetResourceProduction{m.planetResourceProduction}, *err
}

func (m *mockPlanetResourceProductionRepository) Update(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error) {
	m.updateCalled++
	m.updatedPlanetResourceProductions = append(m.updatedPlanetResourceProductions, production)
	return production, m.updateErr
}

type mockPlanetResourceStorageRepository struct {
	repositories.PlanetResourceStorageRepository

	planetResourceStorage persistence.PlanetResourceStorage
	errs                  []error
	calls                 int
	updateErr             error

	createCalled                  int
	createdPlanetResourceStorage  persistence.PlanetResourceStorage
	listForPlanetIds              []uuid.UUID
	listForPlanetCalled           int
	updateCalled                  int
	updatedPlanetResourceStorages []persistence.PlanetResourceStorage
}

func (m *mockPlanetResourceStorageRepository) Create(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error) {
	m.createCalled++
	m.createdPlanetResourceStorage = storage

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.planetResourceStorage, *err
}

func (m *mockPlanetResourceStorageRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceStorage, error) {
	m.listForPlanetCalled++
	m.listForPlanetIds = append(m.listForPlanetIds, planet)

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.PlanetResourceStorage{m.planetResourceStorage}, *err
}

func (m *mockPlanetResourceStorageRepository) Update(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error) {
	m.updateCalled++
	m.updatedPlanetResourceStorages = append(m.updatedPlanetResourceStorages, storage)
	return storage, m.updateErr
}

type mockPlayerRepository struct {
	repositories.PlayerRepository

	player persistence.Player
	err    error

	createCalled         int
	createdPlayer        persistence.Player
	getCalled            int
	getId                uuid.UUID
	listCalled           int
	listForApiUserId     uuid.UUID
	listForApiUserCalled int
	deleteCalled         int
	deleteId             uuid.UUID
}

func (m *mockPlayerRepository) Create(ctx context.Context, tx db.Transaction, player persistence.Player) (persistence.Player, error) {
	m.createCalled++
	m.createdPlayer = player
	return m.player, m.err
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uuid.UUID) (persistence.Player, error) {
	m.getCalled++
	m.getId = id
	return m.player, m.err
}

func (m *mockPlayerRepository) List(ctx context.Context) ([]persistence.Player, error) {
	m.listCalled++
	return []persistence.Player{m.player}, m.err
}

func (m *mockPlayerRepository) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]persistence.Player, error) {
	m.listForApiUserCalled++
	m.listForApiUserId = apiUser
	return []persistence.Player{m.player}, m.err
}

func (m *mockPlayerRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

func getValueToReturnOr[T any](count int, values []T, value T) *T {
	out := getValueToReturn(count, values)
	if out == nil {
		return &value
	}

	return out
}

func getValueToReturn[T any](count int, values []T) *T {
	var out *T
	if count > len(values) {
		count = 0
	}
	if count < len(values) {
		out = &values[count]
	}

	return out
}
