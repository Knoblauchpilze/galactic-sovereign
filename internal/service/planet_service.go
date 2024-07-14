package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type PlanetService interface {
	Create(ctx context.Context, planetDto communication.PlanetDtoRequest) (communication.PlanetDtoResponse, error)
	Get(ctx context.Context, id uuid.UUID) (communication.PlanetDtoResponse, error)
	List(ctx context.Context) ([]communication.PlanetDtoResponse, error)
	ListForPlayer(ctx context.Context, player uuid.UUID) ([]communication.PlanetDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type planetServiceImpl struct {
	conn db.ConnectionPool

	planetRepo repositories.PlanetRepository
}

func NewPlanetService(conn db.ConnectionPool, repos repositories.Repositories) PlanetService {
	return &planetServiceImpl{
		conn:       conn,
		planetRepo: repos.Planet,
	}
}

func (s *planetServiceImpl) Create(ctx context.Context, planetDto communication.PlanetDtoRequest) (communication.PlanetDtoResponse, error) {
	planet := communication.FromPlanetDtoRequest(planetDto)

	createdPlanet, err := s.planetRepo.Create(ctx, planet)
	if err != nil {
		return communication.PlanetDtoResponse{}, err
	}

	out := communication.ToPlanetDtoResponse(createdPlanet)
	return out, nil
}

func (s *planetServiceImpl) Get(ctx context.Context, id uuid.UUID) (communication.PlanetDtoResponse, error) {
	planet, err := s.planetRepo.Get(ctx, id)
	if err != nil {
		return communication.PlanetDtoResponse{}, err
	}

	out := communication.ToPlanetDtoResponse(planet)
	return out, nil
}

func (s *planetServiceImpl) List(ctx context.Context) ([]communication.PlanetDtoResponse, error) {
	planets, err := s.planetRepo.List(ctx)
	if err != nil {
		return []communication.PlanetDtoResponse{}, err
	}

	var out []communication.PlanetDtoResponse
	for _, planet := range planets {
		dto := communication.ToPlanetDtoResponse(planet)
		out = append(out, dto)
	}

	return out, nil
}

func (s *planetServiceImpl) ListForPlayer(ctx context.Context, player uuid.UUID) ([]communication.PlanetDtoResponse, error) {
	planets, err := s.planetRepo.ListForPlayer(ctx, player)
	if err != nil {
		return []communication.PlanetDtoResponse{}, err
	}

	var out []communication.PlanetDtoResponse
	for _, planet := range planets {
		dto := communication.ToPlanetDtoResponse(planet)
		out = append(out, dto)
	}

	return out, nil
}

func (s *planetServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	return s.planetRepo.Delete(ctx, tx, id)
}
