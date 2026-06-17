package drivenadapters

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
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

	testBootstrapDatabase = "postgres"
	testBootstrapUser     = "postgres"
	testBootstrapPassword = "postgres"

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
		postgres.WithDatabase(testBootstrapDatabase),
		postgres.WithUsername(testBootstrapUser),
		postgres.WithPassword(testBootstrapPassword),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err, "Actual err: %v", err)

	return postgresContainer
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
			`ALTER ROLE %s IN DATABASE %s SET search_path = %s`,
			testDatabaseUser,
			testDatabaseName,
			testDatabaseSchema,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	conn := newConnectionForTestContainer(
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

	bootstrapConn := newConnectionForTestContainer(
		t,
		postgresContainer,
		testBootstrapDatabase,
		testBootstrapUser,
		testBootstrapPassword,
	)

	_, err := bootstrapConn.Exec(
		ctx,
		fmt.Sprintf(`CREATE USER %s WITH PASSWORD '%s'`, testDatabaseUser, testDatabasePassword),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = bootstrapConn.Exec(
		ctx,
		fmt.Sprintf(`CREATE DATABASE %s OWNER %s`, testTemplateDatabase, testDatabaseUser),
	)
	require.NoError(t, err, "Actual err: %v", err)

	templateBootstrapConn := newConnectionForTestContainer(
		t,
		postgresContainer,
		testTemplateDatabase,
		testBootstrapUser,
		testBootstrapPassword,
	)
	defer templateBootstrapConn.Close(ctx)

	_, err = templateBootstrapConn.Exec(
		ctx,
		fmt.Sprintf(`CREATE SCHEMA %s AUTHORIZATION %s`, testDatabaseSchema, testDatabaseUser),
	)
	require.NoError(t, err, "Actual err: %v", err)

	_, err = templateBootstrapConn.Exec(
		ctx,
		fmt.Sprintf(
			`ALTER ROLE %s IN DATABASE %s SET search_path = %s`,
			testDatabaseUser,
			testTemplateDatabase,
			testDatabaseSchema,
		),
	)
	require.NoError(t, err, "Actual err: %v", err)

	templateConn := newConnectionForTestContainer(
		t,
		postgresContainer,
		testTemplateDatabase,
		testDatabaseUser,
		testDatabasePassword,
	)

	for _, migrationPath := range testMigrationPaths(t) {
		migration, readErr := os.ReadFile(migrationPath)
		require.NoError(t, readErr, "Actual err: %v", readErr)

		_, err = templateConn.Exec(ctx, string(migration))
		require.NoError(t, err, "Actual err: %v", err)
	}

	templateConn.Close(ctx)

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

func newConnectionForTestContainer(
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

func testMigrationPaths(t *testing.T) []string {
	t.Helper()

	_, currentFilePath, _, ok := runtime.Caller(0)
	require.True(t, ok)

	migrationPaths, err := filepath.Glob(
		filepath.Join(filepath.Dir(currentFilePath), "../../../../database/galactic-sovereign/migrations/*.up.sql"),
	)
	require.NoError(t, err, "Actual err: %v", err)

	slices.SortStableFunc(migrationPaths, func(left string, right string) int {
		leftPrefix, _ := migrationFilePrefix(left)
		rightPrefix, _ := migrationFilePrefix(right)
		return leftPrefix - rightPrefix
	})

	return migrationPaths
}

func migrationFilePrefix(path string) (int, error) {
	prefix := strings.SplitN(filepath.Base(path), "_", 2)[0]
	return strconv.Atoi(prefix)
}

func randFloat(precision int) float64 {
	rounder := math.Pow(10, float64(precision))
	return math.Round(rand.Float64()*rounder) / rounder
}
