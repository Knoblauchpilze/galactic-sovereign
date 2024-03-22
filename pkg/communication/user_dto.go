package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UserDto struct {
	Id       uuid.UUID
	Email    string
	Password string

	CreatedAt time.Time
}

func FromUser(user persistence.User) UserDto {
	return UserDto{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,

		CreatedAt: user.CreatedAt,
	}
}
