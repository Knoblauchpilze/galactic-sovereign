package drivenadapters

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	someTime      = time.Date(2024, 11, 29, 17, 53, 29, 0, time.UTC)
	someOtherTime = time.Date(2026, 6, 1, 8, 20, 15, 0, time.UTC)

	sharedTestContainerSuite = &testContainerSuite{}
)

const (
	testContainerImage = "postgres:18-alpine"

	testDatabaseUser     = "galactic_sovereign_manager"
	testDatabasePassword = "manager_password"
	testDatabaseSchema   = "galactic_sovereign_schema"
	testTemplateDatabase = "db_galactic_sovereign_template"
)

type testContainerSuite struct {
	stateLock           sync.Mutex
	initialized         bool
	container           *postgres.PostgresContainer
	bootstrapConn       db.Connection
	testDatabaseCounter int
}

func TestMain(m *testing.M) {
	code := m.Run()
	sharedTestContainerSuite.teardown()
	os.Exit(code)
}

func createTestContainer(t *testing.T) *postgres.PostgresContainer {
	t.Helper()

	postgresContainer, err := postgres.Run(
		context.Background(),
		testContainerImage,
		postgres.WithDatabase(testTemplateDatabase),
		postgres.WithUsername(testDatabaseUser),
		postgres.WithPassword(testDatabasePassword),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err, "Actual err: %v", err)

	return postgresContainer
}

func (s *testContainerSuite) createConnection(
	t *testing.T,
	postgresContainer *postgres.PostgresContainer,
	database string,
	user string,
	password string,
) db.Connection {
	t.Helper()

	host, err := postgresContainer.Host(context.Background())
	require.NoError(t, err, "Actual err: %v", err)

	port, err := postgresContainer.MappedPort(context.Background(), "5432/tcp")
	require.NoError(t, err, "Actual err: %v", err)

	portValue, err := strconv.ParseUint(port.Port(), 10, 16)
	require.NoError(t, err, "Actual err: %v", err)

	conn, err := db.New(context.Background(), postgresql.Config{
		Host:           host,
		Port:           uint16(portValue),
		Database:       database,
		User:           user,
		Password:       password,
		ConnectTimeout: 5 * time.Second,
	})
	require.NoError(t, err, "Actual err: %v", err)

	return conn
}

func newTestConnection(t *testing.T) db.Connection {
	t.Helper()

	ctx := context.Background()
	sharedTestContainerSuite.ensureInitialized(t)

	testDatabaseName, bootstrapConn, postgresContainer := sharedTestContainerSuite.nextDatabaseContext()

	_, err := bootstrapConn.Exec(
		ctx,
		fmt.Sprintf(
			`CREATE DATABASE %s WITH TEMPLATE %s OWNER %s`,
			testDatabaseName,
			testTemplateDatabase,
			testDatabaseUser,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = bootstrapConn.Exec(
		ctx,
		fmt.Sprintf(
			`ALTER DATABASE %s SET search_path = %s`,
			testDatabaseName,
			testDatabaseSchema,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	conn := sharedTestContainerSuite.createConnection(
		t,
		postgresContainer,
		testDatabaseName,
		testDatabaseUser,
		testDatabasePassword,
	)

	t.Cleanup(func() {
		conn.Close(ctx)

		_, dropErr := bootstrapConn.Exec(
			ctx,
			fmt.Sprintf(`DROP DATABASE IF EXISTS %s WITH (FORCE)`, testDatabaseName),
		)
		require.NoError(t, dropErr, "Actual err: %v", dropErr)
	})

	return conn
}

func (s *testContainerSuite) nextDatabaseContext() (string, db.Connection, *postgres.PostgresContainer) {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	s.testDatabaseCounter++

	return fmt.Sprintf("db_galactic_sovereign_test_%d", s.testDatabaseCounter), s.bootstrapConn, s.container
}

func (s *testContainerSuite) ensureInitialized(t *testing.T) {
	t.Helper()

	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	if s.initialized {
		return
	}

	ctx := context.Background()
	postgresContainer := createTestContainer(t)

	bootstrapConn := s.createConnection(
		t,
		postgresContainer,
		"postgres",
		testDatabaseUser,
		testDatabasePassword,
	)

	templateConn := s.createConnection(
		t,
		postgresContainer,
		testTemplateDatabase,
		testDatabaseUser,
		testDatabasePassword,
	)
	defer templateConn.Close(ctx)

	_, err := templateConn.Exec(
		ctx,
		fmt.Sprintf(`CREATE SCHEMA %s AUTHORIZATION %s`, testDatabaseSchema, testDatabaseUser),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = templateConn.Exec(
		ctx,
		fmt.Sprintf(
			`ALTER DATABASE %s SET search_path = %s`,
			testTemplateDatabase,
			testDatabaseSchema,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	s.runMigrations(t, postgresContainer, testTemplateDatabase)

	s.container = postgresContainer
	s.bootstrapConn = bootstrapConn
	s.initialized = true
}

func (s *testContainerSuite) teardown() {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	if !s.initialized {
		return
	}

	ctx := context.Background()
	s.bootstrapConn.Close(ctx)

	_ = testcontainers.TerminateContainer(s.container)

	s.initialized = false
	s.container = nil
	s.bootstrapConn = nil
	s.testDatabaseCounter = 0
}

func (s *testContainerSuite) runMigrations(
	t *testing.T,
	postgresContainer *postgres.PostgresContainer,
	database string,
) {
	t.Helper()

	databaseURL := s.databaseURL(t, postgresContainer, database)
	sourceURL := s.migrationSourceURL(t)

	migrationRunner, err := migrate.New(sourceURL, databaseURL)
	require.NoError(t, err, "Actual err: %v", err)

	err = migrationRunner.Up()
	if err != nil {
		require.ErrorIs(t, err, migrate.ErrNoChange, "Actual err: %v", err)
	}

	sourceCloseErr, databaseCloseErr := migrationRunner.Close()
	require.NoError(t, sourceCloseErr, "Actual err: %v", sourceCloseErr)
	require.NoError(t, databaseCloseErr, "Actual err: %v", databaseCloseErr)
}

func (s *testContainerSuite) databaseURL(
	t *testing.T,
	postgresContainer *postgres.PostgresContainer,
	database string,
) string {
	_ = s
	t.Helper()

	host, err := postgresContainer.Host(context.Background())
	require.NoError(t, err, "Actual err: %v", err)
	port, err := postgresContainer.MappedPort(context.Background(), "5432/tcp")
	require.NoError(t, err, "Actual err: %v", err)

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		testDatabaseUser,
		testDatabasePassword,
		host,
		port.Port(),
		database,
	)
}

func (s *testContainerSuite) migrationSourceURL(t *testing.T) string {
	_ = s
	t.Helper()

	_, currentFilePath, _, ok := runtime.Caller(0)
	require.True(t, ok)

	migrationsPath := filepath.Join(filepath.Dir(currentFilePath), "../../../../database/galactic-sovereign/migrations")
	return fmt.Sprintf("file://%s", migrationsPath)
}

func randFloat(precision int) float64 {
	rounder := math.Pow(10, float64(precision))
	return math.Round(rand.Float64()*rounder) / rounder
}
