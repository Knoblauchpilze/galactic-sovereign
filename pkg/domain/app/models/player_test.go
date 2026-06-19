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

	t.Run("creates a homeworld belonging to the player", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []uuid.UUID{},
		}

		actual := p.CreateHomeworld(resources)

		assert.Equal(t, p.Id, actual.Player)
		assert.True(t, actual.Homeworld)
		assert.Equal(t, "homeworld", actual.Name)
		assert.Zero(t, actual.Version)
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})

	t.Run("assigns planet when slice is nil", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: nil,
		}

		actual := p.CreateHomeworld(resources)

		assert.Equal(t, p.Id, actual.Player)
		assert.True(t, actual.Homeworld)
		assert.Equal(t, "homeworld", actual.Name)
		assert.Zero(t, actual.Version)
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})

	t.Run("assigns start amount for each resource", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []uuid.UUID{},
		}

		actual := p.CreateHomeworld(resources)

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

		actual := p.CreateHomeworld(resources)

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

		actual := p.CreateHomeworld(resources)

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
}
