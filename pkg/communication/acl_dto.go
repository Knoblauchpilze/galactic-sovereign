package communication

import (
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type AclDtoRequest struct {
	User uuid.UUID `json:"user"`

	Resource    string   `json:"resource"`
	Permissions []string `json:"permissions"`
}

type AclDtoResponse struct {
	Id   uuid.UUID `json:"id"`
	User uuid.UUID `json:"user"`

	Resource    string   `json:"resource"`
	Permissions []string `json:"permissions"`

	CreatedAt time.Time `json:"createdAt"`
}

func FromAclDtoRequest(acl AclDtoRequest) persistence.Acl {
	t := time.Now()
	return persistence.Acl{
		Id:   uuid.New(),
		User: acl.User,

		Resource:    acl.Resource,
		Permissions: acl.Permissions,

		CreatedAt: t,
		UpdatedAt: t,
	}
}

func ToAclDtoResponse(acl persistence.Acl) AclDtoResponse {
	return AclDtoResponse{
		Id:   acl.Id,
		User: acl.User,

		Resource:    acl.Resource,
		Permissions: acl.Permissions,

		CreatedAt: acl.CreatedAt,
	}
}

type AclResponseDto struct {
	Resource    string   `json:"resource"`
	Permissions []string `json:"permissions"`
}
