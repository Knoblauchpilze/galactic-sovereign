package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var defaultPlanetId = uuid.MustParse("5b0efd85-8817-4454-b8f3-7af5d93253a1")
var defaultPlanetName = "my-planet"

var defaultPlanetDtoRequest = communication.PlanetDtoRequest{
	Player: defaultPlayerId,
	Name:   defaultPlanetName,
}
var defaultPlanet = persistence.Planet{
	Id:     defaultPlanetId,
	Player: defaultPlayerId,
	Name:   defaultPlanetName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func TestPlanetService_Create_CallsRepositoryCreate(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	s.Create(context.Background(), defaultPlanetDtoRequest)

	assert.Equal(1, mur.createCalled)
	assert.Equal(defaultPlanetDtoRequest.Player, mur.createdPlanet.Player)
	assert.Equal(defaultPlanetDtoRequest.Name, mur.createdPlanet.Name)
}

func TestPlanetService_Create_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	_, err := s.Create(context.Background(), defaultPlanetDtoRequest)

	assert.Equal(errDefault, err)
}

func TestPlanetService_Create_ReturnsCreatedPlanet(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{
		planet: defaultPlanet,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	actual, err := s.Create(context.Background(), defaultPlanetDtoRequest)

	assert.Nil(err)

	expected := communication.PlanetDtoResponse{
		Id:     defaultPlanet.Id,
		Player: defaultPlanet.Player,
		Name:   defaultPlanet.Name,

		CreatedAt: defaultPlanet.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestPlanetService_Get_CallsRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	s.Get(context.Background(), defaultPlanetId)

	assert.Equal(1, mur.getCalled)
}

func TestPlanetService_Get_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	_, err := s.Get(context.Background(), defaultPlanetId)

	assert.Equal(errDefault, err)
}

func TestPlanetService_Get_ReturnsPlanet(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{
		planet: defaultPlanet,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	actual, err := s.Get(context.Background(), defaultPlanetId)

	assert.Nil(err)
	assert.Equal(defaultPlanetId, mur.getId)

	expected := communication.PlanetDtoResponse{
		Id:     defaultPlanet.Id,
		Player: defaultPlanet.Player,
		Name:   defaultPlanet.Name,

		CreatedAt: defaultPlanet.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestPlanetService_List_CallsRepositoryList(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	s.List(context.Background())

	assert.Equal(1, mur.listCalled)
}

func TestPlanetService_List_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	_, err := s.List(context.Background())

	assert.Equal(errDefault, err)
}

func TestPlanetService_List_ReturnsAllPlanets(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{
		planet: defaultPlanet,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	actual, err := s.List(context.Background())

	assert.Nil(err)
	expected := []communication.PlanetDtoResponse{
		{
			Id:        defaultPlanet.Id,
			Player:    defaultPlanet.Player,
			Name:      defaultPlanet.Name,
			CreatedAt: defaultPlanet.CreatedAt,
		},
	}
	assert.Equal(expected, actual)
}

func TestPlanetService_Delete_CallsRepositoryDelete(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	s.Delete(context.Background(), defaultPlanetId)

	assert.Equal(1, mur.deleteCalled)
}

func TestPlanetService_Delete_CallsTransactionClose(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	s.Delete(context.Background(), defaultPlanetId)

	assert.Equal(1, mc.tx.closeCalled)
}

func TestPlanetService_Delete_WhenCreatingTransactionFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{}
	mc := &mockConnectionPool{
		err: errDefault,
	}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	err := s.Delete(context.Background(), defaultPlanetId)

	assert.Equal(errDefault, err)
}

func TestPlanetService_Delete_DeletesTheRightPlanet(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	s.Delete(context.Background(), defaultPlanetId)

	assert.Equal(defaultPlanetId, mur.deleteId)
}

func TestPlanetService_Delete_WhenPlanetRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	err := s.Delete(context.Background(), defaultPlanetId)

	assert.Equal(errDefault, err)
}

func TestPlanetService_Delete_WhenRepositoriesSucceeds_ExpectSuccess(t *testing.T) {
	assert := assert.New(t)

	mur := &mockPlanetRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Planet: mur,
	}
	s := NewPlanetService(mc, repos)

	err := s.Delete(context.Background(), defaultPlanetId)

	assert.Nil(err)
}
