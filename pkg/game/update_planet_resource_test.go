package game

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var someTime = time.Date(2024, 9, 25, 20, 24, 18, 651387253, time.UTC)
var buildingId = uuid.MustParse("a8d8706e-53c7-480e-8cef-7a42ec963c0c")
var planetId = uuid.MustParse("ab6806b1-722b-4438-afa4-2b23a389d773")

var defaultPlanetResource = persistence.PlanetResource{
	Planet:   planetId,
	Resource: defaultMetalId,
	Amount:   60.0,

	UpdatedAt: time.Date(2024, 9, 13, 14, 51, 55, 651387248, time.UTC),

	Version: 10,
}

var metalProduction = persistence.PlanetResourceProduction{
	Planet:     planetId,
	Building:   &buildingId,
	Resource:   defaultMetalId,
	Production: 31,
	CreatedAt:  someTime,
	UpdatedAt:  someTime,
}
var crystalProduction = persistence.PlanetResourceProduction{
	Planet:     planetId,
	Resource:   defaultCrystalId,
	Building:   &buildingId,
	Production: 26,

	CreatedAt: someTime,
	UpdatedAt: someTime,

	Version: 7,
}

var metalStorage = persistence.PlanetResourceStorage{
	Planet:    planetId,
	Resource:  defaultMetalId,
	Storage:   89,
	CreatedAt: someTime,
	UpdatedAt: someTime,
}
var crystalStorage = persistence.PlanetResourceStorage{
	Planet:   planetId,
	Resource: defaultCrystalId,
	Storage:  78,

	CreatedAt: someTime,
	UpdatedAt: someTime,

	Version: 7,
}

const resourceProductionPerHour = 5.0
const resourceStorage = 1000
const thresholdForResourceEquality = 1e-6

func TestUnit_ToPlanetResourceProductionMap(t *testing.T) {
	assert := assert.New(t)

	in := []persistence.PlanetResourceProduction{metalProduction, crystalProduction}

	actual := toPlanetResourceProductionMap(in)

	assert.Equal(2, len(actual))

	metal, ok := actual[metalProduction.Resource]
	assert.True(ok)
	assert.Equal(metalProduction.Production, metal)

	crystal, ok := actual[crystalProduction.Resource]
	assert.True(ok)
	assert.Equal(crystalProduction.Production, crystal)
}

func TestUnit_ToPlanetResourceProductionMap_whenMultipleProductionsForResource_expectThemToBeAdded(t *testing.T) {
	assert := assert.New(t)

	metalProduction1 := metalProduction
	metalProduction2 := metalProduction
	metalProduction2.Production = 58

	in := []persistence.PlanetResourceProduction{metalProduction1, metalProduction2}

	actual := toPlanetResourceProductionMap(in)

	assert.Equal(1, len(actual))

	metal, ok := actual[metalProduction.Resource]
	assert.True(ok)
	expectedProduction := metalProduction1.Production + metalProduction2.Production
	assert.Equal(expectedProduction, metal)
}

func TestUnit_ToPlanetResourceStorageMap(t *testing.T) {
	assert := assert.New(t)

	in := []persistence.PlanetResourceStorage{metalStorage, crystalStorage}

	actual := toPlanetResourceStorageMap(in)

	assert.Equal(2, len(actual))

	metal, ok := actual[metalStorage.Resource]
	assert.True(ok)
	assert.Equal(metalStorage.Storage, metal)

	crystal, ok := actual[crystalStorage.Resource]
	assert.True(ok)
	assert.Equal(crystalStorage.Storage, crystal)
}

func TestUnit_ToPlanetResourceStorageMap_whenMultipleProductionsForResource_expectThemToBeAdded(t *testing.T) {
	assert := assert.New(t)

	metalStorage1 := metalStorage
	metalStorage2 := metalStorage
	metalStorage2.Storage = 58

	in := []persistence.PlanetResourceStorage{metalStorage1, metalStorage2}

	actual := toPlanetResourceStorageMap(in)

	assert.Equal(1, len(actual))

	metal, ok := actual[metalStorage.Resource]
	assert.True(ok)
	expectedStorage := metalStorage1.Storage + metalStorage2.Storage
	assert.Equal(expectedStorage, metal)
}

func TestUnit_UpdatePlanetResourceAmountToTime_whenTimeInThePast_expectNoUpdate(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()
	expectedUpdatedAt := resource.UpdatedAt

	inThePast := resource.UpdatedAt.Add(-1 * time.Hour)
	updated := updatePlanetResourceAmountToTime(resource, resourceProductionPerHour, resourceStorage, inThePast)

	assert.Equal(resource.Amount, updated.Amount)
	assert.Equal(expectedUpdatedAt, resource.UpdatedAt)
}

func TestUnit_UpdatePlanetResourceAmountToTime_updatesAmount(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	type testCase struct {
		duration       time.Duration
		expectedAmount float64
	}

	testCases := map[string]testCase{
		"1s": {
			duration:       1 * time.Second,
			expectedAmount: 60.001389,
		},
		"41s": {
			duration:       41 * time.Second,
			expectedAmount: 60.056944,
		},
		"1m": {
			duration:       1 * time.Minute,
			expectedAmount: 60.083333,
		},
		"30m": {
			duration:       30 * time.Minute,
			expectedAmount: 62.5,
		},
		"1h": {
			duration:       1 * time.Hour,
			expectedAmount: 65,
		},
		"4h3m2s2ms": {
			duration:       4*time.Hour + 3*time.Minute + 2*time.Second + 1*time.Millisecond,
			expectedAmount: 80.252779,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			oneHourFromNow := resource.UpdatedAt.Add(testCase.duration)
			updated := updatePlanetResourceAmountToTime(resource, resourceProductionPerHour, resourceStorage, oneHourFromNow)

			assert.InDelta(testCase.expectedAmount, updated.Amount, thresholdForResourceEquality)
		})
	}
}

func TestUnit_UpdatePlanetResourceAmountToTime_whenAmountIsAlreadyAboveStorageCapacity_expectNoChange(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	oneHourFromNow := resource.UpdatedAt.Add(1 * time.Hour)
	smallerStorageCapacityAsResourceAlreadyAvailable := resource.Amount - 10
	updated := updatePlanetResourceAmountToTime(resource,
		resourceProductionPerHour,
		smallerStorageCapacityAsResourceAlreadyAvailable,
		oneHourFromNow)

	assert.Equal(resource.Amount, updated.Amount)
	assert.Equal(oneHourFromNow, updated.UpdatedAt)
}

func TestUnit_UpdatePlanetResourceAmountToTime_whenProductionExceedsStorage_expectCapped(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	oneHourFromNow := resource.UpdatedAt.Add(1 * time.Hour)
	storageSufficientToAbsorbHalfAnHourOfProduction := resource.Amount + resourceProductionPerHour/2
	updated := updatePlanetResourceAmountToTime(resource,
		resourceProductionPerHour,
		storageSufficientToAbsorbHalfAnHourOfProduction,
		oneHourFromNow)

	expectedAmount := resource.Amount + resourceProductionPerHour/2
	assert.Equal(expectedAmount, updated.Amount)
}

func TestUnit_UpdatePlanetResourceAmountToTime_updatesUpdatedAt(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	oneHourFromNow := resource.UpdatedAt.Add(1 * time.Hour)
	updated := updatePlanetResourceAmountToTime(resource, resourceProductionPerHour, resourceStorage, oneHourFromNow)

	assert.Equal(oneHourFromNow, updated.UpdatedAt)
}

type verifyMockInteractions func(repositories.Repositories, *assert.Assertions)
type generateRepositoriesMocks func() repositories.Repositories

type planetResourceUpdateTestCase struct {
	until                     time.Time
	generateRepositoriesMocks generateRepositoriesMocks
	expectedError             error
	verifyMockInteractions    verifyMockInteractions
}

func TestUnit_UpdatePlanetResourcesToTime(t *testing.T) {
	tests := map[string]planetResourceUpdateTestCase{
		"whenUpdatingPlanetUntilTime_expectListResourcesForPlanetCalled": {
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.listForPlanetCalled)
				assert.Equal([]uuid.UUID{planetId}, m.listForPlanetIds)
			},
		},
		"whenUpdatingPlanetUntilTime_whenListResourcesForPlanetFails_expectError": {
			generateRepositoriesMocks: func() repositories.Repositories {
				repos := generateDefaultRepositoriesMocks()
				repos.PlanetResource = &mockPlanetResourceRepository{
					err: errDefault,
				}

				return repos
			},
			expectedError: errDefault,
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.listForPlanetCalled)
			},
		},
		"whenUpdatingPlanetUntilTime_expectListResourceProductionsForPlanetCalled": {
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

				assert.Equal(1, m.listForPlanetCalled)
				assert.Equal([]uuid.UUID{planetId}, m.listForPlanetIds)
			},
		},
		"whenUpdatingPlanetUntilTime_whenListResourceProductionsForPlanetFails_expectError": {
			generateRepositoriesMocks: func() repositories.Repositories {
				repos := generateDefaultRepositoriesMocks()
				repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
					errs: []error{errDefault},
				}

				return repos
			},
			expectedError: errDefault,
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

				assert.Equal(1, m.listForPlanetCalled)
			},
		},
		"whenUpdatingPlanetUntilTime_expectListResourceStoragesForPlanetCalled": {
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceStorageRepoIsAMock(repos, assert)

				assert.Equal(1, m.listForPlanetCalled)
				assert.Equal([]uuid.UUID{planetId}, m.listForPlanetIds)
			},
		},
		"whenUpdatingPlanetUntilTime_whenListResourceStoragesForPlanetFails_expectError": {
			generateRepositoriesMocks: func() repositories.Repositories {
				repos := generateDefaultRepositoriesMocks()
				repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
					errs: []error{errDefault},
				}

				return repos
			},
			expectedError: errDefault,
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceStorageRepoIsAMock(repos, assert)

				assert.Equal(1, m.listForPlanetCalled)
			},
		},
		"whenUpdatingPlanetUntilTime_expectResourceAreUpdatedWithCorrectValue": {
			until: defaultPlanetResource.UpdatedAt.Add(2 * time.Minute),
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.updateCalled)
				assert.Equal(1, len(m.updatedPlanetResources))

				actual := m.updatedPlanetResources[0]
				assert.Equal(planetId, actual.Planet)
				assert.Equal(defaultMetalId, actual.Resource)
				expectedAmount := defaultPlanetResource.Amount + 2.0/60.0*float64(metalProduction.Production)
				assert.Equal(expectedAmount, actual.Amount)
				assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
				expectedUpdatedAt := defaultPlanetResource.UpdatedAt.Add(2 * time.Minute)
				assert.Equal(expectedUpdatedAt, actual.UpdatedAt)
				assert.Equal(defaultPlanetResource.Version, actual.Version)
			},
		},
		"whenUpdatingPlanetUntilTime_whenResourceAlreadyFillingStorage_expectNoUpdate": {
			until: defaultPlanetResource.UpdatedAt.Add(2 * time.Minute),
			generateRepositoriesMocks: func() repositories.Repositories {
				repos := generateDefaultRepositoriesMocks()

				resource := defaultPlanetResource
				resource.Amount = 60
				repos.PlanetResource = &mockPlanetResourceRepository{
					planetResource: resource,
				}

				storage := metalStorage
				storage.Storage = 60
				repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
					planetResourceStorage: storage,
				}

				return repos
			},
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.updateCalled)
				assert.Equal(1, len(m.updatedPlanetResources))

				actual := m.updatedPlanetResources[0]
				assert.Equal(planetId, actual.Planet)
				assert.Equal(defaultMetalId, actual.Resource)
				assert.Equal(60.0, actual.Amount)
				assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
				expectedUpdatedAt := defaultPlanetResource.UpdatedAt.Add(2 * time.Minute)
				assert.Equal(expectedUpdatedAt, actual.UpdatedAt)
				assert.Equal(defaultPlanetResource.Version, actual.Version)
			},
		},
		"whenUpdatingPlanetUntilTime_whenStorageNotBigEnoughForWholeProduction_expectPartialUpdate": {
			until: defaultPlanetResource.UpdatedAt.Add(2 * time.Minute),
			generateRepositoriesMocks: func() repositories.Repositories {
				repos := generateDefaultRepositoriesMocks()

				resource := defaultPlanetResource
				resource.Amount = 60
				repos.PlanetResource = &mockPlanetResourceRepository{
					planetResource: resource,
				}

				production := metalProduction
				production.Production = 60
				repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
					planetResourceProduction: production,
				}

				storage := metalStorage
				storage.Storage = 61
				repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
					planetResourceStorage: storage,
				}

				return repos
			},
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.updateCalled)
				assert.Equal(1, len(m.updatedPlanetResources))

				actual := m.updatedPlanetResources[0]
				assert.Equal(planetId, actual.Planet)
				assert.Equal(defaultMetalId, actual.Resource)
				assert.Equal(61.0, actual.Amount)
				assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
				expectedUpdatedAt := defaultPlanetResource.UpdatedAt.Add(2 * time.Minute)
				assert.Equal(expectedUpdatedAt, actual.UpdatedAt)
				assert.Equal(defaultPlanetResource.Version, actual.Version)
			},
		},
		"whenUpdatingPlanetUntilTime_whenStorageInformationNotAvailable_expectNoProductionPossible": {
			until: defaultPlanetResource.UpdatedAt.Add(2 * time.Minute),
			generateRepositoriesMocks: func() repositories.Repositories {
				repos := generateDefaultRepositoriesMocks()

				resource := defaultPlanetResource
				resource.Amount = 60
				repos.PlanetResource = &mockPlanetResourceRepository{
					planetResource: resource,
				}

				production := metalProduction
				production.Production = 1
				repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
					planetResourceProduction: production,
				}

				repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
					planetResourceStorage: crystalStorage,
				}

				return repos
			},
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.updateCalled)
				assert.Equal(1, len(m.updatedPlanetResources))

				actual := m.updatedPlanetResources[0]
				assert.Equal(planetId, actual.Planet)
				assert.Equal(defaultMetalId, actual.Resource)
				assert.Equal(60.0, actual.Amount)
				assert.Equal(defaultPlanetResource.CreatedAt, actual.CreatedAt)
				expectedUpdatedAt := defaultPlanetResource.UpdatedAt.Add(2 * time.Minute)
				assert.Equal(expectedUpdatedAt, actual.UpdatedAt)
				assert.Equal(defaultPlanetResource.Version, actual.Version)
			},
		},
		"whenUpdatingPlanetUntilTime_whenUpdateOfResourceFails_expectError": {
			generateRepositoriesMocks: func() repositories.Repositories {
				repos := generateDefaultRepositoriesMocks()
				repos.PlanetResource = &mockPlanetResourceRepository{
					planetResource: defaultPlanetResource,
					updateErr:      errDefault,
				}

				return repos
			},
			expectedError: errDefault,
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.updateCalled)
			},
		},
		"whenUpdatingPlanetUntilTime_whenResourceIsNotProduced_expectNoUpdate": {
			generateRepositoriesMocks: func() repositories.Repositories {
				repos := generateDefaultRepositoriesMocks()

				planetResource := defaultPlanetResource
				planetResource.Resource = defaultCrystalId

				repos.PlanetResource = &mockPlanetResourceRepository{
					planetResource: planetResource,
					updateErr:      errDefault,
				}

				return repos
			},
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(0, m.updateCalled)
			},
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {

			tx := &mockTransaction{}

			var repos repositories.Repositories
			if testCase.generateRepositoriesMocks != nil {
				repos = testCase.generateRepositoriesMocks()
			} else {
				repos = generateDefaultRepositoriesMocks()
			}

			data := PlanetResourceUpdateData{
				Planet:                       planetId,
				Until:                        testCase.until,
				PlanetResourceRepo:           repos.PlanetResource,
				PlanetResourceProductionRepo: repos.PlanetResourceProduction,
				PlanetResourceStorageRepo:    repos.PlanetResourceStorage,
			}

			err := UpdatePlanetResourcesToTime(context.Background(), tx, data)

			assert := assert.New(t)
			assert.Equal(testCase.expectedError, err)

			if testCase.verifyMockInteractions != nil {
				testCase.verifyMockInteractions(repos, assert)
			}
		})
	}
}

func generatePlanetResource() persistence.PlanetResource {
	resource := defaultPlanetResource
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()

	return resource
}

func generateDefaultRepositoriesMocks() repositories.Repositories {
	return repositories.Repositories{
		PlanetResource: &mockPlanetResourceRepository{
			planetResource: defaultPlanetResource,
		},
		PlanetResourceProduction: &mockPlanetResourceProductionRepository{
			planetResourceProduction: metalProduction,
		},
		PlanetResourceStorage: &mockPlanetResourceStorageRepository{
			planetResourceStorage: metalStorage,
		},
	}
}

func assertPlanetResourceRepoIsAMock(repos repositories.Repositories, assert *assert.Assertions) *mockPlanetResourceRepository {
	m, ok := repos.PlanetResource.(*mockPlanetResourceRepository)
	if !ok {
		assert.Fail("Provided planet resource repository is not a mock")
	}
	return m
}

func assertPlanetResourceProductionRepoIsAMock(repos repositories.Repositories, assert *assert.Assertions) *mockPlanetResourceProductionRepository {
	m, ok := repos.PlanetResourceProduction.(*mockPlanetResourceProductionRepository)
	if !ok {
		assert.Fail("Provided planet resource production repository is not a mock")
	}
	return m
}

func assertPlanetResourceStorageRepoIsAMock(repos repositories.Repositories, assert *assert.Assertions) *mockPlanetResourceStorageRepository {
	m, ok := repos.PlanetResourceStorage.(*mockPlanetResourceStorageRepository)
	if !ok {
		assert.Fail("Provided planet resource storage repository is not a mock")
	}
	return m
}
