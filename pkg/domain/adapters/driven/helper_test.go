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
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	someTime      = time.Date(2024, time.November, 29, 17, 53, 29, 0, time.UTC)
	someOtherTime = time.Date(2026, time.June, 1, 8, 20, 15, 0, time.UTC)

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
		t.Context(),
		testContainerImage,
		postgres.WithDatabase(testTemplateDatabase),
		postgres.WithUsername(testDatabaseUser),
		postgres.WithPassword(testDatabasePassword),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err, "Actual err: %v", err)

	return postgresContainer
}

func newTestConnection(t *testing.T) db.Connection {
	t.Helper()

	sharedTestContainerSuite.ensureInitialized(t)

	testDatabaseName, bootstrapConn, postgresContainer := sharedTestContainerSuite.nextDatabaseContext()

	_, err := bootstrapConn.Exec(
		t.Context(),
		fmt.Sprintf(
			`CREATE DATABASE %s WITH TEMPLATE %s OWNER %s`,
			testDatabaseName,
			testTemplateDatabase,
			testDatabaseUser,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = bootstrapConn.Exec(
		t.Context(),
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
		conn.Close(t.Context())

		_, dropErr := bootstrapConn.Exec(
			context.Background(),
			fmt.Sprintf(`DROP DATABASE IF EXISTS %s WITH (FORCE)`, testDatabaseName),
		)
		require.NoError(t, dropErr, "Actual err: %v", dropErr)
	})

	return conn
}

func (s *testContainerSuite) createConnection(
	t *testing.T,
	postgresContainer *postgres.PostgresContainer,
	database string,
	user string,
	password string,
) db.Connection {
	t.Helper()

	host, err := postgresContainer.Host(t.Context())
	require.NoError(t, err, "Actual err: %v", err)

	port, err := postgresContainer.MappedPort(t.Context(), "5432/tcp")
	require.NoError(t, err, "Actual err: %v", err)

	portValue, err := strconv.ParseUint(port.Port(), 10, 16)
	require.NoError(t, err, "Actual err: %v", err)

	conn, err := db.New(t.Context(), postgresql.Config{
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
	defer templateConn.Close(t.Context())

	_, err := templateConn.Exec(
		t.Context(),
		fmt.Sprintf(`CREATE SCHEMA %s AUTHORIZATION %s`, testDatabaseSchema, testDatabaseUser),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = templateConn.Exec(
		t.Context(),
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

	s.bootstrapConn.Close(context.Background())

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

	host, err := postgresContainer.Host(t.Context())
	require.NoError(t, err, "Actual err: %v", err)
	port, err := postgresContainer.MappedPort(t.Context(), "5432/tcp")
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

func randFloat(t *testing.T, min float64, max float64, precision int) float64 {
	t.Helper()

	scale := math.Pow(10, float64(precision))
	minScaled := int64(math.Ceil(min * scale))
	maxScaled := int64(math.Floor(max * scale))
	if minScaled > maxScaled {
		t.Fatalf(
			"no representable values in range [%f, %f] for precision=%d",
			min,
			max,
			precision,
		)
	}

	valueScaled := minScaled + rand.Int64N(maxScaled-minScaled+1)
	return float64(valueScaled) / scale
}

func assertEqualIgnoringFields[T any](
	t *testing.T,
	actual T,
	expected T,
	ignoredFields ...string,
) {
	t.Helper()

	equal := eassert.EqualsIgnoringFields(actual, expected, ignoredFields...)
	assert.True(t, equal, "Expected actual=%+v and expected=%+v to be equal", actual, expected)
}

func assertContainsIgnoringFields[T any](
	t *testing.T,
	collection []T,
	expected T,
	ignoredFields ...string,
) {
	t.Helper()

	equal := eassert.ContainsIgnoringFields(collection, expected, ignoredFields...)
	assert.True(t, equal, "Expected collection=%+v to contain expected=%+v", collection, expected)
}
