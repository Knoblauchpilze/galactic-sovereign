package domainservices

import (
	"slices"
	"testing"
	"time"

	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	metalResourceId   = uuid.MustParse("b4419b6b-b3bf-4576-aa92-055283addbc8")
	crystalResourceId = uuid.MustParse("cd2ac9aa-9968-4ff5-b746-88f1f810fbb3")
	crystalMineId     = uuid.MustParse("3904d34d-9a7e-47d4-a332-091700e2c5c3")
	metalStorageId    = uuid.MustParse("22b4c0c3-c8e5-4493-89fc-522fdbb0beee")

	lightFighterId = uuid.MustParse("a31de13b-5905-4468-99c5-d1d1e529b36e")
	smallCargoId   = uuid.MustParse("c0978950-601e-4d35-9c7c-28df69d2cd0e")

	t1 = time.Date(2026, time.July, 3, 6, 32, 27, 0, time.UTC)
	t2 = time.Date(2026, time.July, 3, 7, 32, 27, 0, time.UTC)
	t3 = time.Date(2026, time.July, 3, 8, 32, 27, 0, time.UTC)
	t4 = time.Date(2026, time.July, 3, 9, 32, 27, 0, time.UTC)
	t5 = time.Date(2026, time.July, 3, 10, 32, 27, 0, time.UTC)
)

func TestUnit_AdvancePlanetToTime(t *testing.T) {
	t.Run("updates planet to time when no action is running", func(t *testing.T) {
		p := generateTestPlanet()

		initialStorages := slices.Clone(p.Storages)
		initialProductions := slices.Clone(p.Productions)
		initialBuildings := slices.Clone(p.Buildings)
		initialShips := slices.Clone(p.Ships)

		err := AdvancePlanetToTime(&p, t4)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: t4,
			Version:   4,
			Resources: []models.PlanetResource{
				{Resource: metalResourceId, Amount: 1195},
				{Resource: crystalResourceId, Amount: 2120},
			},
			Storages:    initialStorages,
			Productions: initialProductions,
			Buildings:   initialBuildings,
			Ships:       initialShips,
		}
		assert.Equal(t, expected, p)
	})

	t.Run("updates planet to time when building action finishes after requested time", func(t *testing.T) {
		p := generateTestPlanet()
		action := generateTestBuildingAction()
		p.BuildingAction = &action

		initialStorages := slices.Clone(p.Storages)
		initialProductions := slices.Clone(p.Productions)
		initialBuildings := slices.Clone(p.Buildings)
		initialShips := slices.Clone(p.Ships)

		err := AdvancePlanetToTime(&p, t2)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: t2,
			Version:   4,
			Resources: []models.PlanetResource{
				{Resource: metalResourceId, Amount: 1065},
				{Resource: crystalResourceId, Amount: 2040},
			},
			Storages:       initialStorages,
			Productions:    initialProductions,
			Buildings:      initialBuildings,
			Ships:          initialShips,
			BuildingAction: &action,
		}
		assert.Equal(t, expected, p)
	})

	t.Run("applies building action when it finishes before the requested time", func(t *testing.T) {
		p := generateTestPlanet()
		action := generateTestBuildingAction()
		p.BuildingAction = &action

		err := AdvancePlanetToTime(&p, t4)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: t4,
			Version:   6,
			Resources: []models.PlanetResource{
				{Resource: metalResourceId, Amount: 1195},
				{Resource: crystalResourceId, Amount: 3328},
			},
			Storages: []models.PlanetResourceStorage{
				{Resource: metalResourceId, Storage: 15874},
				{Resource: crystalResourceId, Storage: 78941},
			},
			Productions: []models.PlanetResourceProduction{
				{Resource: crystalResourceId, Production: 14},
				{Resource: metalResourceId, Building: &crystalMineId, Production: 65},
				{Resource: crystalResourceId, Building: &crystalMineId, Production: 1234},
			},
			Buildings: []models.PlanetBuilding{
				{Building: crystalMineId, Level: action.DesiredLevel},
				{Building: metalStorageId, Level: 4},
			},
			Ships: []models.PlanetShip{
				{Ship: lightFighterId, Count: 1},
				{Ship: smallCargoId, Count: 0},
			},
			BuildingAction: nil,
		}
		assert.Equal(t, expected, p)
	})

	t.Run("updates planet to time when ship action finishes after requested time", func(t *testing.T) {
		p := generateTestPlanet()
		action := generateTestShipAction(2)
		p.ShipActions = append(p.ShipActions, action)

		initialStorages := slices.Clone(p.Storages)
		initialProductions := slices.Clone(p.Productions)
		initialBuildings := slices.Clone(p.Buildings)
		initialShips := slices.Clone(p.Ships)

		err := AdvancePlanetToTime(&p, t2)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: t2,
			Version:   4,
			Resources: []models.PlanetResource{
				{Resource: metalResourceId, Amount: 1065},
				{Resource: crystalResourceId, Amount: 2040},
			},
			Storages:    initialStorages,
			Productions: initialProductions,
			Buildings:   initialBuildings,
			Ships:       initialShips,
			ShipActions: []models.ShipAction{action},
		}
		assert.Equal(t, expected, p)
	})

	t.Run("applies ship action when it finishes before the requested time", func(t *testing.T) {
		p := generateTestPlanet()
		action := generateTestShipAction(2)
		p.ShipActions = append(p.ShipActions, action)

		err := AdvancePlanetToTime(&p, t4)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: t4,
			Version:   6,
			Resources: []models.PlanetResource{
				{Resource: metalResourceId, Amount: 1195},
				{Resource: crystalResourceId, Amount: 2120},
			},
			Storages: []models.PlanetResourceStorage{
				{Resource: metalResourceId, Storage: 15874},
				{Resource: crystalResourceId, Storage: 3541},
			},
			Productions: []models.PlanetResourceProduction{
				{Resource: crystalResourceId, Production: 14},
				{Resource: metalResourceId, Building: &crystalMineId, Production: 65},
				{Resource: crystalResourceId, Building: &crystalMineId, Production: 26},
			},
			Buildings: []models.PlanetBuilding{
				{Building: crystalMineId, Level: 2},
				{Building: metalStorageId, Level: 4},
			},
			Ships: []models.PlanetShip{
				{Ship: lightFighterId, Count: 2},
				{Ship: smallCargoId, Count: 0},
			},
			BuildingAction: nil,
			ShipActions: []models.ShipAction{
				{
					Id:                 action.Id,
					Ship:               action.Ship,
					Count:              action.Count - 1,
					CreatedAt:          t3,
					NextCompletionAt:   t5,
					UnitCompletionTime: t4.Sub(t3),
				},
			},
		}
		assert.Equal(t, expected, p)
	})

	t.Run("removes ship action when it is completed before the requested time", func(t *testing.T) {
		p := generateTestPlanet()
		action := generateTestShipAction(1)
		p.ShipActions = append(p.ShipActions, action)

		err := AdvancePlanetToTime(&p, t5)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: t5,
			Version:   6,
			Resources: []models.PlanetResource{
				{Resource: metalResourceId, Amount: 1260},
				{Resource: crystalResourceId, Amount: 2160},
			},
			Storages: []models.PlanetResourceStorage{
				{Resource: metalResourceId, Storage: 15874},
				{Resource: crystalResourceId, Storage: 3541},
			},
			Productions: []models.PlanetResourceProduction{
				{Resource: crystalResourceId, Production: 14},
				{Resource: metalResourceId, Building: &crystalMineId, Production: 65},
				{Resource: crystalResourceId, Building: &crystalMineId, Production: 26},
			},
			Buildings: []models.PlanetBuilding{
				{Building: crystalMineId, Level: 2},
				{Building: metalStorageId, Level: 4},
			},
			Ships: []models.PlanetShip{
				{Ship: lightFighterId, Count: 2},
				{Ship: smallCargoId, Count: 0},
			},
			BuildingAction: nil,
			ShipActions:    []models.ShipAction{},
		}
		assert.Equal(t, expected, p)
	})

	t.Run("leaves uncompleted ship action running when one is removed", func(t *testing.T) {
		p := generateTestPlanet()
		action1 := models.ShipAction{
			Id:                 uuid.New(),
			Ship:               smallCargoId,
			Count:              1,
			CreatedAt:          t2,
			NextCompletionAt:   t3,
			UnitCompletionTime: t3.Sub(t2),
		}
		action2 := models.ShipAction{
			Id:                 uuid.New(),
			Ship:               lightFighterId,
			Count:              1,
			CreatedAt:          t4.Add(1 * time.Minute),
			NextCompletionAt:   t4.Add(2 * time.Minute),
			UnitCompletionTime: 1 * time.Minute,
		}
		p.ShipActions = append(p.ShipActions, action1, action2)

		err := AdvancePlanetToTime(&p, t4)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Planet{
			Id:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: t4,
			Version:   6,
			Resources: []models.PlanetResource{
				{Resource: metalResourceId, Amount: 1195},
				{Resource: crystalResourceId, Amount: 2120},
			},
			Storages: []models.PlanetResourceStorage{
				{Resource: metalResourceId, Storage: 15874},
				{Resource: crystalResourceId, Storage: 3541},
			},
			Productions: []models.PlanetResourceProduction{
				{Resource: crystalResourceId, Production: 14},
				{Resource: metalResourceId, Building: &crystalMineId, Production: 65},
				{Resource: crystalResourceId, Building: &crystalMineId, Production: 26},
			},
			Buildings: []models.PlanetBuilding{
				{Building: crystalMineId, Level: 2},
				{Building: metalStorageId, Level: 4},
			},
			Ships: []models.PlanetShip{
				{Ship: lightFighterId, Count: 1},
				{Ship: smallCargoId, Count: 1},
			},
			BuildingAction: nil,
			ShipActions:    []models.ShipAction{action2},
		}
		assert.Equal(t, expected, p)
	})
}

func generateTestPlanet() models.Planet {
	return models.Planet{
		Id:        uuid.New(),
		CreatedAt: t1,
		UpdatedAt: t1,
		Version:   3,
		Resources: []models.PlanetResource{
			{Resource: metalResourceId, Amount: 1000.0},
			{Resource: crystalResourceId, Amount: 2000.0},
		},
		Storages: []models.PlanetResourceStorage{
			{Resource: metalResourceId, Storage: 15874},
			{Resource: crystalResourceId, Storage: 3541},
		},
		Productions: []models.PlanetResourceProduction{
			{Resource: crystalResourceId, Production: 14},
			{Resource: metalResourceId, Building: &crystalMineId, Production: 65},
			{Resource: crystalResourceId, Building: &crystalMineId, Production: 26},
		},
		Buildings: []models.PlanetBuilding{
			{Building: crystalMineId, Level: 2},
			{Building: metalStorageId, Level: 4},
		},
		Ships: []models.PlanetShip{
			{Ship: lightFighterId, Count: 1},
			{Ship: smallCargoId, Count: 0},
		},
	}
}

func generateTestBuildingAction() models.BuildingAction {
	return models.BuildingAction{
		Id:           uuid.New(),
		Building:     crystalMineId,
		DesiredLevel: 3,
		Storages: []models.BuildingActionResourceStorage{
			{Resource: crystalResourceId, Storage: 78941},
		},
		Productions: []models.BuildingActionResourceProduction{
			{Resource: crystalResourceId, Production: 1234},
		},
		CompletedAt: t3,
	}
}

func generateTestShipAction(count int) models.ShipAction {
	return models.ShipAction{
		Id:                 uuid.New(),
		Ship:               lightFighterId,
		Count:              count,
		CreatedAt:          t3,
		NextCompletionAt:   t4,
		UnitCompletionTime: t4.Sub(t3),
	}
}
