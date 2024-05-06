package service

import (
	"context"
	"time"

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
	Login(ctx context.Context, id uuid.UUID) (communication.ApiKeyDtoResponse, error)
	Logout(ctx context.Context, id uuid.UUID) error
}

type userServiceImpl struct {
	conn       db.ConnectionPool
	userRepo   repositories.UserRepository
	apiKeyRepo repositories.ApiKeyRepository

	apiKeyValidity time.Duration
}

func NewUserService(config Config, conn db.ConnectionPool, userRepo repositories.UserRepository, apiKeyRepo repositories.ApiKeyRepository) UserService {
	return &userServiceImpl{
		conn:       conn,
		userRepo:   userRepo,
		apiKeyRepo: apiKeyRepo,

		apiKeyValidity: config.ApiKeyValidity,
	}
}

func (s *userServiceImpl) Create(ctx context.Context, userDto communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	user := communication.FromUserDtoRequest(userDto)

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	out := communication.ToUserDtoResponse(createdUser)
	return out, nil
}

func (s *userServiceImpl) Get(ctx context.Context, id uuid.UUID) (communication.UserDtoResponse, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	out := communication.ToUserDtoResponse(user)
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

	out := communication.ToUserDtoResponse(updated)
	return out, nil
}

func (s *userServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	apiKeys, err := s.apiKeyRepo.GetForUserTx(ctx, tx, id)
	if err != nil {
		return err
	}

	err = s.apiKeyRepo.DeleteTx(ctx, tx, apiKeys)
	if err != nil {
		return err
	}
	err = s.userRepo.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) Login(ctx context.Context, id uuid.UUID) (communication.ApiKeyDtoResponse, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return communication.ApiKeyDtoResponse{}, err
	}

	apiKey := persistence.ApiKey{
		Id:         uuid.New(),
		Key:        uuid.New(),
		ApiUser:    user.Id,
		ValidUntil: time.Now().Add(s.apiKeyValidity),
	}

	createdKey, err := s.apiKeyRepo.Create(ctx, apiKey)
	if err != nil {
		return communication.ApiKeyDtoResponse{}, err
	}

	out := communication.ToApiKeyDtoResponse(createdKey)
	return out, nil
}

func (s *userServiceImpl) Logout(ctx context.Context, id uuid.UUID) error {
	_, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	apiKeys, err := s.apiKeyRepo.GetForUser(ctx, id)
	if err != nil {
		return err
	}

	if len(apiKeys) == 0 {
		return nil
	}

	return s.apiKeyRepo.Delete(ctx, apiKeys)
}
