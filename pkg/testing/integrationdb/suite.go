package integrationdb

import (
	"context"
	"fmt"
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

const (
	containerImage   = "postgres:18-alpine"
	databaseUser     = "galactic_sovereign_manager"
	databasePassword = "manager_password"
	databaseSchema   = "galactic_sovereign_schema"
	templateDatabase = "db_galactic_sovereign_template"
)

// Suite manages a shared Postgres testcontainer lifecycle for DB-backed integration
// tests. A single container is started on demand and shared across all tests in a
// test binary. Each call to NewTestConnection clones a fresh database from the
// migrated template, giving each test full isolation with minimal overhead.
type Suite struct {
	stateLock           sync.Mutex
	initialized         bool
	container           *postgres.PostgresContainer
	bootstrapConn       db.Connection
	testDatabaseCounter int
}

// NewDatabaseSharedContainer creates a new Suite that can be shared across all
// tests in a binary. Typically stored in a package-level variable and passed to
// TestMain for teardown.
func NewDatabaseSharedContainer() *Suite {
	return &Suite{}
}

// NewTestConnection returns a fresh DB connection backed by the shared test
// container. A per-test database is cloned from the migrated template and is
// automatically dropped when the test finishes.
func (s *Suite) NewTestConnection(t *testing.T) db.Connection {
	t.Helper()

	s.ensureInitialized(t)

	testDatabaseName, bootstrapConn, postgresContainer := s.nextDatabaseContext()

	_, err := bootstrapConn.Exec(
		t.Context(),
		fmt.Sprintf(
			`CREATE DATABASE %s WITH TEMPLATE %s OWNER %s`,
			testDatabaseName,
			templateDatabase,
			databaseUser,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = bootstrapConn.Exec(
		t.Context(),
		fmt.Sprintf(
			`ALTER DATABASE %s SET search_path = %s`,
			testDatabaseName,
			databaseSchema,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	conn := s.createConnection(
		t,
		postgresContainer,
		testDatabaseName,
		databaseUser,
		databasePassword,
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

// Teardown shuts down the shared test container. Must be called from TestMain
// after m.Run() to ensure proper container cleanup.
func (s *Suite) Teardown() {
	s.teardown()
}

func (s *Suite) createConnection(
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

func (s *Suite) nextDatabaseContext() (string, db.Connection, *postgres.PostgresContainer) {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()

	s.testDatabaseCounter++

	return fmt.Sprintf("db_galactic_sovereign_test_%d", s.testDatabaseCounter), s.bootstrapConn, s.container
}

func (s *Suite) ensureInitialized(t *testing.T) {
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
		databaseUser,
		databasePassword,
	)

	templateConn := s.createConnection(
		t,
		postgresContainer,
		templateDatabase,
		databaseUser,
		databasePassword,
	)
	defer templateConn.Close(t.Context())

	_, err := templateConn.Exec(
		t.Context(),
		fmt.Sprintf(`CREATE SCHEMA %s AUTHORIZATION %s`, databaseSchema, databaseUser),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = templateConn.Exec(
		t.Context(),
		fmt.Sprintf(
			`ALTER DATABASE %s SET search_path = %s`,
			templateDatabase,
			databaseSchema,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	s.runMigrations(t, postgresContainer, templateDatabase)

	s.container = postgresContainer
	s.bootstrapConn = bootstrapConn
	s.initialized = true
}

func (s *Suite) teardown() {
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

func (s *Suite) runMigrations(
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

func (s *Suite) databaseURL(
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
		databaseUser,
		databasePassword,
		host,
		port.Port(),
		database,
	)
}

func (s *Suite) migrationSourceURL(t *testing.T) string {
	_ = s
	t.Helper()

	_, currentFilePath, _, ok := runtime.Caller(0)
	require.True(t, ok)

	// This file lives at pkg/testing/integrationdb/ — 3 levels up reaches the
	// module root, one more component reaches the migrations directory.
	migrationsPath := filepath.Join(filepath.Dir(currentFilePath), "../../../database/galactic-sovereign/migrations")
	return fmt.Sprintf("file://%s", migrationsPath)
}

func createTestContainer(t *testing.T) *postgres.PostgresContainer {
	t.Helper()

	postgresContainer, err := postgres.Run(
		t.Context(),
		containerImage,
		postgres.WithDatabase(templateDatabase),
		postgres.WithUsername(databaseUser),
		postgres.WithPassword(databasePassword),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err, "Actual err: %v", err)

	return postgresContainer
}
