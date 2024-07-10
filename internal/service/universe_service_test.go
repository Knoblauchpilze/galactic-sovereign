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

var defaultUniverseId = uuid.MustParse("3e7fde5c-ac70-4e5d-bd09-73029725048d")
var defaultUniverseName = "my-universe"

var defaultUniverseDtoRequest = communication.UniverseDtoRequest{
	Name: defaultUniverseName,
}
var defaultUniverse = persistence.Universe{
	Id:   defaultUniverseId,
	Name: defaultUniverseName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func TestUniverseService_Create_CallsRepositoryCreate(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	s.Create(context.Background(), defaultUniverseDtoRequest)

	assert.Equal(1, mur.createCalled)
	assert.Equal(defaultUniverseDtoRequest.Name, mur.createdUniverse.Name)
}

func TestUniverseService_Create_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	_, err := s.Create(context.Background(), defaultUniverseDtoRequest)

	assert.Equal(errDefault, err)
}

func TestUniverseService_Create_ReturnsCreatedUniverse(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{
		universe: defaultUniverse,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	actual, err := s.Create(context.Background(), defaultUniverseDtoRequest)

	assert.Nil(err)

	expected := communication.UniverseDtoResponse{
		Id:   defaultUniverse.Id,
		Name: defaultUniverse.Name,

		CreatedAt: defaultUniverse.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUniverseService_Get_CallsRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	s.Get(context.Background(), defaultUniverseId)

	assert.Equal(1, mur.getCalled)
}

func TestUniverseService_Get_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	_, err := s.Get(context.Background(), defaultUniverseId)

	assert.Equal(errDefault, err)
}

func TestUniverseService_Get_ReturnsUniverse(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{
		universe: defaultUniverse,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	actual, err := s.Get(context.Background(), defaultUniverseId)

	assert.Nil(err)
	assert.Equal(defaultUniverseId, mur.getId)

	expected := communication.UniverseDtoResponse{
		Id:   defaultUniverse.Id,
		Name: defaultUniverse.Name,

		CreatedAt: defaultUniverse.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestUniverseService_List_CallsRepositoryList(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	s.List(context.Background())

	assert.Equal(1, mur.listCalled)
}

func TestUniverseService_List_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	_, err := s.List(context.Background())

	assert.Equal(errDefault, err)
}

func TestUniverseService_List_ReturnsAllUniverses(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{
		universe: defaultUniverse,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	actual, err := s.List(context.Background())

	assert.Nil(err)
	expected := []communication.UniverseDtoResponse{
		{
			Id:        defaultUniverse.Id,
			Name:      defaultUniverse.Name,
			CreatedAt: defaultUniverse.CreatedAt,
		},
	}
	assert.Equal(expected, actual)
}

func TestUniverseService_Delete_CallsRepositoryDelete(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	s.Delete(context.Background(), defaultUniverseId)

	assert.Equal(1, mur.deleteCalled)
}

func TestUniverseService_Delete_CallsTransactionClose(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	s.Delete(context.Background(), defaultUniverseId)

	assert.Equal(1, mc.tx.closeCalled)
}

func TestUniverseService_Delete_WhenCreatingTransactionFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{}
	mc := &mockConnectionPool{
		err: errDefault,
	}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUniverseId)

	assert.Equal(errDefault, err)
}

func TestUniverseService_Delete_DeletesTheRightUniverse(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	s.Delete(context.Background(), defaultUniverseId)

	assert.Equal(defaultUniverseId, mur.deleteId)
}

func TestUniverseService_Delete_WhenUniverseRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUniverseId)

	assert.Equal(errDefault, err)
}

func TestUniverseService_Delete_WhenRepositoriesSucceeds_ExpectSuccess(t *testing.T) {
	assert := assert.New(t)

	mur := &mockUniverseRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Universe: mur,
	}
	s := NewUniverseService(Config{}, mc, repos)

	err := s.Delete(context.Background(), defaultUniverseId)

	assert.Nil(err)
}
