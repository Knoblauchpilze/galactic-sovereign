package service

import (
	"context"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var someTime = time.Date(2024, 8, 17, 14, 31, 13, 651387247, time.UTC)

func TestUnit_PlanetResourceService(t *testing.T) {
	s := ServicePoolTestSuite{
		generateRepositoriesMocks: generatePlanetResourceServiceMocks,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"whenUpdatingPlanetResources_expectPlanetResourcesChangedWithCorrectValues": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlanetResourceService(conn, repos)
					threeMinutesAfterUpdatedAt := defaultPlanetResource.UpdatedAt.Add(3 * time.Minute)
					return s.UpdatePlanetUntil(ctx, defaultPlanetId, threeMinutesAfterUpdatedAt)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))
					actual := m.updatedPlanetResources[0]

					assert.Equal(defaultPlanetId, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					expectedAmount := defaultPlanetResource.Amount + 3.0/60.0*float64(defaultPlanetResourceProduction.Production)
					assert.Equal(expectedAmount, actual.Amount)

					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					expectedUpdatedAt := defaultPlanetResource.UpdatedAt.Add(3 * time.Minute)
					assert.Equal(expectedUpdatedAt, actual.UpdatedAt)
					assert.Equal(defaultPlanetResource.Version, actual.Version)
				},
			},
			"whenUpdatingPlanetResources_whenStorageIsAlreadyFull_expectNoUpdate": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetResourceServiceMocks()

					storage := defaultPlanetResourceStorage
					storage.Storage = int(defaultPlanetResource.Amount) - 10
					repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
						planetResourceStorage: storage,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlanetResourceService(conn, repos)
					threeMinutesAfterUpdatedAt := defaultPlanetResource.UpdatedAt.Add(3 * time.Minute)
					return s.UpdatePlanetUntil(ctx, defaultPlanetId, threeMinutesAfterUpdatedAt)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))
					actual := m.updatedPlanetResources[0]

					assert.Equal(defaultPlanetId, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					assert.Equal(defaultPlanetResource.Amount, actual.Amount)

					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					expectedUpdatedAt := defaultPlanetResource.UpdatedAt.Add(3 * time.Minute)
					assert.Equal(expectedUpdatedAt, actual.UpdatedAt)
					assert.Equal(defaultPlanetResource.Version, actual.Version)
				},
			},
			"whenUpdatingPlanetResources_whenStorageIsNotEnoughToAbsorbAllProduction_expectPartialUpdate": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetResourceServiceMocks()

					resource := defaultPlanetResource
					resource.Amount = 60.0
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: resource,
					}

					production := defaultPlanetResourceProduction
					production.Production = 60
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						planetResourceProduction: production,
					}

					storage := defaultPlanetResourceStorage
					storage.Storage = 62
					repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
						planetResourceStorage: storage,
					}

					return repos
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlanetResourceService(conn, repos)
					threeMinutesAfterUpdatedAt := defaultPlanetResource.UpdatedAt.Add(3 * time.Minute)
					return s.UpdatePlanetUntil(ctx, defaultPlanetId, threeMinutesAfterUpdatedAt)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))
					actual := m.updatedPlanetResources[0]

					assert.Equal(defaultPlanetId, actual.Planet)
					assert.Equal(defaultPlanetResource.Resource, actual.Resource)
					assert.Equal(62.0, actual.Amount)

					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					expectedUpdatedAt := defaultPlanetResource.UpdatedAt.Add(3 * time.Minute)
					assert.Equal(expectedUpdatedAt, actual.UpdatedAt)
					assert.Equal(defaultPlanetResource.Version, actual.Version)
				},
			},
		},

		transactionTestCases: map[string]transactionTestCase{
			"updatePlanetUntil": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlanetResourceService(conn, repos)
					return s.UpdatePlanetUntil(ctx, defaultPlanetId, someTime)
				},
			},
		},

		transactionInteractionTestCases: map[string]transactionInteractionTestCase{
			"whenUpdatingPlanetResources_createsATransactionAndClosesIt": {
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlanetResourceService(conn, repos)
					return s.UpdatePlanetUntil(ctx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(conn db.Connection, assert *require.Assertions) {
					m := assertConnectionIsAMock(conn, assert)

					assert.Equal(1, len(m.txs))
					assert.Equal(1, m.txs[0].closeCalled)
				},
			},
			"whenUpdatingPlanetResources_whenFailureToCreateTransaction_expectPlanetResourcesNotUpdated": {
				generateConnectionMock: func() db.Connection {
					return &mockConnection{
						errs: []error{
							errDefault,
						},
					}
				},
				handler: func(ctx context.Context, conn db.Connection, repos repositories.Repositories) error {
					s := NewPlanetResourceService(conn, repos)
					return s.UpdatePlanetUntil(ctx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyMockInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(0, m.updateCalled)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generatePlanetResourceServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		PlanetResource: &mockPlanetResourceRepository{
			planetResource: defaultPlanetResource,
		},
		PlanetResourceProduction: &mockPlanetResourceProductionRepository{
			planetResourceProduction: defaultPlanetResourceProduction,
		},
		PlanetResourceStorage: &mockPlanetResourceStorageRepository{
			planetResourceStorage: defaultPlanetResourceStorage,
		},
	}
}

func assertConnectionIsAMock(conn db.Connection, assert *require.Assertions) *mockConnection {
	m, ok := conn.(*mockConnection)
	if !ok {
		assert.Fail("Provided connection is not a mock")
	}
	return m
}
