package service

import (
	"context"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/communication"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/persistence"
	"github.com/Knoblauchpilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
)

type PlayerService interface {
	Create(ctx context.Context, playerDto communication.PlayerDtoRequest) (communication.PlayerDtoResponse, error)
	Get(ctx context.Context, id uuid.UUID) (communication.PlayerDtoResponse, error)
	List(ctx context.Context) ([]communication.PlayerDtoResponse, error)
	ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]communication.PlayerDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type playerServiceImpl struct {
	conn db.Connection

	buildingActionRepo repositories.BuildingActionRepository
	planetRepo         repositories.PlanetRepository
	playerRepo         repositories.PlayerRepository
}

func NewPlayerService(conn db.Connection, repos repositories.Repositories) PlayerService {
	return &playerServiceImpl{
		conn:               conn,
		buildingActionRepo: repos.BuildingAction,
		planetRepo:         repos.Planet,
		playerRepo:         repos.Player,
	}
}

func (s *playerServiceImpl) Create(ctx context.Context, playerDto communication.PlayerDtoRequest) (communication.PlayerDtoResponse, error) {
	player := communication.FromPlayerDtoRequest(playerDto)

	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return communication.PlayerDtoResponse{}, err
	}
	defer tx.Close(ctx)

	createdPlayer, err := s.playerRepo.Create(ctx, tx, player)
	if err != nil {
		return communication.PlayerDtoResponse{}, err
	}

	planet := persistence.Planet{
		Id:        uuid.New(),
		Player:    createdPlayer.Id,
		Name:      "homeworld",
		Homeworld: true,
		CreatedAt: createdPlayer.CreatedAt,
		UpdatedAt: createdPlayer.UpdatedAt,
	}

	_, err = s.planetRepo.Create(ctx, tx, planet)
	if err != nil {
		return communication.PlayerDtoResponse{}, err
	}

	out := communication.ToPlayerDtoResponse(createdPlayer)
	return out, nil
}

func (s *playerServiceImpl) Get(ctx context.Context, id uuid.UUID) (communication.PlayerDtoResponse, error) {
	player, err := s.playerRepo.Get(ctx, id)
	if err != nil {
		return communication.PlayerDtoResponse{}, err
	}

	out := communication.ToPlayerDtoResponse(player)
	return out, nil
}

func (s *playerServiceImpl) List(ctx context.Context) ([]communication.PlayerDtoResponse, error) {
	players, err := s.playerRepo.List(ctx)
	if err != nil {
		return []communication.PlayerDtoResponse{}, err
	}

	var out []communication.PlayerDtoResponse
	for _, player := range players {
		dto := communication.ToPlayerDtoResponse(player)
		out = append(out, dto)
	}

	return out, nil
}

func (s *playerServiceImpl) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]communication.PlayerDtoResponse, error) {
	players, err := s.playerRepo.ListForApiUser(ctx, apiUser)
	if err != nil {
		return []communication.PlayerDtoResponse{}, err
	}

	var out []communication.PlayerDtoResponse
	for _, player := range players {
		dto := communication.ToPlayerDtoResponse(player)
		out = append(out, dto)
	}

	return out, nil
}

func (s *playerServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	err = s.buildingActionRepo.DeleteForPlayer(ctx, tx, id)
	if err != nil {
		return err
	}

	err = s.planetRepo.DeleteForPlayer(ctx, tx, id)
	if err != nil {
		return err
	}

	return s.playerRepo.Delete(ctx, tx, id)
}
