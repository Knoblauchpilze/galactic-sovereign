package drivenadapters

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_ShipRepository_Get(t *testing.T) {
	repo, conn := newTestShipRepository(t)

	t.Run("gets a ship", func(t *testing.T) {
		ship := insertTestShip(t, conn)

		actual, err := repo.Get(t.Context(), ship.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, ship, actual)
	})

	t.Run("gets a ship with costs", func(t *testing.T) {
		ship := insertTestShip(t, conn, addShipCost)

		actual, err := repo.Get(t.Context(), ship.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, ship, actual)
	})

	t.Run("returns error when ship does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func newTestShipRepository(t *testing.T) (*ShipRepository, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewShipRepository(conn), conn
}

func insertTestShip(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Ship),
) models.Ship {
	t.Helper()

	ship := models.Ship{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-ship-%s", uuid.NewString()),
		CreatedAt: someTime,
		// This is intentional: the details (e.g. costs) are returned as empty
		// slices by the adapter
		Costs: []models.ShipCost{},
	}

	sqlQuery := `INSERT INTO ship (id, name, created_at)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		ship.Id,
		ship.Name,
		ship.CreatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &ship)
	}

	return ship
}

func addShipCost(t *testing.T, conn db.Connection, s *models.Ship) {
	t.Helper()

	cost := models.ShipCost{
		Resource:              metalResourceId,
		Cost:                  rand.Intn(897),
		BuildTimeHoursPerUnit: 0.0004,
	}

	sqlQuery := `INSERT INTO ship_cost (ship, resource, cost)
		VALUES ($1, $2, $3)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		s.Id,
		cost.Resource,
		cost.Cost,
	)
	require.NoError(t, err, "Actual err: %v", err)

	s.Costs = append(s.Costs, cost)
}
