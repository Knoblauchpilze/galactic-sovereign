package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/adapters/driven/database"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/domain/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	testServerHost = "localhost"
)

var (
	oberonUniverseId = uuid.MustParse("9682f17b-f5f0-4eda-a747-2537d2151837")

	metalMineId = uuid.MustParse("d176e82d-f2ca-4611-996b-c4804096caef")
	shipyardId  = uuid.MustParse("58d75842-6dc0-4ac0-b36d-55f91b8d060d")

	lightFighterId = uuid.MustParse("a31de13b-5905-4468-99c5-d1d1e529b36e")
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

// urlFor builds a URL under the test server's base path.
// Segments are joined with '/' and appended after the base path.
func urlFor(baseUrl string, segments ...string) string {
	path := ""
	for _, s := range segments {
		path += "/" + s
	}
	return fmt.Sprintf("%s%s", baseUrl, path)
}

func newTestServerConfig() server.Config {
	return server.Config{
		BasePath:        "/v1/galactic-sovereign",
		ShutdownTimeout: 500 * time.Millisecond,
	}
}

// asyncStartServer binds the server on an OS-assigned port, serves until test
// cleanup and returns the base URL to reach the server.
func asyncStartServer(t *testing.T, s HttpServer) string {
	t.Helper()

	listener, err := s.Bind(0)
	require.NoError(t, err, "Actual err: %v", err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- s.Serve(ctx, listener)
	}()

	t.Cleanup(func() {
		cancel()
		err := <-done
		require.NoError(t, err, "Actual err: %v", err)
	})

	defaultConfig := newTestServerConfig()
	port := uint16(listener.Addr().(*net.TCPAddr).Port)
	return fmt.Sprintf(
		"http://%s:%d%s",
		testServerHost,
		port,
		defaultConfig.BasePath,
	)
}

func doGet[T any](t *testing.T, url string) T {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err, "GET %s: %v", url, err)
	defer resp.Body.Close() // nolint:errcheck
	require.Equal(t, http.StatusOK, resp.StatusCode, "GET %s returned %d", url, resp.StatusCode)

	return decodeResponseBody[T](t, resp.Body)
}

func doPost[T any](t *testing.T, url string, body any) T {
	t.Helper()

	var payload io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		require.NoError(t, err, "Actual err: %v", err)

		payload = bytes.NewReader(raw)
	}

	resp, err := http.Post(url, "application/json", payload) // nolint:noctx
	require.NoError(t, err, "POST %s: %v", url, err)
	defer resp.Body.Close() // nolint:errcheck
	require.Equal(t, http.StatusCreated, resp.StatusCode, "POST %s returned %d", url, resp.StatusCode)

	return decodeResponseBody[T](t, resp.Body)
}

func doDelete(t *testing.T, url string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err, "Actual err: %v", err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "DELETE %s: %v", url, err)
	defer resp.Body.Close() // nolint:errcheck
	require.Equal(t, http.StatusNoContent, resp.StatusCode, "DELETE %s returned %d", url, resp.StatusCode)
}

func decodeResponseBody[T any](t *testing.T, body io.ReadCloser) T {
	t.Helper()

	raw, err := io.ReadAll(body)
	require.NoError(t, err, "Actual err: %v", err)

	var envelope rest.ResponseEnvelope[T]
	err = json.Unmarshal(raw, &envelope)
	require.NoError(t, err, "Actual err: %v", err)

	return envelope.Details
}

// addPlanetResources credits the given resource amount to the planet, allowing
// tests to afford actions (e.g. ships) that starting resources can't cover.
func addPlanetResources(
	t *testing.T,
	conn database.Connection,
	planet uuid.UUID,
	resource uuid.UUID,
	amount int,
) {
	t.Helper()

	_, err := conn.Exec(
		t.Context(),
		`UPDATE planet_resource SET amount = amount + $1 WHERE planet = $2 AND resource = $3`,
		amount,
		planet,
		resource,
	)
	require.NoError(t, err, "Actual err: %v", err)
}

// bumpPlanetBuilding bumps the given building to the desired level on the
// planet, allowing tests to perform actions (e.g. ship creation) that the
// starting planet would not allow.
func bumpPlanetBuilding(
	t *testing.T,
	conn database.Connection,
	planet uuid.UUID,
	building uuid.UUID,
	level int,
) {
	t.Helper()

	_, err := conn.Exec(
		t.Context(),
		`UPDATE planet_building SET level = $1 WHERE planet = $2 AND building = $3`,
		level,
		planet,
		building,
	)
	require.NoError(t, err, "Actual err: %v", err)
}

// bumpPlanetShip sets the count of the given ship on a planet to the desired
// value, allowing tests to perform actions (e.g. fleet creation) that the
// starting planet would not allow.
func bumpPlanetShip(
	t *testing.T,
	conn database.Connection,
	planet uuid.UUID,
	ship uuid.UUID,
	count int,
) {
	t.Helper()

	_, err := conn.Exec(
		t.Context(),
		`INSERT INTO planet_ship (planet, ship, count)
			VALUES ($1, $2, $3)
		ON CONFLICT (planet, ship) DO UPDATE SET
			count = excluded.count`,
		planet,
		ship,
		count,
	)
	require.NoError(t, err, "Actual err: %v", err)
}

// insertTestPlanet allows to create a new planet for a player with random coordinates
// Coordinates might collide with existing planets but this should be rare. If it's too
// frequent, a more reliable system to pick coordinates should be implemented.
func insertTestPlanet(
	t *testing.T,
	conn database.Connection,
	player uuid.UUID,
) models.Planet {
	t.Helper()

	planet := models.Planet{
		Id:        uuid.New(),
		Player:    player,
		Name:      fmt.Sprintf("my-planet-%s", uuid.NewString()),
		Homeworld: false,
		Coordinate: models.Coordinate{
			Galaxy:      rand.Intn(3),
			SolarSystem: rand.Intn(267),
			Position:    rand.Intn(7),
		},
		Fields:      1 + rand.Intn(211),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     0,
		Resources:   []models.PlanetResource{},
		Storages:    []models.PlanetResourceStorage{},
		Productions: []models.PlanetResourceProduction{},
		Buildings:   []models.PlanetBuilding{},
		Ships:       []models.PlanetShip{},
		ShipActions: []models.ShipAction{},
	}

	sqlQuery := `INSERT INTO planet (id, player, name, fields, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(
		t.Context(),
		sqlQuery,
		planet.Id,
		planet.Player,
		planet.Name,
		planet.Fields,
		planet.CreatedAt,
		planet.UpdatedAt,
		planet.Version,
	)
	require.NoError(t, err, "Actual err: %v", err)

	sqlQuery = `INSERT INTO planet_coordinate (planet, universe, galaxy, solar_system, position)
		SELECT
			$1,
			universe,
			$2,
			$3,
			$4
		FROM
			player
		WHERE
			id = $5
		ON CONFLICT (planet) DO UPDATE SET
			galaxy = excluded.galaxy,
			solar_system = excluded.solar_system,
			position = excluded.position`
	_, err = conn.Exec(
		t.Context(),
		sqlQuery,
		planet.Id,
		planet.Coordinate.Galaxy,
		planet.Coordinate.SolarSystem,
		planet.Coordinate.Position,
		planet.Player,
	)
	require.NoError(t, err, "Actual err: %v", err)

	return planet
}
