package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_Player_CreateHomeworld(t *testing.T) {
	resources := []Resource{
		{
			Id:              metalResourceId,
			StartAmount:     145,
			StartStorage:    123,
			StartProduction: 985,
		},
		{
			Id:              crystalResourceId,
			StartAmount:     325,
			StartStorage:    421,
			StartProduction: 752,
		},
	}

	buildings := []Building{
		{
			Id: uuid.New(),
			Costs: []BuildingCost{
				{
					Resource: metalResourceId,
					Cost:     550,
					Progress: 26.4,
				},
			},
		},
		{
			Id: uuid.New(),
			Storages: []BuildingResourceStorage{
				{
					Resource: metalResourceId,
					Scale:    1500.0,
					Progress: 26.4,
				},
			},
		},
	}

	t.Run("creates a homeworld belonging to the player", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []uuid.UUID{},
		}

		actual := p.CreateHomeworld(resources, buildings)

		assert.Equal(t, p.Id, actual.Player)
		assert.True(t, actual.Homeworld)
		assert.Equal(t, "homeworld", actual.Name)
		assert.Zero(t, actual.Version)
		assert.Equal(t, actual.Id, p.Homeworld)
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})

	t.Run("assigns planet when slice is nil", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: nil,
		}

		actual := p.CreateHomeworld(resources, buildings)

		assert.Equal(t, p.Id, actual.Player)
		assert.True(t, actual.Homeworld)
		assert.Equal(t, "homeworld", actual.Name)
		assert.Zero(t, actual.Version)
		assert.Equal(t, actual.Id, p.Homeworld)
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})

	t.Run("assigns start amount for each resource", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []uuid.UUID{},
		}

		actual := p.CreateHomeworld(resources, buildings)

		expected := []PlanetResource{
			{
				Resource: metalResourceId,
				Amount:   145,
			},
			{
				Resource: crystalResourceId,
				Amount:   325,
			},
		}
		assert.Equal(t, expected, actual.Resources)
	})

	t.Run("assigns start storage for each resource", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []uuid.UUID{},
		}

		actual := p.CreateHomeworld(resources, buildings)

		expected := []PlanetResourceStorage{
			{
				Resource: metalResourceId,
				Storage:  123,
			},
			{
				Resource: crystalResourceId,
				Storage:  421,
			},
		}
		assert.Equal(t, expected, actual.Storages)
	})

	t.Run("assigns start production for each resource", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []uuid.UUID{},
		}

		actual := p.CreateHomeworld(resources, buildings)

		expected := []PlanetResourceProduction{
			{
				Resource:   metalResourceId,
				Production: 985,
			},
			{
				Resource:   crystalResourceId,
				Production: 752,
			},
		}
		assert.Equal(t, expected, actual.Productions)
	})

	t.Run("creates each building with level 0", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []uuid.UUID{},
		}

		actual := p.CreateHomeworld(resources, buildings)

		expected := []PlanetBuilding{
			{
				Building: buildings[0].Id,
				Level:    0,
			},
			{
				Building: buildings[1].Id,
				Level:    0,
			},
		}
		assert.Equal(t, expected, actual.Buildings)
	})
}
