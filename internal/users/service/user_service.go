package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, user communication.UserDtoRequest) (communication.UserDtoResponse, error)
	Get(ctx context.Context, id uuid.UUID) (communication.UserDtoResponse, error)
	List(ctx context.Context) ([]uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, user communication.UserDtoRequest) (communication.UserDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userServiceImpl struct {
	conn       db.ConnectionPool
	userRepo   repositories.UserRepository
	apiKeyRepo repositories.ApiKeyRepository
}

func NewUserService(conn db.ConnectionPool, userRepo repositories.UserRepository, apiKeyRepo repositories.ApiKeyRepository) UserService {
	return &userServiceImpl{
		conn:       conn,
		userRepo:   userRepo,
		apiKeyRepo: apiKeyRepo,
	}
}

func (s *userServiceImpl) Create(ctx context.Context, userDto communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	user := communication.FromUserDtoRequest(userDto)

	apiKey := persistence.ApiKey{
		// TODO: Reactivate this to force a failure in the transaction
		//Id:      uuid.MustParse("d5826856-63f8-41c8-a643-8bb80d5feb78"),
		Id:      uuid.New(),
		Key:     uuid.New(),
		ApiUser: user.Id,
		Enabled: true,
	}

	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}
	defer tx.Close(ctx)

	createdUser, err := s.userRepo.TransactionalCreate(ctx, tx, user)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}
	createdKey, err := s.apiKeyRepo.TransactionalCreate(ctx, tx, apiKey)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	out := communication.ToUserDtoResponse(createdUser, []uuid.UUID{createdKey.Key})
	return out, nil
}

func (s *userServiceImpl) Get(ctx context.Context, id uuid.UUID) (communication.UserDtoResponse, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	// TODO: Fetch the API keys
	out := communication.ToUserDtoResponse(user, []uuid.UUID{})
	return out, nil
}

func (s *userServiceImpl) List(ctx context.Context) ([]uuid.UUID, error) {
	return s.userRepo.List(ctx)
}

func (s *userServiceImpl) Update(ctx context.Context, id uuid.UUID, userDto communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	user.Email = userDto.Email
	user.Password = userDto.Password

	updated, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	// TODO: Fetch the API keys
	out := communication.ToUserDtoResponse(updated, []uuid.UUID{})
	return out, nil
}

func (s *userServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
