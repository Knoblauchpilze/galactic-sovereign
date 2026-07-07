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

	t1 = time.Date(2026, time.July, 3, 6, 32, 27, 0, time.UTC)
	t2 = time.Date(2026, time.July, 3, 7, 32, 27, 0, time.UTC)
	t3 = time.Date(2026, time.July, 3, 8, 32, 27, 0, time.UTC)
	t4 = time.Date(2026, time.July, 3, 9, 32, 27, 0, time.UTC)
)

func TestUnit_AdvancePlanetToTime(t *testing.T) {
	t.Run("updates planet to time when no building action is running", func(t *testing.T) {
		p := generateTestPlanet()

		initialBuildings := slices.Clone(p.Buildings)
		initialStorages := slices.Clone(p.Storages)
		initialProductions := slices.Clone(p.Productions)

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
		}
		assert.Equal(t, expected, p)
	})

	t.Run("updates planet to time when building action finishes after requested time", func(t *testing.T) {
		p := generateTestPlanet()
		action := generateTestBuildingAction(p)
		p.BuildingAction = &action

		initialBuildings := slices.Clone(p.Buildings)
		initialStorages := slices.Clone(p.Storages)
		initialProductions := slices.Clone(p.Productions)

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
			BuildingAction: &action,
		}
		assert.Equal(t, expected, p)
	})

	t.Run("applies building action when it finishes before the requested time", func(t *testing.T) {
		p := generateTestPlanet()
		action := generateTestBuildingAction(p)
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
			BuildingAction: nil,
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
	}
}

func generateTestBuildingAction(p models.Planet) models.BuildingAction {
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
