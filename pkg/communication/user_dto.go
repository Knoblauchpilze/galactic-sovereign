package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UserDtoRequest struct {
	// https://stackoverflow.com/questions/18635671/how-to-define-multiple-name-tags-in-a-struct
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UserDtoResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`

	CreatedAt time.Time `json:"createdAt"`
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

func ToUserDtoResponse(user persistence.User) UserDtoResponse {
	return UserDtoResponse{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,

		CreatedAt: user.CreatedAt,
	}
}
