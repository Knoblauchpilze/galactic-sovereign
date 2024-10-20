package communication

import (
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/google/uuid"
)

type LimitDtoRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LimitDtoResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func FromLimitDtoRequest(limit LimitDtoRequest) persistence.Limit {
	t := time.Now()
	return persistence.Limit{
		Id:    uuid.New(),
		Name:  limit.Name,
		Value: limit.Value,

		CreatedAt: t,
		UpdatedAt: t,
	}
}

func ToLimitDtoResponse(limit persistence.Limit) LimitDtoResponse {
	return LimitDtoResponse{
		Name:  limit.Name,
		Value: limit.Value,
	}
}

type UserLimitDtoRequest struct {
	Name string    `json:"name"`
	User uuid.UUID `json:"user"`

	Limits []LimitDtoRequest `json:"limits"`
}

type UserLimitDtoResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	User uuid.UUID `json:"user"`

	Limits []LimitDtoResponse `json:"limits"`

	CreatedAt time.Time `json:"createdAt"`
}

func FromUserLimitDtoRequest(userLimit UserLimitDtoRequest) persistence.UserLimit {
	var limits []persistence.Limit
	for _, limit := range userLimit.Limits {
		limits = append(limits, FromLimitDtoRequest(limit))
	}

	t := time.Now()
	return persistence.UserLimit{
		Id:   uuid.New(),
		Name: userLimit.Name,
		User: userLimit.User,

		Limits: limits,

		CreatedAt: t,
		UpdatedAt: t,
	}
}

func ToUserLimitDtoResponse(userLimit persistence.UserLimit) UserLimitDtoResponse {
	var limitsResponseDtos []LimitDtoResponse
	for _, limit := range userLimit.Limits {
		limitsResponseDtos = append(limitsResponseDtos, ToLimitDtoResponse(limit))
	}

	return UserLimitDtoResponse{
		Id:   userLimit.Id,
		Name: userLimit.Name,
		User: userLimit.User,

		Limits: limitsResponseDtos,

		CreatedAt: userLimit.CreatedAt,
	}
}
