package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
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
