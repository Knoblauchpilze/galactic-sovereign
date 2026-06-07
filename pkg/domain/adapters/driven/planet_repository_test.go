package drivenadapters

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	drivenports "github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/ports/driven"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	crystalResourceId = uuid.MustParse("cd2ac9aa-9968-4ff5-b746-88f1f810fbb3")
	crystalMineId     = uuid.MustParse("3904d34d-9a7e-47d4-a332-091700e2c5c3")
	metalStorageId    = uuid.MustParse("22b4c0c3-c8e5-4493-89fc-522fdbb0beee")
)

func TestIT_PlanetRepository_Create(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)

	t.Run("creates a planet", func(t *testing.T) {
		player, _ := insertTestPlayerInUniverse(t, conn)

		planet := models.Planet{
			Id:          uuid.New(),
			Player:      player.Id,
			Name:        fmt.Sprintf("my-planet-%s", uuid.NewString()),
			Homeworld:   false,
			CreatedAt:   someTime,
			UpdatedAt:   someOtherTime,
			Version:     3,
			Resources:   []models.PlanetResource{},
			Storages:    []models.PlanetResourceStorage{},
			Productions: []models.PlanetResourceProduction{},
			Buildings:   []models.PlanetBuilding{},
		}

		err := repo.Create(context.Background(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertPlanetExists(t, conn, planet.Id)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, actual)
	})

	t.Run("returns error when planet with same id already exists", func(t *testing.T) {
		planet, player, _ := insertTestPlanetForPlayer(t, conn)

		duplicatedPlanet := models.Planet{
			Id:        planet.Id,
			Player:    player.Id,
			Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
			Homeworld: false,
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
			Version:   4,
		}

		err := repo.Create(context.Background(), duplicatedPlanet)
		assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	})

	t.Run("creates homeworld and marks it as such", func(t *testing.T) {

		player, _ := insertTestPlayerInUniverse(t, conn)

		planet := models.Planet{
			Id:          uuid.New(),
			Player:      player.Id,
			Name:        fmt.Sprintf("my-planet-%s", uuid.NewString()),
			Homeworld:   true,
			CreatedAt:   someTime,
			UpdatedAt:   someOtherTime,
			Resources:   []models.PlanetResource{},
			Storages:    []models.PlanetResourceStorage{},
			Productions: []models.PlanetResourceProduction{},
			Buildings:   []models.PlanetBuilding{},
		}

		err := repo.Create(context.Background(), planet)
		require.NoError(t, err, "Actual err: %v", err)
		assertPlanetExists(t, conn, planet.Id)
		assertPlanetIsHomeworld(t, conn, planet.Id, planet.Player)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, actual)
	})

	t.Run("returns error when player already has a homeworld", func(t *testing.T) {
		_, player, _ := insertTestHomeworldPlanetForPlayer(t, conn)

		planet := models.Planet{
			Id:        uuid.New(),
			Player:    player.Id,
			Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
			Homeworld: true,
			CreatedAt: someTime,
			UpdatedAt: someOtherTime,
		}

		err := repo.Create(context.Background(), planet)
		assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
		assertPlanetDoesNotExist(t, conn, planet.Id)
	})
}

func TestIT_PlanetRepository_Get(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	t.Run("gets planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, planet, actual)
	})

	t.Run("gets planet with resources", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, planet)
	})

	t.Run("gets planet with storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, planet)
	})

	t.Run("gets planet with productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, planet)
	})

	t.Run("gets planet with productions for building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, planet)
	})

	t.Run("gets planet with building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)

		actual, err := repo.Get(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assert.Equal(t, actual, planet)
	})

	t.Run("returns error when planet does not exist", func(t *testing.T) {
		id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
		_, err := repo.Get(context.Background(), id)

		assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
	})
}

func TestIT_PlanetRepository_List(t *testing.T) {
	repo, conn := newTestPlanetRepositoryWithContainer(t)
	defer conn.Close(context.Background())
	p1, player1, _ := insertTestPlanetForPlayer(t, conn)
	p2 := insertTestPlanet(t, conn, player1.Id)
	p3, player2, _ := insertTestPlanetForPlayer(t, conn)
	p4 := insertTestPlanet(t, conn, player2.Id)
	p5 := insertTestPlanet(t, conn, player2.Id, addPlanetResource)
	p6 := insertTestPlanet(t, conn, player2.Id, addPlanetStorage)
	p7 := insertTestPlanet(t, conn, player2.Id, addPlanetProduction)
	p8 := insertTestPlanet(t, conn, player2.Id, addPlanetProductionForBuilding)
	p9 := insertTestPlanet(t, conn, player2.Id, addPlanetBuilding)

	actual, err := repo.List(context.Background())
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 9)
	assert.Contains(t, actual, p1)
	assert.Contains(t, actual, p2)
	assert.Contains(t, actual, p3)
	assert.Contains(t, actual, p4)
	assert.Contains(t, actual, p5)
	assert.Contains(t, actual, p6)
	assert.Contains(t, actual, p7)
	assert.Contains(t, actual, p8)
	assert.Contains(t, actual, p9)
}

func newTestPlanetRepositoryWithContainer(t *testing.T) (drivenports.ForManagingPlanets, db.Connection) {
	t.Helper()

	ctx := context.Background()
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	require.NoError(t, err, "Actual err: %v", err)

	t.Cleanup(func() {
		require.NoError(t, pgContainer.Terminate(context.Background()))
	})

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err, "Actual err: %v", err)

	mappedPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err, "Actual err: %v", err)

	portInt, err := strconv.ParseUint(mappedPort.Port(), 10, 16)
	require.NoError(t, err, "Actual err: %v", err)
	port := uint16(portInt)

	postgresConfig := postgresql.Config{
		Host:           host,
		Port:           port,
		Database:       "postgres",
		User:           "postgres",
		Password:       "postgres",
		ConnectTimeout: dbTestConfig.ConnectTimeout,
	}

	bootstrapConn, err := db.New(ctx, postgresConfig)
	require.NoError(t, err, "Actual err: %v", err)
	t.Cleanup(func() {
		bootstrapConn.Close(context.Background())
	})

	require.NoError(t, bootstrapDatabaseForPlanetTests(ctx, bootstrapConn, host, port))

	adminConfig := postgresql.Config{
		Host:           host,
		Port:           port,
		Database:       "db_galactic_sovereign",
		User:           "galactic_sovereign_admin",
		Password:       "admin_password",
		ConnectTimeout: dbTestConfig.ConnectTimeout,
	}
	adminConn, err := db.New(ctx, adminConfig)
	require.NoError(t, err, "Actual err: %v", err)
	t.Cleanup(func() {
		adminConn.Close(context.Background())
	})

	require.NoError(t, runAllUpMigrationsForPlanetTests(ctx, adminConn))

	managerConfig := postgresql.Config{
		Host:           host,
		Port:           port,
		Database:       "db_galactic_sovereign",
		User:           "galactic_sovereign_manager",
		Password:       "manager_password",
		ConnectTimeout: dbTestConfig.ConnectTimeout,
	}

	conn, err := db.New(ctx, managerConfig)
	require.NoError(t, err, "Actual err: %v", err)
	return NewPlanetRepository(conn), conn
}

func bootstrapDatabaseForPlanetTests(ctx context.Context, conn db.Connection, host string, port uint16) error {
	statements := []string{
		"CREATE USER galactic_sovereign_admin WITH CREATEDB PASSWORD 'admin_password'",
		"CREATE USER galactic_sovereign_manager WITH PASSWORD 'manager_password'",
		"CREATE USER galactic_sovereign_user WITH PASSWORD 'user_password'",
		"GRANT galactic_sovereign_user TO galactic_sovereign_manager",
		"GRANT galactic_sovereign_manager TO galactic_sovereign_admin",
		"CREATE DATABASE db_galactic_sovereign OWNER galactic_sovereign_admin",
		"REVOKE ALL ON DATABASE db_galactic_sovereign FROM public",
		"GRANT CONNECT ON DATABASE db_galactic_sovereign TO galactic_sovereign_user",
	}

	for _, stmt := range statements {
		if _, err := conn.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	setupConfig := postgresql.Config{
		Host:           host,
		Port:           port,
		Database:       "db_galactic_sovereign",
		User:           "postgres",
		Password:       "postgres",
		ConnectTimeout: dbTestConfig.ConnectTimeout,
	}

	setupConn, err := db.New(ctx, setupConfig)
	if err != nil {
		return err
	}
	defer setupConn.Close(context.Background())

	dbStatements := []string{
		"CREATE SCHEMA galactic_sovereign_schema AUTHORIZATION galactic_sovereign_admin",
		"ALTER ROLE galactic_sovereign_admin IN DATABASE db_galactic_sovereign SET search_path = galactic_sovereign_schema",
		"ALTER ROLE galactic_sovereign_manager IN DATABASE db_galactic_sovereign SET search_path = galactic_sovereign_schema",
		"ALTER ROLE galactic_sovereign_user IN DATABASE db_galactic_sovereign SET search_path = galactic_sovereign_schema",
		"GRANT USAGE ON SCHEMA galactic_sovereign_schema TO galactic_sovereign_user",
		"GRANT CREATE ON SCHEMA galactic_sovereign_schema TO galactic_sovereign_admin",
		"ALTER DEFAULT PRIVILEGES FOR ROLE galactic_sovereign_admin GRANT SELECT ON TABLES TO galactic_sovereign_user",
		"ALTER DEFAULT PRIVILEGES FOR ROLE galactic_sovereign_admin GRANT INSERT, UPDATE, DELETE ON TABLES TO galactic_sovereign_manager",
	}

	for _, stmt := range dbStatements {
		if _, err := setupConn.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}

func runAllUpMigrationsForPlanetTests(ctx context.Context, conn db.Connection) error {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to resolve current file path")
	}

	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", ".."))
	migrationsDir := filepath.Join(repoRoot, "database", "galactic-sovereign", "migrations")

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	type migrationFile struct {
		version int
		name    string
	}

	upFiles := []migrationFile{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		versionPart, _, found := strings.Cut(name, "_")
		if !found {
			continue
		}

		version, err := strconv.Atoi(versionPart)
		if err != nil {
			return err
		}

		upFiles = append(upFiles, migrationFile{version: version, name: name})
	}

	sort.Slice(upFiles, func(i int, j int) bool {
		if upFiles[i].version == upFiles[j].version {
			return upFiles[i].name < upFiles[j].name
		}
		return upFiles[i].version < upFiles[j].version
	})

	for _, migration := range upFiles {
		migrationPath := filepath.Join(migrationsDir, migration.name)
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return err
		}

		if _, err := conn.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("applying migration %s: %w", migration.name, err)
		}
	}

	return nil
}

func TestIT_PlanetRepository_ListForPlayer(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())
	p1, _, _ := insertTestPlanetForPlayer(t, conn)
	p2, player1, _ := insertTestPlanetForPlayer(t, conn)
	p3 := insertTestPlanet(t, conn, player1.Id, addPlanetResource)
	p4 := insertTestPlanet(t, conn, player1.Id, addPlanetStorage)
	p5 := insertTestPlanet(t, conn, player1.Id, addPlanetProduction)
	p6 := insertTestPlanet(t, conn, player1.Id, addPlanetProductionForBuilding)
	p7 := insertTestPlanet(t, conn, player1.Id, addPlanetBuilding)

	actual, err := repo.ListForPlayer(context.Background(), player1.Id)
	require.NoError(t, err, "Actual err: %v", err)

	assert.GreaterOrEqual(t, len(actual), 6)
	assert.Contains(t, actual, p2)
	assert.Contains(t, actual, p3)
	assert.Contains(t, actual, p4)
	assert.Contains(t, actual, p5)
	assert.Contains(t, actual, p6)
	assert.Contains(t, actual, p7)
	for _, planet := range actual {
		assert.NotEqual(t, planet.Id, p1.Id)
	}
}

func TestIT_PlanetRepository_Delete(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	t.Run("deletes planet", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
	})

	t.Run("deletes planet with resources", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with resources", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetResource)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with storages", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetStorage)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with productions", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetProduction)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with productions for building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with productions for building", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetBuildingDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes homeworld with building", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn, addPlanetBuilding)

		err := repo.Delete(context.Background(), planet.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
		assertPlanetBuildingDoesNotExist(t, conn, planet.Id)
	})

	t.Run("succeeds when the planet does not exist", func(t *testing.T) {
		nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

		err := repo.Delete(context.Background(), nonExistingId)
		require.NoError(t, err, "Actual err: %v", err)
	})
}

func TestIT_PlanetRepository_DeleteForPlayer(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	t.Run("deletes all planets", func(t *testing.T) {
		p1, player, _ := insertTestHomeworldPlanetForPlayer(t, conn)
		p2 := insertTestPlanet(t, conn, player.Id)

		err := repo.DeleteForPlayer(context.Background(), player.Id)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, p1.Id)
		assertPlanetIsNotHomeworld(t, conn, p1.Id)
		assertPlanetDoesNotExist(t, conn, p2.Id)
	})

	t.Run("deletes homeworld", func(t *testing.T) {
		planet, _, _ := insertTestHomeworldPlanetForPlayer(t, conn)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetIsNotHomeworld(t, conn, planet.Id)
	})

	t.Run("deletes planet with resources", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetResource)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetResourceDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with storages", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetStorage)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetStorageDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with productions", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProduction)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with productions for building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetProductionForBuilding)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetProductionDoesNotExist(t, conn, planet.Id)
	})

	t.Run("deletes planet with building", func(t *testing.T) {
		planet, _, _ := insertTestPlanetForPlayer(t, conn, addPlanetBuilding)

		err := repo.DeleteForPlayer(context.Background(), planet.Player)
		require.NoError(t, err, "Actual err: %v", err)

		assertPlanetDoesNotExist(t, conn, planet.Id)
		assertPlanetBuildingDoesNotExist(t, conn, planet.Id)
	})
}

func TestIT_PlanetRepository_CreationDeletionWorkflow(t *testing.T) {
	repo, conn := newTestPlanetRepository(t)
	defer conn.Close(context.Background())

	type testCase struct {
		name   string
		planet models.Planet
	}

	testCases := []testCase{
		{
			name: "planet",
			planet: models.Planet{
				Id:          uuid.New(),
				Name:        fmt.Sprintf("my-planet-%s", uuid.NewString()),
				Homeworld:   false,
				CreatedAt:   someTime,
				UpdatedAt:   someOtherTime,
				Version:     4,
				Resources:   []models.PlanetResource{},
				Storages:    []models.PlanetResourceStorage{},
				Productions: []models.PlanetResourceProduction{},
				Buildings:   []models.PlanetBuilding{},
			},
		},
		{
			name: "homeworld",
			planet: models.Planet{
				Id:          uuid.New(),
				Name:        fmt.Sprintf("my-homeworld-%s", uuid.NewString()),
				Homeworld:   true,
				CreatedAt:   someTime,
				UpdatedAt:   someOtherTime,
				Version:     6,
				Resources:   []models.PlanetResource{},
				Storages:    []models.PlanetResourceStorage{},
				Productions: []models.PlanetResourceProduction{},
				Buildings:   []models.PlanetBuilding{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			player, _ := insertTestPlayerInUniverse(t, conn)

			tc.planet.Player = player.Id

			func() {
				err := repo.Create(context.Background(), tc.planet)
				require.NoError(t, err, "Actual err: %v", err)
			}()

			func() {
				planetFromDb, err := repo.Get(context.Background(), tc.planet.Id)
				require.NoError(t, err, "Actual err: %v", err)

				assert.Equal(t, tc.planet, planetFromDb)
			}()

			func() {
				err := repo.Delete(context.Background(), tc.planet.Id)
				require.NoError(t, err, "Actual err: %v", err)
			}()

			assertPlanetDoesNotExist(t, conn, tc.planet.Id)
		})
	}
}

func newTestPlanetRepository(t *testing.T) (drivenports.ForManagingPlanets, db.Connection) {
	t.Helper()
	conn := newTestConnection(t)
	return NewPlanetRepository(conn), conn
}

func insertTestPlanet(
	t *testing.T,
	conn db.Connection,
	player uuid.UUID,
	modifiers ...func(*testing.T, db.Connection, *models.Planet),
) models.Planet {
	t.Helper()

	planet := models.Planet{
		Id:        uuid.New(),
		Player:    player,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
		Version:   7,
		// This is intentional: the details (e.g. resources, etc.) are returned as empty
		// slices by the adapter
		Resources:   []models.PlanetResource{},
		Storages:    []models.PlanetResourceStorage{},
		Productions: []models.PlanetResourceProduction{},
		Buildings:   []models.PlanetBuilding{},
	}

	sqlQuery := `INSERT INTO planet (id, player, name, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		planet.Id,
		planet.Player,
		planet.Name,
		planet.CreatedAt,
		planet.UpdatedAt,
		planet.Version,
	)
	require.NoError(t, err, "Actual err: %v", err)

	for _, modifier := range modifiers {
		modifier(t, conn, &planet)
	}

	return planet
}

func addPlanetHomeworld(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	sqlQuery := `INSERT INTO homeworld (player, planet) VALUES ($1, $2)`
	_, err := conn.Exec(context.Background(), sqlQuery, p.Player, p.Id)
	require.NoError(t, err, "Actual err: %v", err)

	p.Homeworld = true
}

func addPlanetResource(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	resource := models.PlanetResource{
		Resource: crystalResourceId,
		// Amount is stored with 5 decimals in the DB
		Amount:    randFloat(5),
		CreatedAt: someTime,
		UpdatedAt: someOtherTime,
	}

	sqlQuery := `INSERT INTO planet_resource (planet, resource, amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		resource.Resource,
		resource.Amount,
		resource.CreatedAt,
		resource.UpdatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Resources = append(p.Resources, resource)
}

func addPlanetStorage(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	storage := models.PlanetResourceStorage{
		Resource:  crystalResourceId,
		Storage:   6233,
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	sqlQuery := `INSERT INTO planet_resource_storage (planet, resource, storage, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		storage.Resource,
		storage.Storage,
		storage.CreatedAt,
		storage.UpdatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Storages = append(p.Storages, storage)
}

func addPlanetProductionForBuilding(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	production := models.PlanetResourceProduction{
		Building:   &crystalMineId,
		Resource:   metalResourceId,
		Production: rand.Intn(784152),
		CreatedAt:  someTime,
		UpdatedAt:  someOtherTime,
	}

	sqlQuery := `INSERT INTO planet_resource_production
		(planet, building, resource, production, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		production.Building,
		production.Resource,
		production.Production,
		production.CreatedAt,
		production.UpdatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Productions = append(p.Productions, production)
}

func addPlanetProduction(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	production := models.PlanetResourceProduction{
		Building:   nil,
		Resource:   metalResourceId,
		Production: rand.Intn(6589),
		CreatedAt:  someTime,
		UpdatedAt:  someOtherTime,
	}

	sqlQuery := `INSERT INTO planet_resource_production
		(planet, building, resource, production, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		production.Building,
		production.Resource,
		production.Production,
		production.CreatedAt,
		production.UpdatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Productions = append(p.Productions, production)
}

func addPlanetBuilding(t *testing.T, conn db.Connection, p *models.Planet) {
	t.Helper()

	building := models.PlanetBuilding{
		Building:  metalStorageId,
		Level:     0,
		CreatedAt: someTime,
		UpdatedAt: someTime,
	}

	sqlQuery := `INSERT INTO planet_building (planet, building, level, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(
		context.Background(),
		sqlQuery,
		p.Id,
		building.Building,
		building.Level,
		building.CreatedAt,
		building.UpdatedAt,
	)
	require.NoError(t, err, "Actual err: %v", err)

	p.Buildings = append(p.Buildings, building)
}

func insertTestPlanetForPlayer(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Planet),
) (models.Planet, models.Player, models.Universe) {
	t.Helper()

	player, universe := insertTestPlayerInUniverse(t, conn)
	planet := insertTestPlanet(t, conn, player.Id, modifiers...)
	return planet, player, universe
}

func insertTestHomeworldPlanetForPlayer(
	t *testing.T,
	conn db.Connection,
	modifiers ...func(*testing.T, db.Connection, *models.Planet),
) (models.Planet, models.Player, models.Universe) {
	t.Helper()

	player, universe := insertTestPlayerInUniverse(t, conn)
	modifiers = append(modifiers, addPlanetHomeworld)
	planet := insertTestPlanet(t, conn, player.Id, modifiers...)
	return planet, player, universe
}

func assertPlanetExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT id FROM planet WHERE id = $1`
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, id, value)
}

func assertPlanetDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(id) FROM planet WHERE id = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, id)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetIsHomeworld(t *testing.T, conn db.Connection, planet uuid.UUID, player uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1 AND player = $2`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet, player)
	require.NoError(t, err, "Actual err: %v", err)
	require.Equal(t, 1, value)
}

func assertPlanetIsNotHomeworld(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(*) FROM homeworld WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetResourceDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(resource) FROM planet_resource WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetStorageDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	t.Helper()

	sqlQuery := `SELECT COUNT(resource) FROM planet_resource_storage WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetProductionDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(resource) FROM planet_resource_production WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}

func assertPlanetBuildingDoesNotExist(t *testing.T, conn db.Connection, planet uuid.UUID) {
	sqlQuery := `SELECT COUNT(building) FROM planet_building WHERE planet = $1`
	value, err := db.QueryOne[int](context.Background(), conn, sqlQuery, planet)
	require.NoError(t, err, "Actual err: %v", err)
	require.Zero(t, value)
}
