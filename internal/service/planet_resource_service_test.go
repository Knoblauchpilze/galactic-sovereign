package service

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test_PlanetResourceService(t *testing.T) {
	s := ServiceTransactionTestSuite{
		generateRepositoriesMock: generateValidPlanetResourceServiceMocks,

		interactionTestCases: map[string]serviceTransactionInteractionTestCase{
			"whenUpdatingPlanetUntilTime_expectListResourcesForPlanetCalled": {
				handler: func(ctx context.Context, tx db.Transaction, repos repositories.Repositories) error {
					s := NewPlanetResourceService(repos)
					return s.UpdatePlanetUntil(ctx, tx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal([]uuid.UUID{defaultPlanetId}, m.listForPlanetIds)
				},
			},
			"whenUpdatingPlanetUntilTime_whenListResourcesForPlanetFails_expectError": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidPlanetResourceServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, tx db.Transaction, repos repositories.Repositories) error {
					s := NewPlanetResourceService(repos)
					return s.UpdatePlanetUntil(ctx, tx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"whenUpdatingPlanetUntilTime_expectListResourceProductionsForPlanetCalled": {
				handler: func(ctx context.Context, tx db.Transaction, repos repositories.Repositories) error {
					s := NewPlanetResourceService(repos)
					return s.UpdatePlanetUntil(ctx, tx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal([]uuid.UUID{defaultPlanetId}, m.listForPlanetIds)
				},
			},
			"whenUpdatingPlanetUntilTime_whenListResourceProductionsForPlanetFails_expectError": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidPlanetResourceServiceMocks()
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, tx db.Transaction, repos repositories.Repositories) error {
					s := NewPlanetResourceService(repos)
					return s.UpdatePlanetUntil(ctx, tx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
				},
			},
			"whenUpdatingPlanetUntilTime_expectResourceAreUpdatedWithCorrectValue": {
				generateTransactionMock: func() db.Transaction {
					twoMinutesAfterUpdatedAt := defaultPlanetResource.UpdatedAt.Add(2 * time.Minute)

					return &mockTransaction{
						timeStamp: twoMinutesAfterUpdatedAt,
					}
				},
				handler: func(ctx context.Context, tx db.Transaction, repos repositories.Repositories) error {
					s := NewPlanetResourceService(repos)
					return s.UpdatePlanetUntil(ctx, tx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
					assert.Equal(1, len(m.updatedPlanetResources))

					actual := m.updatedPlanetResources[0]
					assert.Equal(defaultPlanetId, actual.Planet)
					assert.Equal(metalResourceId, actual.Resource)
					expectedAmount := defaultPlanetResource.Amount + 2.0/60.0*float64(defaultPlanetResourceProduction.Production)
					assert.Equal(expectedAmount, actual.Amount)
					assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
					expectedUpdatedAt := defaultPlanetResource.UpdatedAt.Add(2 * time.Minute)
					assert.Equal(expectedUpdatedAt, actual.UpdatedAt)
					assert.Equal(defaultPlanetResource.Version, actual.Version)
				},
			},
			"whenUpdatingPlanetUntilTime_whenUpdateOfResourceFails_expectError": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidPlanetResourceServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: defaultPlanetResource,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, tx db.Transaction, repos repositories.Repositories) error {
					s := NewPlanetResourceService(repos)
					return s.UpdatePlanetUntil(ctx, tx, defaultPlanetId, someTime)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.updateCalled)
				},
			},
			"whenUpdatingPlanetUntilTime_whenResourceIsNotProduced_expectNoUpdate": {
				generateRepositoriesMock: func() repositories.Repositories {
					repos := generateValidPlanetResourceServiceMocks()

					planetResource := defaultPlanetResource
					planetResource.Resource = crystalResourceId

					repos.PlanetResource = &mockPlanetResourceRepository{
						planetResource: planetResource,
						updateErr:      errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, tx db.Transaction, repos repositories.Repositories) error {
					s := NewPlanetResourceService(repos)
					return s.UpdatePlanetUntil(ctx, tx, defaultPlanetId, someTime)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(0, m.updateCalled)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generateValidPlanetResourceServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		PlanetResource: &mockPlanetResourceRepository{
			planetResource: defaultPlanetResource,
		},
		PlanetResourceProduction: &mockPlanetResourceProductionRepository{
			planetResourceProduction: defaultPlanetResourceProduction,
		},
	}
}
