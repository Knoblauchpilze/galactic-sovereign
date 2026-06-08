package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	someTime = time.Date(2026, 6, 8, 8, 22, 35, 0, time.UTC)
)

func TestUnit_Building_CreateBuildingAction(t *testing.T) {
	buildingId := uuid.New()
	p := Planet{
		Id: uuid.New(),
		Buildings: []PlanetBuilding{
			{
				Building: buildingId,
				Level:    4,
			},
		},
	}

	t.Run("correctly calculates action costs", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs: []BuildingCost{
				{
					Resource: uuid.New(),
					Cost:     36,
					Progress: 1.5,
				},
				{
					Resource: uuid.New(),
					Cost:     78,
					Progress: 1.7,
				},
			},
			Productions: []BuildingResourceProduction{},
			Storages:    []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p, b.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := BuildingAction{
			// The identifier is generated
			Id:           action.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,

			CreatedAt: action.CreatedAt,
			// CompletedAt time.Time

			Version: 0,

			Costs: []BuildingActionCost{
				{
					Resource: b.Costs[0].Resource,
					Amount:   182,
				},
				{
					Resource: b.Costs[1].Resource,
					Amount:   651,
				},
			},
			Storages:    []BuildingActionResourceStorage{},
			Productions: []BuildingActionResourceProduction{},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("correctly calculates action resource productions", func(t *testing.T) {
		b := Building{
			Id:        buildingId,
			Name:      "test-building",
			CreatedAt: someTime,
			Costs:     []BuildingCost{},
			Productions: []BuildingResourceProduction{
				{
					Resource: uuid.New(),
					Base:     74,
					Progress: 4.5,
				},
				{
					Resource: uuid.New(),
					Base:     98,
					Progress: 2.4,
				},
			},
			Storages: []BuildingResourceStorage{},
		}

		action, err := b.CreateBuildingAction(p, b.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := BuildingAction{
			Id:           action.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,

			CreatedAt: action.CreatedAt,
			// CompletedAt time.Time

			Version: 0,

			Costs: []BuildingActionCost{},
			Productions: []BuildingActionResourceProduction{
				{
					Resource:   b.Productions[0].Resource,
					Production: 682754,
				},
				{
					Resource:   b.Productions[1].Resource,
					Production: 39016,
				},
			},
			Storages: []BuildingActionResourceStorage{},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("correctly calculates action resource storages", func(t *testing.T) {
		b := Building{
			Id:          buildingId,
			Name:        "test-building",
			CreatedAt:   someTime,
			Costs:       []BuildingCost{},
			Productions: []BuildingResourceProduction{},
			Storages: []BuildingResourceStorage{
				{
					Resource: uuid.New(),
					Base:     8904,
					Scale:    2,
					Progress: 2.2,
				},
				{
					Resource: uuid.New(),
					Base:     312,
					Scale:    1.2,
					Progress: 1.1,
				},
			},
		}

		action, err := b.CreateBuildingAction(p, b.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := BuildingAction{
			Id:           action.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,

			CreatedAt: action.CreatedAt,
			// CompletedAt time.Time

			Version: 0,

			Costs:       []BuildingActionCost{},
			Productions: []BuildingActionResourceProduction{},
			Storages: []BuildingActionResourceStorage{
				{
					Resource: b.Storages[0].Resource,
					Storage:  917112,
				},
				{
					Resource: b.Storages[1].Resource,
					Storage:  312,
				},
			},
		}
		assert.Equal(t, expected, action)
	})

	t.Run("returns error when building does not exist on planet", func(t *testing.T) {
		b := Building{Id: uuid.New()}

		_, err := b.CreateBuildingAction(p, b.Id)

		assert.ErrorIs(t, err, errBuildingNotFound)
	})
}

// func TestUnit_DetermineBuildingActionResourceProduction(t *testing.T) {
// 	assert := assert.New(t)

// 	action := persistence.BuildingAction{
// 		Id:           uuid.MustParse("7f548f48-2bac-46f0-b655-56487472b5db"),
// 		DesiredLevel: 10,
// 	}
// 	baseProductions := []persistence.BuildingResourceProduction{
// 		{
// 			Building: defaultId,
// 			Resource: defaultMetalId,
// 			Base:     21,
// 			Progress: 1.2,
// 		},
// 		{
// 			Building: defaultId,
// 			Resource: defaultCrystalId,
// 			Base:     27,
// 			Progress: 1.3,
// 		},
// 	}

// 	productions := DetermineBuildingActionResourceProduction(action, baseProductions)

// 	assert.Equal(2, len(productions))
// 	expectedResourceProduction := persistence.BuildingActionResourceProduction{
// 		Action:     action.Id,
// 		Resource:   defaultMetalId,
// 		Production: 1300,
// 	}
// 	assert.Equal(expectedResourceProduction, productions[0])
// 	expectedResourceProduction = persistence.BuildingActionResourceProduction{
// 		Action:     action.Id,
// 		Resource:   defaultCrystalId,
// 		Production: 3722,
// 	}
// 	assert.Equal(expectedResourceProduction, productions[1])
// }

// func TestUnit_DetermineBuildingActionResourceStorage(t *testing.T) {
// 	assert := assert.New(t)

// 	action := persistence.BuildingAction{
// 		Id:           uuid.MustParse("7f548f48-2bac-46f0-b655-56487472b5db"),
// 		DesiredLevel: 10,
// 	}
// 	baseStorages := []persistence.BuildingResourceStorage{
// 		{
// 			Building: defaultId,
// 			Resource: defaultMetalId,
// 			Base:     25,
// 			Scale:    527.78,
// 			Progress: 3.174,
// 		},
// 		{
// 			Building: defaultId,
// 			Resource: defaultCrystalId,
// 			Base:     1736,
// 			Scale:    1045.78,
// 			Progress: 1.995,
// 		},
// 	}

// 	storages := DetermineBuildingActionResourceStorage(action, baseStorages)

// 	assert.Equal(2, len(storages))
// 	expectedResourceStorage := persistence.BuildingActionResourceStorage{
// 		Action:   action.Id,
// 		Resource: defaultMetalId,
// 		Storage:  1_369_185_075,
// 	}
// 	assert.Equal(expectedResourceStorage, storages[0])
// 	expectedResourceStorage = persistence.BuildingActionResourceStorage{
// 		Action:   action.Id,
// 		Resource: defaultCrystalId,
// 		Storage:  1_813_087_080,
// 	}
// 	assert.Equal(expectedResourceStorage, storages[1])
// }

// func TestUnit_ConsolidateBuildingActionLevel_WhenNoBuilding_SetsDefault(t *testing.T) {
// 	assert := assert.New(t)

// 	action := persistence.BuildingAction{
// 		Building: defaultId,
// 	}
// 	buildings := []persistence.PlanetBuilding{}

// 	actual := ConsolidateBuildingActionLevel(action, buildings)

// 	assert.Equal(0, actual.CurrentLevel)
// 	assert.Equal(1, actual.DesiredLevel)
// }

// func TestUnit_ConsolidateBuildingActionLevel_WhenBuildingExists_SetsCorrectLevel(t *testing.T) {
// 	assert := assert.New(t)

// 	action := persistence.BuildingAction{
// 		Building: defaultId,
// 	}
// 	buildings := []persistence.PlanetBuilding{
// 		{
// 			Building: defaultId,
// 			Level:    26,
// 		},
// 	}

// 	actual := ConsolidateBuildingActionLevel(action, buildings)

// 	assert.Equal(26, actual.CurrentLevel)
// 	assert.Equal(27, actual.DesiredLevel)
// }

// func TestUnit_ConsolidateBuildingActionCompletionTime_WhenResourceNotFound_ExpectError(t *testing.T) {
// 	assert := assert.New(t)

// 	action := persistence.BuildingAction{
// 		Building: defaultId,
// 	}
// 	resources := []persistence.Resource{
// 		{
// 			Id:   defaultMetalId,
// 			Name: "metal",
// 		},
// 	}
// 	costs := []persistence.BuildingActionCost{
// 		{
// 			Resource: defaultMetalId,
// 			Amount:   1250,
// 		},
// 		{
// 			Resource: defaultCrystalId,
// 			Amount:   3750,
// 		},
// 	}

// 	_, err := ConsolidateBuildingActionCompletionTime(action, resources, costs)

// 	assert.True(errors.IsErrorWithCode(err, NoSuchResource))
// }

// func TestUnit_ConsolidateBuildingActionCompletionTime(t *testing.T) {
// 	assert := assert.New(t)

// 	action := persistence.BuildingAction{
// 		Building: defaultId,
// 	}
// 	costs := []persistence.BuildingActionCost{
// 		{
// 			Resource: defaultMetalId,
// 			Amount:   1250,
// 		},
// 		{
// 			Resource: defaultCrystalId,
// 			Amount:   3750,
// 		},
// 	}

// 	actual, err := ConsolidateBuildingActionCompletionTime(action, defaultResources, costs)

// 	assert.Nil(err)
// 	expectedCompletionTime := actual.CreatedAt.Add(2 * time.Hour)
// 	assert.Equal(expectedCompletionTime, actual.CompletedAt)
// }
