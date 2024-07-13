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

var defaultPlayerId = uuid.MustParse("f16455b7-514c-44b1-847f-ba2cf99c765b")
var defaultPlayerName = "my-player"

var defaultPlayerDtoRequest = communication.PlayerDtoRequest{
	ApiUser:  defaultUserId,
	Universe: defaultUniverseId,
	Name:     defaultPlayerName,
}
var defaultPlayer = persistence.Player{
	Id:       defaultPlayerId,
	ApiUser:  defaultUserId,
	Universe: defaultUniverseId,
	Name:     defaultPlayerName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func TestPlayerService_Create_CallsRepositoryCreate(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	s.Create(context.Background(), defaultPlayerDtoRequest)

	assert.Equal(1, mpr.createCalled)
	assert.Equal(defaultPlayerDtoRequest.ApiUser, mpr.createdPlayer.ApiUser)
	assert.Equal(defaultPlayerDtoRequest.Universe, mpr.createdPlayer.Universe)
	assert.Equal(defaultPlayerDtoRequest.Name, mpr.createdPlayer.Name)
}

func TestPlayerService_Create_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	_, err := s.Create(context.Background(), defaultPlayerDtoRequest)

	assert.Equal(errDefault, err)
}

func TestPlayerService_Create_ReturnsCreatedPlayer(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{
		player: defaultPlayer,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	actual, err := s.Create(context.Background(), defaultPlayerDtoRequest)

	assert.Nil(err)

	expected := communication.PlayerDtoResponse{
		Id:       defaultPlayer.Id,
		ApiUser:  defaultPlayer.ApiUser,
		Universe: defaultPlayer.Universe,
		Name:     defaultPlayer.Name,

		CreatedAt: defaultPlayer.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestPlayerService_Get_CallsRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	s.Get(context.Background(), defaultPlayerId)

	assert.Equal(1, mpr.getCalled)
}

func TestPlayerService_Get_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	_, err := s.Get(context.Background(), defaultPlayerId)

	assert.Equal(errDefault, err)
}

func TestPlayerService_Get_ReturnsPlayer(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{
		player: defaultPlayer,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	actual, err := s.Get(context.Background(), defaultPlayerId)

	assert.Nil(err)
	assert.Equal(defaultPlayerId, mpr.getId)

	expected := communication.PlayerDtoResponse{
		Id:       defaultPlayer.Id,
		ApiUser:  defaultPlayer.ApiUser,
		Universe: defaultPlayer.Universe,
		Name:     defaultPlayer.Name,

		CreatedAt: defaultPlayer.CreatedAt,
	}
	assert.Equal(expected, actual)
}

func TestPlayerService_List_CallsRepositoryList(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	s.List(context.Background())

	assert.Equal(1, mpr.listCalled)
}

func TestPlayerService_List_WhenRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	_, err := s.List(context.Background())

	assert.Equal(errDefault, err)
}

func TestPlayerService_List_ReturnsAllPlayers(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{
		player: defaultPlayer,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	actual, err := s.List(context.Background())

	assert.Nil(err)
	expected := []communication.PlayerDtoResponse{
		{
			Id:        defaultPlayer.Id,
			ApiUser:   defaultPlayer.ApiUser,
			Universe:  defaultPlayer.Universe,
			Name:      defaultPlayer.Name,
			CreatedAt: defaultPlayer.CreatedAt,
		},
	}
	assert.Equal(expected, actual)
}

func TestPlayerService_Delete_CallsRepositoryDelete(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	s.Delete(context.Background(), defaultPlayerId)

	assert.Equal(1, mpr.deleteCalled)
}

func TestPlayerService_Delete_CallsTransactionClose(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	s.Delete(context.Background(), defaultPlayerId)

	assert.Equal(1, mc.tx.closeCalled)
}

func TestPlayerService_Delete_WhenCreatingTransactionFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{}
	mc := &mockConnectionPool{
		err: errDefault,
	}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	err := s.Delete(context.Background(), defaultPlayerId)

	assert.Equal(errDefault, err)
}

func TestPlayerService_Delete_DeletesTheRightPlayer(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	s.Delete(context.Background(), defaultPlayerId)

	assert.Equal(defaultPlayerId, mpr.deleteId)
}

func TestPlayerService_Delete_WhenPlayerRepositoryFails_ExpectError(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{
		err: errDefault,
	}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	err := s.Delete(context.Background(), defaultPlayerId)

	assert.Equal(errDefault, err)
}

func TestPlayerService_Delete_WhenRepositoriesSucceeds_ExpectSuccess(t *testing.T) {
	assert := assert.New(t)

	mpr := &mockPlayerRepository{}
	mc := &mockConnectionPool{}
	repos := repositories.Repositories{
		Player: mpr,
	}
	s := NewPlayerService(mc, repos)

	err := s.Delete(context.Background(), defaultPlayerId)

	assert.Nil(err)
}
