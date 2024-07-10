package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type PlayerService interface {
	Create(ctx context.Context, playerDto communication.PlayerDtoRequest) (communication.PlayerDtoResponse, error)
	Get(ctx context.Context, id uuid.UUID) (communication.PlayerDtoResponse, error)
	List(ctx context.Context) ([]communication.PlayerDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type playerServiceImpl struct {
	conn db.ConnectionPool

	playerRepo repositories.PlayerRepository
}

func NewPlayerService(config Config, conn db.ConnectionPool, repos repositories.Repositories) PlayerService {
	return &playerServiceImpl{
		conn:       conn,
		playerRepo: repos.Player,
	}
}

func (s *playerServiceImpl) Create(ctx context.Context, playerDto communication.PlayerDtoRequest) (communication.PlayerDtoResponse, error) {
	player := communication.FromPlayerDtoRequest(playerDto)

	createdPlayer, err := s.playerRepo.Create(ctx, player)
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

func (s *playerServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	return s.playerRepo.Delete(ctx, tx, id)
}
