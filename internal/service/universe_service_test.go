package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	eassert "github.com/Knoblauchpilze/easy-assert/assert"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_UniverseService_Create(t *testing.T) {
	id := uuid.New()
	universeDtoRequest := communication.UniverseDtoRequest{
		Name: fmt.Sprintf("my-universe-%s", id),
	}
	beforeInsertion := time.Now()

	service, conn := newTestUniverseService(t)
	out, err := service.Create(context.Background(), universeDtoRequest)

	assert.Nil(t, err)

	assert.Equal(t, universeDtoRequest.Name, out.Name)
	eassert.AreTimeCloserThan(beforeInsertion, out.CreatedAt, 1*time.Second)
	assertUniverseExists(t, conn, out.Id)
}

func TestIT_UniverseService_Create_WhenNameAlreadyExists_ExpectFailure(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)
	universeDtoRequest := communication.UniverseDtoRequest{
		Name: universe.Name,
	}

	_, err := service.Create(context.Background(), universeDtoRequest)

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestIT_UniverseService_Get_WithNoDataRegistered(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)

	actual, err := service.Get(context.Background(), universe.Id)

	assert.Nil(t, err)
	assert.Equal(t, universe.Id, actual.Id)
	assert.Equal(t, universe.Name, actual.Name)
	assert.True(t, eassert.AreTimeCloserThan(universe.CreatedAt, actual.CreatedAt, 1*time.Second))
}

func TestIT_UniverseService_Get_WithResources(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)
	resource := insertTestResource(t, conn)

	actual, err := service.Get(context.Background(), universe.Id)

	assert.Nil(t, err)
	assert.Equal(t, universe.Id, actual.Id)

	expected := communication.ResourceDtoResponse{
		Id:   resource.Id,
		Name: resource.Name,
	}
	assert.True(t, eassert.ContainsIgnoringFields(actual.Resources, expected, "CreatedAt"))
}

func TestIT_UniverseService_Get_WithBuildings(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)
	building := insertTestBuilding(t, conn)

	actual, err := service.Get(context.Background(), universe.Id)

	assert.Nil(t, err)
	assert.Equal(t, universe.Id, actual.Id)

	expected := communication.FullBuildingDtoResponse{
		BuildingDtoResponse: communication.BuildingDtoResponse{
			Id:   building.Id,
			Name: building.Name,
		},
	}
	assert.True(t, eassert.ContainsIgnoringFields(actual.Buildings, expected, "CreatedAt"))
}

func TestIT_UniverseService_Get_WhenBuildingHasCosts_ExpectThemToBeReturned(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)
	building := insertTestBuilding(t, conn)
	cost, resource := insertTestBuildingCost(t, conn, building.Id)

	actual, err := service.Get(context.Background(), universe.Id)

	assert.Nil(t, err)
	assert.Equal(t, universe.Id, actual.Id)

	expected := communication.FullBuildingDtoResponse{
		BuildingDtoResponse: communication.BuildingDtoResponse{
			Id:   building.Id,
			Name: building.Name,
		},
		Costs: []communication.BuildingCostDtoResponse{
			{
				Building: building.Id,
				Resource: resource.Id,
				Cost:     cost.Cost,
				Progress: cost.Progress,
			},
		},
	}
	assert.True(t, eassert.ContainsIgnoringFields(actual.Buildings, expected, "CreatedAt"))
}

func TestIT_UniverseService_Get_WhenBuildingHasProductions_ExpectThemToBeReturned(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)
	building := insertTestBuilding(t, conn)
	prod, resource := insertTestBuildingResourceProduction(t, conn, building.Id)

	actual, err := service.Get(context.Background(), universe.Id)

	assert.Nil(t, err)
	assert.Equal(t, universe.Id, actual.Id)

	expected := communication.FullBuildingDtoResponse{
		BuildingDtoResponse: communication.BuildingDtoResponse{
			Id:   building.Id,
			Name: building.Name,
		},
		Productions: []communication.BuildingResourceProductionDtoResponse{
			{
				Building: building.Id,
				Resource: resource.Id,
				Base:     prod.Base,
				Progress: prod.Progress,
			},
		},
	}
	assert.True(t, eassert.ContainsIgnoringFields(actual.Buildings, expected, "CreatedAt"))
}

func TestIT_UniverseService_Get_WhenBuildingHasStorages_ExpectThemToBeReturned(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)
	building := insertTestBuilding(t, conn)
	storage, resource := insertTestBuildingResourceStorage(t, conn, building.Id)

	actual, err := service.Get(context.Background(), universe.Id)

	assert.Nil(t, err)
	assert.Equal(t, universe.Id, actual.Id)

	expected := communication.FullBuildingDtoResponse{
		BuildingDtoResponse: communication.BuildingDtoResponse{
			Id:   building.Id,
			Name: building.Name,
		},
		Storages: []communication.BuildingResourceStorageDtoResponse{
			{
				Building: building.Id,
				Resource: resource.Id,
				Base:     storage.Base,
				Scale:    storage.Scale,
				Progress: storage.Progress,
			},
		},
	}
	assert.True(t, eassert.ContainsIgnoringFields(actual.Buildings, expected, "CreatedAt"))
}

func TestIT_UniverseService_List(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)

	out, err := service.List(context.Background())

	assert.Nil(t, err)
	expected := communication.UniverseDtoResponse{
		Id:   universe.Id,
		Name: universe.Name,
	}
	assert.True(t, eassert.ContainsIgnoringFields(out, expected, "CreatedAt"))
}

func TestIT_UniverseService_Delete(t *testing.T) {
	service, conn := newTestUniverseService(t)
	universe := insertTestUniverse(t, conn)

	err := service.Delete(context.Background(), universe.Id)

	assert.Nil(t, err)
	assertUniverseDoesNotExist(t, conn, universe.Id)
}

func TestIT_UniverseService_Delete_WhenUniverseDoesNotExist_ExpectSuccess(t *testing.T) {
	nonExistingId := uuid.MustParse("00000000-0000-1221-0000-000000000000")

	service, _ := newTestUniverseService(t)
	err := service.Delete(context.Background(), nonExistingId)

	assert.Nil(t, err)
}

func newTestUniverseService(t *testing.T) (UniverseService, db.Connection) {
	conn := newTestConnection(t)

	repos := repositories.Repositories{
		Universe:                   repositories.NewUniverseRepository(conn),
		Resource:                   repositories.NewResourceRepository(),
		Building:                   repositories.NewBuildingRepository(),
		BuildingCost:               repositories.NewBuildingCostRepository(),
		BuildingResourceProduction: repositories.NewBuildingResourceProductionRepository(),
		BuildingResourceStorage:    repositories.NewBuildingResourceStorageRepository(),
	}

	return NewUniverseService(conn, repos), conn
}

func assertUniverseExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, "SELECT id FROM universe WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertUniverseDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[int](context.Background(), conn, "SELECT COUNT(id) FROM universe WHERE id = $1", id)
	require.Nil(t, err)
	require.Zero(t, value)
}

func insertTestUniverse(t *testing.T, conn db.Connection) persistence.Universe {
	someTime := time.Date(2024, 12, 8, 10, 10, 22, 0, time.UTC)

	universe := persistence.Universe{
		Id:        uuid.New(),
		Name:      fmt.Sprintf("my-universe-%s", uuid.NewString()),
		CreatedAt: someTime,
	}

	sqlQuery := `INSERT INTO universe (id, name, created_at) VALUES ($1, $2, $3) RETURNING updated_at`
	updatedAt, err := db.QueryOne[time.Time](
		context.Background(),
		conn,
		sqlQuery,
		universe.Id,
		universe.Name,
		universe.CreatedAt,
	)
	require.Nil(t, err)

	universe.UpdatedAt = updatedAt

	return universe
}

func assertResourceRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockResourceRepository {
	m, ok := repos.Resource.(*mockResourceRepository)
	if !ok {
		assert.Fail("Provided resource repository is not a mock")
	}
	return m
}
