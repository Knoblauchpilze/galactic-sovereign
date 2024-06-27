package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type AuthService interface {
	Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationDtoResponse, error)
}

type authServiceImpl struct {
	conn       db.ConnectionPool
	aclRepo    repositories.AclRepository
	apiKeyRepo repositories.ApiKeyRepository
	userLimit  repositories.UserLimitRepository
}

func NewAuthService(conn db.ConnectionPool, repos repositories.Repositories) AuthService {
	return &authServiceImpl{
		conn:       conn,
		aclRepo:    repos.Acl,
		apiKeyRepo: repos.ApiKey,
		userLimit:  repos.UserLimit,
	}
}

func (s *authServiceImpl) Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationDtoResponse, error) {
	var out communication.AuthorizationDtoResponse

	key, err := s.apiKeyRepo.GetForKey(ctx, apiKey)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return out, errors.NewCode(UserNotAuthenticated)
		}

		return out, err
	}

	if key.ValidUntil.Before(time.Now()) {
		return out, errors.NewCode(AuthenticationExpired)
	}

	tx, err := s.conn.StartTransaction(ctx)
	if err != nil {
		return out, err
	}
	defer tx.Close(ctx)

	ids, err := s.aclRepo.GetForUser(ctx, tx, key.ApiUser)
	if err != nil {
		return out, err
	}

	for _, id := range ids {
		acl, err := s.aclRepo.Get(ctx, tx, id)
		if err != nil {
			return out, err
		}

		dto := communication.ToAclDtoResponse(acl)
		out.Acls = append(out.Acls, dto)
	}

	ids, err = s.userLimit.GetForUser(ctx, tx, key.ApiUser)
	if err != nil {
		return out, err
	}

	for _, id := range ids {
		userLimit, err := s.userLimit.Get(ctx, tx, id)
		if err != nil {
			return out, err
		}

		dto := communication.ToUserLimitDtoResponse(userLimit)
		out.Limits = append(out.Limits, dto.Limits...)
	}

	return out, nil
}
