package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnit_Player_CreateHomeworld(t *testing.T) {
	u := Universe{
		Id:        uuid.New(),
		Resources: sampleResources(),
		Buildings: sampleBuildings(),
	}

	t.Run("creates a homeworld belonging to the player", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: []uuid.UUID{},
		}

		actual := p.CreateHomeworld(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.True(t, actual.Homeworld)
		assert.Equal(t, "homeworld", actual.Name)
		assert.Equal(t, actual.Id, p.Homeworld)
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})

	t.Run("assigns planet when slice is nil", func(t *testing.T) {
		p := Player{
			Id:      uuid.New(),
			Planets: nil,
		}

		actual := p.CreateHomeworld(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.True(t, actual.Homeworld)
		assert.Equal(t, "homeworld", actual.Name)
		assert.Equal(t, actual.Id, p.Homeworld)
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})
}

func TestUnit_Player_Colonize(t *testing.T) {
	u := Universe{
		Id:        uuid.New(),
		Resources: sampleResources(),
		Buildings: sampleBuildings(),
	}

	t.Run("creates a planet belonging to the player", func(t *testing.T) {
		p := Player{
			Id:        uuid.New(),
			Homeworld: uuid.New(),
			Planets:   []uuid.UUID{},
		}

		actual := p.Colonize(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.False(t, actual.Homeworld)
		assert.Equal(t, "colony", actual.Name)
		assert.NotEqual(t, actual.Id, p.Homeworld)
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})

	t.Run("assigns planet when slice is nil", func(t *testing.T) {
		p := Player{
			Id:        uuid.New(),
			Homeworld: uuid.New(),
			Planets:   nil,
		}

		actual := p.Colonize(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.False(t, actual.Homeworld)
		assert.Equal(t, "colony", actual.Name)
		assert.NotEqual(t, actual.Id, p.Homeworld)
		assert.Equal(t, []uuid.UUID{actual.Id}, p.Planets)
	})

	t.Run("assigns planet when multiple planets already exist", func(t *testing.T) {
		homeworldId := uuid.New()
		p := Player{
			Id:        uuid.New(),
			Homeworld: homeworldId,
			Planets:   []uuid.UUID{homeworldId},
		}

		actual := p.Colonize(u)

		assert.Equal(t, p.Id, actual.Player)
		assert.False(t, actual.Homeworld)
		assert.Equal(t, "colony", actual.Name)
		assert.NotEqual(t, actual.Id, p.Homeworld)
		assert.Equal(t, []uuid.UUID{homeworldId, actual.Id}, p.Planets)
	})
}
