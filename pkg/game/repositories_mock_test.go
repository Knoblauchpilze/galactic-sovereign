package game

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type mockTransaction struct {
	db.Transaction

	timeStamp time.Time

	closeCalled int
}

func (m *mockTransaction) Close(ctx context.Context) {
	m.closeCalled++
}

func (m *mockTransaction) TimeStamp() time.Time {
	return m.timeStamp
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
	deleteForPlanetCalled  int
	deleteForPlanetId      uuid.UUID
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

func (m *mockPlanetResourceRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet
	return m.err
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
	deleteForPlanetCalled            int
	deleteForPlanetId                uuid.UUID
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

func (m *mockPlanetResourceProductionRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
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
	deleteForPlanetCalled         int
	deleteForPlanetId             uuid.UUID
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

func (m *mockPlanetResourceStorageRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
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
