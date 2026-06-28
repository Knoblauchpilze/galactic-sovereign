package models

import (
	"slices"
	"testing"
	"time"

	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	crystalMineId = uuid.MustParse("3904d34d-9a7e-47d4-a332-091700e2c5c3")
	metalMineId   = uuid.MustParse("d176e82d-f2ca-4611-996b-c4804096caef")
)

func TestUnit_Planet_AddBuildingAction(t *testing.T) {
	t.Run("returns error when planet already has an action", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		actionId := uuid.New()
		p.BuildingAction = &BuildingAction{Id: actionId}

		b := generateTestBuilding(t)

		err := p.AddBuildingAction(b)

		assert.ErrorIs(t, domainerrors.ErrActionAlreadyInProgress, err, "Actual err: %v", err)
		require.NotNil(t, p.BuildingAction)
		assert.Equal(t, actionId, p.BuildingAction.Id)
		assert.Equal(t, 3, p.Version)
	})

	t.Run("returns error when planet does not have enough resources", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)
		p.Resources = []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   189,
			},
			{
				Resource: crystalResourceId,
				// Needed value: 651
				Amount: 650,
			},
		}

		b := generateTestBuilding(t, withBuildingCost)

		err := p.AddBuildingAction(b)

		assert.ErrorIs(t, err, domainerrors.ErrNotEnoughResources, "Actual err: %v", err)
		assert.Nil(t, p.BuildingAction)
		assert.Equal(t, 3, p.Version)
	})

	t.Run("returns error when building does not exist on planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)

		b := Building{Id: uuid.New()}

		err := p.AddBuildingAction(b)

		assert.ErrorIs(t, domainerrors.ErrBuildingNotFound, err, "Actual err: %v", err)
		assert.Nil(t, p.BuildingAction)
		assert.Equal(t, 3, p.Version)
	})

	t.Run("assigns building action to planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, p.BuildingAction)

		completionTime := 1199520 * time.Millisecond
		expectedAction := &BuildingAction{
			Id:           p.BuildingAction.Id,
			Planet:       p.Id,
			Building:     b.Id,
			CurrentLevel: p.Buildings[0].Level,
			DesiredLevel: p.Buildings[0].Level + 1,
			CreatedAt:    someTime,
			CompletedAt:  someTime.Add(completionTime),
			Costs: []BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   182,
				},
				{
					Resource: crystalResourceId,
					Amount:   651,
				},
			},
			Storages: []BuildingActionResourceStorage{
				{
					Resource: metalResourceId,
					Storage:  917112,
				},
				{
					Resource: crystalResourceId,
					Storage:  312,
				},
			},
			Productions: []BuildingActionResourceProduction{
				{
					Resource:   metalResourceId,
					Production: 682754,
				},
				{
					Resource:   crystalResourceId,
					Production: 39016,
				},
			},
		}
		assert.Equal(t, expectedAction, p.BuildingAction)
	})

	t.Run("deducts action costs from the available planet resources", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		initialResources := slices.Clone(p.Resources)

		err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, p.BuildingAction)

		expectedMetalAmount := initialResources[0].Amount - float64(p.BuildingAction.Costs[0].Amount)
		expectedCrystalAmount := initialResources[1].Amount - float64(p.BuildingAction.Costs[1].Amount)
		expectedResources := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   expectedMetalAmount,
			},
			{
				Resource: crystalResourceId,
				Amount:   expectedCrystalAmount,
			},
		}
		assert.Equal(t, expectedResources, p.Resources)
	})

	t.Run("bumps version by one", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.UpdatedAt = someTime
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		initialVersion := p.Version

		err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, p.BuildingAction)

		assert.Equal(t, initialVersion+1, p.Version)
	})

	t.Run("does not bump updated at field", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.UpdatedAt = someTime
		b := generateTestBuilding(t, withBuildingCost, withBuildingProduction, withBuildingStorage)

		err := p.AddBuildingAction(b)
		require.NoError(t, err, "Actual err: %v", err)
		require.NotNil(t, p.BuildingAction)

		assert.Equal(t, someTime, p.UpdatedAt)
	})
}

func TestUnit_Planet_CancelBuildingAction(t *testing.T) {
	t.Run("returns error when planet does not have an action", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)

		err := p.CancelBuildingAction()

		assert.ErrorIs(t, domainerrors.ErrNoActionInProgress, err, "Actual err: %v", err)
	})

	t.Run("resets building action in planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.BuildingAction = &BuildingAction{Id: uuid.New()}

		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, p.BuildingAction)
	})

	t.Run("adds back action costs to the available planet resources", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)
		p.Resources = []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   1000,
			},
			{
				Resource: crystalResourceId,
				Amount:   2000,
			},
		}

		p.BuildingAction = &BuildingAction{
			Id: uuid.New(),
			Costs: []BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   36,
				},
				{
					Resource: crystalResourceId,
					Amount:   178,
				},
			},
		}

		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   1036,
			},
			{
				Resource: crystalResourceId,
				Amount:   2178,
			},
		}
		assert.Equal(t, expected, p.Resources)
	})

	t.Run("adds back action costs when resources are not present on the planet", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding)
		p.Resources = []PlanetResource{}
		p.BuildingAction = &BuildingAction{
			Id: uuid.New(),
			Costs: []BuildingActionCost{
				{
					Resource: metalResourceId,
					Amount:   36,
				},
				{
					Resource: crystalResourceId,
					Amount:   178,
				},
			},
		}

		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   36,
			},
			{
				Resource: crystalResourceId,
				Amount:   178,
			},
		}
		assert.Equal(t, expected, p.Resources)
	})

	t.Run("bumps version by one", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.BuildingAction = &BuildingAction{Id: uuid.New()}

		initialVersion := p.Version

		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, initialVersion+1, p.Version)
	})

	t.Run("does not bump updated at field", func(t *testing.T) {
		p := generateTestPlanet(t, withPlanetBuilding, withManyResources)
		p.BuildingAction = &BuildingAction{Id: uuid.New()}

		err := p.CancelBuildingAction()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, someTime, p.UpdatedAt)
	})
}

func TestUnit_Planet_UpdateToTime(t *testing.T) {
	someTime := time.Date(2026, time.June, 25, 20, 19, 37, 0, time.UTC)
	someTimeLater := someTime.Add(1*time.Hour + 2*time.Minute + 3*time.Second)

	t.Run("does not update resource when no production is defined", func(t *testing.T) {
		p := Planet{
			Resources: []PlanetResource{{Resource: crystalResourceId, Amount: 36}},
			UpdatedAt: someTime,
		}

		err := p.UpdateToTime(someTimeLater)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Len(t, p.Resources, 1)
		assert.Equal(t, 36.0, p.Resources[0].Amount)
	})

	t.Run("does not change resource value when already over the storage capacity", func(t *testing.T) {
		p := Planet{
			Resources:   []PlanetResource{{Resource: crystalResourceId, Amount: 36}},
			Storages:    []PlanetResourceStorage{{Resource: crystalResourceId, Storage: 30}},
			Productions: []PlanetResourceProduction{{Resource: crystalResourceId, Production: 30}},
			UpdatedAt:   someTime,
		}

		err := p.UpdateToTime(someTimeLater)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Len(t, p.Resources, 1)
		assert.Equal(t, 36.0, p.Resources[0].Amount)
	})

	t.Run("does not change resource when update time is before already computed time", func(t *testing.T) {
		p := Planet{
			Resources:   []PlanetResource{{Resource: crystalResourceId, Amount: 36}},
			Storages:    []PlanetResourceStorage{{Resource: crystalResourceId, Storage: 300}},
			Productions: []PlanetResourceProduction{{Resource: crystalResourceId, Production: 30}},
			UpdatedAt:   someTimeLater,
		}

		err := p.UpdateToTime(someTime)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Len(t, p.Resources, 1)
		assert.Equal(t, 36.0, p.Resources[0].Amount)
	})

	t.Run("caps resource at storage capacity", func(t *testing.T) {
		p := Planet{
			Resources:   []PlanetResource{{Resource: crystalResourceId, Amount: 36}},
			Storages:    []PlanetResourceStorage{{Resource: crystalResourceId, Storage: 45}},
			Productions: []PlanetResourceProduction{{Resource: crystalResourceId, Production: 30}},
			UpdatedAt:   someTime,
		}

		err := p.UpdateToTime(someTimeLater)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Len(t, p.Resources, 1)
		assert.Equal(t, 45.0, p.Resources[0].Amount)
	})

	t.Run("adds full production value when storage is sufficient", func(t *testing.T) {
		p := Planet{
			Resources:   []PlanetResource{{Resource: crystalResourceId, Amount: 36}},
			Storages:    []PlanetResourceStorage{{Resource: crystalResourceId, Storage: 300}},
			Productions: []PlanetResourceProduction{{Resource: crystalResourceId, Production: 30}},
			UpdatedAt:   someTime,
		}

		err := p.UpdateToTime(someTimeLater)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Len(t, p.Resources, 1)
		assert.Equal(t, 67.025, p.Resources[0].Amount)
	})

	t.Run("adds all production for a resource", func(t *testing.T) {
		p := Planet{
			Resources: []PlanetResource{{Resource: crystalResourceId, Amount: 36}},
			Storages:  []PlanetResourceStorage{{Resource: crystalResourceId, Storage: 300}},
			Productions: []PlanetResourceProduction{
				{Resource: crystalResourceId, Production: 30},
				{Resource: crystalResourceId, Production: 45, Building: &crystalMineId},
			},
			UpdatedAt: someTime,
		}

		err := p.UpdateToTime(someTimeLater)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Len(t, p.Resources, 1)
		assert.Equal(t, 113.5625, p.Resources[0].Amount)
	})

	t.Run("keeps resource at 0 when no storage is defined for it", func(t *testing.T) {
		p := Planet{
			Resources: []PlanetResource{{Resource: crystalResourceId, Amount: 36}},
			Storages:  nil,
			Productions: []PlanetResourceProduction{
				{Resource: crystalResourceId, Production: 30},
				{Resource: crystalResourceId, Production: 45, Building: &crystalMineId},
			},
			UpdatedAt: someTime,
		}

		err := p.UpdateToTime(someTimeLater)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Len(t, p.Resources, 1)
		assert.Equal(t, 36.0, p.Resources[0].Amount)
	})

	t.Run("bumps updated at field", func(t *testing.T) {
		p := Planet{UpdatedAt: someTime}

		err := p.UpdateToTime(someTimeLater)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, someTimeLater, p.UpdatedAt)
	})

	t.Run("does not bump updated at field when update time is in the past", func(t *testing.T) {
		p := Planet{UpdatedAt: someTimeLater}

		err := p.UpdateToTime(someTime)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, someTimeLater, p.UpdatedAt)
	})

	t.Run("returns error when building action finishes before update time", func(t *testing.T) {
		t1 := time.Date(2026, time.June, 26, 8, 30, 50, 0, time.UTC)
		t2 := time.Date(2026, time.June, 26, 8, 31, 50, 0, time.UTC)
		t3 := time.Date(2026, time.June, 26, 8, 32, 50, 0, time.UTC)

		p := Planet{
			BuildingAction: &BuildingAction{
				CompletedAt: t2,
			},
			UpdatedAt: t1,
		}

		err := p.UpdateToTime(t3)

		assert.ErrorIs(t, domainerrors.ErrPlanetNotUpToDate, err, "Actual err is: %v", err)
	})
}

func TestUnit_Planet_ApplyAction(t *testing.T) {
	t1 := time.Date(2026, time.June, 26, 8, 41, 30, 0, time.UTC)
	t2 := time.Date(2026, time.June, 26, 8, 42, 30, 0, time.UTC)

	t.Run("returns error when no action is in progress", func(t *testing.T) {
		p := Planet{BuildingAction: nil}

		err := p.ApplyAction()

		assert.ErrorIs(t, domainerrors.ErrNoActionInProgress, err, "Actual err: %v", err)
	})

	t.Run("returns error when planet update time is not matching action completion time", func(t *testing.T) {
		p := Planet{
			BuildingAction: &BuildingAction{
				CompletedAt: t2,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()

		assert.ErrorIs(t, domainerrors.ErrActionNotCompleted, err, "Actual err: %v", err)
	})

	t.Run("applies resource production changes when applicable", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: metalMineId, Level: 1},
			},
			Productions: []PlanetResourceProduction{
				{Resource: metalResourceId, Production: 50},
				{Resource: metalResourceId, Production: 30, Building: &metalMineId},
			},
			BuildingAction: &BuildingAction{
				Building:     metalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				Productions: []BuildingActionResourceProduction{
					{Resource: metalResourceId, Production: 45},
				},
				CompletedAt: t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResourceProduction{
			{Resource: metalResourceId, Production: 50},
			{Resource: metalResourceId, Production: 45, Building: &metalMineId},
		}
		assert.Equal(t, expected, p.Productions)
	})

	t.Run("registers new resource production for building when applicable", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: metalMineId, Level: 1},
			},
			Productions: []PlanetResourceProduction{
				{Resource: metalResourceId, Production: 30},
			},
			BuildingAction: &BuildingAction{
				Building:     metalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				Productions: []BuildingActionResourceProduction{
					{Resource: metalResourceId, Production: 45},
				},
				CompletedAt: t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResourceProduction{
			{Resource: metalResourceId, Production: 30},
			{Resource: metalResourceId, Production: 45, Building: &metalMineId},
		}
		assert.Equal(t, expected, p.Productions)
	})

	t.Run("does not change resource production when action has no effect", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: metalMineId, Level: 1},
			},
			Productions: []PlanetResourceProduction{
				{Resource: metalResourceId, Production: 30},
			},
			BuildingAction: &BuildingAction{
				Building:     metalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				Productions:  []BuildingActionResourceProduction{},
				CompletedAt:  t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResourceProduction{
			{Resource: metalResourceId, Production: 30},
		}
		assert.Equal(t, expected, p.Productions)
	})

	t.Run("applies resource storage changes when applicable", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: metalMineId, Level: 1},
			},
			Storages: []PlanetResourceStorage{
				{Resource: metalResourceId, Storage: 1000},
			},
			BuildingAction: &BuildingAction{
				Building:     metalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				Storages: []BuildingActionResourceStorage{
					{Resource: metalResourceId, Storage: 2000},
				},
				CompletedAt: t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResourceStorage{
			{Resource: metalResourceId, Storage: 2000},
		}
		assert.Equal(t, expected, p.Storages)
	})

	t.Run("does not change resource storage when action has no effect", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: metalMineId, Level: 1},
			},
			Storages: []PlanetResourceStorage{
				{Resource: metalResourceId, Storage: 1000},
			},
			BuildingAction: &BuildingAction{
				Building:     metalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				Storages:     []BuildingActionResourceStorage{},
				CompletedAt:  t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetResourceStorage{
			{Resource: metalResourceId, Storage: 1000},
		}
		assert.Equal(t, expected, p.Storages)
	})

	t.Run("updates building to desired level", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: crystalMineId, Level: 1},
			},
			BuildingAction: &BuildingAction{
				Building:     crystalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				CompletedAt:  t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		expected := []PlanetBuilding{
			{Building: crystalMineId, Level: 2},
		}
		assert.Equal(t, expected, p.Buildings)
	})

	t.Run("bumps version by one", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: crystalMineId, Level: 1},
			},
			Version: 0,
			BuildingAction: &BuildingAction{
				Building:     crystalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				CompletedAt:  t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, 1, p.Version)
	})

	t.Run("does not change updated at field", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: crystalMineId, Level: 1},
			},
			Version: 0,
			BuildingAction: &BuildingAction{
				Building:     crystalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				CompletedAt:  t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, t1, p.UpdatedAt)
	})

	t.Run("removes building action from planet", func(t *testing.T) {
		p := Planet{
			Buildings: []PlanetBuilding{
				{Building: crystalMineId, Level: 1},
			},
			Version: 0,
			BuildingAction: &BuildingAction{
				Building:     crystalMineId,
				CurrentLevel: 1,
				DesiredLevel: 2,
				CompletedAt:  t1,
			},
			UpdatedAt: t1,
		}

		err := p.ApplyAction()
		require.NoError(t, err, "Actual err: %v", err)

		assert.Nil(t, p.BuildingAction)
	})
}

func generateTestPlanet(
	t *testing.T,
	modifiers ...func(*testing.T, *Planet),
) Planet {
	t.Helper()

	p := Planet{
		Id:        uuid.New(),
		UpdatedAt: someTime,
		Version:   3,
	}

	for _, modifier := range modifiers {
		modifier(t, &p)
	}

	return p
}

func withManyResources(t *testing.T, p *Planet) {
	t.Helper()

	// High enough values to not have to worry about costs
	p.Resources = []PlanetResource{
		{
			Resource: metalResourceId,
			Amount:   999999,
		},
		{
			Resource: crystalResourceId,
			Amount:   999999,
		},
	}
}

func withPlanetBuilding(t *testing.T, p *Planet) {
	p.Buildings = []PlanetBuilding{
		{
			Building: buildingId,
			Level:    4,
		},
	}
}
