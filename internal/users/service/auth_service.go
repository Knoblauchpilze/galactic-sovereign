package service

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type AuthService interface {
	Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationResponseDto, error)
}

type authServiceImpl struct {
	conn       db.ConnectionPool
	userRepo   repositories.UserRepository
	apiKeyRepo repositories.ApiKeyRepository
}

func NewAuthService(conn db.ConnectionPool, userRepo repositories.UserRepository, apiKeyRepo repositories.ApiKeyRepository) AuthService {
	return &authServiceImpl{
		conn:       conn,
		userRepo:   userRepo,
		apiKeyRepo: apiKeyRepo,
	}
}

func (s *authServiceImpl) Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationResponseDto, error) {
	out := communication.AuthorizationResponseDto{
		Acls: []communication.AclResponseDto{
			{
				Resource:   "r1",
				Permission: "DELETE",
			},
		},
		Limits: []communication.LimitResponseDto{
			{
				Name:  "l1",
				Value: "-1",
			},
		},
	}

	return out, errors.NewCode(errors.NotImplementedCode)
}
