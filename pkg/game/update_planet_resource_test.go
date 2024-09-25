package game

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
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

var defaultPlanetResourceProduction = persistence.PlanetResourceProduction{
	Planet:     planetId,
	Building:   &buildingId,
	Resource:   defaultMetalId,
	Production: 31,
	CreatedAt:  someTime,
	UpdatedAt:  someTime,
}
var resourceProductionPerHour = 5.0

var thresholdForResourceEquality = 1e-6

func TestUpdatePlanetResourceAmountToTime_whenTimeInThePast_expectNoUpdate(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	inThePast := resource.UpdatedAt.Add(-1 * time.Hour)
	updated := updatePlanetResourceAmountToTime(resource, resourceProductionPerHour, inThePast)

	assert.Equal(resource.Amount, updated.Amount)
}

func TestUpdatePlanetResourceAmountToTime_updatesAmount(t *testing.T) {
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
			updated := updatePlanetResourceAmountToTime(resource, resourceProductionPerHour, oneHourFromNow)

			assert.InDelta(testCase.expectedAmount, updated.Amount, thresholdForResourceEquality)
		})
	}
}

func TestUpdatePlanetResourceAmountToTime_updatesUpdatedAt(t *testing.T) {
	assert := assert.New(t)

	resource := generatePlanetResource()

	oneHourFromNow := resource.UpdatedAt.Add(1 * time.Hour)
	updated := updatePlanetResourceAmountToTime(resource, resourceProductionPerHour, oneHourFromNow)

	assert.Equal(oneHourFromNow, updated.UpdatedAt)
}

type verifyMockInteractions func(repositories.Repositories, *assert.Assertions)
type generateRepositoriesMock func() repositories.Repositories

type planetResourceUpdateTestCase struct {
	until                    time.Time
	generateRepositoriesMock generateRepositoriesMock
	expectedError            error
	verifyMockInteractions   verifyMockInteractions
}

func Test_PlanetResourceService(t *testing.T) {
	tests := map[string]planetResourceUpdateTestCase{
		"whenUpdatingPlanetUntilTime_expectListResourcesForPlanetCalled": {
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.listForPlanetCalled)
				assert.Equal([]uuid.UUID{planetId}, m.listForPlanetIds)
			},
		},
		"whenUpdatingPlanetUntilTime_whenListResourcesForPlanetFails_expectError": {
			generateRepositoriesMock: func() repositories.Repositories {
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
			generateRepositoriesMock: func() repositories.Repositories {
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
		"whenUpdatingPlanetUntilTime_expectResourceAreUpdatedWithCorrectValue": {
			until: defaultPlanetResource.UpdatedAt.Add(2 * time.Minute),
			verifyMockInteractions: func(repos repositories.Repositories, assert *assert.Assertions) {
				m := assertPlanetResourceRepoIsAMock(repos, assert)

				assert.Equal(1, m.updateCalled)
				assert.Equal(1, len(m.updatedPlanetResources))

				actual := m.updatedPlanetResources[0]
				assert.Equal(planetId, actual.Planet)
				assert.Equal(defaultMetalId, actual.Resource)
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
			generateRepositoriesMock: func() repositories.Repositories {
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
			if testCase.generateRepositoriesMock != nil {
				repos = testCase.generateRepositoriesMock()
			} else {
				repos = generateDefaultRepositoriesMocks()
			}

			data := PlanetResourceUpdateData{
				Planet:                       planetId,
				Until:                        testCase.until,
				PlanetResourceRepo:           repos.PlanetResource,
				PlanetResourceProductionRepo: repos.PlanetResourceProduction,
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
			planetResourceProduction: defaultPlanetResourceProduction,
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
