package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UserDtoRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UserDtoResponse struct {
	Id       uuid.UUID
	Email    string
	Password string

	ApiKeys []uuid.UUID

	CreatedAt time.Time
}

func ToUserDtoResponse(user persistence.User, apiKeys []uuid.UUID) UserDtoResponse {
	if apiKeys == nil {
		apiKeys = make([]uuid.UUID, 0)
	}

	return UserDtoResponse{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,

		ApiKeys: apiKeys,

		CreatedAt: user.CreatedAt,
	}
}

func FromUserDtoRequest(user UserDtoRequest) persistence.User {
	t := time.Now()
	return persistence.User{
		Id:       uuid.New(),
		Email:    user.Email,
		Password: user.Password,

		CreatedAt: t,
		UpdatedAt: t,
	}
}
