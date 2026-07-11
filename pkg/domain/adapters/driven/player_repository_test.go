package drivenadapters

import (
	"fmt"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	domainerrors "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models/errors"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_PlayerRepository_Create(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)

	t.Run("creates a player with a homeworld", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		planetId := uuid.New()

		player := models.Player{
			Id:        uuid.New(),
			ApiUser:   uuid.New(),
			Universe:  universe.Id,
			Name:      fmt.Sprintf("player-%s", uuid.NewString()),
			CreatedAt: someTime,
			Homeworld: planetId,
			Planets:   []uuid.UUID{planetId},
		}
		planet := models.Planet{
			Id:        planetId,
			Player:    player.Id,
			Name:      fmt.Sprintf("planetr-%s", uuid.NewString()),
			Homeworld: true,
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
			Version:   0,
			Resources: []models.PlanetResource{
				{
					Resource: metalResourceId,
					Amount:   389,
				},
			},
			Storages: []models.PlanetResourceStorage{
				{
					Resource: crystalResourceId,
					Storage:  74812,
				},
			},
			Productions: []models.PlanetResourceProduction{
				{
					Resource:   metalResourceId,
					Production: 14587,
				},
			},
			Buildings: []models.PlanetBuilding{
				{
					Building: metalMineId,
					Level:    2,
				},
			},
		}

		err := repo.Create(t.Context(), player, planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertPlayerExists(t, conn, player.Id)
		assertPlanetExists(t, conn, planet.Id)

		actual, err := repo.Get(t.Context(), player.Id)
		expectedPlayer := player
		expectedPlayer.Planets = []uuid.UUID{planet.Id}
		require.NoError(t, err, "Actual err: %v", err)
		assert.Equal(t, expectedPlayer, actual)

		actualPlanet := loadPlanetFromDb(t, conn, planet.Id)
		assert.Equal(t, planet, actualPlanet)
	})

	t.Run("does not create planet for a player", func(t *testing.T) {
		universe := insertTestUniverse(t, conn)

		player := models.Player{
			Id:        uuid.New(),
			ApiUser:   uuid.New(),
			Universe:  universe.Id,
			Name:      fmt.Sprintf("player-%s", uuid.NewString()),
			CreatedAt: someTime,
			Planets:   []uuid.UUID{uuid.New()},
		}
		planet := models.Planet{
			Id:          uuid.New(),
			Player:      player.Id,
			Name:        fmt.Sprintf("planetr-%s", uuid.NewString()),
			Homeworld:   true,
			CreatedAt:   someTime,
			UpdatedAt:   someOtherTime,
			Version:     0,
			Resources:   []models.PlanetResource{},
			Storages:    []models.PlanetResourceStorage{},
			Productions: []models.PlanetResourceProduction{},
			Buildings:   []models.PlanetBuilding{},
		}

		err := repo.Create(t.Context(), player, planet)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlayerExists(t, conn, player.Id)
		assertPlanetDoesNotExist(t, conn, player.Planets[0])
	})

	t.Run("returns error when player with same name already exists", func(t *testing.T) {
		player, universe := insertTestPlayerInUniverse(t, conn)

		newPlayer := models.Player{
			Id:        uuid.New(),
			ApiUser:   uuid.New(),
			Universe:  universe.Id,
			Name:      player.Name,
			CreatedAt: someTime,
		}
		planet := models.Planet{
			Id:          uuid.New(),
			Player:      player.Id,
			Name:        fmt.Sprintf("planetr-%s", uuid.NewString()),
			Homeworld:   true,
			CreatedAt:   someTime,
			UpdatedAt:   someOtherTime,
			Version:     0,
			Resources:   []models.PlanetResource{},
			Storages:    []models.PlanetResourceStorage{},
			Productions: []models.PlanetResourceProduction{},
			Buildings:   []models.PlanetBuilding{},
		}

		err := repo.Create(t.Context(), newPlayer, planet)

		assert.Equal(t, domainerrors.ErrNameAlreadyTaken, err, "Actual err: %v", err)
		assertPlayerDoesNotExist(t, conn, newPlayer.Id)
	})
}

func TestIT_PlayerRepository_Get(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)

	t.Run("gets a player", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)

		actual, err := repo.Get(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, player, actual)
	})

	t.Run("gets a player with planets", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn, addPlayerPlanet)

		actual, err := repo.Get(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, player, actual)
	})

	t.Run("gets a player with a homeworld", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)
		planet := insertTestPlanet(t, conn, player.Id, addPlanetHomeworld)

		actual, err := repo.Get(t.Context(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		expected := models.Player{
			Id:        player.Id,
			ApiUser:   player.ApiUser,
			Universe:  player.Universe,
			Name:      player.Name,
			CreatedAt: player.CreatedAt,
			Version:   player.Version,
			Homeworld: planet.Id,
			Planets:   []uuid.UUID{planet.Id},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("returns error when player does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(t.Context(), id)

		assert.ErrorIs(t, err, domainerrors.ErrNotFound, "Actual err: %v", err)
	})
}

func TestIT_PlayerRepository_List(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)
	p1, universe := insertTestPlayerInUniverse(t, conn)
	p2 := insertTestPlayer(t, conn, universe.Id)
	p3 := insertTestPlayer(t, conn, universe.Id, addPlayerPlanet)

	p4, _ := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, p4.Id, addPlanetHomeworld)

	actual, err := repo.List(t.Context())
	require.NoError(t, err, "Actual err: %v", err)

	// The additional resources are players from the seed data
	assert.Contains(t, actual, p1)
	assert.Contains(t, actual, p2)
	assert.Contains(t, actual, p3)

	expectedP4 := p4
	expectedP4.Homeworld = planet.Id
	expectedP4.Planets = []uuid.UUID{planet.Id}
	assert.Contains(t, actual, expectedP4)
}

func TestIT_PlayerRepository_ListForApiUser(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)

	t.Run("lists player for an API user", func(t *testing.T) {
		p1, universe := insertTestPlayerInUniverse(t, conn)
		insertTestPlayer(t, conn, universe.Id)

		actual, err := repo.ListForApiUser(t.Context(), p1.ApiUser)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Player{p1}, actual)
	})

	t.Run("lists player with a homeworld", func(t *testing.T) {
		p1, universe := insertTestPlayerInUniverse(t, conn)
		planet := insertTestPlanet(t, conn, p1.Id, addPlanetHomeworld)

		insertTestPlayer(t, conn, universe.Id)

		actual, err := repo.ListForApiUser(t.Context(), p1.ApiUser)
		require.NoError(t, err, "Actual err: %v", err)

		expectedP1 := p1
		expectedP1.Homeworld = planet.Id
		expectedP1.Planets = []uuid.UUID{planet.Id}
		assert.Equal(t, []models.Player{expectedP1}, actual)
	})

	t.Run("lists player with planets for an API user", func(t *testing.T) {
		p1, universe := insertTestPlayerInUniverse(t, conn, addPlayerPlanet)
		insertTestPlayer(t, conn, universe.Id)

		actual, err := repo.ListForApiUser(t.Context(), p1.ApiUser)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, []models.Player{p1}, actual)
	})
}

func TestIT_PlayerRepository_Delete(t *testing.T) {
	repo, conn := newTestPlayerRepository(t)

	t.Run("deletes a player", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)

		err := repo.Delete(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlayerDoesNotExist(t, conn, player.Id)
	})

	t.Run("deletes a player with planet", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)

		planet := insertTestPlanet(t, conn, player.Id)
		player.Planets = []uuid.UUID{planet.Id}

		err := repo.Delete(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlayerDoesNotExist(t, conn, player.Id)
		assertPlanetDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes a player with a building action", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)
		planet := insertTestPlanet(t, conn, player.Id, addPlanetBuildingAction)
		player.Planets = []uuid.UUID{planet.Id}

		err := repo.Delete(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlayerDoesNotExist(t, conn, player.Id)
		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertBuildingActionDoesNotExist(t, conn, planet.BuildingAction.Id)
	})

	t.Run("succeeds when the player does not exist", func(t *testing.T) {
		player := models.Player{Id: uuid.New()}

		err := repo.Delete(t.Context(), player)
		require.NoError(t, err, "Actual err: %v", err)
	})
}

func newTestPlayerRepository(t *testing.T) (drivenports.ForManagingPlayers, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewPlayerRepository(conn), conn
}

func insertTestPlayer(
	t *testing.T,
	conn db.Connection,
	universe uuid.UUID,
	modifiers ...func(*testing.T, db.Connection, *models.Player),
) models.Player {
	t.Helper()

	player := models.Player{
		Id:        uuid.New(),
		ApiUser:   uuid.New(),
		Universe:  universe,
		Name:      fmt.Sprintf("my-player-%s", uuid.NewString()),
		CreatedAt: someTime,
		// This is intentional: the details (e.g. planets, etc.) are returned as empty
		// slices by the adapter
		Planets: []uuid.UUID{},
	}

	sqlQuery := `INSERT INTO player (id, api_user, universe, name, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		player.Id,
		player.ApiUser,
		player.Universe,
		player.Name,
		player.CreatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &player)
	}

	return player
}

func insertTestPlayerInUniverse(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Player),
) (models.Player, models.Universe) {
	universe := insertTestUniverse(t, conn)
	player := insertTestPlayer(t, conn, universe.Id, modifiers...)
	return player, universe
}

func addPlayerPlanet(t *testing.T, conn db.Connection, p *models.Player) {
	t.Helper()

	planetId := uuid.New()

	sqlQuery := `INSERT INTO planet (id, player, name, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		planetId,
		p.Id,
		fmt.Sprintf("my-planet-%s", planetId.String()),
		someTime,
		someOtherTime,
		8,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Planets = append(p.Planets, planetId)
}

func assertPlayerExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT id FROM player WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, id, value)
}

func assertPlayerDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM player WHERE id = $1`
	value, err := db.QueryOne[int](t.Context(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
